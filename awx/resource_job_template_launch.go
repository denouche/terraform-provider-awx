/*
*TBD*

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
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	awx "github.com/denouche/goawx/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
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
			"limit": {
				Type:        schema.TypeString,
				Required:    false,
				Optional:    true,
				ForceNew:    true,
				Description: "List of comma delimited hosts to limit job execution. Required ask_limit_on_launch set on job_template.",
			},
			"inventory_id": {
				Type:        schema.TypeInt,
				Required:    false,
				Optional:    true,
				Computed:    true,
				Description: "Override Inventory ID. Required ask_inventory_on_launch set on job_template.",
				ForceNew:    true,
			},
			"extra_vars": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Override job template variables. YAML or JSON values are supported.",
				ForceNew:    true,
				StateFunc:   normalizeJsonYaml,
			},
			"wait_for_completion": {
				Type:        schema.TypeBool,
				Required:    false,
				Optional:    true,
				Default:     false,
				Description: "Resource creation will wait for job completion.",
				ForceNew:    true,
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
		},
	}
}

func statusInstanceState(ctx context.Context, svc *awx.JobService, id int) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		output, err := svc.GetJob(id, map[string]string{})
		return output, output.Status, err
	}
}

func jobTemplateLaunchWait(ctx context.Context, svc *awx.JobService, job *awx.JobLaunch, timeout time.Duration) error {

	stateConf := &retry.StateChangeConf{
		Pending:    []string{"new", "pending", "waiting", "running"},
		Target:     []string{"successful"},
		Refresh:    statusInstanceState(ctx, svc, job.ID),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)

	return err
}

// JobTemplateLaunchData provides payload data used by the JobTemplateLaunch method
type JobTemplateLaunchData struct {
	Limit       string `json:"limit,omitempty"`
	InventoryID int    `json:"inventory,omitempty"`
	ExtraVars   string `json:"extra_vars,omitempty"`
}

func resourceJobTemplateLaunchCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*awx.AWX)
	awxService := client.JobTemplateService
	awxJobService := client.JobService

	jobTemplateID := d.Get("job_template_id").(int)
	_, err := awxService.GetJobTemplateByID(jobTemplateID, make(map[string]string))
	if err != nil {
		return buildDiagNotFoundFail("job template", jobTemplateID, err)
	}

	data := JobTemplateLaunchData{
		Limit:       d.Get("limit").(string),
		InventoryID: d.Get("inventory_id").(int),
		ExtraVars:   d.Get("extra_vars").(string),
	}

	var iData map[string]interface{}
	idata, _ := json.Marshal(data)
	json.Unmarshal(idata, &iData)

	res, err := awxService.Launch(jobTemplateID, iData, map[string]string{})
	if err != nil {
		log.Printf("Failed to create Template Launch %v", err)
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create JobTemplate",
			Detail:   fmt.Sprintf("JobTemplateLaunch with template ID %d, failed to create %s", d.Get("job_template_id").(int), err.Error()),
		})
		return diags
	}

	// return resourceJobRead(ctx, d, m)
	d.SetId(strconv.Itoa(res.ID))

	if d.Get("wait_for_completion").(bool) {
		err = jobTemplateLaunchWait(ctx, awxJobService, res, d.Timeout(schema.TimeoutCreate))
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "JobTemplate execution failure",
				Detail:   fmt.Sprintf("JobTemplateLaunch with ID %d and template ID %d, failed to complete %s", res.ID, d.Get("job_template_id").(int), err.Error()),
			})
		}
	}
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
