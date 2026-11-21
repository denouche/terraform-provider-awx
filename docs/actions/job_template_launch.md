---
layout: "awx"
page_title: "AWX: awx_job_template_launch (action)"
sidebar_current: "docs-awx-action-job_template_launch"
description: |-
  Launches a job from an AWX job template without storing anything in state.
---

# awx_job_template_launch (action)

Launches a job from an AWX job template. Requires Terraform 1.14 or later.

Unlike the `awx_job_template_launch` resource, the action keeps nothing in state, so there is nothing to destroy or replace later. AWX can purge old jobs without affecting it.

## Example Usage

Run a job template whenever a value changes, by triggering the action from the resource that tracks it:

```hcl
action "awx_job_template_launch" "cluster_upgrade" {
  config {
    job_template_id = awx_job_template.cluster_upgrade.id
    extra_vars      = jsonencode({ kubernetes_version = var.kubernetes_version })
    timeout         = "300m"
  }
}

resource "terraform_data" "kubernetes_version_tracker" {
  input = var.kubernetes_version

  lifecycle {
    action_trigger {
      events  = [after_update]
      actions = [action.awx_job_template_launch.cluster_upgrade]
    }
  }
}
```

Or run it on demand:

```shell
terraform apply -invoke=action.awx_job_template_launch.cluster_upgrade
```

## Argument Reference

The following arguments are supported:

* `job_template_id` - (Required) Job template ID.
* `limit` - (Optional) List of comma delimited hosts to limit job execution. Requires `ask_limit_on_launch` on the job template.
* `inventory_id` - (Optional) Override inventory ID. Requires `ask_inventory_on_launch` on the job template.
* `extra_vars` - (Optional) Override job template variables. YAML or JSON values are supported.
* `wait_for_completion` - (Optional) Wait for the job to finish, and fail the action if the job fails. Defaults to `true`, unlike the resource.
* `timeout` - (Optional) How long to wait for the job, as a Go duration such as `"90m"`. Defaults to `"20m"`.

## Failed jobs are not retried

Terraform saves the triggering resource's new state before it invokes the action. If the job fails, the apply fails, but the next apply finds no changes and does not invoke the action again. Re-run the job with `terraform apply -invoke=...` once the cause is fixed.

The `awx_job_template_launch` resource behaves differently: a failed launch is tainted and re-created on the next apply.
