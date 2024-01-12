/*
*TBD*

Example Usage

```hcl
resource "awx_workflow_job_template_schedule" "default" {
  workflow_job_template_id      = awx_workflow_job_template.default.id

  name                      = "schedule-test"
  rrule                     = "DTSTART;TZID=Europe/Paris:20211214T120000 RRULE:INTERVAL=1;FREQ=DAILY"
  extra_data                = <<EOL
organization_name: testorg
EOL
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

func resourceWorkflowJobTemplateSchedule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceWorkflowJobTemplateScheduleCreate,
		ReadContext:   resourceWorkflowJobTemplateScheduleRead,
		UpdateContext: resourceWorkflowJobTemplateScheduleUpdate,
		DeleteContext: resourceScheduleDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"rrule": {
				Type:     schema.TypeString,
				Required: true,
			},
			"workflow_job_template_id": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "The workflow_job_template id for this schedule",
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"inventory": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Inventory applied as a prompt, assuming job template prompts for inventory (id, default=``)",
			},
			"extra_data": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Extra data to be pass for the schedule (YAML format)",
			},
		},
	}
}

func resourceWorkflowJobTemplateScheduleCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*awx.AWX)
	awxService := client.WorkflowJobTemplateScheduleService

	workflowJobTemplateID := d.Get("workflow_job_template_id").(int)

	result, err := awxService.CreateWorkflowJobTemplateSchedule(workflowJobTemplateID, map[string]interface{}{
		"name":        d.Get("name").(string),
		"rrule":       d.Get("rrule").(string),
		"description": d.Get("description").(string),
		"enabled":     d.Get("enabled").(bool),
		"inventory":   AtoipOr(d.Get("inventory").(string), nil),
		"extra_data":  unmarshalYaml(d.Get("extra_data").(string)),
	}, map[string]string{})
	if err != nil {
		log.Printf("Fail to Create Schedule for WorkflowJobTemplate %d: %v", workflowJobTemplateID, err)
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create Schedule",
			Detail:   fmt.Sprintf("Schedule failed to create %s", err.Error()),
		})
		return diags
	}

	d.SetId(strconv.Itoa(result.ID))
	return resourceScheduleRead(ctx, d, m)
}

func resourceWorkflowJobTemplateScheduleUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*awx.AWX)
	awxService := client.ScheduleService
	id, diags := convertStateIDToNummeric("Update Schedule", d)
	if diags.HasError() {
		return diags
	}

	params := make(map[string]string)
	_, err := awxService.GetByID(id, params)
	if err != nil {
		return buildDiagNotFoundFail("schedule", id, err)
	}

	_, err = awxService.Update(id, map[string]interface{}{
		"name":                  d.Get("name").(string),
		"rrule":                 d.Get("rrule").(string),
		"workflow_job_template": d.Get("workflow_job_template_id").(int),
		"description":           d.Get("description").(string),
		"enabled":               d.Get("enabled").(bool),
		"inventory":             d.Get("inventory").(string),
		"extra_data":            unmarshalYaml(d.Get("extra_data").(string)),
	}, map[string]string{})
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to update Schedule",
			Detail:   fmt.Sprintf("Schedule with name %s failed to update %s", d.Get("name").(string), err.Error()),
		})
		return diags
	}

	return resourceScheduleRead(ctx, d, m)
}

func resourceWorkflowJobTemplateScheduleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*awx.AWX)
	awxService := client.ScheduleService
	id, diags := convertStateIDToNummeric("Read schedule", d)
	if diags.HasError() {
		return diags
	}

	res, err := awxService.GetByID(id, make(map[string]string))
	if err != nil {
		return buildDiagNotFoundFail("schedule", id, err)
	}
	d = setWorkflowJobTemplateScheduleResourceData(d, res)
	return nil
}

func setWorkflowJobTemplateScheduleResourceData(d *schema.ResourceData, r *awx.Schedule) *schema.ResourceData {
	d.Set("name", r.Name)
	d.Set("rrule", r.Rrule)
	// Map api to state
	d.Set("workflow_job_template_id", r.UnifiedJobTemplate)
	d.Set("description", r.Description)
	d.Set("enabled", r.Enabled)
	d.Set("inventory", r.Inventory)
	d.Set("extra_data", marshalYaml(r.ExtraData))
	d.SetId(strconv.Itoa(r.ID))
	return d
}
