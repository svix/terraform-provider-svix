package internal

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	svix "github.com/svix/svix-webhooks/go"
	"github.com/svix/svix-webhooks/go/models"
)

var _ resource.Resource = &ApiTokenResource{}

func NewApiTokenResource() resource.Resource {
	return &ApiTokenResource{}
}

type ApiTokenResource struct {
	state appState
}
type ApiTokenResourceModel struct {
	EnvironmentId types.String      `tfsdk:"environment_id"`
	Name          types.String      `tfsdk:"name"`
	Scopes        types.List        `tfsdk:"scopes"`
	Token         types.String      `tfsdk:"token"`
	Id            types.String      `tfsdk:"id"`
	CreatedAt     timetypes.RFC3339 `tfsdk:"created_at"`
	ExpiresAt     timetypes.RFC3339 `tfsdk:"expires_at"`
}

func (r *ApiTokenResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "svix_api_token"
}

func (r *ApiTokenResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	state, ok := req.ProviderData.(appState)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *svix.Svix, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}
	r.state = state
}

func (r *ApiTokenResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"environment_id": schema.StringAttribute{
				Required:    true,
				Description: ENV_ID_DESC,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"scopes": schema.ListAttribute{
				Optional:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
				},
			},
			// non modifiable fields
			"token":      schema.StringAttribute{Sensitive: true, Computed: true, Description: "The api token"},
			"id":         schema.StringAttribute{Computed: true},
			"created_at": schema.StringAttribute{Computed: true, CustomType: timetypes.RFC3339Type{}},
			"expires_at": schema.StringAttribute{Computed: true, CustomType: timetypes.RFC3339Type{}},
		},
	}
}

func (r *ApiTokenResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// load state/plan
	var data ApiTokenResourceModel
	var env_id string
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("environment_id"), &env_id)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// create svix client
	svx, err := r.state.clientWithEnvId(env_id)
	if err != nil {
		resp.Diagnostics.AddError(UNABLE_TO_CREATE_SVIX_CLIENT, err.Error())
		return
	}

	// call api
	var scopes []string
	resp.Diagnostics.Append(data.Scopes.ElementsAs(ctx, &scopes, false)...)
	if resp.Diagnostics.HasError() {
		return
	}
	res, err := svx.Management.Authentication.CreateApiToken(
		ctx,
		models.ApiTokenIn{
			Name:   data.Name.ValueString(),
			Scopes: scopes,
		},
		&svix.ManagementAuthenticationCreateApiTokenOptions{
			IdempotencyKey: randStr32(),
		},
	)

	if err != nil {
		logSvixError(&resp.Diagnostics, err, "Unable to create api token")
		return
	}

	// set the state
	setCreateState(ctx, resp, "environment_id", env_id)
	setCreateState(ctx, resp, "name", res.Name)
	setCreateState(ctx, resp, "scopes", res.Scopes)
	setCreateState(ctx, resp, "token", res.Token)
	setCreateState(ctx, resp, "id", res.Id)
	setCreateState(ctx, resp, "created_at", timetypes.NewRFC3339TimeValue(res.CreatedAt))
	setCreateState(ctx, resp, "expires_at", timetypes.NewRFC3339TimePointerValue(res.ExpiresAt))
}

func (r *ApiTokenResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// TODO(10202) read the actual state of the API token
	var stateData ApiTokenResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &stateData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// set the state
	setReadState(ctx, resp, "environment_id", stateData.EnvironmentId)
	setReadState(ctx, resp, "name", stateData.Name)
	setReadState(ctx, resp, "scopes", stateData.Scopes)
	setReadState(ctx, resp, "token", stateData.Token)
	setReadState(ctx, resp, "id", stateData.Id)
	setReadState(ctx, resp, "created_at", stateData.CreatedAt)
	setReadState(ctx, resp, "expires_at", stateData.ExpiresAt)
}

func (r *ApiTokenResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// load state/plan
	var env_id, key_id string
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("environment_id"), &env_id)...)
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &key_id)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// create svix client
	svx, err := r.state.clientWithEnvId(env_id)
	if err != nil {
		resp.Diagnostics.AddError(UNABLE_TO_CREATE_SVIX_CLIENT, err.Error())
		return
	}

	// call api
	err = svx.Management.Authentication.ExpireApiToken(
		ctx, key_id, models.ApiTokenExpireIn{
			Expiry: ptr(int32(0)),
		},
		&svix.ManagementAuthenticationExpireApiTokenOptions{
			IdempotencyKey: randStr32(),
		},
	)

	if err != nil {
		logSvixError(&resp.Diagnostics, err, "Unable to expire api token")
		return
	}

}

func (r *ApiTokenResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError("Terraform tried to update the `svix_api_token` resource, this should not be possible. please contact the developers", "")
}
