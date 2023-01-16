/*
Use this data source to query credential type by ID.

Example Usage

```hcl
data "awx_credential_type" "default" {
    id = 1
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

func dataSourceCredentialTypeByID() *schema.Resource {
    return &schema.Resource{
        ReadContext: dataSourceCredentialTypeByIDRead,
        Schema: map[string]*schema.Schema{
            "id": {
                Type:     schema.TypeInt,
                Required: true,
            },
            "name": {
                Type:     schema.TypeString,
                Computed: true,
            },
            "description": {
                Type:     schema.TypeString,
                Computed: true,
            },
            "kind": {
                Type:     schema.TypeString,
                Computed: true,
            },
            "inputs": {
                Type:     schema.TypeString,
                Computed: true,
            },
            "injectors": {
                Type:     schema.TypeString,
                Computed: true,
            },
        },
    }
}

func dataSourceCredentialTypeByIDRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
    var diags diag.Diagnostics

    client := m.(*awx.AWX)
    id := d.Get("id").(int)

    if id == 0 {
        diags = append(diags, diag.Diagnostic{
            Severity: diag.Error,
            Summary:  "Get: Missing Parameters",
            Detail:   "id parameter is required.",
        })
        return diags
    }

    credType, err := client.CredentialTypeService.GetCredentialTypeByID(id, map[string]string{})
    if err != nil {
        diags = append(diags, diag.Diagnostic{
            Severity: diag.Error,
            Summary:  "Unable to fetch credential type",
            Detail:   fmt.Sprintf("Unable to fetch credential type with ID: %d. Error: %s", id, err.Error()),
        })
    }

    d.Set("name", credType.Name)
    d.Set("description", credType.Description)
    d.Set("kind", credType.Kind)
    d.Set("inputs", credType.Inputs)
    d.Set("injectors", credType.Injectors)
    d.SetId(strconv.Itoa(id))

    return diags
}
