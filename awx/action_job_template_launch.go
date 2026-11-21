package awx

import (
	"context"
	"fmt"
	"time"

	awx "github.com/denouche/goawx/client"
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/action/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const defaultJobTemplateLaunchActionTimeout = 20 * time.Minute

// jobTemplateLaunchAction launches a job template without keeping anything in
// state. Unlike the awx_job_template_launch resource, there is nothing to
// destroy or replace later, so AWX purging the job does not affect it.
type jobTemplateLaunchAction struct {
	client *awx.AWX
}

var (
	_ action.Action              = (*jobTemplateLaunchAction)(nil)
	_ action.ActionWithConfigure = (*jobTemplateLaunchAction)(nil)
)

// NewJobTemplateLaunchAction returns the awx_job_template_launch action.
func NewJobTemplateLaunchAction() action.Action {
	return &jobTemplateLaunchAction{}
}

type jobTemplateLaunchActionModel struct {
	JobTemplateID     types.Int64  `tfsdk:"job_template_id"`
	Limit             types.String `tfsdk:"limit"`
	InventoryID       types.Int64  `tfsdk:"inventory_id"`
	ExtraVars         types.String `tfsdk:"extra_vars"`
	WaitForCompletion types.Bool   `tfsdk:"wait_for_completion"`
	Timeout           types.String `tfsdk:"timeout"`
}

func (a *jobTemplateLaunchAction) Metadata(_ context.Context, req action.MetadataRequest, resp *action.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_job_template_launch"
}

func (a *jobTemplateLaunchAction) Schema(_ context.Context, _ action.SchemaRequest, resp *action.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Launches a job from an AWX job template. Nothing is stored in state.",
		Attributes: map[string]schema.Attribute{
			"job_template_id": schema.Int64Attribute{
				Required:    true,
				Description: "Job template ID",
			},
			"limit": schema.StringAttribute{
				Optional:    true,
				Description: "List of comma delimited hosts to limit job execution. Required ask_limit_on_launch set on job_template.",
			},
			"inventory_id": schema.Int64Attribute{
				Optional:    true,
				Description: "Override Inventory ID. Required ask_inventory_on_launch set on job_template.",
			},
			"extra_vars": schema.StringAttribute{
				Optional:    true,
				Description: "Override job template variables. YAML or JSON values are supported.",
			},
			"wait_for_completion": schema.BoolAttribute{
				Optional:    true,
				Description: "Wait for the job to finish, failing the action if the job fails. Defaults to true.",
			},
			"timeout": schema.StringAttribute{
				Optional:    true,
				Description: "How long to wait for the job when wait_for_completion is true, as a Go duration such as \"90m\". Defaults to \"20m\".",
			},
		},
	}
}

func (a *jobTemplateLaunchAction) Configure(_ context.Context, req action.ConfigureRequest, resp *action.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*awx.AWX)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected provider data",
			fmt.Sprintf("Expected *client.AWX, got %T.", req.ProviderData),
		)
		return
	}
	a.client = client
}

func (a *jobTemplateLaunchAction) Invoke(ctx context.Context, req action.InvokeRequest, resp *action.InvokeResponse) {
	var config jobTemplateLaunchActionModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if a.client == nil {
		resp.Diagnostics.AddError(
			"AWX provider is not configured",
			"The provider configuration was not known when the action ran.",
		)
		return
	}

	timeout := defaultJobTemplateLaunchActionTimeout
	if !config.Timeout.IsNull() {
		parsed, err := time.ParseDuration(config.Timeout.ValueString())
		if err != nil {
			resp.Diagnostics.AddAttributeError(
				path.Root("timeout"), "Invalid timeout",
				fmt.Sprintf("%q is not a Go duration: %s", config.Timeout.ValueString(), err),
			)
			return
		}
		timeout = parsed
	}

	jobTemplateID := int(config.JobTemplateID.ValueInt64())
	if _, err := a.client.JobTemplateService.GetJobTemplateByID(jobTemplateID, map[string]string{}); err != nil {
		resp.Diagnostics.AddError(
			"Unable to fetch job template",
			fmt.Sprintf("Unable to load job template with id %d: got %s", jobTemplateID, err),
		)
		return
	}

	job, err := launchJobTemplate(a.client.JobTemplateService, jobTemplateID, JobTemplateLaunchData{
		Limit:       config.Limit.ValueString(),
		InventoryID: int(config.InventoryID.ValueInt64()),
		ExtraVars:   config.ExtraVars.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to launch job template",
			fmt.Sprintf("Launching job template %d failed: %s", jobTemplateID, err),
		)
		return
	}
	sendProgress(resp, fmt.Sprintf("Launched job %d from job template %d", job.ID, jobTemplateID))

	if !config.WaitForCompletion.IsNull() && !config.WaitForCompletion.ValueBool() {
		return
	}

	err = waitForJob(ctx, a.client.JobService, job.ID, timeout, func(status string) {
		sendProgress(resp, fmt.Sprintf("Job %d is %s", job.ID, status))
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Job template execution failure",
			fmt.Sprintf("Job %d from job template %d did not succeed: %s", job.ID, jobTemplateID, err),
		)
	}
}

func sendProgress(resp *action.InvokeResponse, message string) {
	if resp.SendProgress != nil {
		resp.SendProgress(action.InvokeProgressEvent{Message: message})
	}
}
