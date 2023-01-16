/*
Associate a job template and a credential.

Example Usage

```hcl
resource "awx_job_template_credential" "baseconfig" {
    job_template_id = data.awx_job_template.baseconfig.id
    credential_id   = data.awx_credential_machine.pi_connection.id
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

func resourceJobTemplateCredentials() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceJobTemplateCredentialsCreate,
		DeleteContext: resourceJobTemplateCredentialsDelete,
		ReadContext:   resourceJobTemplateCredentialsRead,

		Schema: map[string]*schema.Schema{
			"job_template_id": {
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

func resourceJobTemplateCredentialsCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*awx.AWX)
	awxService := client.JobTemplateService
	jobTemplateID := d.Get("job_template_id").(int)
	_, err := awxService.GetJobTemplateByID(jobTemplateID, make(map[string]string))
	if err != nil {
		return buildDiagNotFoundFail("job template", jobTemplateID, err)
	}

	result, err := awxService.AssociateCredentials(jobTemplateID, map[string]interface{}{
		"id": d.Get("credential_id").(int),
	}, map[string]string{})

	if err != nil {
		return buildDiagnosticsMessage("Create: JobTemplate not AssociateCredentials", "Fail to add credentials with Id %v, for Template ID %v, got error: %s", d.Get("credential_id").(int), jobTemplateID, err.Error())
	}

	d.SetId(strconv.Itoa(result.ID))
	return diags
}

func resourceJobTemplateCredentialsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	return diags
}

func resourceJobTemplateCredentialsDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*awx.AWX)
	awxService := client.JobTemplateService
	jobTemplateID := d.Get("job_template_id").(int)
	res, err := awxService.GetJobTemplateByID(jobTemplateID, make(map[string]string))
	if err != nil {
		return buildDiagNotFoundFail("job template", jobTemplateID, err)
	}

	_, err = awxService.DisAssociateCredentials(res.ID, map[string]interface{}{
		"id": d.Get("credential_id").(int),
	}, map[string]string{})
	if err != nil {
		return buildDiagDeleteFail("JobTemplate DisAssociateCredentials", fmt.Sprintf("DisAssociateCredentials %v, from JobTemplateID %v got %s ", d.Get("credential_id").(int), d.Get("job_template_id").(int), err.Error()))
	}

	d.SetId("")
	return diags
}
