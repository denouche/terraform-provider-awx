/*
Use this data source to list execution environments.

Example Usage

```hcl
data "awx_execution_environment" "default" {}

data "awx_execution_environment" "default" {
    name = "Default"
}

data "awx_execution_environment" "default" {
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

func dataSourceExecutionEnvironment() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceExecutionEnvironmentsRead,
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

func dataSourceExecutionEnvironmentsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*awx.AWX)
	params := make(map[string]string)
	if groupName, okName := d.GetOk("name"); okName {
		params["name"] = groupName.(string)
	}

	if groupID, okGroupID := d.GetOk("id"); okGroupID {
		params["id"] = strconv.Itoa(groupID.(int))
	}

	executionEnvironments, _, err := client.ExecutionEnvironmentsService.ListExecutionEnvironments(params)
	if err != nil {
		return buildDiagnosticsMessage(
			"Get: Fail to fetch execution environment",
			"Fail to find the execution environment got: %s",
			err.Error(),
		)
	}
	if len(executionEnvironments) > 1 {
		return buildDiagnosticsMessage(
			"Get: find more than one element",
			"The query returns more than one execution environment, %d",
			len(executionEnvironments),
		)
		return diags
	}

	executionEnvironment := executionEnvironments[0]
	d = setExecutionEnvironmentsResourceData(d, executionEnvironment)
	return diags
}
