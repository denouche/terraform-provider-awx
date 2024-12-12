/*
This resource lets you attach or detach an instance group to an organization

Example Usage

```hcl
resource "awx_organization" "default" {
    name = "default"
}

resource "awx_instance_group" "default" {
    name = "my-default"
}

resource "awx_organization_instance_group" "default" {
    organization_id   = awx_organization.default.id
    instance_group_id = awx_instance_group.default.id
}
```

*/
package awx

import (
	"context"
	"fmt"
	"strconv"

	awx "github.com/denouche/goawx/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceOrganizationInstanceGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceOrganizationInstanceGroupCreate,
		DeleteContext: resourceOrganizationInstanceGroupDelete,
        ReadContext:   resourceOrganizationInstanceGroupRead,

		Schema: map[string]*schema.Schema{
			"organization_id": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"instance_group_id": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceOrganizationInstanceGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	return diags
}

func resourceOrganizationInstanceGroupCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*awx.AWX)
	awxService := client.OrganizationsService
	OrganizationID := d.Get("organization_id").(int)
	_, err := awxService.GetOrganizationsByID(OrganizationID, make(map[string]string))
	if err != nil {
		return buildDiagNotFoundFail("organization", OrganizationID, err)
	}

	result, err := awxService.AssociateInstanceGroups(OrganizationID, map[string]interface{}{
		"id": d.Get("instance_group_id").(int),
	}, map[string]string{})

	if err != nil {
		return buildDiagnosticsMessage("Create: Organization not AssociateInstanceGroups", "Fail to associate Instance Group %v with Organization %v, got error: %s", d.Get("instance_group_id").(int), d.Get("organization_id").(int), err.Error())
	}

	d.SetId(strconv.Itoa(result.ID))
	return diags
}

func resourceOrganizationInstanceGroupDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*awx.AWX)
	awxService := client.OrganizationsService
	OrganizationID := d.Get("organization_id").(int)
	res, err := awxService.GetOrganizationsByID(OrganizationID, make(map[string]string))
	if err != nil {
		return buildDiagNotFoundFail("organization", OrganizationID, err)
	}

	_, err = awxService.DisAssociateInstanceGroups(res.ID, map[string]interface{}{
		"id": d.Get("instance_group_id").(int),
	}, map[string]string{})
	if err != nil {
		return buildDiagDeleteFail("Organization DisAssociateInstanceGroups", fmt.Sprintf("Fail to disassociate Instance Group %v from Organization %v, got error: %s ", d.Get("instance_group_id").(int), d.Get("organization_id").(int), err.Error()))
	}

	d.SetId("")
	return diags
}
