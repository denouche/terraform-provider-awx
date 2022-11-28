/*
*TBD*

Example Usage

```hcl
resource "awx_job_template" "my_template" {
  name           = "my_template"
  job_type       = "run"
  inventory_id   = awx_inventory.myinv.id
  project_id     = awx_project.myproj.id
  playbook       = "ansible/myplaybook.yml"
}

data "awx_project_role" "template_admin_role" {
  name         = "Admin"
   project_id 	= awx_project.my_template.id
}

*/
package awx

import (
	"context"
	"strconv"

	awx "github.com/denouche/goawx/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceJobTemplateRole() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceJobTemplateRoleRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"job_template_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
		},
	}
}

func dataSourceJobTemplateRoleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*awx.AWX)
	params := make(map[string]string)

	job_template_id := d.Get("job_template_id").(int)
	job_template, err := client.JobTemplateService.GetJobTemplateByID(job_template_id, params)
	if err != nil {
		return buildDiagnosticsMessage(
			"Get: Fail to fetch Job Template",
			"Fail to find the Job Template, got: %s",
			err.Error(),
		)
	}

	roleslist := []*awx.ApplyRole{
		job_template.SummaryFields.ObjectRoles.AdminRole,
		job_template.SummaryFields.ObjectRoles.ReadRole,
		job_template.SummaryFields.ObjectRoles.ExecuteRole,
	}

	if roleName, okName := d.GetOk("name"); okName {
		name := roleName.(string)

		for _, v := range roleslist {
			if v != nil && name == v.Name {
				d = setJobTemplateRoleData(d, v)
				return diags
			}
		}
	}

	if roleID, okID := d.GetOk("id"); okID {
		id := roleID.(int)
		for _, v := range roleslist {
			if v != nil && id == v.ID {
				d = setJobTemplateRoleData(d, v)
				return diags
			}
		}
	}

	return buildDiagnosticsMessage(
		"Failed to fetch job template role - Not Found",
		"The ob template role was not found %s", roleslist,
	)
}

func setJobTemplateRoleData(d *schema.ResourceData, r *awx.ApplyRole) *schema.ResourceData {
	d.Set("name", r.Name)
	d.SetId(strconv.Itoa(r.ID))
	return d
}
