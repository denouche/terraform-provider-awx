---
layout: "awx"
page_title: "AWX: awx_organization"
sidebar_current: "docs-awx-resource-organization"
description: |-
  *TBD*
---

# awx_organization

*TBD*

## Example Usage

```hcl
resource "awx_organization" "default" {
    name            = "acc-test"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) 
* `custom_virtualenv` - (Optional) Local absolute file path containing a custom Python virtualenv to use
* `default_environment` - (Optional) The default execution environment for jobs run by this organization.
* `description` - (Optional) 
* `max_hosts` - (Optional) Maximum number of hosts allowed to be managed by this organization

