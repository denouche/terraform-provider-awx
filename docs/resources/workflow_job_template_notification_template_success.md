---
layout: "awx"
page_title: "AWX: awx_workflow_job_template_notification_template_success"
sidebar_current: "docs-awx-resource-workflow_job_template_notification_template_success"
description: |-
  *TBD*
---

# awx_workflow_job_template_notification_template_success

*TBD*

## Example Usage

```hcl
resource "awx_workflow_job_template_notification_template_success" "baseconfig" {
  workflow_job_template_id   = awx_workflow_job_template.baseconfig.id
  notification_template_id   = awx_notification_template.default.id
}
```

## Argument Reference

The following arguments are supported:

* `notification_template_id` - (Required, ForceNew) 
* `workflow_job_template_id` - (Required, ForceNew) 

