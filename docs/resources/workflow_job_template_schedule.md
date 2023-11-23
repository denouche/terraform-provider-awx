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
    workflow_job_template_id  = awx_workflow_job_template.default.id

    name                      = "schedule-test"
    rrule                     = "DTSTART;TZID=Europe/Paris:20211214T120000 RRULE:INTERVAL=1;FREQ=DAILY"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required)
* `rrule` - (Required)
* `workflow_job_template_id` - (Required)
* `description` - (Optional)
* `inventory` - (Optional)
* `timezone` - (Optional)
* `extra_data` - (Optional)
