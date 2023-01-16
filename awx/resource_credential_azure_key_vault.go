/*
Create an Azure credential.

Example Usage

```hcl
resource "awx_credential_azure_key_vault" "default" {
    name                = "acc-test"
    organization_id     = 1
    url                 = "..."
    client              = "..."
    secret              = "..."
    tenant              = "..."
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

func resourceCredentialAzureKeyVault() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCredentialAzureKeyVaultCreate,
		ReadContext:   resourceCredentialAzureKeyVaultRead,
		UpdateContext: resourceCredentialAzureKeyVaultUpdate,
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
			"url": {
				Type:     schema.TypeString,
				Required: true,
			},
			"client": {
				Type:     schema.TypeString,
				Required: true,
			},
			"secret": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"tenant": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceCredentialAzureKeyVaultCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	var err error

	newCredential := map[string]interface{}{
		"name":            d.Get("name").(string),
		"description":     d.Get("description").(string),
		"organization":    d.Get("organization_id").(int),
		"credential_type": 19, // Azure Key Vault
		"inputs": map[string]interface{}{
			"url":    d.Get("url").(string),
			"client": d.Get("client").(string),
			"secret": d.Get("secret").(string),
			"tenant": d.Get("tenant").(string),
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
	resourceCredentialAzureKeyVaultRead(ctx, d, m)

	return diags
}

func resourceCredentialAzureKeyVaultRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
	d.Set("organization_id", cred.OrganizationID)
	d.Set("url", cred.Inputs["url"])
	d.Set("client", cred.Inputs["client"])
	d.Set("secret", d.Get("secret").(string))
	d.Set("tenant", cred.Inputs["tenant"])

	return diags
}

func resourceCredentialAzureKeyVaultUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	keys := []string{
		"name",
		"description",
		"url",
		"client",
		//"secret",
		"tenant",
	}

	if d.HasChanges(keys...) {
		var err error

		id, _ := strconv.Atoi(d.Id())
		updatedCredential := map[string]interface{}{
			"name":            d.Get("name").(string),
			"description":     d.Get("description").(string),
			"organization":    d.Get("organization_id").(int),
			"credential_type": 19, // Azure Key Vault
			"inputs": map[string]interface{}{
				"url":    d.Get("url").(string),
				"client": d.Get("client").(string),
				"secret": d.Get("secret").(string),
				"tenant": d.Get("tenant").(string),
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

	return resourceCredentialAzureKeyVaultRead(ctx, d, m)
}
