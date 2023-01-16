/*
Create a machine credential.

Example Usage

```hcl
resource "awx_credential_machine" "default" {
    name                = "acc-test"
    organization_id     = 1
    username            = "..."
    password            = "..."
}
*/
package awx

import (
	"context"
	"fmt"
	"strconv"

	awx "github.com/denouche/goawx/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCredentialMachine() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCredentialMachineCreate,
		ReadContext:   resourceCredentialMachineRead,
		UpdateContext: resourceCredentialMachineUpdate,
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
			"username": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"ssh_key_data": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"ssh_public_key_data": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ssh_key_unlock": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"become_method": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"become_username": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"become_password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
		},
	}
}

func resourceCredentialMachineCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	var err error

	newCredential := map[string]interface{}{
		"name":            d.Get("name").(string),
		"description":     d.Get("description").(string),
		"organization":    d.Get("organization_id").(int),
		"credential_type": 1, // SSH
		"inputs": map[string]interface{}{
			"username":            d.Get("username").(string),
			"password":            d.Get("password").(string),
			"ssh_key_data":        d.Get("ssh_key_data").(string),
			"ssh_public_key_data": d.Get("ssh_public_key_data").(string),
			"ssh_key_unlock":      d.Get("ssh_key_unlock").(string),
			"become_method":       d.Get("become_method").(string),
			"become_username":     d.Get("become_username").(string),
			"become_password":     d.Get("become_password").(string),
		},
	}

	client := m.(*awx.AWX)
	cred, err := client.CredentialsService.CreateCredentials(newCredential, map[string]string{})
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create new credentials",
			Detail:   fmt.Sprintf("Unable to create new credentials: %s", err.Error()),
		})
		return diags
	}

	d.SetId(strconv.Itoa(cred.ID))
	resourceCredentialMachineRead(ctx, d, m)

	return diags
}

func resourceCredentialMachineRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(*awx.AWX)
	id, _ := strconv.Atoi(d.Id())
	cred, err := client.CredentialsService.GetCredentialsByID(id, map[string]string{})
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to fetch credentials",
			Detail:   fmt.Sprintf("Unable to credentials with id %d: %s", id, err.Error()),
		})
		return diags
	}

	d.Set("name", cred.Name)
	d.Set("description", cred.Description)
	d.Set("username", cred.Inputs["username"])
	d.Set("password", cred.Inputs["password"])
	d.Set("ssh_key_data", cred.Inputs["ssh_key_data"])
	d.Set("ssh_public_key_data", cred.Inputs["ssh_public_key_data"])
	d.Set("ssh_key_unlock", cred.Inputs["ssh_key_unlock"])
	d.Set("become_method", cred.Inputs["become_method"])
	d.Set("become_username", cred.Inputs["become_username"])
	d.Set("become_password", cred.Inputs["become_password"])
	d.Set("organization_id", cred.OrganizationID)

	return diags
}

func resourceCredentialMachineUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	keys := []string{
		"name",
		"description",
		"username",
		"password",
		"ssh_key_data",
		"ssh_public_key_data",
		"ssh_key_unlock",
		"become_method",
		"become_username",
		"become_password",
		"organization_id",
		"team_id",
		"owner_id",
	}

	if d.HasChanges(keys...) {
		var err error

		id, _ := strconv.Atoi(d.Id())
		updatedCredential := map[string]interface{}{
			"name":            d.Get("name").(string),
			"description":     d.Get("description").(string),
			"organization":    d.Get("organization_id").(int),
			"credential_type": 1, // SSH
			"inputs": map[string]interface{}{
				"username":            d.Get("username").(string),
				"password":            d.Get("password").(string),
				"ssh_key_data":        d.Get("ssh_key_data").(string),
				"ssh_public_key_data": d.Get("ssh_public_key_data").(string),
				"ssh_key_unlock":      d.Get("ssh_key_unlock").(string),
				"become_method":       d.Get("become_method").(string),
				"become_username":     d.Get("become_username").(string),
				"become_password":     d.Get("become_password").(string),
			},
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

	return resourceCredentialMachineRead(ctx, d, m)
}
