/*
Create a notification template for an organization.

The notification_type field can take different values:
* email
* grafana
* irc
* mattermost
* pagerduty
* rocketchat
* slack
* twilio
* webhook

Example Usage

```hcl
resource "awx_notification_template" "default" {
    name            = "notification_template-test"
    organization_id = 1
    notification_type = "success"
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

	awx "github.com/denouche/goawx/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceNotificationTemplate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNotificationTemplateCreate,
		ReadContext:   resourceNotificationTemplateRead,
		UpdateContext: resourceNotificationTemplateUpdate,
		DeleteContext: resourceNotificationTemplateDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"organization_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"notification_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"notification_configuration": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
		},
	}
}

func resourceNotificationTemplateCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*awx.AWX)
	awxService := client.NotificationTemplatesService

	notificationConfigurationStr := d.Get("notification_configuration").(string)
	notificationConfigurationMap := make(map[string]interface{})
	err := json.Unmarshal([]byte(notificationConfigurationStr), &notificationConfigurationMap)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create NotificationTemplate",
			Detail:   fmt.Sprintf("error while unmarshal notification_configuration: %s", err.Error()),
		})
		return diags
	}

	result, err := awxService.Create(map[string]interface{}{
		"name":                       d.Get("name").(string),
		"description":                d.Get("description").(string),
		"organization":               d.Get("organization_id").(string),
		"notification_type":          d.Get("notification_type").(string),
		"notification_configuration": notificationConfigurationMap,
	}, map[string]string{})
	if err != nil {
		log.Printf("Fail to Create notification_template %v", err)
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create NotificationTemplate",
			Detail:   fmt.Sprintf("NotificationTemplate failed to create %s", err.Error()),
		})
		return diags
	}

	d.SetId(strconv.Itoa(result.ID))
	return resourceNotificationTemplateRead(ctx, d, m)
}

func resourceNotificationTemplateUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*awx.AWX)
	awxService := client.NotificationTemplatesService
	id, diags := convertStateIDToNummeric("Update NotificationTemplate", d)
	if diags.HasError() {
		return diags
	}

	params := make(map[string]string)
	_, err := awxService.GetByID(id, params)
	if err != nil {
		return buildDiagNotFoundFail("notification_template", id, err)
	}

	notificationConfigurationStr := d.Get("notification_configuration").(string)
	notificationConfigurationMap := make(map[string]interface{})
	err = json.Unmarshal([]byte(notificationConfigurationStr), &notificationConfigurationMap)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create NotificationTemplate",
			Detail:   fmt.Sprintf("error while unmarshal notification_configuration: %s", err.Error()),
		})
		return diags
	}

	_, err = awxService.Update(id, map[string]interface{}{
		"name":                       d.Get("name").(string),
		"description":                d.Get("description").(string),
		"organization":               d.Get("organization_id").(string),
		"notification_type":          d.Get("notification_type").(string),
		"notification_configuration": notificationConfigurationMap,
	}, map[string]string{})
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to update NotificationTemplate",
			Detail:   fmt.Sprintf("notification_template with name %s failed to update %s", d.Get("name").(string), err.Error()),
		})
		return diags
	}

	return resourceNotificationTemplateRead(ctx, d, m)
}

func resourceNotificationTemplateRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*awx.AWX)
	awxService := client.NotificationTemplatesService
	id, diags := convertStateIDToNummeric("Read notification_template", d)
	if diags.HasError() {
		return diags
	}

	res, err := awxService.GetByID(id, make(map[string]string))
	if err != nil {
		return buildDiagNotFoundFail("notification_template", id, err)

	}
	d = setNotificationTemplateResourceData(d, res)
	return nil
}

func resourceNotificationTemplateDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*awx.AWX)
	awxService := client.NotificationTemplatesService
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

func setNotificationTemplateResourceData(d *schema.ResourceData, r *awx.NotificationTemplate) *schema.ResourceData {
	d.Set("name", r.Name)
	d.Set("description", r.Description)
	d.Set("organization", r.Organization)
	d.Set("notification_type", r.NotificationType)
	d.Set("notification_configuration", r.NotificationConfiguration)
	d.SetId(strconv.Itoa(r.ID))
	return d
}
