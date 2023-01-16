/*
Set LDAP definition for a team.

Example Usage

```hcl
data "awx_organization" "default" {
  name = "Default"
}

resource "awx_team" "admin_team" {
  name = "Admins"
  organization_id = data.awx_organization.default.id
}

resource "awx_settings_ldap_team_map" "admin_team_map" {
  name         = resource.awx_team.admin_team.name
  users        = ["CN=MyTeam,OU=Groups,DC=example,DC=com"]
  organization = data.awx_organization.default.name
  remove       = true
}
```

*/

package awx

import (
    "context"
    "encoding/json"
    "sync"
    "time"

    awx "github.com/denouche/goawx/client"
    "github.com/hashicorp/terraform-plugin-sdk/v2/diag"
    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var ldapTeamMapAccessMutex sync.Mutex

func resourceSettingsLDAPTeamMap() *schema.Resource {
    return &schema.Resource{
        CreateContext: resourceSettingsLDAPTeamMapCreate,
        ReadContext:   resourceSettingsLDAPTeamMapRead,
        DeleteContext: resourceSettingsLDAPTeamMapDelete,
        UpdateContext: resourceSettingsLDAPTeamMapUpdate,

        Schema: map[string]*schema.Schema{
            "name": {
                Type:        schema.TypeString,
                Required:    true,
                Description: "Name of this Team",
            },
            "users": {
                Type: schema.TypeList,
                Elem: &schema.Schema{
                    Type: schema.TypeString,
                },
                Optional:    true,
                Description: "Group DNs to map to this team",
            },
            "organization": {
                Type:        schema.TypeString,
                Required:    true,
                Description: "Name of the team organization",
            },
            "remove": {
                Type:        schema.TypeBool,
                Optional:    true,
                Default:     false,
                Description: "When True, a user who is not a member of the given groups will be removed from the team",
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

type team_map_entry struct {
    UserDNs      interface{} `json:"users"`
    Organization string      `json:"organization"`
    Remove       bool        `json:"remove"`
}

type teammap map[string]team_map_entry

func resourceSettingsLDAPTeamMapCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

    ldapTeamMapAccessMutex.Lock()
    defer ldapTeamMapAccessMutex.Unlock()

    client := m.(*awx.AWX)
    awxService := client.SettingService

    res, err := awxService.GetSettingsBySlug("ldap", make(map[string]string))
    if err != nil {
        return buildDiagnosticsMessage(
            "Create: failed to fetch settings",
            "Failed to fetch any ldap setting, got: %s", err.Error(),
        )
    }

    /*return buildDiagnosticsMessage(
    	"returning as desired",
    	"Data: %v", res,
    )*/
    tmaps := make(teammap)
    err = json.Unmarshal((*res)["AUTH_LDAP_TEAM_MAP"], &tmaps)
    if err != nil {
        return buildDiagnosticsMessage(
            "Create: failed to parse AUTH_LDAP_TEAM_MAP setting",
            "Failed to parse AUTH_LDAP_TEAM_MAP setting, got: %s with input %s", err.Error(), (*res)["AUTH_LDAP_TEAM_MAP"],
        )
    }

    name := d.Get("name").(string)

    _, ok := tmaps[name]
    if ok {
        return buildDiagnosticsMessage(
            "Create: team map already exists",
            "Map for ldap to team map %v already exists", d.Id(),
        )
    }

    newtmap := team_map_entry{
        UserDNs:      d.Get("users").([]interface{}),
        Organization: d.Get("organization").(string),
        Remove:       d.Get("remove").(bool),
    }

    tmaps[name] = newtmap

    payload := map[string]interface{}{
        "AUTH_LDAP_TEAM_MAP": tmaps,
    }

    _, err = awxService.UpdateSettings("ldap", payload, make(map[string]string))
    if err != nil {
        return buildDiagnosticsMessage(
            "Create: team map not created",
            "failed to save team map data, got: %s", err.Error(),
        )
    }

    d.SetId(name)
    return resourceSettingsLDAPTeamMapRead(ctx, d, m)
}

func resourceSettingsLDAPTeamMapUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

    ldapTeamMapAccessMutex.Lock()
    defer ldapTeamMapAccessMutex.Unlock()

    client := m.(*awx.AWX)
    awxService := client.SettingService

    res, err := awxService.GetSettingsBySlug("ldap", make(map[string]string))
    if err != nil {
        return buildDiagnosticsMessage(
            "Update: Unable to fetch settings",
            "Unable to load settings with slug ldap: got %s", err.Error(),
        )
    }

    tmaps := make(teammap)
    err = json.Unmarshal((*res)["AUTH_LDAP_TEAM_MAP"], &tmaps)
    if err != nil {
        return buildDiagnosticsMessage(
            "Update: failed to parse AUTH_LDAP_TEAM_MAP setting",
            "Failed to parse AUTH_LDAP_TEAM_MAP setting, got: %s", err.Error(),
        )
    }

    id := d.Id()
    name := d.Get("name").(string)
    organization := d.Get("organization").(string)
    users := d.Get("users").([]interface{})
    remove := d.Get("remove").(bool)

    if name != id {
        tmaps[name] = tmaps[id]
        delete(tmaps, id)
    }

    utmap := tmaps[name]
    utmap.UserDNs = users
    utmap.Organization = organization
    utmap.Remove = remove
    tmaps[name] = utmap

    payload := map[string]interface{}{
        "AUTH_LDAP_TEAM_MAP": tmaps,
    }

    _, err = awxService.UpdateSettings("ldap", payload, make(map[string]string))
    if err != nil {
        return buildDiagnosticsMessage(
            "Update: team map not created",
            "failed to save team map data, got: %s", err.Error(),
        )
    }

    d.SetId(name)
    return resourceSettingsLDAPTeamMapRead(ctx, d, m)
}

func resourceSettingsLDAPTeamMapRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
    var diags diag.Diagnostics
    client := m.(*awx.AWX)
    awxService := client.SettingService

    res, err := awxService.GetSettingsBySlug("ldap", make(map[string]string))
    if err != nil {
        return buildDiagnosticsMessage(
            "Unable to fetch settings",
            "Unable to load settings with slug ldap: got %s",
            err.Error(),
        )
    }
    tmaps := make(teammap)
    err = json.Unmarshal((*res)["AUTH_LDAP_TEAM_MAP"], &tmaps)
    if err != nil {
        return buildDiagnosticsMessage(
            "Unable to parse AUTH_LDAP_TEAM_MAP",
            "Unable to parse AUTH_LDAP_TEAM_MAP, got: %s", err.Error(),
        )
    }
    mapdef, ok := tmaps[d.Id()]
    if !ok {
        return buildDiagnosticsMessage(
            "Unable to fetch ldap team map",
            "Unable to load ldap team map %v: not found", d.Id(),
        )
    }

    /*return buildDiagnosticsMessage(
    	"returning as desired",
    	"Data: %v %T", mapdef.UserDNs, mapdef.UserDNs,
    )*/

    var users []string
    switch tt := mapdef.UserDNs.(type) {
    case string:
        users = []string{tt}
    case []string:
        users = tt
    case []interface{}:
        for _, v := range tt {
            if dn, ok := v.(string); ok {
                users = append(users, dn)
            }
        }
    }

    d.Set("name", d.Id())
    d.Set("users", users)
    d.Set("organization", mapdef.Organization)
    d.Set("remove", mapdef.Remove)
    return diags
}

func resourceSettingsLDAPTeamMapDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

    ldapTeamMapAccessMutex.Lock()
    defer ldapTeamMapAccessMutex.Unlock()

    var diags diag.Diagnostics
    client := m.(*awx.AWX)
    awxService := client.SettingService

    res, err := awxService.GetSettingsBySlug("ldap", make(map[string]string))
    if err != nil {
        return buildDiagnosticsMessage(
            "Delete: Unable to fetch settings",
            "Unable to load settings with slug ldap: got %s", err.Error(),
        )
    }

    tmaps := make(teammap)
    err = json.Unmarshal((*res)["AUTH_LDAP_TEAM_MAP"], &tmaps)
    if err != nil {
        return buildDiagnosticsMessage(
            "Delete: failed to parse AUTH_LDAP_TEAM_MAP setting",
            "Failed to parse AUTH_LDAP_TEAM_MAP setting, got: %s", err.Error(),
        )
    }

    id := d.Id()
    delete(tmaps, id)

    payload := map[string]interface{}{
        "AUTH_LDAP_TEAM_MAP": tmaps,
    }

    _, err = awxService.UpdateSettings("ldap", payload, make(map[string]string))
    if err != nil {
        return buildDiagnosticsMessage(
            "Delete: team map not created",
            "failed to save team map data, got: %s", err.Error(),
        )
    }
    d.SetId("")
    return diags
}
