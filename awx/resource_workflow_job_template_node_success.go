/*
Add a node when the previous step is successful.

Example Usage

```hcl
resource "random_uuid" "workflow_node_k3s_uuid" {}

resource "awx_workflow_job_template_node_success" "k3s" {
    workflow_job_template_id        = data.awx_workflow_job_template.default.id
    workflow_job_template_node_id   = data.awx_workflow_job_template_node.default.id
    unified_job_template_id         = data.awx_job_template.k3s.id
    identifier                      = random_uuid.workflow_node_k3s_uuid.result
    inventory_id                   = data.awx_inventory.default.id
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

func resourceWorkflowJobTemplateNodeSuccess() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceWorkflowJobTemplateNodeSuccessCreate,
		ReadContext:   resourceWorkflowJobTemplateNodeRead,
		UpdateContext: resourceWorkflowJobTemplateNodeUpdate,
		DeleteContext: resourceWorkflowJobTemplateNodeDelete,
		Schema:        workflowJobNodeSchema,
	}
}

func resourceWorkflowJobTemplateNodeSuccessCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := m.(*awx.AWX)
	awxService := client.WorkflowJobTemplateNodeSuccessService
	return createNodeForWorkflowJob(awxService, ctx, d, m)
}
