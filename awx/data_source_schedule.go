/*
*TBD*

Example Usage

```hcl
data "awx_schedule" "default" {
  name            = "private_services"
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

func dataSourceSchedule() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSchedulesRead,
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

func dataSourceSchedulesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*awx.AWX)
	params := make(map[string]string)
	if groupName, okName := d.GetOk("name"); okName {
		params["name"] = groupName.(string)
	}

	if groupID, okID := d.GetOk("id"); okID {
		params["id"] = strconv.Itoa(groupID.(int))
	}

	if len(params) == 0 {
		return buildDiagnosticsMessage(
			"Get: Missing Parameters",
			"Please use one of the selectors (name or id)",
		)
	}

	schedules, _, err := client.ScheduleService.List(params)
	if err != nil {
		return buildDiagnosticsMessage(
			"Get: Fail to fetch Schedule Group",
			"Fail to find the group got: %s",
			err.Error(),
		)
	}
	if len(schedules) > 1 {
		return buildDiagnosticsMessage(
			"Get: find more than one Element",
			"The Query Returns more than one Group, %d",
			len(schedules),
		)
	}
	if len(schedules) == 0 {
		return buildDiagnosticsMessage(
			"Get: Schedule does not exist",
			"The Query Returns no Schedule matching filter %v",
			params,
		)
	}

	schedule := schedules[0]
	d = setScheduleResourceData(d, schedule)
	return diags
}
