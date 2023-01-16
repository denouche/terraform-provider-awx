/*
Use this data source to list notification template.

Example Usage

```hcl
data "awx_notification_template" "default" {}

data "awx_notification_template" "default" {
    name = "private_services"
}
```

*/
package awx

import (
	"context"
	"strconv"

	awx "github.com/denouche/goawx/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceNotificationTemplate() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNotificationTemplatesRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func dataSourceNotificationTemplatesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*awx.AWX)
	params := make(map[string]string)
	if groupName, okName := d.GetOk("name"); okName {
		params["name"] = groupName.(string)
	}

	if groupID, okID := d.GetOk("id"); okID {
		params["id"] = strconv.Itoa(groupID.(int))
	}

	notificationTemplates, _, err := client.NotificationTemplatesService.List(params)
	if err != nil {
		return buildDiagnosticsMessage(
			"Get: Fail to fetch NotificationTemplate",
			"Fail to find the group got: %s",
			err.Error(),
		)
	}
	if len(notificationTemplates) > 1 {
		return buildDiagnosticsMessage(
			"Get: find more than one Element",
			"The Query Returns more than one Group, %d",
			len(notificationTemplates),
		)
	}

	notificationTemplate := notificationTemplates[0]
	d = setNotificationTemplateResourceData(d, notificationTemplate)
	return diags
}
