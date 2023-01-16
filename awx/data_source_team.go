/*
Use this data source to list teams.

Example Usage

```hcl
data "awx_team" "default" {}

data "awx_team" "default" {
    name = "Default"
}

data "awx_team" "default" {
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

func dataSourceTeam() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTeamsRead,
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

func dataSourceTeamsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*awx.AWX)
	params := make(map[string]string)
	if teamName, okName := d.GetOk("name"); okName {
		params["name"] = teamName.(string)
	}

	if teamID, okTeamID := d.GetOk("id"); okTeamID {
		params["id"] = strconv.Itoa(teamID.(int))
	}

	teams, _, err := client.TeamService.ListTeams(params)
	if err != nil {
		return buildDiagnosticsMessage(
			"Get: Fail to fetch Team",
			"Fail to find the team got: %s",
			err.Error(),
		)
	}
	if len(teams) > 1 {
		return buildDiagnosticsMessage(
			"Get: find more than one Element",
			"The Query Returns more than one team, %d",
			len(teams),
		)
	}

	team := teams[0]
	entitlements, _, err := client.TeamService.ListTeamRoleEntitlements(team.ID, make(map[string]string))
	if err != nil {
		return buildDiagnosticsMessage(
			"Get: Failed to fetch team role entitlements",
			"Fail to retrieve team role entitlements got: %s",
			err.Error(),
		)
	}

	d = setTeamResourceData(d, team, entitlements)
	return diags
}
