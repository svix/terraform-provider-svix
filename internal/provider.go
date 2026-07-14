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
}

type SvixProviderModel struct {
	ServerUrl types.String `tfsdk:"server_url"`
	Token     types.String `tfsdk:"token"`
}

func (p *SvixProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "svix"
	resp.Version = Version
}

func (p *SvixProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"server_url": schema.StringAttribute{
				MarkdownDescription: "Svix server URL. Required for most resources, but not for `svix_autoconfig` " +
					"(credentials are embedded in the AutoConfig token). Can also be set via the `SVIX_SERVER_URL` environment variable.",
				Optional: true,
			},
			"token": schema.StringAttribute{
				MarkdownDescription: "API token. Required for most resources, but not for `svix_autoconfig` " +
					"(credentials are embedded in the AutoConfig token). Can also be set via the `SVIX_TOKEN` environment variable.",
				Optional:  true,
				Sensitive: true,
			},
		},
	}
}

func (p *SvixProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	token := os.Getenv("SVIX_TOKEN")
	server_url := os.Getenv("SVIX_SERVER_URL")
	var data SvixProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if data.Token.ValueString() != "" {
		token = data.Token.ValueString()
	}
	if data.ServerUrl.ValueString() != "" {
		server_url = data.ServerUrl.ValueString()
	}

	var parsed url.URL
	if server_url != "" {
		u, err := url.Parse(server_url)
		if err != nil {
			resp.Diagnostics.AddError("Unable to parse endpoint url", err.Error())
			return
		}
		parsed = *u
	}

	_, debug := os.LookupEnv("SVIX_DEBUG")

	state := appState{
		token:     token,
		serverUrl: parsed,
		debug:     debug,
	}

	resp.DataSourceData = state
	resp.ResourceData = state
}

func (p *SvixProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewApiTokenResource,
		NewAutoConfigResource,
		NewEnvironmentResource,
		NewEnvironmentSettingsResource,
		NewEventTypeOpenapiImportResource,
		NewEventTypeResource,
		NewOperationalWebhooksEndpoint,
		NewSvixIngestSourceResource,
		NewIngestEndpointResource,
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

func New() func() provider.Provider {
	return func() provider.Provider {
		return &SvixProvider{}
	}
}

type appState struct {
	token     string
	serverUrl url.URL
	debug     bool
}

var userAgentSuffix = fmt.Sprintf("tf-provider-v%s", Version)

func (s *appState) requireAuth() error {
	if s.token == "" {
		return fmt.Errorf(
			"API token is required for this resource; set the provider token attribute or the SVIX_TOKEN environment variable",
		)
	}
	if s.serverUrl.String() == "" {
		return fmt.Errorf(
			"server URL is required for this resource; set the provider server_url attribute or the SVIX_SERVER_URL environment variable",
		)
	}
	return nil
}

// get the default client without an envId suffixed
func (s *appState) DefaultSvixClient() (*svix.Svix, error) {
	if err := s.requireAuth(); err != nil {
		return nil, err
	}
	svx, err := svix.New(s.token, &svix.SvixOptions{ServerUrl: &s.serverUrl, Debug: s.debug})
	if err != nil {
		return nil, err
	}
	err = svx.SetUserAgentSuffix(userAgentSuffix)
	if err != nil {
		return nil, err
	}
	return svx, nil
}

// create a new svix client with the envId suffixed on the token
func (s *appState) ClientWithEnvId(envId string) (*svix.Svix, error) {
	if err := s.requireAuth(); err != nil {
		return nil, err
	}
	bearerToken := fmt.Sprintf("%s|%s", s.token, envId)
	svx, err := svix.New(bearerToken, &svix.SvixOptions{ServerUrl: &s.serverUrl, Debug: s.debug})
	if err != nil {
		return nil, err
	}
	err = svx.SetUserAgentSuffix(userAgentSuffix)
	if err != nil {
		return nil, err
	}
	return svx, nil
}

// create a new internal svix client with the envId suffixed on the token
func (s *appState) InternalClientWithEnvId(envId string) (*svix_internal.InternalSvix, error) {
	if err := s.requireAuth(); err != nil {
		return nil, err
	}
	bearerToken := fmt.Sprintf("%s|%s", s.token, envId)
	svx, err := svix_internal.New(bearerToken, &s.serverUrl, s.debug, &userAgentSuffix)
	if err != nil {
		return nil, err
	}
	return svx, nil
}

// get the default internal svix client without an envId suffixed
func (s *appState) InternalDefaultSvixClient() (*svix_internal.InternalSvix, error) {
	if err := s.requireAuth(); err != nil {
		return nil, err
	}
	svx, err := svix_internal.New(s.token, &s.serverUrl, s.debug, &userAgentSuffix)
	if err != nil {
		return nil, err
	}
	return svx, nil
}
