---
layout: "awx"
page_title: "AWX: awx_setting"
sidebar_current: "docs-awx-resource-setting"
description: |-
  This resource configure generic AWX settings.
Please note that resource deletion only delete object from terraform state and do not reset setting to his initial value.
---

# awx_setting

This resource configure generic AWX settings.
Please note that resource deletion only delete object from terraform state and do not reset setting to his initial value.

See available settings list here: https://docs.ansible.com/ansible-tower/latest/html/towerapi/api_ref.html#/Settings/Settings_settings_update

## Example Usage

```hcl
resource "awx_setting" "social_auth_saml_technical_contact" {
  name  = "SOCIAL_AUTH_SAML_TECHNICAL_CONTACT"
  value = <<EOF
  {
    "givenName": "Myorg",
    "emailAddress": "test@foo.com"
  }
  EOF
}

resource "awx_setting" "social_auth_saml_sp_entity_id" {
  name  = "SOCIAL_AUTH_SAML_SP_ENTITY_ID"
  value = "test"
}

resource "awx_setting" "schedule_max_jobs" {
  name  = "SCHEDULE_MAX_JOBS"
  value = 15
}

resource "awx_setting" "remote_host_headers" {
  name  = "REMOTE_HOST_HEADERS"
  value = <<EOF
  [
    "HTTP_X_FORWARDED_FOR",
    "REMOTE_ADDR",
    "REMOTE_HOST"
  ]
  EOF
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of setting to modify
* `value` - (Required) Value to be modified for given setting.

