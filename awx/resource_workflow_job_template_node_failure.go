/*
Add a node when the previous step is failed.

Example Usage

```hcl
resource "random_uuid" "workflow_node_k3s_uuid" {}

resource "awx_workflow_job_template_node_failure" "k3s" {
    workflow_job_template_id        = awx_workflow_job_template.default.id
    workflow_job_template_node_id   = awx_workflow_job_template_node.default.id
    unified_job_template_id         = awx_job_template.k3s.id
    identifier                      = random_uuid.workflow_node_k3s_uuid.result
    inventory_id                   = awx_inventory.default.id
}
```

*/

package awx

import (
	"context"

	awx "github.com/denouche/goawx/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceWorkflowJobTemplateNodeFailure() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceWorkflowJobTemplateNodeFailureCreate,
		ReadContext:   resourceWorkflowJobTemplateNodeRead,
		UpdateContext: resourceWorkflowJobTemplateNodeUpdate,
		DeleteContext: resourceWorkflowJobTemplateNodeDelete,
		Schema:        workflowJobNodeSchema,
	}
}

func resourceWorkflowJobTemplateNodeFailureCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*awx.AWX)
	awxService := client.WorkflowJobTemplateNodeFailureService
	return createNodeForWorkflowJob(awxService, ctx, d, m)
}
