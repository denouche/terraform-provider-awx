/*
Create a team.

Example Usage

```hcl
data "awx_organization" "default" {
    name = "Default"
}

data "awx_inventory" "myinv" {
    name = "My Inventory"
}

data "awx_inventory_role" "myinv_admins" {
    name         = "Admin"
    inventory_id = data.awx_inventory.myinv.id
}

data "awx_project" "myproj" {
    name = "My Project"
}

data "awx_project_role" "myproj_admins" {
    name        = "Admin"
    project_id  = data.awx_project.myproj.id
}

resource "awx_team" "admins_team" {
    name             = "admins-team"
    organization_id  = data.awx_organization.default.id

    role_entitlement {
        role_id = data.awx_inventory_role.myinv_admins.id
    }
    role_entitlement {
        role_id = data.awx_project_role.myproj_admins.id
    }
}
```

*/

package awx

import (
	"context"
	"fmt"
	"strconv"
	"time"

	awx "github.com/denouche/goawx/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceTeam() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTeamCreate,
		ReadContext:   resourceTeamRead,
		DeleteContext: resourceTeamDelete,
		UpdateContext: resourceTeamUpdate,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of this Team",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Optional description of this Team.",
			},
			"organization_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Numeric ID of the Team organization",
			},
			"role_entitlement": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Set of role IDs of the role entitlements",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"role_id": {
							Type:     schema.TypeInt,
							Required: true,
						},
					},
				},
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
	}
}

func resourceTeamCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*awx.AWX)
	awxService := client.TeamService

	orgID := d.Get("organization_id").(int)
	teamName := d.Get("name").(string)
	_, res, err := awxService.ListTeams(map[string]string{
		"name":         teamName,
		"organization": strconv.Itoa(orgID),
	},
	)
	if err != nil {
		return buildDiagnosticsMessage("Create: Fail to find Team", "Fail to find Team %s Organization ID %v, %s", teamName, orgID, err.Error())
	}
	if len(res.Results) >= 1 {
		return buildDiagnosticsMessage("Create: Already exist", "Team with name %s  already exists in the Organization ID %v", teamName, orgID)
	}

	result, err := awxService.CreateTeam(map[string]interface{}{
		"name":         teamName,
		"description":  d.Get("description").(string),
		"organization": d.Get("organization_id").(int),
	}, map[string]string{})
	if err != nil {
		return buildDiagnosticsMessage("Create: Team not created", "Team with name %s  in the Organization ID %v not created, %s", teamName, orgID, err.Error())
	}

	d.SetId(strconv.Itoa(result.ID))

	if rent, entOk := d.GetOk("role_entitlement"); entOk {
		entset := rent.(*schema.Set).List()
		err := roleEntitlementUpdate(m, result.ID, entset, false)
		if err != nil {
			return buildDiagnosticsMessage(
				"Create: team role entitlement not created",
				"Role entitlement for team %s not created: %s", teamName, err.Error(),
			)
		}
	}

	return resourceTeamRead(ctx, d, m)
}

func roleEntitlementUpdate(m interface{}, team_id int, roles []interface{}, remove bool) error {
	client := m.(*awx.AWX)
	awxService := client.TeamService

	for _, v := range roles {
		emap := v.(map[string]interface{})
		payload := map[string]interface{}{
			"id": emap["role_id"],
		}
		if remove {
			payload["disassociate"] = true // presence of key triggers removal
		}

		_, err := awxService.UpdateTeamRoleEntitlement(team_id, payload, make(map[string]string))
		if err != nil {
			return err
		}
	}
	return nil
}

func resourceTeamUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*awx.AWX)
	awxService := client.TeamService

	id, diags := convertStateIDToNummeric("Update Team", d)
	if diags.HasError() {
		return diags
	}
	if d.HasChange("role_entitlement") {
		oi, ni := d.GetChange("role_entitlement")
		if oi == nil {
			oi = new(schema.Set)
		}
		if ni == nil {
			ni = new(schema.Set)
		}
		oe := oi.(*schema.Set)
		ne := ni.(*schema.Set)

		remove := oe.Difference(ne).List()
		add := ne.Difference(oe).List()

		err := roleEntitlementUpdate(m, id, remove, true)
		if err != nil {
			return buildDiagnosticsMessage(
				"Update: Failed To Update Team Role Entitlement",
				"Failed to remove team role entitlement: got %s", err.Error(),
			)
		}
		err = roleEntitlementUpdate(m, id, add, false)
		if err != nil {
			return buildDiagnosticsMessage(
				"Update: Failed To Update Team Role Entitlement",
				"Failed to add team role entitlement: got %s", err.Error(),
			)
		}
		//d.SetPartial("role_entitlemen")
	}
	_, err := awxService.UpdateTeam(id, map[string]interface{}{
		"name":         d.Get("name").(string),
		"description":  d.Get("description").(string),
		"organization": d.Get("organization_id").(int),
	}, map[string]string{})
	if err != nil {
		return buildDiagnosticsMessage("Update: Failed To Update Team", "Fail to get Team with ID %v, got %s", id, err.Error())
	}
	d.Partial(false)
	return resourceTeamRead(ctx, d, m)
}

func resourceTeamRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*awx.AWX)
	awxService := client.TeamService

	id, diags := convertStateIDToNummeric("Read Team", d)
	if diags.HasError() {
		return diags
	}

	team, err := awxService.GetTeamByID(id, make(map[string]string))
	if err != nil {
		return buildDiagNotFoundFail("team", id, err)
	}
	entitlements, _, err := awxService.ListTeamRoleEntitlements(id, make(map[string]string))
	if err != nil {
		return buildDiagNotFoundFail("team roles", id, err)
	}

	d = setTeamResourceData(d, team, entitlements)
	return diags
}

func resourceTeamDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	digMessagePart := "Team"
	client := m.(*awx.AWX)
	awxService := client.TeamService

	id, diags := convertStateIDToNummeric("Delete Team", d)
	if diags.HasError() {
		return diags
	}

	if _, err := awxService.DeleteTeam(id); err != nil {
		return buildDiagDeleteFail(digMessagePart, fmt.Sprintf("TeamID %v, got %s ", id, err.Error()))
	}
	d.SetId("")
	return diags
}

func setTeamResourceData(d *schema.ResourceData, r *awx.Team, e []*awx.ApplyRole) *schema.ResourceData {
	d.Set("name", r.Name)
	d.Set("description", r.Description)
	d.Set("organization_id", r.Organization)

	var entlist []interface{}
	for _, v := range e {
		elem := make(map[string]interface{})
		elem["role_id"] = v.ID
		entlist = append(entlist, elem)
	}
	f := schema.HashResource(&schema.Resource{
		Schema: map[string]*schema.Schema{
			"role_id": {Type: schema.TypeInt},
		}})

	ent := schema.NewSet(f, entlist)

	d.Set("role_entitlement", ent)

	d.SetId(strconv.Itoa(r.ID))
	return d
}
