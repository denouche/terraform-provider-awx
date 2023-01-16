/*
Create a credential type.

Example Usage

```hcl
resource "awx_credential_type" "default" {
    name      = "acc-test"
    injectors = <<YAML
...
YAML
    inputs    = <<YAML
...
YAML
}
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

func resourceCredentialType() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCredentialTypeCreate,
		ReadContext:   resourceCredentialTypeRead,
		UpdateContext: resourceCredentialTypeUpdate,
		DeleteContext: CredentialTypeServiceDeleteByID,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of this credential type.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Optional description of this credential type.",
			},
			"kind": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "cloud",
				Description: "Choices cloud or net",
			},
			"inputs": {
				Type:     schema.TypeString,
				Required: true,
			},
			"injectors": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceCredentialTypeCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	var err error

	inputs := d.Get("inputs").(string)
	inputs_map := make(map[string]interface{})
	inputs_jsonerr := json.Unmarshal([]byte(inputs), &inputs_map)

	if inputs_jsonerr != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create new Credential Type",
			Detail:   fmt.Sprintf("Unable to create new credential type: %s", inputs_jsonerr.Error()),
		})
		return diags
	}

	injectors := d.Get("injectors").(string)
	injectors_map := make(map[string]interface{})
	injectors_jsonerr := json.Unmarshal([]byte(injectors), &injectors_map)

	if injectors_jsonerr != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create new Credential Type",
			Detail:   fmt.Sprintf("Unable to create new credential type: %s", injectors_jsonerr.Error()),
		})
		return diags
	}

	newCredentialType := map[string]interface{}{
		"name":        d.Get("name").(string),
		"description": d.Get("description").(string),
		"kind":        d.Get("kind").(string),
		"inputs":      inputs_map,
		"injectors":   injectors_map,
	}

	client := m.(*awx.AWX)
	credtype, err := client.CredentialTypeService.CreateCredentialType(newCredentialType, map[string]string{})
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create new credential type",
			Detail:   fmt.Sprintf("Unable to create new credential type: %s", err.Error()),
		})
		return diags
	}

	d.SetId(strconv.Itoa(credtype.ID))
	resourceCredentialTypeRead(ctx, d, m)

	return diags
}

func resourceCredentialTypeRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(*awx.AWX)
	id, _ := strconv.Atoi(d.Id())
	credtype, err := client.CredentialTypeService.GetCredentialTypeByID(id, map[string]string{})
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to fetch credential type",
			Detail:   fmt.Sprintf("Unable to credential type with id %d: %s", id, err.Error()),
		})
		return diags
	}

	d.Set("name", credtype.Name)
	d.Set("description", credtype.Description)
	d.Set("kind", credtype.Kind)
	d.Set("inputs", credtype.Inputs)
	d.Set("injectors", credtype.Injectors)

	return diags
}

func resourceCredentialTypeUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	keys := []string{
		"name",
		"description",
		"kind",
		"inputs",
		"injectors",
	}

	if d.HasChanges(keys...) {
		var err error

		inputs := d.Get("inputs").(string)
		inputs_map := make(map[string]interface{})
		inputs_jsonerr := json.Unmarshal([]byte(inputs), &inputs_map)
		if inputs_jsonerr != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Unable to create new credential",
				Detail:   fmt.Sprintf("Unable to update credential type: %s", inputs_jsonerr.Error()),
			})
			return diags
		}

		injectors := d.Get("injectors").(string)
		injectors_map := make(map[string]interface{})
		injectors_jsonerr := json.Unmarshal([]byte(injectors), &injectors_map)
		if injectors_jsonerr != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Unable to create new credential",
				Detail:   fmt.Sprintf("Unable to update credential type: %s", injectors_jsonerr.Error()),
			})
			return diags
		}

		id, _ := strconv.Atoi(d.Id())
		updatedCredentialType := map[string]interface{}{
			"name":        d.Get("name").(string),
			"description": d.Get("description").(string),
			"kind":        d.Get("kind").(string),
			"inputs":      inputs_map,
			"injectors":   injectors_map,
		}

		client := m.(*awx.AWX)
		_, err = client.CredentialTypeService.UpdateCredentialTypeByID(id, updatedCredentialType, map[string]string{})
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Unable to update existing credential type",
				Detail:   fmt.Sprintf("Unable to update existing credential type with id %d: %s", id, err.Error()),
			})
			return diags
		}
	}

	return resourceCredentialTypeRead(ctx, d, m)
}
