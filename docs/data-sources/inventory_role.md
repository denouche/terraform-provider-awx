---
layout: "awx"
page_title: "AWX: awx_inventory_role"
sidebar_current: "docs-awx-datasource-inventory_role"
description: |-
  *TBD*
---

# awx_inventory_role

*TBD*

## Example Usage

```hcl
resource "awx_inventory" "myinv" {
  name = "My Inventory"
  ...
}

data "awx_inventory_role" "inv_admin_role" {
  name         = "Admin"
  inventory_id = data.awx_inventory.myinv.id
}
```

## Argument Reference

The following arguments are supported:

* `inventory_id` - (Required) 
* `id` - (Optional) 
* `name` - (Optional) 

