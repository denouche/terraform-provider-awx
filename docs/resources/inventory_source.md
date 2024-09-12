---
layout: "awx"
page_title: "AWX: awx_inventory_source"
sidebar_current: "docs-awx-resource-inventory_source"
description: |-
  *TBD*
---

# awx_inventory_source

*TBD*

## Example Usage

```hcl
data "awx_organization" "default" {
  name = "Default"
}

data "awx_inventory" "default" {
  name            = "private_services"
  organization_id = data.awx_organization.default.id
}

data "awx_project" "default" {
  name = "Default"
}

resource "awx_inventory_source" "inventory_source" {
  name              = "hosts"
  inventory_id      = data.awx_inventory.default.id
  source            = "scm"
  source_project_id = data.awx_project.default.id
  source_path       = "inventory.yml"
}
```

## Argument Reference

The following arguments are supported:

* `inventory_id` - (Required, ForceNew) 
* `name` - (Required) 
* `credential_id` - (Optional) 
* `description` - (Optional) 
* `enabled_value` - (Optional) 
* `enabled_var` - (Optional) 
* `group_by` - (Optional) 
* `host_filter` - (Optional) 
* `instance_filters` - (Optional) 
* `overwrite_vars` - (Optional) 
* `overwrite` - (Optional) 
* `source_path` - (Optional) 
* `source_project_id` - (Optional) 
* `source_regions` - (Optional) 
* `source_vars` - (Optional) 
* `source` - (Optional) 
* `update_cache_timeout` - (Optional) 
* `update_on_launch` - (Optional) 
* `verbosity` - (Optional) 

