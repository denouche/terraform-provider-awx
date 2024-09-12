---
layout: "awx"
page_title: "AWX: awx_workflow_job_template_node_failure"
sidebar_current: "docs-awx-resource-workflow_job_template_node_failure"
description: |-
  *TBD*
---

# awx_workflow_job_template_node_failure

*TBD*

## Example Usage

```hcl
resource "random_uuid" "workflow_node_k3s_uuid" {}

resource "awx_workflow_job_template_node_failure" "k3s" {
  workflow_job_template_id      = awx_workflow_job_template.default.id
  workflow_job_template_node_id = awx_workflow_job_template_node.default.id
  unified_job_template_id       = awx_job_template.k3s.id
  inventory_id                  = awx_inventory.default.id
  identifier                    = random_uuid.workflow_node_k3s_uuid.result
}
```

## Argument Reference

The following arguments are supported:

* `identifier` - (Required) 
* `unified_job_template_id` - (Required) 
* `workflow_job_template_id` - (Required) 
* `workflow_job_template_node_id` - (Required) The workflow_job_template_node id from with the new node will start
* `all_parents_must_converge` - (Optional) 
* `diff_mode` - (Optional) 
* `extra_data` - (Optional) 
* `inventory_id` - (Optional) Inventory applied as a prompt, assuming job template prompts for inventory.
* `job_tags` - (Optional) 
* `job_type` - (Optional) 
* `limit` - (Optional) 
* `scm_branch` - (Optional) 
* `skip_tags` - (Optional) 
* `verbosity` - (Optional) 

