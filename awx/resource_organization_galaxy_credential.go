/*
*TBD*

Example Usage

```hcl
resource "awx_organization_galaxy_credential" "baseconfig" {
  organization_id = awx_organization.baseconfig.id
  credential_id   = awx_credential_machine.pi_connection.id
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

func resourceOrganizationsGalaxyCredentials() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceOrganizationsGalaxyCredentialsCreate,
		DeleteContext: resourceOrganizationsGalaxyCredentialsDelete,
		ReadContext:   resourceOrganizationsGalaxyCredentialsRead,

		Schema: map[string]*schema.Schema{

			"organization_id": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"credential_id": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceOrganizationsGalaxyCredentialsCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*awx.AWX)
	awxService := client.OrganizationsService
	OrganizationID := d.Get("organization_id").(int)
	_, err := awxService.GetOrganizationsByID(OrganizationID, make(map[string]string))
	if err != nil {
		return buildDiagNotFoundFail("organization", OrganizationID, err)
	}

	result, err := awxService.AssociateGalaxyCredentials(OrganizationID, map[string]interface{}{
		"id": d.Get("credential_id").(int),
	}, map[string]string{})

	if err != nil {
		return buildDiagnosticsMessage("Create: Organization not AssociateGalaxyCredentials", "Fail to add Galaxy credentials with Id %v, for Organization ID %v, got error: %s", d.Get("credential_id").(int), OrganizationID, err.Error())
	}

	d.SetId(strconv.Itoa(result.ID))
	return diags
}

func resourceOrganizationsGalaxyCredentialsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	return diags
}

func resourceOrganizationsGalaxyCredentialsDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*awx.AWX)
	awxService := client.OrganizationsService
	OrganizationID := d.Get("organization_id").(int)
	res, err := awxService.GetOrganizationsByID(OrganizationID, make(map[string]string))
	if err != nil {
		return buildDiagNotFoundFail("organization", OrganizationID, err)
	}

	_, err = awxService.DisAssociateGalaxyCredentials(res.ID, map[string]interface{}{
		"id": d.Get("credential_id").(int),
	}, map[string]string{})
	if err != nil {
		return buildDiagDeleteFail("Organization DisAssociateGalaxyCredentials", fmt.Sprintf("DisAssociateGalaxyCredentials %v, from OrganizationID %v got %s ", d.Get("credential_id").(int), d.Get("organization_id").(int), err.Error()))
	}

	d.SetId("")
	return diags
}
