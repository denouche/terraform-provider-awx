/*
Add a notification on the job template in case of error.

Example Usage

```hcl
resource "awx_job_template_notification_template_error" "baseconfig" {
    job_template_id            = data.awx_job_template.baseconfig.id
    notification_template_id    = data.awx_notification_template.default.id
}
```

*/

package awx

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceJobTemplateNotificationTemplateError() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceJobTemplateNotificationTemplateCreateForType("error"),
		DeleteContext: resourceJobTemplateNotificationTemplateDeleteForType("error"),
		ReadContext:   resourceJobTemplateNotificationTemplateRead,

		Schema: map[string]*schema.Schema{
			"job_template_id": {
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
