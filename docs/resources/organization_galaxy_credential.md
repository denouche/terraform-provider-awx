---
layout: "awx"
page_title: "AWX: awx_organization_galaxy_credential"
sidebar_current: "docs-awx-resource-organization_galaxy_credential"
description: |-
  *TBD*
---

# awx_organization_galaxy_credential

*TBD*

## Example Usage

```hcl
resource "awx_organization_galaxy_credential" "baseconfig" {
  organization_id = awx_organization.baseconfig.id
  credential_id   = awx_credential_machine.pi_connection.id
}
```

## Argument Reference

The following arguments are supported:

* `credential_id` - (Required, ForceNew) 
* `organization_id` - (Required, ForceNew) 

