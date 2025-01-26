/*
*TBD*

Example Usage

```hcl
resource "awx_workflow_job_template_launch" "now" {
  workflow_job_template_id = someid
  wait_for_completion = true
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

func resourceWorkflowJobTeamplateLaunch() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceWorkflowJobTeamplateLaunchCreate,
		ReadContext:   resourceWorkflowJobRead,
		DeleteContext: resourceWorkflowJobDelete,

		Schema: map[string]*schema.Schema{
			"workflow_job_template_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Workflow job template ID",
				ForceNew:    true,
			},
			"extra_vars": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Override workflow job template variables. YAML or JSON values are supported.",
				ForceNew:    true,
				StateFunc:   normalizeJsonYaml,
			},
			"wait_for_completion": {
				Type:        schema.TypeBool,
				Required:    false,
				Optional:    true,
				Default:     false,
				Description: "Resource creation will wait for workflow job completion.",
				ForceNew:    true,
			},
			"limit": {
				Type:        schema.TypeString,
				Required:    false,
				Optional:    true,
				ForceNew:    true,
				Description: "List of comma delimited hosts to limit workflow job execution. Required ask_limit_on_launch set on workflow_job_template.",
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
		},
	}
}

func statusInstanceWorkflowState(ctx context.Context, svc *awx.WorkflowJobService, id int) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {

		output, err := svc.GetWorkflowJob(id, map[string]string{})

		return output, output.Status, err
	}
}

func workflowJobTemplateLaunchWait(ctx context.Context, svc *awx.WorkflowJobService, job *awx.JobLaunch, timeout time.Duration) error {

	stateConf := &retry.StateChangeConf{
		Pending:    []string{"new", "pending", "waiting", "running"},
		Target:     []string{"successful"},
		Refresh:    statusInstanceWorkflowState(ctx, svc, job.ID),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)

	return err
}

// WorkflowJobTemplateLaunchData provides payload data used by the WorkflowJobTemplateLaunch method
type WorkflowJobTemplateLaunchData struct {
	ExtraVars string `json:"extra_vars,omitempty"`
	Limit     string `json:"limit,omitempty"`
}

func resourceWorkflowJobTeamplateLaunchCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*awx.AWX)
	awxService := client.WorkflowJobTemplateService
	awxWorkflowJobService := client.WorkflowJobService

	workflowJobTemplateID := d.Get("workflow_job_template_id").(int)
	_, err := awxService.GetWorkflowJobTemplateByID(workflowJobTemplateID, make(map[string]string))
	if err != nil {
		return buildDiagNotFoundFail("Workflow job template", workflowJobTemplateID, err)
	}

	data := WorkflowJobTemplateLaunchData{
		ExtraVars: d.Get("extra_vars").(string),
		Limit:     d.Get("limit").(string),
	}

	var iData map[string]interface{}
	idata, _ := json.Marshal(data)
	json.Unmarshal(idata, &iData)

	res, err := awxService.Launch(workflowJobTemplateID, iData, map[string]string{})
	if err != nil {
		log.Printf("Failed to create Workflow Template Launch %v", err)
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create workflowJobTemplate",
			Detail:   fmt.Sprintf("WorkflowJobTemplateLaunch with template ID %d, failed to create %s", d.Get("workflow_job_template_id").(int), err.Error()),
		})
		return diags
	}

	// return resourceWorkflowJobRead(ctx, d, m)
	d.SetId(strconv.Itoa(res.ID))
	if d.Get("wait_for_completion").(bool) { // Print the full structure of the result object
		err = workflowJobTemplateLaunchWait(ctx, awxWorkflowJobService, res, d.Timeout(schema.TimeoutCreate))
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "WorkflowJobTemplate execution failure",
				Detail:   fmt.Sprintf("WorkflowJobTemplateLaunch with ID %d and Workflow template ID %d, failed to complete %s", res.ID, d.Get("workflow_job_template_id").(int), err.Error()),
			})
		}
	}
	return diags
}

func resourceWorkflowJobRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	return diags
}

func resourceWorkflowJobDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*awx.AWX)
	awxService := client.WorkflowJobService
	jobID, diags := convertStateIDToNummeric("Delete Workflow Job", d)
	_, err := awxService.GetWorkflowJob(jobID, map[string]string{})
	if err != nil {
		return buildDiagNotFoundFail("Workflow job", jobID, err)
	}

	d.SetId("")
	return diags
}
