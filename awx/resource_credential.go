/*
Create a credential.

Example Usage

```hcl
resource "awx_credential" "default" {
    name                = "acc-test"
    organization_id     = 1
    credential_type_id  = 1
    inputs              = <<YAML
...
YAML
}
```

*/

package awx

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	awx "github.com/denouche/goawx/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCredential() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCredentialCreate,
		ReadContext:   resourceCredentialRead,
		UpdateContext: resourceCredentialUpdate,
		DeleteContext: CredentialsServiceDeleteByID,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"organization_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"credential_type_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Specify the type of credential you want to create. Refer to the Ansible Tower documentation for details on each type",
			},
			"inputs": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
		},
	}
}

func resourceCredentialCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	var err error

	inputs := d.Get("inputs").(string)
	inputs_map := make(map[string]interface{})
	jsonerr := json.Unmarshal([]byte(inputs), &inputs_map)

	if jsonerr != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create new credential",
			Detail:   fmt.Sprintf("Unable to create new credential: %s", jsonerr.Error()),
		})
		return diags
	}

	newCredential := map[string]interface{}{
		"name":            d.Get("name").(string),
		"description":     d.Get("description").(string),
		"organization":    d.Get("organization_id").(int),
		"credential_type": d.Get("credential_type_id").(int),
		"inputs":          inputs_map,
	}

	client := m.(*awx.AWX)
	cred, err := client.CredentialsService.CreateCredentials(newCredential, map[string]string{})
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create new credential",
			Detail:   fmt.Sprintf("Unable to create new credential: %s", err.Error()),
		})
		return diags
	}

	d.SetId(strconv.Itoa(cred.ID))
	resourceCredentialRead(ctx, d, m)

	return diags
}

func resourceCredentialRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(*awx.AWX)
	id, _ := strconv.Atoi(d.Id())
	cred, err := client.CredentialsService.GetCredentialsByID(id, map[string]string{})
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to fetch credential",
			Detail:   fmt.Sprintf("Unable to credential with id %d: %s", id, err.Error()),
		})
		return diags
	}

	d.Set("name", cred.Name)
	d.Set("description", cred.Description)
	d.Set("organization_id", cred.OrganizationID)
	d.Set("inputs", cred.Inputs)

	return diags
}

func resourceCredentialUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	keys := []string{
		"name",
		"description",
		"organization_id",
		"inputs",
	}

	if d.HasChanges(keys...) {
		var err error

		inputs := d.Get("inputs").(string)
		inputs_map := make(map[string]interface{})
		jsonerr := json.Unmarshal([]byte(inputs), &inputs_map)

		if jsonerr != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Unable to create new credential",
				Detail:   fmt.Sprintf("Unable to create new credential: %s", jsonerr.Error()),
			})
			return diags
		}

		id, _ := strconv.Atoi(d.Id())
		updatedCredential := map[string]interface{}{
			"name":            d.Get("name").(string),
			"description":     d.Get("description").(string),
			"organization":    d.Get("organization_id").(int),
			"credential_type": d.Get("credential_type_id"),
			"inputs":          inputs_map,
		}

		client := m.(*awx.AWX)
		_, err = client.CredentialsService.UpdateCredentialsByID(id, updatedCredential, map[string]string{})
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Unable to update existing credentials",
				Detail:   fmt.Sprintf("Unable to update existing credentials with id %d: %s", id, err.Error()),
			})
			return diags
		}
	}

	return resourceCredentialSCMRead(ctx, d, m)
}
