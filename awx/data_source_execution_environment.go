/*
*TBD*

Example Usage

```hcl
data "awx_execution_environment" "default" {
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

	if len(params) == 0 {
		return buildDiagnosticsMessage(
			"Get: Missing Parameters",
			"Please use one of the selectors (name or group_id)",
		)
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
	}
	if len(executionEnvironments) == 0 {
		return buildDiagnosticsMessage(
			"Get: Execution Environment does not exist",
			"The Query Returns no Execution Environment matching filter %v",
			params,
		)
	}

	executionEnvironment := executionEnvironments[0]
	d = setExecutionEnvironmentsResourceData(d, executionEnvironment)
	return diags
}
