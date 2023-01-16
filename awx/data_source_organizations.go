/*
Use this data source to list organizations.

Example Usage

```hcl
data "awx_organizations" "all_orgs" {}
```

*/
package awx

import (
	"context"
	"strconv"
	"time"

	awx "github.com/denouche/goawx/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceOrganizations() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceOrganizationsRead,
		Schema: map[string]*schema.Schema{
			"organizations": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceOrganizationsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*awx.AWX)

	parsedOrgs := make([]map[string]interface{}, 0)

	orgs, err := client.OrganizationsService.ListOrganizations(map[string]string{})
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to fetch organizations",
			Detail:   "Unable to fetch organizations from AWX API",
		})
		return diags
	}
	for _, c := range orgs {
		parsedOrgs = append(parsedOrgs, map[string]interface{}{
			"id":   c.ID,
			"name": c.Name,
		})
	}

	err = d.Set("organizations", parsedOrgs)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}
