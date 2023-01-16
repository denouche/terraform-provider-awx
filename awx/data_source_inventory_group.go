/*
Use this data source to list inventory groups for a specified inventory.

Example Usage

```hcl
data "awx_inventory" "default" {
    id = 1
}

data "awx_inventory_group" "default" {
    name         = "k3sPrimary"
    inventory_id = data.awx_inventory.default.id
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

func dataSourceInventoryGroup() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceInventoryGroupRead,
		Schema: map[string]*schema.Schema{
			"inventory_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
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

func dataSourceInventoryGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*awx.AWX)
	params := make(map[string]string)

	inventoryId := d.Get("inventory_id").(int)

	if inventoryId == 0 {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Get: Missing Parameters",
			Detail:   "inventory_id parameter is required.",
		})
		return diags
	}

	if groupName, okName := d.GetOk("name"); okName {
		params["name"] = groupName.(string)
	}

	if groupID, okGroupID := d.GetOk("id"); okGroupID {
		params["id"] = strconv.Itoa(groupID.(int))
	}

	inventoryID := d.Get("inventory_id").(int)
	groups, _, err := client.InventoryGroupService.ListInventoryGroups(inventoryID, params)
	if err != nil {
		return buildDiagnosticsMessage(
			"Get: Fail to fetch Inventory Group",
			"Fail to find the group got: %s",
			err.Error(),
		)
	}
	if len(groups) > 1 {
		return buildDiagnosticsMessage(
			"Get: find more than one Element",
			"The Query Returns more than one Group, %d",
			len(groups),
		)
		return diags
	}

	group := groups[0]
	d = setInventoryGroupResourceData(d, group)
	return diags
}
