---
layout: "awx"
page_title: "AWX: awx_credential_gitlab"
sidebar_current: "docs-awx-resource-credential_gitlab"
description: |-
  *TBD*
---

# awx_credential_gitlab

*TBD*

## Example Usage

```hcl
resource "awx_credential_gitlab" "credential" {
   organization_id = awx_organization.default.id
   name            = "awx-scm-credential"
   description 	   = "test"
   token           = "My_TOKEN"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) 
* `token` - (Required) 
* `description` - (Optional) 
* `organization_id` - (Optional) 

