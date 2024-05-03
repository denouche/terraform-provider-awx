/*
*TBD*

Example Usage

```hcl
*TBD*
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

func resourceInstanceGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceInstanceGroupCreate,
		ReadContext:   resourceInstanceGroupRead,
		UpdateContext: resourceInstanceGroupUpdate,
		DeleteContext: resourceInstanceGroupDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"is_container_group": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"policy_instance_minimum": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  nil,
			},
			"policy_instance_percentage": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  nil,
			},
			"pod_spec_override": {
				Type:      schema.TypeString,
				Optional:  true,
				Default:   "",
				StateFunc: normalizeJsonYaml,
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceInstanceGroupCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := m.(*awx.AWX)
	awxService := client.InstanceGroupsService

	result, err := awxService.CreateInstanceGroup(map[string]interface{}{
		"name":                       d.Get("name").(string),
		"policy_instance_minimum":    d.Get("policy_instance_minimum").(int),
		"is_container_group":         d.Get("is_container_group").(bool),
		"policy_instance_percentage": d.Get("policy_instance_percentage").(int),
		"pod_spec_override":          d.Get("pod_spec_override").(string),
	}, map[string]string{})
	if err != nil {
		return buildDiagCreateFail(diagElementInstanceGroupTitle, err)
	}

	d.SetId(strconv.Itoa(result.ID))
	return resourceInstanceGroupRead(ctx, d, m)

}

func resourceInstanceGroupUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*awx.AWX)
	awxService := client.InstanceGroupsService
	id, diags := convertStateIDToNummeric(diagElementInstanceGroupTitle, d)
	if diags.HasError() {
		return diags
	}

	updateData := map[string]interface{}{
		"name":              d.Get("name").(string),
		"is_container_group": d.Get("is_container_group").(bool),
		"pod_spec_override": d.Get("pod_spec_override").(string),
	}

	if v, ok := d.GetOk("policy_instance_minimum"); ok {
		updateData["policy_instance_minimum"] = v.(int)
	}

	if v, ok := d.GetOk("policy_instance_percentage"); ok {
		updateData["policy_instance_percentage"] = v.(int)
	}

	_, err := awxService.UpdateInstanceGroup(id, updateData, nil)
	if err != nil {
		return buildDiagUpdateFail(diagElementInstanceGroupTitle, id, err)
	}

	return resourceInstanceGroupRead(ctx, d, m)
}

func resourceInstanceGroupDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*awx.AWX)
	awxService := client.InstanceGroupsService

	id, diags := convertStateIDToNummeric(diagElementInstanceGroupTitle, d)
	if diags.HasError() {
		return diags
	}

	if _, err := awxService.DeleteInstanceGroup(id); err != nil {
		return buildDiagDeleteFail(
			diagElementInstanceGroupTitle,
			fmt.Sprintf("ID: %v, got %s ",
				id, err.Error()))
	}
	d.SetId("")
	return nil
}

func resourceInstanceGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*awx.AWX)
	awxService := client.InstanceGroupsService

	id, diags := convertStateIDToNummeric(diagElementInstanceGroupTitle, d)
	if diags.HasError() {
		return diags
	}

	res, err := awxService.GetInstanceGroupByID(id, make(map[string]string))
	if err != nil {
		return buildDiagNotFoundFail(diagElementInstanceGroupTitle, id, err)
	}
	d = setInstanceGroupResourceData(d, res)
	return diags
}

func setInstanceGroupResourceData(d *schema.ResourceData, r *awx.InstanceGroup) *schema.ResourceData {
	d.Set("name", r.Name)
	d.Set("is_container_group", r.IsContainerGroup)
	d.Set("pod_spec_override", normalizeJsonYaml(r.PodSpecOverride))

	d.SetId(strconv.Itoa(r.ID))
	return d
}
