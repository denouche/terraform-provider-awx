/*
*TBD*

Example Usage

```hcl
resource "awx_workflow_job_template_notification_template_error" "baseconfig" {
  workflow_job_template_id   = awx_workflow_job_template.baseconfig.id
  notification_template_id   = awx_notification_template.default.id
}
```

*/
package awx

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceWorkflowJobTemplateNotificationTemplateError() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceWorkflowJobTemplateNotificationTemplateCreateForType("error"),
		DeleteContext: resourceWorkflowJobTemplateNotificationTemplateDeleteForType("error"),
		ReadContext:   resourceWorkflowJobTemplateNotificationTemplateRead,

		Schema: map[string]*schema.Schema{
			"workflow_job_template_id": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"notification_template_id": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
		},
	}
}
