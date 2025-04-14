package internal

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	svix_internal "github.com/svix/svix-webhooks/go/internalapi"
	"github.com/svix/svix-webhooks/go/models"
)

var _ resource.Resource = &ApiTokenResource{}

func NewEnvironmentResource() resource.Resource {
	return &EnvironmentResource{}
}

type EnvironmentResource struct {
	state appState
}

type EnvironmentResourceModel struct {
	CreatedAt timetypes.RFC3339 `tfsdk:"created_at"`
	UpdatedAt timetypes.RFC3339 `tfsdk:"updated_at"`
	Id        types.String      `tfsdk:"id"`
	Region    types.String      `tfsdk:"region"`
	Name      types.String      `tfsdk:"name"`
	Type      types.String      `tfsdk:"type"`
}

func (r *EnvironmentResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	state, ok := req.ProviderData.(appState)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected appState, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}
	r.state = state
}

func (r *EnvironmentResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "svix_environment"
}

func (r *EnvironmentResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{Required: true, Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
				stringvalidator.LengthAtMost(256),
			}},
			"type": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.OneOf("development", "production"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				}},
			// non modifiable fields
			"id":         schema.StringAttribute{Computed: true},
			"region":     schema.StringAttribute{Computed: true},
			"created_at": schema.StringAttribute{Computed: true, CustomType: timetypes.RFC3339Type{}},
			"updated_at": schema.StringAttribute{Computed: true, CustomType: timetypes.RFC3339Type{}},
		},
	}
}

func (r *EnvironmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// load state/plan
	var data EnvironmentResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// create svix client
	svx, err := r.state.internalDefaultSvixClient()
	if err != nil {
		resp.Diagnostics.AddError(UNABLE_TO_CREATE_SVIX_CLIENT, err.Error())
		return
	}

	// call api
	res, err := svx.Management.Environment.Create(
		ctx, models.EnvironmentModelIn{
			Name: data.Name.ValueString(),
			Type: models.EnvironmentType(data.Type.ValueString()),
		},
		&svix_internal.ManagementEnvironmentCreateOptions{IdempotencyKey: randStr32()},
	)
	if err != nil {
		logSvixError(&resp.Diagnostics, err, "Failed to create environment")
		return
	}

	// set the state
	setCreateState(ctx, resp, "id", res.Id)
	setCreateState(ctx, resp, "name", res.Name)
	setCreateState(ctx, resp, "type", res.Type)
	setCreateState(ctx, resp, "region", res.Region)
	setCreateState(ctx, resp, "created_at", timetypes.NewRFC3339TimeValue(res.CreatedAt))
	setCreateState(ctx, resp, "updated_at", timetypes.NewRFC3339TimeValue(res.UpdatedAt))
}

func (r *EnvironmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// load state/plan
	var env_id string
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &env_id)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// create svix client
	svx, err := r.state.internalDefaultSvixClient()
	if err != nil {
		resp.Diagnostics.AddError(UNABLE_TO_CREATE_SVIX_CLIENT, err.Error())
		return
	}

	// call api
	res, err := svx.Management.Environment.Get(ctx, env_id)
	if err != nil {
		logSvixError(&resp.Diagnostics, err, "Failed to read environment")
		return
	}

	// set the state
	setReadState(ctx, resp, "id", res.Id)
	setReadState(ctx, resp, "name", res.Name)
	setReadState(ctx, resp, "type", res.Type)
	setReadState(ctx, resp, "region", res.Region)
	setReadState(ctx, resp, "created_at", timetypes.NewRFC3339TimeValue(res.CreatedAt))
	setReadState(ctx, resp, "updated_at", timetypes.NewRFC3339TimeValue(res.UpdatedAt))
}

func (r *EnvironmentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// load state/plan
	var data EnvironmentResourceModel
	var env_id string
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &env_id)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// create svix client
	svx, err := r.state.internalDefaultSvixClient()
	if err != nil {
		resp.Diagnostics.AddError(UNABLE_TO_CREATE_SVIX_CLIENT, err.Error())
		return
	}

	// call api
	res, err := svx.Management.Environment.Update(ctx, env_id, models.EnvironmentModelUpdate{
		Name: data.Name.ValueString(),
	})
	if err != nil {
		logSvixError(&resp.Diagnostics, err, "Failed to update environment")
		return
	}

	// set the state
	setUpdateState(ctx, resp, "id", res.Id)
	setUpdateState(ctx, resp, "name", res.Name)
	setUpdateState(ctx, resp, "type", res.Type)
	setUpdateState(ctx, resp, "region", res.Region)
	setUpdateState(ctx, resp, "created_at", timetypes.NewRFC3339TimeValue(res.CreatedAt))
	setUpdateState(ctx, resp, "updated_at", timetypes.NewRFC3339TimeValue(res.UpdatedAt))
}

func (r *EnvironmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// load state/plan
	var env_id string
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &env_id)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// create svix client
	svx, err := r.state.internalDefaultSvixClient()
	if err != nil {
		resp.Diagnostics.AddError(UNABLE_TO_CREATE_SVIX_CLIENT, err.Error())
		return
	}

	// call api

	err = svx.Management.Environment.Delete(ctx, env_id)
	if err != nil {
		logSvixError(&resp.Diagnostics, err, "Failed to delete environment")
		return
	}
}
