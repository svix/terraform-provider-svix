package internal

import (
	"context"
	"net/url"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	svix "github.com/svix/svix-webhooks/go"
)

var _ provider.Provider = &SvixProvider{}
var _ provider.ProviderWithFunctions = &SvixProvider{}
var _ provider.ProviderWithEphemeralResources = &SvixProvider{}

type SvixProvider struct {
	version string
}

type SvixProviderModel struct {
	ServerUrl types.String `tfsdk:"server_url"`
	Token     types.String `tfsdk:"token"`
}

func (p *SvixProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "svix"
	resp.Version = p.version
}

func (p *SvixProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"server_url": schema.StringAttribute{
				MarkdownDescription: "Svix server url",
				Optional:            true,
			},
			"token": schema.StringAttribute{
				MarkdownDescription: "Api token",
				Optional:            true,
				Sensitive:           true,
			},
		},
	}
}

func (p *SvixProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	token := os.Getenv("SVIX_TOKEN")
	server_url := os.Getenv("SVIX_SERVER_URL")
	var data SvixProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if data.Token.ValueString() != "" {
		token = data.Token.ValueString()
	}
	if data.ServerUrl.ValueString() != "" {
		server_url = data.ServerUrl.ValueString()
	}

	if token == "" {
		resp.Diagnostics.AddError(
			"Missing API Token Configuration",
			"While configuring the provider, the API token was not found in "+
				"the SVIX_TOKEN environment variable or provider "+
				"configuration block token attribute.",
		)
	}
	if server_url == "" {
		resp.Diagnostics.AddError(
			"Missing Server URL Configuration",
			"While configuring the provider, the API token was not found in "+
				"the SVIX_SERVER_URL environment variable or provider "+
				"configuration block server_url attribute.",
		)
	}
	if resp.Diagnostics.HasError() {
		return
	}

	// Example client configuration for data sources and resources
	url, err := url.Parse(server_url)
	if err != nil {
		resp.Diagnostics.AddError("Unable to parse endpoint url", err.Error())
	}
	svx, err := svix.New(token, &svix.SvixOptions{ServerUrl: url})
	if err != nil {
		resp.Diagnostics.AddError("Unable to create svix client", err.Error())
	}
	if resp.Diagnostics.HasError() {
		return
	}
	resp.DataSourceData = svx
	resp.ResourceData = svx

}

func (p *SvixProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewEventTypeResource,
		NewOperationalWebhooksEndpoint,
	}
}
func (p *SvixProvider) EphemeralResources(ctx context.Context) []func() ephemeral.EphemeralResource {
	return nil
}

func (p *SvixProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return nil
}

func (p *SvixProvider) Functions(ctx context.Context) []func() function.Function {
	return nil
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &SvixProvider{
			version: version,
		}
	}
}
