/*
Create a credential input source.

Example Usage

```hcl
resource "awx_credential_input_source" "default" {
    input_field_name     = "acc-test"
    target              = 1
    source              = 1
}
```

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

func resourceCredentialInputSource() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCredentialInputSourceCreate,
		ReadContext:   resourceCredentialInputSourceRead,
		UpdateContext: resourceCredentialInputSourceUpdate,
		DeleteContext: resourceCredentialInputSourceDelete,
		Schema: map[string]*schema.Schema{
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"input_field_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"target": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"source": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"metadata": {
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			},
		},
	}
}

func resourceCredentialInputSourceCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	var err error

	newSourceInput := map[string]interface{}{
		"description":       d.Get("description").(string),
		"input_field_name":  d.Get("input_field_name").(string),
		"target_credential": d.Get("target").(int),
		"source_credential": d.Get("source").(int),
		"metadata":          d.Get("metadata").(map[string]interface{}),
	}

	client := m.(*awx.AWX)
	cred, err := client.CredentialInputSourceService.CreateCredentialInputSource(newSourceInput, map[string]string{})
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create new credentials",
			Detail:   fmt.Sprintf("Unable to create new credentials: %s", err.Error()),
		})
		return diags
	}

	d.SetId(strconv.Itoa(cred.ID))
	resourceCredentialInputSourceRead(ctx, d, m)

	return diags
}

func resourceCredentialInputSourceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(*awx.AWX)
	id, _ := strconv.Atoi(d.Id())
	inputSource, err := client.CredentialInputSourceService.GetCredentialInputSourceByID(id, map[string]string{})
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to fetch credentials",
			Detail:   fmt.Sprintf("Unable to credentials with id %d: %s", id, err.Error()),
		})
		return diags
	}

	d.Set("description", inputSource.Description)
	d.Set("input_field_name", inputSource.InputFieldName)
	d.Set("target", inputSource.TargetCredential)
	d.Set("source", inputSource.SourceCredential)
	d.Set("metadata", inputSource.Metadata)

	return diags
}

func resourceCredentialInputSourceUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	keys := []string{
		"description",
		"input_field_name",
		"target",
		"source",
		"metadata",
	}

	if d.HasChanges(keys...) {
		var err error

		id, _ := strconv.Atoi(d.Id())
		updatedSourceInput := map[string]interface{}{
			"description":       d.Get("description").(string),
			"input_field_name":  d.Get("input_field_name").(string),
			"target_credential": d.Get("target").(int),
			"source_credential": d.Get("source").(int),
			"metadata":          d.Get("metadata").(map[string]interface{}),
		}

		client := m.(*awx.AWX)
		_, err = client.CredentialInputSourceService.UpdateCredentialInputSourceByID(id, updatedSourceInput, map[string]string{})
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Unable to update existing credentials",
				Detail:   fmt.Sprintf("Unable to update existing credentials with id %d: %s", id, err.Error()),
			})
			return diags
		}
	}

	return resourceCredentialInputSourceRead(ctx, d, m)
}

func resourceCredentialInputSourceDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	id, _ := strconv.Atoi(d.Id())
	client := m.(*awx.AWX)
	err := client.CredentialInputSourceService.DeleteCredentialInputSourceByID(id, map[string]string{})
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to delete existing credentials",
			Detail:   fmt.Sprintf("Unable to delete existing credentials with id %d: %s", id, err.Error()),
		})
	}

	return diags
}
