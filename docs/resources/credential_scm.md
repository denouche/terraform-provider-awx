---
layout: "awx"
page_title: "AWX: awx_credential_scm"
sidebar_current: "docs-awx-resource-credential_scm"
description: |-
  *TBD*
---

# awx_credential_scm

*TBD*

## Example Usage

```hcl
data "awx_organization" "default" {
  name = "Default"
}

resource "awx_credential_scm" "credential" {
  name            = "SCM Token"
  description     = "Token to access repo at [...]"
  organization_id = data.awx_organization.default.id
  username        = "Username"
  password        = "Password"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) 
* `organization_id` - (Required) 
* `description` - (Optional) 
* `password` - (Optional) 
* `ssh_key_data` - (Optional) 
* `ssh_key_unlock` - (Optional) 
* `username` - (Optional) 

