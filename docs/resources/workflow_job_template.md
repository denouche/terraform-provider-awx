---
layout: "awx"
page_title: "AWX: awx_workflow_job_template"
sidebar_current: "docs-awx-resource-workflow_job_template"
description: |-
  *TBD*
---

# awx_workflow_job_template

*TBD*

## Example Usage

```hcl
resource "awx_workflow_job_template" "default" {
  name            = "workflow-job"
  organization_id = var.organization_id
  inventory_id    = awx_inventory.default.id
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of this workflow job template. (string, required)
* `allow_simultaneous` - (Optional) 
* `ask_inventory_on_launch` - (Optional) 
* `ask_limit_on_launch` - (Optional) 
* `ask_scm_branch_on_launch` - (Optional) 
* `ask_variables_on_launch` - (Optional) 
* `description` - (Optional) Optional description of this workflow job template.
* `inventory_id` - (Optional) Inventory applied as a prompt, assuming job template prompts for inventory. (id, default=``)
* `limit` - (Optional) 
* `organization_id` - (Optional) The organization used to determine access to this template. (id, default=``)
* `scm_branch` - (Optional) 
* `survey_enabled` - (Optional) 
* `variables` - (Optional) 
* `webhook_credential` - (Optional) 
* `webhook_service` - (Optional) 

