---
layout: "awx"
page_title: "AWX: awx_user"
sidebar_current: "docs-awx-resource-user"
description: |-
  *TBD*
---

# awx_user

*TBD*

## Example Usage

```hcl
data "awx_organization" "default" {
  name = "Default"
}

data "awx_organization_role" "orga_read" {
  name            = "Read"
  organization_id = awx_organization.default.id
}

resource "awx_user" "my_user" {
  username = "my_user"
  password = "my_password"
  role_entitlement {
    role_id = data.awx_organization_role.orga_read.id
  }
}
```

## Argument Reference

The following arguments are supported:

* `password` - (Required) 
* `username` - (Required) 
* `email` - (Optional) 
* `first_name` - (Optional) 
* `is_superuser` - (Optional) 
* `is_system_auditor` - (Optional) 
* `last_name` - (Optional) 
* `role_entitlement` - (Optional) Set of role IDs of the role entitlements

The `role_entitlement` object supports the following:

* `role_id` - (Required) ************************* Please input Description for Schema ************************* 

