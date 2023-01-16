/*
This resource configure generic AWX settings.
Please note that resource deletion only delete object from terraform state and do not reset setting to his initial value.

See available settings list here: https://docs.ansible.com/ansible-tower/latest/html/towerapi/api_ref.html#/Settings/Settings_settings_update

Example Usage

```hcl
resource "awx_setting" "social_auth_saml_technical_contact" {
  name  = "SOCIAL_AUTH_SAML_TECHNICAL_CONTACT"
  value = <<EOF
  {
    "givenName": "Myorg",
    "emailAddress": "test@foo.com"
  }
  EOF
}

resource "awx_setting" "social_auth_saml_sp_entity_id" {
  name  = "SOCIAL_AUTH_SAML_SP_ENTITY_ID"
  value = "test"
}

resource "awx_setting" "schedule_max_jobs" {
  name  = "SCHEDULE_MAX_JOBS"
  value = 15
}

resource "awx_setting" "remote_host_headers" {
  name  = "REMOTE_HOST_HEADERS"
  value = <<EOF
  [
    "HTTP_X_FORWARDED_FOR",
    "REMOTE_ADDR",
    "REMOTE_HOST"
  ]
  EOF
}
```

*/

package awx

import (
	"context"
	"encoding/json"
	"time"

	awx "github.com/denouche/goawx/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceSetting() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSettingUpdate,
		ReadContext:   resourceSettingRead,
		DeleteContext: resourceSettingDelete,
		UpdateContext: resourceSettingUpdate,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of setting to modify",
			},
			"value": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Value to be modified for given setting.",
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
	}
}

type setting map[string]string

func resourceSettingUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*awx.AWX)
	awxService := client.SettingService

	_, err := awxService.GetSettingsBySlug("all", make(map[string]string))
	if err != nil {
		return buildDiagnosticsMessage(
			"Create: failed to fetch settings",
			"Failed to fetch setting, got: %s", err.Error(),
		)
	}

	var map_decoded map[string]interface{}
	var array_decoded []interface{}
	var formatted_value interface{}

	name := d.Get("name").(string)
	value := d.Get("value").(string)

	// Attempt to unmarshall string into a map
	err = json.Unmarshal([]byte(value), &map_decoded)

	if err != nil {
		// Attempt to unmarshall string into an array
		err = json.Unmarshal([]byte(value), &array_decoded)

		if err != nil {
			formatted_value = value
		} else {
			formatted_value = array_decoded
		}
	} else {
		formatted_value = map_decoded
	}

	payload := map[string]interface{}{
		name: formatted_value,
	}

	_, err = awxService.UpdateSettings("all", payload, make(map[string]string))
	if err != nil {
		return buildDiagnosticsMessage(
			"Create: setting not created",
			"failed to save setting data, got: %s, %s", err.Error(), value,
		)
	}

	d.SetId(name)
	return resourceSettingRead(ctx, d, m)
}

func resourceSettingRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*awx.AWX)
	awxService := client.SettingService

	_, err := awxService.GetSettingsBySlug("all", make(map[string]string))
	if err != nil {
		return buildDiagnosticsMessage(
			"Unable to fetch settings",
			"Unable to load settings with slug all: got %s", err.Error(),
		)
	}

	d.Set("name", d.Id())
	d.Set("value", d.Get("value").(string))
	return diags
}

func resourceSettingDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	d.SetId("")
	return diags
}
