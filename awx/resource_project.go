/*
Create a project.

Example Usage

```hcl
data "awx_organization" "default" {
    name = "Default"
}

resource "awx_project" "base_service_config" {
    name                 = "base-service-configuration"
    organization_id      = data.awx_organization.default.id
    scm_type             = "git"
    scm_url              = "https://github.com/nolte/ansible_playbook-baseline-online-server"
    scm_update_on_launch = true
}
```

*/

package awx

import (
	"context"
	"fmt"
	"strconv"
	"time"

	awx "github.com/denouche/goawx/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceProject() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceProjectCreate,
		ReadContext:   resourceProjectRead,
		DeleteContext: resourceProjectDelete,
		UpdateContext: resourceProjectUpdate,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of this project",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Optional description of this project.",
			},
			"local_path": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Local path (relative to PROJECTS_ROOT) containing playbooks and related files for this project.",
			},
			"scm_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "One of \"\" (manual), git, hg, svn",
			},
			"scm_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "",
			},
			"scm_credential_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Numeric ID of the scm used credential",
			},
			"scm_branch": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Specific branch, tag or commit to checkout.",
			},
			"scm_clean": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"scm_delete_on_update": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"organization_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Numeric ID of the project organization",
			},
			"scm_update_on_launch": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"scm_update_cache_timeout": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
	}
}

func resourceProjectCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*awx.AWX)
	awxService := client.ProjectService

	orgID := d.Get("organization_id").(int)
	projectName := d.Get("name").(string)
	_, res, err := awxService.ListProjects(map[string]string{
		"name":         projectName,
		"organization": strconv.Itoa(orgID),
	},
	)
	if err != nil {
		return buildDiagnosticsMessage("Create: Fail to find Project", "Fail to find Project %s Organization ID %v, %s", projectName, orgID, err.Error())
	}
	if len(res.Results) >= 1 {
		return buildDiagnosticsMessage("Create: Always exist", "Project with name %s  already exists in the Organization ID %v", projectName, orgID)
	}
	credentials := ""
	if d.Get("scm_credential_id").(int) > 0 {
		credentials = strconv.Itoa(d.Get("scm_credential_id").(int))
	}
	result, err := awxService.CreateProject(map[string]interface{}{
		"name":                 projectName,
		"description":          d.Get("description").(string),
		"local_path":           d.Get("local_path").(string),
		"scm_type":             d.Get("scm_type").(string),
		"scm_url":              d.Get("scm_url").(string),
		"scm_branch":           d.Get("scm_branch").(string),
		"scm_clean":            d.Get("scm_clean").(bool),
		"scm_delete_on_update": d.Get("scm_delete_on_update").(bool),
		"organization":         d.Get("organization_id").(int),
		"credential":           credentials,

		"scm_update_on_launch":     d.Get("scm_update_on_launch").(bool),
		"scm_update_cache_timeout": d.Get("scm_update_cache_timeout").(int),
	}, map[string]string{})
	if err != nil {
		return buildDiagnosticsMessage("Create: Project not created", "Project with name %s  in the Organization ID %v not created, %s", projectName, orgID, err.Error())
	}

	d.SetId(strconv.Itoa(result.ID))
	return resourceProjectRead(ctx, d, m)
}

func resourceProjectUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*awx.AWX)
	awxService := client.ProjectService

	id, diags := convertStateIDToNummeric("Update Project", d)
	if diags.HasError() {
		return diags
	}
	credentials := ""
	if d.Get("scm_credential_id").(int) > 0 {
		credentials = strconv.Itoa(d.Get("scm_credential_id").(int))
	}

	data := map[string]interface{}{
		"name":                     d.Get("name").(string),
		"description":              d.Get("description").(string),
		"scm_type":                 d.Get("scm_type").(string),
		"scm_url":                  d.Get("scm_url").(string),
		"scm_branch":               d.Get("scm_branch").(string),
		"scm_clean":                d.Get("scm_clean").(bool),
		"scm_delete_on_update":     d.Get("scm_delete_on_update").(bool),
		"credential":               credentials,
		"organization":             d.Get("organization_id").(int),
		"scm_update_on_launch":     d.Get("scm_update_on_launch").(bool),
		"scm_update_cache_timeout": d.Get("scm_update_cache_timeout").(int),
	}

	// Cannot change local_path for git-based projects
	if d.Get("local_path").(string) != "" && d.Get("scm_type").(string) != "git" {
		data["local_path"] = d.Get("local_path").(string)
	}

	_, err := awxService.UpdateProject(id, data, map[string]string{})
	if err != nil {
		return buildDiagnosticsMessage("Update: Fail To Update Project", "Fail to get Project with ID %v, got %s", id, err.Error())
	}
	return resourceProjectRead(ctx, d, m)
}

func resourceProjectRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*awx.AWX)
	awxService := client.ProjectService

	id, diags := convertStateIDToNummeric("Read Project", d)
	if diags.HasError() {
		return diags
	}

	res, err := awxService.GetProjectByID(id, make(map[string]string))
	if err != nil {
		return buildDiagNotFoundFail("project", id, err)
	}
	d = setProjectResourceData(d, res)
	return diags
}

func resourceProjectDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	digMessagePart := "Project"
	client := m.(*awx.AWX)
	awxService := client.ProjectService
	var jobID int
	var finished time.Time
	id, diags := convertStateIDToNummeric("Delete Project", d)
	if diags.HasError() {
		return diags
	}

	res, err := awxService.GetProjectByID(id, make(map[string]string))
	if err != nil {
		d.SetId("")
		return buildDiagNotFoundFail("project", id, err)
	}

	if res.SummaryFields.CurrentJob["id"] != nil {
		jobID = int(res.SummaryFields.CurrentJob["id"].(float64))
	} else if res.SummaryFields.LastJob["id"] != nil {
		jobID = int(res.SummaryFields.LastJob["id"].(float64))
	}
	if jobID != 0 {
		_, err = client.ProjectUpdatesService.ProjectUpdateCancel(jobID)
		if err != nil {
			return buildDiagnosticsMessage(
				"Delete: Fail to canel Job",
				"Fail to canel the Job %v for Project with ID %v, got %s",
				jobID, id, err.Error(),
			)
		}
	}
	// check if finished is 0
	for finished.IsZero() {
		prj, _ := client.ProjectUpdatesService.ProjectUpdateGet(jobID)
		finished = prj.Finished
		time.Sleep(1 * time.Second)
	}

	if _, err = awxService.DeleteProject(id); err != nil {
		return buildDiagDeleteFail(digMessagePart, fmt.Sprintf("ProjectID %v, got %s ", id, err.Error()))
	}
	d.SetId("")
	return diags
}

func setProjectResourceData(d *schema.ResourceData, r *awx.Project) *schema.ResourceData {
	d.Set("name", r.Name)
	d.Set("description", r.Description)
	d.Set("scm_type", r.ScmType)
	d.Set("scm_url", r.ScmURL)
	d.Set("scm_branch", r.ScmBranch)
	d.Set("scm_clean", r.ScmClean)
	d.Set("scm_delete_on_update", r.ScmDeleteOnUpdate)
	d.Set("organization_id", r.Organization)

	id, err := strconv.Atoi(r.Credential)
	if err == nil {
		d.Set("scm_credential_id", id)
	}
	d.Set("scm_update_on_launch", r.ScmUpdateOnLaunch)
	d.Set("scm_update_cache_timeout", r.ScmUpdateCacheTimeout)

	d.SetId(strconv.Itoa(r.ID))
	return d
}
