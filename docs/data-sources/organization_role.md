---
layout: "awx"
page_title: "AWX: awx_organization_role"
sidebar_current: "docs-awx-datasource-organization_role"
description: |-
  *TBD*
---

# awx_organization_role

*TBD*

## Example Usage

```hcl
resource "awx_organization" "myorg" {
  name = "My AWX Org"
  ...
}

data "awx_organization_role" "org_admins" {
  name            = "Admin"
  organization_id = resource.awx_organization.myorg.id
}
```

## Argument Reference

The following arguments are supported:

* `organization_id` - (Required) 
* `id` - (Optional) 
* `name` - (Optional) 

