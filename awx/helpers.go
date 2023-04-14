package awx

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	awx "github.com/denouche/goawx/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"gopkg.in/yaml.v2"
)

const (
	diagElementInventoryGroupTitle  = "Inventory Group"
	diagElementInventorySourceTitle = "Inventory Source"
	diagElementInventoryTitle       = "Inventory"
	diagElementHostTitle            = "Host"
)

func resourceJobTemplateDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*awx.AWX)
	awxService := client.JobTemplateService
	id, diags := convertStateIDToNummeric("Delete JobTemplate", d)
	if diags.HasError() {
		return diags
	}
	_, err := awxService.DeleteJobTemplate(id)
	if err != nil {
		return buildDiagDeleteFail(
			"JobTemplate",
			fmt.Sprintf(
				"JobTemplateID %v, got %s ",
				id, err.Error(),
			),
		)
	}
	d.SetId("")
	return nil
}
func buildDiagCreateFail(tfMethode string, err error) diag.Diagnostics {
	return buildDiagnosticsMessage(
		fmt.Sprintf("Unable to create %s", tfMethode),
		"Unable to create %s got %s",
		tfMethode, err.Error(),
	)
}
func buildDiagUpdateFail(tfMethode string, id int, err error) diag.Diagnostics {
	return buildDiagnosticsMessage(
		fmt.Sprintf("Unable to update %s", tfMethode),
		"Unable to update %s with id %d: got %s",
		tfMethode, id, err.Error(),
	)
}

func buildDiagNotFoundFail(tfMethode string, id int, err error) diag.Diagnostics {
	return buildDiagnosticsMessage(
		fmt.Sprintf("Unable to fetch %s", tfMethode),
		"Unable to load %s with id %d: got %s",
		tfMethode, id, err.Error(),
	)
}

func buildDiagDeleteFail(tfMethode, details string) diag.Diagnostics {
	return buildDiagnosticsMessage(
		buildDiagDeleteFailSummary(tfMethode),
		buildDiagDeleteFailDetails(tfMethode, details),
	)
}
func buildDiagDeleteFailSummary(tfMethode string) string {
	return fmt.Sprintf("%s delete failed", tfMethode)
}
func buildDiagDeleteFailDetails(tfMethode, detailsString string) string {
	return fmt.Sprintf("Fail to delete %s, %s", tfMethode, detailsString)
}

func convertStateIDToNummeric(tfElement string, d *schema.ResourceData) (int, diag.Diagnostics) {
	var diags diag.Diagnostics
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return id, buildDiagnosticsMessage(
			fmt.Sprintf("%s, State ID Not Converted", tfElement),
			"Value in State %s is`t nummeric, %s",
			d.Id(), err.Error(),
		)
	}
	return id, diags
}

func buildDiagnosticsMessage(diagSummary, diagDetails string, detailsVars ...interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	diags = append(diags, diag.Diagnostic{
		Severity: diag.Error,
		Summary:  diagSummary,
		Detail:   fmt.Sprintf(diagDetails, detailsVars...),
	})
	return diags
}

func CredentialsServiceDeleteByID(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	id, _ := strconv.Atoi(d.Id())
	client := m.(*awx.AWX)
	err := client.CredentialsService.DeleteCredentialsByID(id, map[string]string{})
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to delete existing credentials",
			Detail:   fmt.Sprintf("Unable to delete existing credentials with id %d: %s", id, err.Error()),
		})
	}

	return diags
}
func CredentialTypeServiceDeleteByID(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	id, _ := strconv.Atoi(d.Id())
	client := m.(*awx.AWX)
	err := client.CredentialTypeService.DeleteCredentialTypeByID(id, map[string]string{})
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to delete existing credential type",
			Detail:   fmt.Sprintf("Unable to delete existing credential type with id %d: %s", id, err.Error()),
		})
	}

	return diags
}

func normalizeJsonYaml(s interface{}) string {
	result := string("")
	if j, ok := normalizeJsonOk(s); ok {
		result = j
	} else if y, ok := normalizeYamlOk(s); ok {
		result = y
	} else {
		result = s.(string)
	}
	return result
}
func normalizeJsonOk(s interface{}) (string, bool) {
	if s == nil || s == "" {
		return "", true
	}
	var j interface{}
	err := json.Unmarshal([]byte(s.(string)), &j)
	if err != nil {
		return fmt.Sprintf("Error parsing JSON: %s", err), false
	}
	b, _ := json.Marshal(j)
	return string(b[:]), true
}

func normalizeYamlOk(s interface{}) (string, bool) {
	if s == nil || s == "" {
		return "", true
	}
	var j interface{}
	err := yaml.Unmarshal([]byte(s.(string)), &j)
	if err != nil {
		return fmt.Sprintf("Error parsing YAML: %s", err), false
	}
	b, _ := yaml.Marshal(j)
	return string(b[:]), true
}

func unmarshalYaml(str string) map[string]interface{} {
	asMap := map[string]interface{}{}
	err := yaml.Unmarshal([]byte(str), &asMap)
	if err != nil {
		asMap = nil
	}
	return asMap
}

func marshalYaml(v interface{}) string {
	extraDataBytes, err := yaml.Marshal(v)
	if err != nil {
		return string(extraDataBytes)
	}
	return ""
}
