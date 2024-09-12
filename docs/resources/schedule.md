---
layout: "awx"
page_title: "AWX: awx_schedule"
sidebar_current: "docs-awx-resource-schedule"
description: |-
  *TBD*
---

# awx_schedule

*TBD*

## Example Usage

```hcl
resource "awx_schedule" "default" {
  name                      = "schedule-test"
  rrule                     = "DTSTART;TZID=Europe/Paris:20211214T120000 RRULE:INTERVAL=1;FREQ=DAILY"
  unified_job_template_id   = awx_job_template.baseconfig.id
  extra_data                = <<EOL
organization_name: testorg
EOL
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) 
* `rrule` - (Required) 
* `unified_job_template_id` - (Required) 
* `description` - (Optional) 
* `enabled` - (Optional) 
* `extra_data` - (Optional) Extra data to be pass for the schedule (YAML format)
* `inventory` - (Optional) 

