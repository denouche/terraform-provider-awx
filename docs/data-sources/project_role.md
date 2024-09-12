---
layout: "awx"
page_title: "AWX: awx_project_role"
sidebar_current: "docs-awx-datasource-project_role"
description: |-
  *TBD*
---

# awx_project_role

*TBD*

## Example Usage

```hcl
resource "awx_project" "myproj" {
  name = "My AWX Project"
  ...
}

data "awx_project_role" "proj_admins" {
  name       = "Admin"
  project_id = resource.awx_project.myproj.id
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required) 
* `id` - (Optional) 
* `name` - (Optional) 

