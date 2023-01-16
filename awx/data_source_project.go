/*
Use this data source to list projects.

Example Usage

```hcl
data "awx_project" "default" {}

data "awx_project" "default" {
    name = "Default"
}
```

*/
package awx

import (
    "context"
    "strconv"

    awx "github.com/denouche/goawx/client"
    "github.com/hashicorp/terraform-plugin-sdk/v2/diag"
    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceProject() *schema.Resource {
    return &schema.Resource{
        ReadContext: dataSourceProjectsRead,
        Schema: map[string]*schema.Schema{
            "id": {
                Type:     schema.TypeInt,
                Optional: true,
                Computed: true,
            },
            "name": {
                Type:     schema.TypeString,
                Optional: true,
                Computed: true,
            },
        },
    }
}

func dataSourceProjectsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
    var diags diag.Diagnostics
    client := m.(*awx.AWX)
    params := make(map[string]string)
    if groupName, okName := d.GetOk("name"); okName {
        params["name"] = groupName.(string)
    }

    if groupID, okGroupID := d.GetOk("id"); okGroupID {
        params["id"] = strconv.Itoa(groupID.(int))
    }

    projects, _, err := client.ProjectService.ListProjects(params)
    if err != nil {
        return buildDiagnosticsMessage(
            "Get: Fail to fetch Inventory Group",
            "Fail to find the group got: %s",
            err.Error(),
        )
    }
    if len(projects) > 1 {
        return buildDiagnosticsMessage(
            "Get: find more than one Element",
            "The Query Returns more than one Group, %d",
            len(projects),
        )
    }

    Project := projects[0]
    d = setProjectResourceData(d, Project)
    return diags
}
