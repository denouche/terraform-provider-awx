
# awx Provider

Ansible Tower Provider for handle Tower Projects with [rest](https://docs.ansible.com/ansible-tower/latest/html/towerapi/api_ref.html)

## Example Usage

```hcl
provider "awx" {
  hostname = "http://localhost:8078"
  username = "test"
  password = "changeme"
}
```

## Argument Reference

The following arguments are supported:

* `hostname` - (Optional) The API endpoint for AWX. Defaults to `"http://localhost"`.
* `username` - (Optional) The username for API access. Defaults to `"admin"`.
* `password` - (Optional) The password for API access. Defaults to `"password"`.
* `insecure` - (Optional) Whether to check the TLS certificate. Defaults to `false`.
