/*
*TBD*

Example Usage

```hcl
data "awx_organization" "default" {
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

func dataSourceOrganization() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceOrganizationRead,
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

func dataSourceOrganizationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
	organizations, err := client.OrganizationsService.ListOrganizations(params)
	if err != nil {
		return buildDiagnosticsMessage(
			"Get: Fail to fetch organization",
			"Fail to find the organization got: %s",
			err.Error(),
		)
	}
	if len(organizations) > 1 {
		return buildDiagnosticsMessage(
			"Get: find more than one Element",
			"The Query Returns more than one organization, %d",
			len(organizations),
		)
	}
	if len(organizations) == 0 {
		return buildDiagnosticsMessage(
			"Get: Organization does not exist",
			"The Query Returns no Organization matching filter %v",
			params,
		)
	}

	organization := organizations[0]
	d = setOrganizationsResourceData(d, organization)
	return diags
}
