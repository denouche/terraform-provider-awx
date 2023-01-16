package awx

/*
Base for workflow job template notification template.
*/

import (
	"context"
	"strconv"

	awx "github.com/denouche/goawx/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func getResourceWorkflowJobTemplateNotificationTemplateAssociateFuncForType(client *awx.WorkflowJobTemplateNotificationTemplatesService, typ string) func(workflowJobTemplateID int, notificationTemplateID int) (*awx.NotificationTemplate, error) {
	switch typ {
	case "error":
		return client.AssociateWorkflowJobTemplateNotificationTemplatesError
	case "success":
		return client.AssociateWorkflowJobTemplateNotificationTemplatesSuccess
	case "started":
		return client.AssociateWorkflowJobTemplateNotificationTemplatesStarted
	}
	return nil
}

func getResourceWorkflowJobTemplateNotificationTemplateDisassociateFuncForType(client *awx.WorkflowJobTemplateNotificationTemplatesService, typ string) func(workflowJobTemplateID int, notificationTemplateID int) (*awx.NotificationTemplate, error) {
	switch typ {
	case "error":
		return client.DisassociateWorkflowJobTemplateNotificationTemplatesError
	case "success":
		return client.DisassociateWorkflowJobTemplateNotificationTemplatesSuccess
	case "started":
		return client.DisassociateWorkflowJobTemplateNotificationTemplatesStarted
	}
	return nil
}

func resourceWorkflowJobTemplateNotificationTemplateCreateForType(typ string) func(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return func(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
		var diags diag.Diagnostics
		client := m.(*awx.AWX)
		awxWorkflowJobTemplateService := client.WorkflowJobTemplateService
		workflowJobTemplateID := d.Get("workflow_job_template_id").(int)
		_, err := awxWorkflowJobTemplateService.GetWorkflowJobTemplateByID(workflowJobTemplateID, make(map[string]string))
		if err != nil {
			return buildDiagNotFoundFail("workflow job template", workflowJobTemplateID, err)
		}

		awxWorkflowJobTemplateNotifService := client.WorkflowJobTemplateNotificationTemplatesService
		notificationTemplateID := d.Get("notification_template_id").(int)
		associationFunc := getResourceWorkflowJobTemplateNotificationTemplateAssociateFuncForType(awxWorkflowJobTemplateNotifService, typ)
		if associationFunc == nil {
			return buildDiagnosticsMessage("Create: WorkflowJobTemplate not AssociateWorkflowJobTemplateNotificationTemplates", "Fail to find association function for notification_template type %s", typ)
		}

		result, err := associationFunc(workflowJobTemplateID, notificationTemplateID)
		if err != nil {
			return buildDiagnosticsMessage("Create: WorkflowJobTemplate not AssociateWorkflowJobTemplateNotificationTemplates", "Fail to associate notification_template credentials with ID %v, for workflow_job_template ID %v, got error: %s", notificationTemplateID, workflowJobTemplateID, err.Error())
		}

		d.SetId(strconv.Itoa(result.ID))
		return diags
	}
}

func resourceWorkflowJobTemplateNotificationTemplateRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	return diags
}

func resourceWorkflowJobTemplateNotificationTemplateDeleteForType(typ string) func(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return func(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
		var diags diag.Diagnostics
		client := m.(*awx.AWX)
		awxWorkflowJobTemplateService := client.WorkflowJobTemplateService
		workflowJobTemplateID := d.Get("workflow_job_template_id").(int)
		_, err := awxWorkflowJobTemplateService.GetWorkflowJobTemplateByID(workflowJobTemplateID, make(map[string]string))
		if err != nil {
			return buildDiagNotFoundFail("workflow job template", workflowJobTemplateID, err)
		}

		awxWorkflowJobTemplateNotifService := client.WorkflowJobTemplateNotificationTemplatesService
		notificationTemplateID := d.Get("notification_template_id").(int)
		disassociationFunc := getResourceWorkflowJobTemplateNotificationTemplateDisassociateFuncForType(awxWorkflowJobTemplateNotifService, typ)
		if disassociationFunc == nil {
			return buildDiagnosticsMessage("Create: WorkflowJobTemplate not DisassociateWorkflowJobTemplateNotificationTemplates", "Fail to find disassociation function for notification_template type %s", typ)
		}

		_, err = disassociationFunc(workflowJobTemplateID, notificationTemplateID)
		if err != nil {
			return buildDiagnosticsMessage("Create: WorkflowJobTemplate not DisassociateWorkflowJobTemplateNotificationTemplates", "Fail to associate notification_template credentials with ID %v, for job_template ID %v, got error: %s", notificationTemplateID, workflowJobTemplateID, err.Error())
		}

		d.SetId("")
		return diags
	}
}
