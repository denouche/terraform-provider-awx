/*
Use this data source to list workflow job templates.

Example Usage

```hcl
data "awx_workflow_job_template" "default" {}

data "awx_workflow_job_template" "default" {
    name = "Default"
}

data "awx_workflow_job_template" "default" {
    id = 1
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

func dataSourceWorkflowJobTemplate() *schema.Resource {
    return &schema.Resource{
        ReadContext: dataSourceWorkflowJobTemplateRead,
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

func dataSourceWorkflowJobTemplateRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
    var diags diag.Diagnostics
    client := m.(*awx.AWX)
    params := make(map[string]string)
    if groupName, okName := d.GetOk("name"); okName {
        params["name"] = groupName.(string)
    }

    if groupID, okGroupID := d.GetOk("id"); okGroupID {
        params["id"] = strconv.Itoa(groupID.(int))
    }

    workflowJobTemplate, _, err := client.WorkflowJobTemplateService.ListWorkflowJobTemplates(params)
    if err != nil {
        return buildDiagnosticsMessage(
            "Get: Fail to fetch Inventory Group",
            "Fail to find the group got: %s",
            err.Error(),
        )
    }
    if groupName, okName := d.GetOk("name"); okName {
        for _, template := range workflowJobTemplate {
            if template.Name == groupName {
                d = setWorkflowJobTemplateResourceData(d, template)
                return diags
            }
        }
    }
    if _, okGroupID := d.GetOk("id"); okGroupID {
        if len(workflowJobTemplate) != 1 {
            return buildDiagnosticsMessage(
                "Get: find more than one Element",
                "The Query Returns more than one Group, %d",
                len(workflowJobTemplate),
            )
        }
        d = setWorkflowJobTemplateResourceData(d, workflowJobTemplate[0])
        return diags
    }
    return buildDiagnosticsMessage(
        "Get: find more than one Element",
        "The Query Returns more than one Group, %d",
        len(workflowJobTemplate),
    )
}
