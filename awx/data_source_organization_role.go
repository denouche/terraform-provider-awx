/*
Use this data source to list organization roles for a specified organization.

Example Usage

```hcl
resource "awx_organization" "myorg" {
    name = "My AWX Org"
}

data "awx_organization_role" "org_admins" {
    name            = "Admin"
    organization_id = resource.awx_organization.myorg.id
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

func dataSourceOrganizationRole() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceOrganizationRolesRead,
		Schema: map[string]*schema.Schema{
			"organization_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
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

func dataSourceOrganizationRolesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*awx.AWX)
	params := make(map[string]string)

	orgId := d.Get("organization_id").(int)
	if orgId == 0 {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Get: Missing Parameters",
			Detail:   "organization_id parameter is required.",
		})
		return diags
	}

	organization, err := client.OrganizationsService.GetOrganizationsByID(orgId, params)
	if err != nil {
		return buildDiagnosticsMessage(
			"Get: Fail to fetch organization role",
			"Fail to find the organization role got: %s",
			err.Error(),
		)
	}

	roleslist := []*awx.ApplyRole{
		organization.SummaryFields.ObjectRoles.AdhocRole,
		organization.SummaryFields.ObjectRoles.AdminRole,
		organization.SummaryFields.ObjectRoles.ApprovalRole,
		organization.SummaryFields.ObjectRoles.AuditorRole,
		organization.SummaryFields.ObjectRoles.CredentialAdminRole,
		organization.SummaryFields.ObjectRoles.ExecuteRole,
		organization.SummaryFields.ObjectRoles.InventoryAdminRole,
		organization.SummaryFields.ObjectRoles.JobTemplateAdminRole,
		organization.SummaryFields.ObjectRoles.MemberRole,
		organization.SummaryFields.ObjectRoles.NotificationAdminRole,
		organization.SummaryFields.ObjectRoles.ProjectAdminRole,
		organization.SummaryFields.ObjectRoles.ReadRole,
		organization.SummaryFields.ObjectRoles.UpdateRole,
		organization.SummaryFields.ObjectRoles.UseRole,
		organization.SummaryFields.ObjectRoles.WorkflowAdminRole,
	}

	if roleID, okID := d.GetOk("id"); okID {
		id := roleID.(int)
		for _, v := range roleslist {
			if v != nil && id == v.ID {
				d = setOrganizationRoleData(d, v)
				return diags
			}
		}
	}

	if roleName, okName := d.GetOk("name"); okName {
		name := roleName.(string)

		for _, v := range roleslist {
			if v != nil && name == v.Name {
				d = setOrganizationRoleData(d, v)
				return diags
			}
		}
	}

	return buildDiagnosticsMessage(
		"Failed to fetch organization role - Not Found",
		"The organization role was not found",
	)
}

func setOrganizationRoleData(d *schema.ResourceData, r *awx.ApplyRole) *schema.ResourceData {
	d.Set("name", r.Name)
	d.SetId(strconv.Itoa(r.ID))
	return d
}
