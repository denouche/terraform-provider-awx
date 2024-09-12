---
layout: "awx"
page_title: "AWX: awx_workflow_job_template_schedule"
sidebar_current: "docs-awx-resource-workflow_job_template_schedule"
description: |-
  *TBD*
---

# awx_workflow_job_template_schedule

*TBD*

## Example Usage

```hcl
resource "awx_workflow_job_template_schedule" "default" {
  workflow_job_template_id      = awx_workflow_job_template.default.id

  name                      = "schedule-test"
  rrule                     = "DTSTART;TZID=Europe/Paris:20211214T120000 RRULE:INTERVAL=1;FREQ=DAILY"
  extra_data                = <<EOL
organization_name: testorg
EOL
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) 
* `rrule` - (Required) 
* `workflow_job_template_id` - (Required) The workflow_job_template id for this schedule
* `description` - (Optional) 
* `enabled` - (Optional) 
* `extra_data` - (Optional) Extra data to be pass for the schedule (YAML format)
* `inventory` - (Optional) Inventory applied as a prompt, assuming job template prompts for inventory (id, default=``)
* `unified_job_template_id` - (Optional) 

