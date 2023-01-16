/*
Create a sourec control credential.

Example Usage

```hcl
resource "awx_credential_scm" "default" {
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

func resourceCredentialSCM() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCredentialSCMCreate,
		ReadContext:   resourceCredentialSCMRead,
		UpdateContext: resourceCredentialSCMUpdate,
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
			"ssh_key_unlock": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
		},
	}
}

func resourceCredentialSCMCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	var err error

	newCredential := map[string]interface{}{
		"name":            d.Get("name").(string),
		"description":     d.Get("description").(string),
		"organization":    d.Get("organization_id").(int),
		"credential_type": 2, // Source Controll
		"inputs": map[string]interface{}{
			"username":       d.Get("username").(string),
			"password":       d.Get("password").(string),
			"ssh_key_data":   d.Get("ssh_key_data").(string),
			"ssh_key_unlock": d.Get("ssh_key_unlock").(string),
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
	resourceCredentialSCMRead(ctx, d, m)

	return diags
}

func resourceCredentialSCMRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
	d.Set("ssh_key_unlock", cred.Inputs["ssh_key_unlock"])
	d.Set("organization_id", cred.OrganizationID)

	return diags
}

func resourceCredentialSCMUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	keys := []string{
		"name",
		"description",
		"username",
		"password",
		"ssh_key_data",
		"ssh_key_unlock",
		"organization_id",
	}

	if d.HasChanges(keys...) {
		var err error

		id, _ := strconv.Atoi(d.Id())
		updatedCredential := map[string]interface{}{
			"name":            d.Get("name").(string),
			"description":     d.Get("description").(string),
			"organization":    d.Get("organization_id").(int),
			"credential_type": 2, // Source Controll
			"inputs": map[string]interface{}{
				"username":       d.Get("username").(string),
				"password":       d.Get("password").(string),
				"ssh_key_data":   d.Get("ssh_key_data").(string),
				"ssh_key_unlock": d.Get("ssh_key_unlock").(string),
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

	return resourceCredentialSCMRead(ctx, d, m)
}
