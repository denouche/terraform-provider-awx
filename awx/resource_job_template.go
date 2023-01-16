/*
Create a job template.

Example Usage

```hcl
data "awx_inventory" "default" {
    name            = "my-inventory"
    organization_id = data.awx_organization.default.id
}

data "awx_project" "default" {
    name            = "my-project"
    organization_id = data.awx_organization.default.id
}

resource "awx_job_template" "baseconfig" {
    name           = "baseconfig"
    job_type       = "run"
    inventory_id   = data.awx_inventory.default.id
    project_id     = data.awx_project.default.id
    playbook       = "master-configure-system.yml"
}
```

*/

package awx

import (
	"context"
	"fmt"
	"log"
	"strconv"

	awx "github.com/denouche/goawx/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceJobTemplate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceJobTemplateCreate,
		ReadContext:   resourceJobTemplateRead,
		UpdateContext: resourceJobTemplateUpdate,
		DeleteContext: resourceJobTemplateDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			// Run, Check, Scan
			"job_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "One of: run, check, scan",
			},
			"inventory_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"project_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"playbook": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"forks": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},
			"limit": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			//0,1,2,3,4,5
			"verbosity": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "One of 0,1,2,3,4,5",
			},
			"extra_vars": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"job_tags": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"force_handlers": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"skip_tags": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"start_at_task": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"timeout": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},
			"use_fact_cache": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"host_config_key": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"ask_diff_mode_on_launch": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"ask_limit_on_launch": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"ask_tags_on_launch": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"ask_verbosity_on_launch": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"ask_inventory_on_launch": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"ask_variables_on_launch": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"ask_credential_on_launch": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"survey_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"become_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"diff_mode": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"ask_skip_tags_on_launch": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"allow_simultaneous": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"custom_virtualenv": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"ask_job_type_on_launch": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"execution_environment": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
		},
	}
}

func resourceJobTemplateCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*awx.AWX)
	awxService := client.JobTemplateService

	result, err := awxService.CreateJobTemplate(map[string]interface{}{
		"name":                     d.Get("name").(string),
		"description":              d.Get("description").(string),
		"job_type":                 d.Get("job_type").(string),
		"inventory":                AtoipOr(d.Get("inventory_id").(string), nil),
		"project":                  d.Get("project_id").(int),
		"playbook":                 d.Get("playbook").(string),
		"forks":                    d.Get("forks").(int),
		"limit":                    d.Get("limit").(string),
		"verbosity":                d.Get("verbosity").(int),
		"extra_vars":               d.Get("extra_vars").(string),
		"job_tags":                 d.Get("job_tags").(string),
		"force_handlers":           d.Get("force_handlers").(bool),
		"skip_tags":                d.Get("skip_tags").(string),
		"start_at_task":            d.Get("start_at_task").(string),
		"timeout":                  d.Get("timeout").(int),
		"use_fact_cache":           d.Get("use_fact_cache").(bool),
		"host_config_key":          d.Get("host_config_key").(string),
		"ask_diff_mode_on_launch":  d.Get("ask_diff_mode_on_launch").(bool),
		"ask_variables_on_launch":  d.Get("ask_variables_on_launch").(bool),
		"ask_limit_on_launch":      d.Get("ask_limit_on_launch").(bool),
		"ask_tags_on_launch":       d.Get("ask_tags_on_launch").(bool),
		"ask_skip_tags_on_launch":  d.Get("ask_skip_tags_on_launch").(bool),
		"ask_job_type_on_launch":   d.Get("ask_job_type_on_launch").(bool),
		"ask_verbosity_on_launch":  d.Get("ask_verbosity_on_launch").(bool),
		"ask_inventory_on_launch":  d.Get("ask_inventory_on_launch").(bool),
		"ask_credential_on_launch": d.Get("ask_credential_on_launch").(bool),
		"survey_enabled":           d.Get("survey_enabled").(bool),
		"become_enabled":           d.Get("become_enabled").(bool),
		"diff_mode":                d.Get("diff_mode").(bool),
		"allow_simultaneous":       d.Get("allow_simultaneous").(bool),
		"custom_virtualenv":        AtoipOr(d.Get("custom_virtualenv").(string), nil),
		"execution_environment":    AtoipOr(d.Get("execution_environment").(string), nil),
	}, map[string]string{})
	if err != nil {
		log.Printf("Fail to Create Template %v", err)
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create JobTemplate",
			Detail:   fmt.Sprintf("JobTemplate with name %s in the project id %d, failed to create %s", d.Get("name").(string), d.Get("project_id").(int), err.Error()),
		})
		return diags
	}

	d.SetId(strconv.Itoa(result.ID))
	return resourceJobTemplateRead(ctx, d, m)
}

func resourceJobTemplateUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*awx.AWX)
	awxService := client.JobTemplateService
	id, diags := convertStateIDToNummeric("Update JobTemplate", d)
	if diags.HasError() {
		return diags
	}

	params := make(map[string]string)
	_, err := awxService.GetJobTemplateByID(id, params)
	if err != nil {
		return buildDiagNotFoundFail("job template", id, err)
	}

	_, err = awxService.UpdateJobTemplate(id, map[string]interface{}{
		"name":                     d.Get("name").(string),
		"description":              d.Get("description").(string),
		"job_type":                 d.Get("job_type").(string),
		"inventory":                AtoipOr(d.Get("inventory_id").(string), nil),
		"project":                  d.Get("project_id").(int),
		"playbook":                 d.Get("playbook").(string),
		"forks":                    d.Get("forks").(int),
		"limit":                    d.Get("limit").(string),
		"verbosity":                d.Get("verbosity").(int),
		"extra_vars":               d.Get("extra_vars").(string),
		"job_tags":                 d.Get("job_tags").(string),
		"force_handlers":           d.Get("force_handlers").(bool),
		"skip_tags":                d.Get("skip_tags").(string),
		"start_at_task":            d.Get("start_at_task").(string),
		"timeout":                  d.Get("timeout").(int),
		"use_fact_cache":           d.Get("use_fact_cache").(bool),
		"host_config_key":          d.Get("host_config_key").(string),
		"ask_diff_mode_on_launch":  d.Get("ask_diff_mode_on_launch").(bool),
		"ask_variables_on_launch":  d.Get("ask_variables_on_launch").(bool),
		"ask_limit_on_launch":      d.Get("ask_limit_on_launch").(bool),
		"ask_tags_on_launch":       d.Get("ask_tags_on_launch").(bool),
		"ask_skip_tags_on_launch":  d.Get("ask_skip_tags_on_launch").(bool),
		"ask_job_type_on_launch":   d.Get("ask_job_type_on_launch").(bool),
		"ask_verbosity_on_launch":  d.Get("ask_verbosity_on_launch").(bool),
		"ask_inventory_on_launch":  d.Get("ask_inventory_on_launch").(bool),
		"ask_credential_on_launch": d.Get("ask_credential_on_launch").(bool),
		"survey_enabled":           d.Get("survey_enabled").(bool),
		"become_enabled":           d.Get("become_enabled").(bool),
		"diff_mode":                d.Get("diff_mode").(bool),
		"allow_simultaneous":       d.Get("allow_simultaneous").(bool),
		"custom_virtualenv":        AtoipOr(d.Get("custom_virtualenv").(string), nil),
		"execution_environment":    AtoipOr(d.Get("execution_environment").(string), nil),
	}, map[string]string{})
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to update JobTemplate",
			Detail:   fmt.Sprintf("JobTemplate with name %s in the project id %d failed to update %s", d.Get("name").(string), d.Get("project_id").(int), err.Error()),
		})
		return diags
	}

	return resourceJobTemplateRead(ctx, d, m)
}

func resourceJobTemplateRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*awx.AWX)
	awxService := client.JobTemplateService
	id, diags := convertStateIDToNummeric("Read JobTemplate", d)
	if diags.HasError() {
		return diags
	}

	res, err := awxService.GetJobTemplateByID(id, make(map[string]string))
	if err != nil {
		return buildDiagNotFoundFail("job template", id, err)

	}
	d = setJobTemplateResourceData(d, res)
	return nil
}

func setJobTemplateResourceData(d *schema.ResourceData, r *awx.JobTemplate) *schema.ResourceData {
	d.Set("allow_simultaneous", r.AllowSimultaneous)
	d.Set("ask_credential_on_launch", r.AskCredentialOnLaunch)
	d.Set("ask_job_type_on_launch", r.AskJobTypeOnLaunch)
	d.Set("ask_limit_on_launch", r.AskLimitOnLaunch)
	d.Set("ask_skip_tags_on_launch", r.AskSkipTagsOnLaunch)
	d.Set("ask_tags_on_launch", r.AskTagsOnLaunch)
	d.Set("ask_variables_on_launch", r.AskVariablesOnLaunch)
	d.Set("description", r.Description)
	d.Set("extra_vars", normalizeJsonYaml(r.ExtraVars))
	d.Set("force_handlers", r.ForceHandlers)
	d.Set("forks", r.Forks)
	d.Set("host_config_key", r.HostConfigKey)
	d.Set("inventory_id", r.Inventory)
	d.Set("job_tags", r.JobTags)
	d.Set("job_type", r.JobType)
	d.Set("diff_mode", r.DiffMode)
	d.Set("custom_virtualenv", r.CustomVirtualenv)
	d.Set("limit", r.Limit)
	d.Set("name", r.Name)
	d.Set("become_enabled", r.BecomeEnabled)
	d.Set("use_fact_cache", r.UseFactCache)
	d.Set("playbook", r.Playbook)
	d.Set("project_id", r.Project)
	d.Set("skip_tags", r.SkipTags)
	d.Set("start_at_task", r.StartAtTask)
	d.Set("survey_enabled", r.SurveyEnabled)
	d.Set("verbosity", r.Verbosity)
	d.SetId(strconv.Itoa(r.ID))
	return d
}
