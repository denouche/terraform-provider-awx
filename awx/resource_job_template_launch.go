/*
Launch a job template.

Example Usage

```hcl
data "awx_inventory" "default" {
    name            = "private_services"
    organization_id = data.awx_organization.default.id
}

resource "awx_job_template" "baseconfig" {
    name           = "baseconfig"
    job_type       = "run"
    inventory_id   = data.awx_inventory.default.id
    project_id     = awx_project.base_service_config.id
    playbook       = "master-configure-system.yml"
    become_enabled = true
}

resource "awx_job_template_launch" "now" {
    job_template_id = awx_job_template.baseconfig.id
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

func resourceJobTemplateLaunch() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceJobTemplateLaunchCreate,
		ReadContext:   resourceJobRead,
		DeleteContext: resourceJobDelete,

		Schema: map[string]*schema.Schema{
			"job_template_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Job template ID",
				ForceNew:    true,
			},
		},
	}
}

func resourceJobTemplateLaunchCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*awx.AWX)
	awxService := client.JobTemplateService
	jobTemplateID := d.Get("job_template_id").(int)
	_, err := awxService.GetJobTemplateByID(jobTemplateID, make(map[string]string))
	if err != nil {
		return buildDiagNotFoundFail("job template", jobTemplateID, err)
	}

	res, err := awxService.Launch(jobTemplateID, map[string]interface{}{}, map[string]string{})
	if err != nil {
		log.Printf("Failed to create Template Launch %v", err)
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create JobTemplate",
			Detail:   fmt.Sprintf("JobTemplate with name %s in the project id %d, failed to create %s", d.Get("name").(string), d.Get("project_id").(int), err.Error()),
		})
		return diags
	}

	// return resourceJobRead(ctx, d, m)
	d.SetId(strconv.Itoa(res.ID))
	return diags
}

func resourceJobRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	return diags
}

func resourceJobDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*awx.AWX)
	awxService := client.JobService
	jobID, diags := convertStateIDToNummeric("Delete Job", d)
	_, err := awxService.GetJob(jobID, map[string]string{})
	if err != nil {
		return buildDiagNotFoundFail("job", jobID, err)
	}

	d.SetId("")
	return diags
}
