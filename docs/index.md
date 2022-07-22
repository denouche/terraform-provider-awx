
# awx Provider

Ansible Tower Provider for handle Tower Projects with [rest](https://docs.ansible.com/ansible-tower/latest/html/towerapi/api_ref.html)

## Example Usage

Using username and password:
```hcl
provider "awx" {
    hostname = "http://localhost:8078"
    username = "test"
    password = "changeme"
}
```

Using token:
```hcl
provider "awx" {
  hostname = "http://localhost:8078"
  token    = "awxtoken"
}
```

> ⚠️ Be careful, if you set both token and username/password the token will have the precedence.

## Argument Reference

The following arguments are supported:

* `hostname` - (Optional) The API endpoint for AWX. Defaults to `"http://localhost"`.
* `username` - (Optional) The username for API access. Defaults to `"admin"`.
* `password` - (Optional) The password for API access. Defaults to `"password"`.
* `token`    - (Optional) The AWX token for API access. Defaults to empty.
* `insecure` - (Optional) Whether to check the TLS certificate. Defaults to `false`.
