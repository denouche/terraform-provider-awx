package awx

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// frameworkProvider serves the parts of the provider that need
// terraform-plugin-framework, such as actions. main.go muxes it with the
// SDKv2 provider, so its schema must stay identical to Provider()'s.
type frameworkProvider struct{}

var (
	_ provider.Provider            = (*frameworkProvider)(nil)
	_ provider.ProviderWithActions = (*frameworkProvider)(nil)
)

// NewFrameworkProvider returns the plugin-framework half of the provider.
func NewFrameworkProvider() provider.Provider {
	return &frameworkProvider{}
}

type frameworkProviderModel struct {
	Hostname types.String `tfsdk:"hostname"`
	Insecure types.Bool   `tfsdk:"insecure"`
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
	Token    types.String `tfsdk:"token"`
}

func (p *frameworkProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "awx"
}

func (p *frameworkProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"hostname": schema.StringAttribute{
				Optional: true,
			},
			"insecure": schema.BoolAttribute{
				Optional:    true,
				Description: "Disable SSL verification of API calls",
			},
			"username": schema.StringAttribute{
				Optional: true,
			},
			"password": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
			},
			"token": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
			},
		},
	}
}

// stringOrEnv mirrors schema.EnvDefaultFunc in the SDKv2 provider.
func stringOrEnv(v types.String, env, fallback string) string {
	if !v.IsNull() && !v.IsUnknown() {
		return v.ValueString()
	}
	if e := os.Getenv(env); e != "" {
		return e
	}
	return fallback
}

func (p *frameworkProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config frameworkProviderModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Values that are unknown until apply (for example from another
	// resource) leave the client unset; actions report it when invoked.
	if config.Hostname.IsUnknown() || config.Username.IsUnknown() || config.Password.IsUnknown() ||
		config.Token.IsUnknown() || config.Insecure.IsUnknown() {
		return
	}

	client, err := newAWXClient(
		stringOrEnv(config.Hostname, "AWX_HOSTNAME", "http://localhost"),
		stringOrEnv(config.Username, "AWX_USERNAME", "admin"),
		stringOrEnv(config.Password, "AWX_PASSWORD", "password"),
		stringOrEnv(config.Token, "AWX_TOKEN", ""),
		config.Insecure.ValueBool(),
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create AWX client",
			"Unable to auth user against AWX API: check the hostname, username and password",
		)
		return
	}

	resp.ActionData = client
}

func (p *frameworkProvider) Resources(context.Context) []func() resource.Resource {
	return nil
}

func (p *frameworkProvider) DataSources(context.Context) []func() datasource.DataSource {
	return nil
}

func (p *frameworkProvider) Actions(context.Context) []func() action.Action {
	return []func() action.Action{
		NewJobTemplateLaunchAction,
	}
}
