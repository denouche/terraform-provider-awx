/*
*TBD*

Example Usage

```hcl
resource "awx_schedule" "default" {
  name                      = "schedule-test"
  rrule                     = "DTSTART;TZID=Europe/Paris:20211214T120000 RRULE:INTERVAL=1;FREQ=DAILY"
  unified_job_template_id   = awx_job_template.baseconfig.id
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

func resourceSchedule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceScheduleCreate,
		ReadContext:   resourceScheduleRead,
		UpdateContext: resourceScheduleUpdate,
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
			"unified_job_template_id": {
				Type:     schema.TypeInt,
				Required: true,
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
				Type:     schema.TypeInt,
				Optional: true,
			},
			"extra_data": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Extra data to be pass for the schedule (YAML format)",
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceScheduleCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*awx.AWX)
	awxService := client.ScheduleService

	scheduleData := map[string]interface{}{
		"name":                 d.Get("name").(string),
		"rrule":                d.Get("rrule").(string),
		"unified_job_template": d.Get("unified_job_template_id").(int),
		"description":          d.Get("description").(string),
		"enabled":              d.Get("enabled").(bool),
		"extra_data":           unmarshalYaml(d.Get("extra_data").(string)),
	}
	if _, ok := d.GetOk("inventory"); ok {
		scheduleData["inventory"] = d.Get("inventory").(int)
	}

	result, err := awxService.Create(scheduleData, map[string]string{})
	if err != nil {
		log.Printf("Fail to Create Schedule %v", err)
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

func resourceScheduleUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

	scheduleData := map[string]interface{}{
		"name":                 d.Get("name").(string),
		"rrule":                d.Get("rrule").(string),
		"unified_job_template": d.Get("unified_job_template_id").(int),
		"description":          d.Get("description").(string),
		"enabled":              d.Get("enabled").(bool),
		"extra_data":           unmarshalYaml(d.Get("extra_data").(string)),
	}
	if _, ok := d.GetOk("inventory"); ok {
		scheduleData["inventory"] = d.Get("inventory").(int)
	}

	_, err = awxService.Update(id, scheduleData, map[string]string{})
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

func resourceScheduleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
	d = setScheduleResourceData(d, res)
	return nil
}

func resourceScheduleDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*awx.AWX)
	awxService := client.ScheduleService
	id, diags := convertStateIDToNummeric(diagElementHostTitle, d)
	if diags.HasError() {
		return diags
	}

	if _, err := awxService.Delete(id); err != nil {
		return buildDiagDeleteFail(
			diagElementHostTitle,
			fmt.Sprintf("id %v, got %s ",
				id, err.Error()))
	}
	d.SetId("")
	return nil
}

func setScheduleResourceData(d *schema.ResourceData, r *awx.Schedule) *schema.ResourceData {
	d.Set("name", r.Name)
	d.Set("rrule", r.Rrule)
	d.Set("unified_job_template_id", r.UnifiedJobTemplate)
	d.Set("description", r.Description)
	d.Set("enabled", r.Enabled)
	d.Set("inventory", r.Inventory)
	d.Set("extra_data", marshalYaml(r.ExtraData))
	d.SetId(strconv.Itoa(r.ID))
	return d
}
