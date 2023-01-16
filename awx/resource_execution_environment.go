/*
Create an execution environment.

Example Usage

```hcl
resource "awx_execution_environment" "default" {
    name  = "acc-test"
    image = "..."
}
```

*/

package awx

import (
	"context"
	"fmt"
	"log"
	"strconv"

	awx "github.com/denouche/goawx/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceExecutionEnvironment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceExecutionEnvironmentsCreate,
		ReadContext:   resourceExecutionEnvironmentsRead,
		UpdateContext: resourceExecutionEnvironmentsUpdate,
		DeleteContext: resourceExecutionEnvironmentsDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"image": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"organization": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"credential": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
		},
	}
}

func resourceExecutionEnvironmentsCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*awx.AWX)
	awxService := client.ExecutionEnvironmentsService

	result, err := awxService.CreateExecutionEnvironment(map[string]interface{}{
		"name":         d.Get("name").(string),
		"image":        d.Get("image").(string),
		"description":  d.Get("description").(string),
		"organization": AtoipOr(d.Get("organization").(string), nil),
		"credential":   AtoipOr(d.Get("credential").(string), nil),
	}, map[string]string{})
	if err != nil {
		log.Printf("Fail to Create ExecutionEnvironment %v", err)
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create ExecutionEnvironments",
			Detail:   fmt.Sprintf("ExecutionEnvironments with name, failed to create %s", d.Get("name").(string), err.Error()),
		})
		return diags
	}

	d.SetId(strconv.Itoa(result.ID))
	return resourceExecutionEnvironmentsRead(ctx, d, m)
}

func resourceExecutionEnvironmentsUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*awx.AWX)
	awxService := client.ExecutionEnvironmentsService
	id, diags := convertStateIDToNummeric("Update ExecutionEnvironments", d)
	if diags.HasError() {
		return diags
	}

	params := make(map[string]string)

	_, err := awxService.GetExecutionEnvironmentByID(id, params)
	if err != nil {
		return buildDiagNotFoundFail("ExecutionEnvironments", id, err)
	}

	_, err = awxService.UpdateExecutionEnvironment(id, map[string]interface{}{
		"name":         d.Get("name").(string),
		"image":        d.Get("image").(string),
		"description":  d.Get("description").(string),
		"organization": AtoipOr(d.Get("organization").(string), nil),
		"credential":   AtoipOr(d.Get("credential").(string), nil),
	}, map[string]string{})
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to update ExecutionEnvironments",
			Detail:   fmt.Sprintf("ExecutionEnvironments with name %s failed to update %s", d.Get("name").(string), err.Error()),
		})
		return diags
	}

	return resourceExecutionEnvironmentsRead(ctx, d, m)
}

func resourceExecutionEnvironmentsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*awx.AWX)
	awxService := client.ExecutionEnvironmentsService
	id, diags := convertStateIDToNummeric("Read ExecutionEnvironments", d)
	if diags.HasError() {
		return diags
	}

	res, err := awxService.GetExecutionEnvironmentByID(id, make(map[string]string))
	if err != nil {
		return buildDiagNotFoundFail("ExecutionEnvironment", id, err)

	}
	d = setExecutionEnvironmentsResourceData(d, res)
	return nil
}

func resourceExecutionEnvironmentsDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	digMessagePart := "ExecutionEnvironment"
	client := m.(*awx.AWX)
	awxService := client.ExecutionEnvironmentsService
	id, diags := convertStateIDToNummeric("Delete ExecutionEnvironment", d)
	if diags.HasError() {
		return diags
	}

	if _, err := awxService.DeleteExecutionEnvironment(id); err != nil {
		return buildDiagDeleteFail(digMessagePart, fmt.Sprintf("ExecutionEnvironmentID %v, got %s ", id, err.Error()))
	}
	d.SetId("")
	return diags
}

func setExecutionEnvironmentsResourceData(d *schema.ResourceData, r *awx.ExecutionEnvironment) *schema.ResourceData {
	d.Set("name", r.Name)
	d.Set("image", r.Image)
	d.Set("description", r.Description)
	d.Set("organization", r.Organization)
	d.Set("credential", r.Credential)
	d.SetId(strconv.Itoa(r.ID))
	return d
}
