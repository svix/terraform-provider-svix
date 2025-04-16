package internal

import (
	"context"
	"fmt"
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
	svix_internal "github.com/svix/svix-webhooks/go/internalapi"
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
			"While configuring the provider, the Server URL was not found in "+
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
		return
	}

	appState := appState{
		token:     token,
		serverUrl: *url,
	}

	resp.DataSourceData = appState
	resp.ResourceData = appState

}

func (p *SvixProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewApiTokenResource,
		NewEnvironmentResource,
		NewEventTypeOpenapiImportResource,
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

type appState struct {
	token     string
	serverUrl url.URL
}

// get the default client without an envId suffixed
func (s *appState) DefaultSvixClient() (*svix.Svix, error) {
	svx, err := svix.New(s.token, &svix.SvixOptions{ServerUrl: &s.serverUrl, Debug: true})
	if err != nil {
		return nil, err
	}
	return svx, nil
}

// create a new svix client with the envId suffixed on the token
func (s *appState) ClientWithEnvId(envId string) (*svix.Svix, error) {
	bearerToken := fmt.Sprintf("%s|%s", s.token, envId)
	svx, err := svix.New(bearerToken, &svix.SvixOptions{ServerUrl: &s.serverUrl, Debug: true})
	if err != nil {
		return nil, err
	}
	return svx, nil

}

// create a new internal svix client with the envId suffixed on the token
func (s *appState) InternalClientWithEnvId(envId string) (*svix_internal.InternalSvix, error) {
	bearerToken := fmt.Sprintf("%s|%s", s.token, envId)
	svx, err := svix_internal.New(bearerToken, &s.serverUrl, true)
	if err != nil {
		return nil, err
	}
	return svx, nil

}

// get the default internal svix client without an envId suffixed
func (s *appState) InternalDefaultSvixClient() (*svix_internal.InternalSvix, error) {
	svx, err := svix_internal.New(s.token, &s.serverUrl, true)
	if err != nil {
		return nil, err
	}
	return svx, nil
}
