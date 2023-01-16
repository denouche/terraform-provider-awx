/*
Create an inventory for an organization.

Example Usage

```hcl
data "awx_organization" "default" {
    name = "Default"
}

resource "awx_inventory" "default" {
    name            = "acc-test"
    organization_id = data.awx_organization.default.id
    variables       = <<YAML
---
system_supporters:
  - pi
YAML
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

func resourceInventory() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceInventoryCreate,
		ReadContext:   resourceInventoryRead,
		DeleteContext: resourceInventoryDelete,
		UpdateContext: resourceInventoryUpdate,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"organization_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"kind": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"host_filter": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"variables": {
				Type:      schema.TypeString,
				Optional:  true,
				Default:   "",
				StateFunc: normalizeJsonYaml,
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceInventoryCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*awx.AWX)
	awxService := client.InventoriesService

	result, err := awxService.CreateInventory(map[string]interface{}{
		"name":         d.Get("name").(string),
		"organization": d.Get("organization_id").(string),
		"description":  d.Get("description").(string),
		"kind":         d.Get("kind").(string),
		"host_filter":  d.Get("host_filter").(string),
		"variables":    d.Get("variables").(string),
	}, map[string]string{})
	if err != nil {
		return buildDiagCreateFail(diagElementInventoryTitle, err)
	}

	d.SetId(strconv.Itoa(result.ID))
	return resourceInventoryRead(ctx, d, m)

}

func resourceInventoryUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*awx.AWX)
	awxService := client.InventoriesService
	id, diags := convertStateIDToNummeric(diagElementInventoryTitle, d)
	if diags.HasError() {
		return diags
	}
	_, err := awxService.UpdateInventory(id, map[string]interface{}{
		"name":         d.Get("name").(string),
		"organization": d.Get("organization_id").(string),
		"description":  d.Get("description").(string),
		"kind":         d.Get("kind").(string),
		"host_filter":  d.Get("host_filter").(string),
		"variables":    d.Get("variables").(string),
	}, nil)
	if err != nil {
		return buildDiagUpdateFail(diagElementInventoryTitle, id, err)
	}

	return resourceInventoryRead(ctx, d, m)

}

func resourceInventoryRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*awx.AWX)
	awxService := client.InventoriesService
	id, err := strconv.Atoi(d.Id())
	id, diags := convertStateIDToNummeric(diagElementInventoryTitle, d)
	if diags.HasError() {
		return diags
	}
	r, err := awxService.GetInventory(id, map[string]string{})
	if err != nil {
		return buildDiagNotFoundFail(diagElementInventoryTitle, id, err)
	}
	d = setInventoryResourceData(d, r)
	return nil
}

func resourceInventoryDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*awx.AWX)
	awxService := client.InventoriesService
	id, diags := convertStateIDToNummeric(diagElementInventoryTitle, d)
	if diags.HasError() {
		return diags
	}
	if _, err := awxService.DeleteInventory(id); err != nil {
		return buildDiagDeleteFail(
			diagElementInventoryTitle,
			fmt.Sprintf(
				"%s %v, got %s ",
				diagElementInventoryTitle, id, err.Error(),
			),
		)
	}
	d.SetId("")
	return nil
}

func setInventoryResourceData(d *schema.ResourceData, r *awx.Inventory) *schema.ResourceData {
	d.Set("name", r.Name)
	d.Set("organization_id", strconv.Itoa(r.Organization))
	d.Set("description", r.Description)
	d.Set("kind", r.Kind)
	d.Set("host_filter", r.HostFilter)
	d.Set("variables", normalizeJsonYaml(r.Variables))
	d.SetId(strconv.Itoa(r.ID))
	return d
}
