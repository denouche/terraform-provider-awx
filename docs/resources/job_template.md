---
layout: "awx"
page_title: "AWX: awx_job_template"
sidebar_current: "docs-awx-resource-job_template"
description: |-
  *TBD*
---

# awx_job_template

*TBD*

## Example Usage

```hcl
data "awx_inventory" "default" {
  name            = "private_services"
  organization_id = data.awx_organization.default.id
}

resource "awx_job_template" "baseconfig" {
  name           = "baseconfig"
  job_type       = "run"
  inventory_id   = data.awx_inventory.default.id
  project_id     = awx_project.base_service_config.id
  playbook       = "master-configure-system.yml"
  become_enabled = true
}
```

## Argument Reference

The following arguments are supported:

* `inventory_id` - (Required) 
* `job_type` - (Required) One of: run, check, scan
* `name` - (Required) 
* `project_id` - (Required) 
* `allow_simultaneous` - (Optional) 
* `ask_credential_on_launch` - (Optional) 
* `ask_diff_mode_on_launch` - (Optional) 
* `ask_inventory_on_launch` - (Optional) 
* `ask_job_type_on_launch` - (Optional) 
* `ask_limit_on_launch` - (Optional) 
* `ask_skip_tags_on_launch` - (Optional) 
* `ask_tags_on_launch` - (Optional) 
* `ask_variables_on_launch` - (Optional) 
* `ask_verbosity_on_launch` - (Optional) 
* `become_enabled` - (Optional) 
* `custom_virtualenv` - (Optional) 
* `description` - (Optional) 
* `diff_mode` - (Optional) 
* `execution_environment` - (Optional) 
* `extra_vars` - (Optional) 
* `force_handlers` - (Optional) 
* `forks` - (Optional) 
* `host_config_key` - (Optional) 
* `job_tags` - (Optional) 
* `limit` - (Optional) 
* `playbook` - (Optional) 
* `skip_tags` - (Optional) 
* `start_at_task` - (Optional) 
* `survey_enabled` - (Optional) 
* `timeout` - (Optional) 
* `use_fact_cache` - (Optional) 
* `verbosity` - (Optional) One of 0,1,2,3,4,5

