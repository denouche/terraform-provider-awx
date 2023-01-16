/*
Create a workflow job template.

Example Usage

```hcl
data "awx_inventory" "default" {
    name = "default"
}

resource "awx_workflow_job_template" "default" {
    name            = "workflow-job"
    organization_id = var.organization_id
    inventory_id    = data.awx_inventory.default.id
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

func resourceWorkflowJobTemplate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceWorkflowJobTemplateCreate,
		ReadContext:   resourceWorkflowJobTemplateRead,
		UpdateContext: resourceWorkflowJobTemplateUpdate,
		DeleteContext: resourceWorkflowJobTemplateDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of this workflow job template. (string, required)",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Optional description of this workflow job template.",
			},
			"variables": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "",
				StateFunc:   normalizeJsonYaml,
			},
			"organization_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "The organization used to determine access to this template. (id, default=``)",
			},
			"survey_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"allow_simultaneous": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"ask_variables_on_launch": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"inventory_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Inventory applied as a prompt, assuming job template prompts for inventory. (id, default=``)",
				Default:     "",
			},
			"limit": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"scm_branch": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"ask_inventory_on_launch": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"ask_scm_branch_on_launch": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"ask_limit_on_launch": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"webhook_service": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"webhook_credential": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
		},
	}
}

func resourceWorkflowJobTemplateCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*awx.AWX)
	awxService := client.WorkflowJobTemplateService

	result, err := awxService.CreateWorkflowJobTemplate(map[string]interface{}{
		"name":                     d.Get("name").(string),
		"description":              d.Get("description").(string),
		"organization":             d.Get("organization_id").(int),
		"inventory":                AtoipOr(d.Get("inventory_id").(string), nil),
		"extra_vars":               d.Get("variables").(string),
		"survey_enabled":           d.Get("survey_enabled").(bool),
		"allow_simultaneous":       d.Get("allow_simultaneous").(bool),
		"ask_variables_on_launch":  d.Get("ask_variables_on_launch").(bool),
		"limit":                    d.Get("limit").(string),
		"scm_branch":               d.Get("scm_branch").(string),
		"ask_inventory_on_launch":  d.Get("ask_inventory_on_launch").(bool),
		"ask_scm_branch_on_launch": d.Get("ask_scm_branch_on_launch").(bool),
		"ask_limit_on_launch":      d.Get("ask_limit_on_launch").(bool),
		"webhook_service":          d.Get("webhook_service").(string),
		"webhook_credential":       d.Get("webhook_credential").(string),
	}, map[string]string{})
	if err != nil {
		log.Printf("Fail to Create Template %v", err)
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create WorkflowJobTemplate",
			Detail:   fmt.Sprintf("WorkflowJobTemplate with name %s failed to create %s", d.Get("name").(string), err.Error()),
		})
		return diags
	}

	d.SetId(strconv.Itoa(result.ID))
	return resourceWorkflowJobTemplateRead(ctx, d, m)
}

func resourceWorkflowJobTemplateUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*awx.AWX)
	awxService := client.WorkflowJobTemplateService
	id, diags := convertStateIDToNummeric("Update WorkflowJobTemplate", d)
	if diags.HasError() {
		return diags
	}

	params := make(map[string]string)
	_, err := awxService.GetWorkflowJobTemplateByID(id, params)
	if err != nil {
		return buildDiagNotFoundFail("job Workflow template", id, err)
	}

	_, err = awxService.UpdateWorkflowJobTemplate(id, map[string]interface{}{
		"name":                     d.Get("name").(string),
		"description":              d.Get("description").(string),
		"organization":             d.Get("organization_id").(int),
		"inventory":                AtoipOr(d.Get("inventory_id").(string), nil),
		"extra_vars":               d.Get("variables").(string),
		"survey_enabled":           d.Get("survey_enabled").(bool),
		"allow_simultaneous":       d.Get("allow_simultaneous").(bool),
		"ask_variables_on_launch":  d.Get("ask_variables_on_launch").(bool),
		"limit":                    d.Get("limit").(string),
		"scm_branch":               d.Get("scm_branch").(string),
		"ask_inventory_on_launch":  d.Get("ask_inventory_on_launch").(bool),
		"ask_scm_branch_on_launch": d.Get("ask_scm_branch_on_launch").(bool),
		"ask_limit_on_launch":      d.Get("ask_limit_on_launch").(bool),
		"webhook_service":          d.Get("webhook_service").(string),
		"webhook_credential":       d.Get("webhook_credential").(string),
	}, map[string]string{})
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to update WorkflowJobTemplate",
			Detail:   fmt.Sprintf("WorkflowJobTemplate with name %s in the project id %d failed to update %s", d.Get("name").(string), d.Get("project_id").(int), err.Error()),
		})
		return diags
	}

	return resourceWorkflowJobTemplateRead(ctx, d, m)
}

func resourceWorkflowJobTemplateRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*awx.AWX)
	awxService := client.WorkflowJobTemplateService
	id, diags := convertStateIDToNummeric("Read WorkflowJobTemplate", d)
	if diags.HasError() {
		return diags
	}

	res, err := awxService.GetWorkflowJobTemplateByID(id, make(map[string]string))
	if err != nil {
		return buildDiagNotFoundFail("workflow job template", id, err)

	}
	d = setWorkflowJobTemplateResourceData(d, res)
	return nil
}

func resourceWorkflowJobTemplateDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*awx.AWX)
	awxService := client.WorkflowJobTemplateService
	id, diags := convertStateIDToNummeric(diagElementHostTitle, d)
	if diags.HasError() {
		return diags
	}

	if _, err := awxService.DeleteWorkflowJobTemplate(id); err != nil {
		return buildDiagDeleteFail(
			diagElementHostTitle,
			fmt.Sprintf("id %v, got %s ",
				id, err.Error()))
	}
	d.SetId("")
	return nil
}

func setWorkflowJobTemplateResourceData(d *schema.ResourceData, r *awx.WorkflowJobTemplate) *schema.ResourceData {

	d.Set("name", r.Name)
	d.Set("description", r.Description)
	d.Set("organization_id", strconv.Itoa(r.Organization))
	d.Set("inventory_id", strconv.Itoa(r.Inventory))
	d.Set("survey_enabled", r.SurveyEnabled)
	d.Set("allow_simultaneous", r.AllowSimultaneous)
	d.Set("ask_variables_on_launch", r.AskVariablesOnLaunch)
	d.Set("limit", r.Limit)
	d.Set("scm_branch", r.ScmBranch)
	d.Set("ask_inventory_on_launch", r.AskInventoryOnLaunch)
	d.Set("ask_scm_branch_on_launch", r.AskScmBranchOnLaunch)
	d.Set("ask_limit_on_launch", r.AskLimitOnLaunch)
	d.Set("webhook_service", r.WebhookService)
	d.Set("webhook_credential", r.WebhookCredential)
	d.Set("variables", normalizeJsonYaml(r.ExtraVars))

	d.SetId(strconv.Itoa(r.ID))
	return d
}
