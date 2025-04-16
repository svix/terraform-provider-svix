package internal

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	svix "github.com/svix/svix-webhooks/go"
	"github.com/svix/svix-webhooks/go/models"
)

var _ resource.Resource = &EventTypeResource{}

func NewEventTypeResource() resource.Resource {
	return &EventTypeResource{}
}

type EventTypeResource struct {
	state appState
}

type EventTypeResourceModel struct {
	EnvironmentId types.String         `tfsdk:"environment_id"`
	Archived      types.Bool           `tfsdk:"archived"`
	CreatedAt     timetypes.RFC3339    `tfsdk:"created_at"`
	Deprecated    types.Bool           `tfsdk:"deprecated"`
	Description   types.String         `tfsdk:"description"`
	FeatureFlag   types.String         `tfsdk:"feature_flag"`
	GroupName     types.String         `tfsdk:"group_name"`
	Name          types.String         `tfsdk:"name"`
	Schemas       jsontypes.Normalized `tfsdk:"schemas"`
	UpdatedAt     timetypes.RFC3339    `tfsdk:"updated_at"`
}

func (r *EventTypeResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "svix_event_type"
}

func (r *EventTypeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"environment_id": schema.StringAttribute{
				Required:    true,
				Description: ENV_ID_DESC,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"archived": schema.BoolAttribute{Computed: true, Optional: true, Default: booldefault.StaticBool(false)},
			"created_at": schema.StringAttribute{
				Computed:   true,
				CustomType: timetypes.RFC3339Type{},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"deprecated":  schema.BoolAttribute{Computed: true, Optional: true, Default: booldefault.StaticBool(false)},
			"description": schema.StringAttribute{Required: true},
			"feature_flag": schema.StringAttribute{Optional: true, Validators: []validator.String{
				stringvalidator.LengthAtMost(256),
				stringvalidator.RegexMatches(saneStringRegex(), "String must match against `^[a-zA-Z0-9\\-_.]+$`"),
			}},
			"group_name": schema.StringAttribute{Optional: true, Validators: []validator.String{
				stringvalidator.LengthAtMost(256),
				stringvalidator.RegexMatches(saneStringRegex(), "String must match against `^[a-zA-Z0-9\\-_.]+$`"),
			}},
			"name": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(256),
					stringvalidator.RegexMatches(saneStringRegex(), "String must match against `^[a-zA-Z0-9\\-_.]+$`"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"schemas": schema.StringAttribute{Optional: true, CustomType: jsontypes.NormalizedType{}},
			"updated_at": schema.StringAttribute{
				Computed:   true,
				CustomType: timetypes.RFC3339Type{},
			},
		},
	}
}

func (r *EventTypeResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *EventTypeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// load state/plan
	var data EventTypeResourceModel
	var envId string
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("environment_id"), &envId)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// create svix client
	svx, err := r.state.ClientWithEnvId(envId)
	if err != nil {
		resp.Diagnostics.AddError(UNABLE_TO_CREATE_SVIX_CLIENT, err.Error())
		return
	}

	// create the EventTypeIn struct
	var schema *map[string]any
	if data.Schemas.IsNull() {
		schema = nil
	} else {
		resp.Diagnostics.Append(data.Schemas.Unmarshal(&schema)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}
	eventTypeIn := models.EventTypeIn{
		Name:        data.Name.ValueString(),
		Description: data.Description.ValueString(),
		Archived:    boolOrNil(data.Archived),
		Deprecated:  boolOrNil(data.Deprecated),
		Schemas:     schema,
		FeatureFlag: strOrNil(data.FeatureFlag),
		GroupName:   strOrNil(data.GroupName),
	}
	reqOpts := svix.EventTypeCreateOptions{
		IdempotencyKey: randStr32(),
	}

	// call api
	res, err := svx.EventType.Create(ctx, eventTypeIn, &reqOpts)
	if err != nil {
		logSvixError(&resp.Diagnostics, err, "Failed to create event type")
		return
	}

	// save state
	var schemasJson *string
	if res.Schemas != nil {
		jsonV, err := json.Marshal(res.Schemas)
		if err != nil {
			resp.Diagnostics.AddAttributeError(
				path.Root("schemas"),
				"Failed to marshal a map[string]any to a string",
				err.Error(),
			)
			return
		}
		schemasJson = ptr(string(jsonV))
	}

	setCreateState(ctx, resp, rp("environment_id"), envId)
	setCreateState(ctx, resp, rp("archived"), res.Archived)
	setCreateState(ctx, resp, rp("created_at"), timetypes.NewRFC3339TimeValue(res.CreatedAt))
	setCreateState(ctx, resp, rp("deprecated"), res.Deprecated)
	setCreateState(ctx, resp, rp("description"), res.Description)
	setCreateState(ctx, resp, rp("feature_flag"), res.FeatureFlag)
	setCreateState(ctx, resp, rp("group_name"), res.GroupName)
	setCreateState(ctx, resp, rp("name"), res.Name)
	setCreateState(ctx, resp, rp("schemas"), jsontypes.NewNormalizedPointerValue(schemasJson))
	setCreateState(ctx, resp, rp("updated_at"), timetypes.NewRFC3339TimeValue(res.CreatedAt))
}

func (r *EventTypeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// load state/plan
	var envId, name string
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("environment_id"), &envId)...)
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("name"), &name)...)

	// create svix client
	svx, err := r.state.ClientWithEnvId(envId)
	if err != nil {
		resp.Diagnostics.AddError(UNABLE_TO_CREATE_SVIX_CLIENT, err.Error())
		return
	}

	// call api
	res, err := svx.EventType.Get(ctx, name)
	if err != nil {
		logSvixError(&resp.Diagnostics, err, "Failed to read event type")
		return
	}

	// save state
	var schemasJson *string
	if res.Schemas != nil {
		jsonV, err := json.Marshal(res.Schemas)
		if err != nil {
			resp.Diagnostics.AddAttributeError(
				path.Root("schemas"),
				"Failed to marshal a map[string]any to a string",
				err.Error(),
			)
			return
		}
		schemasJson = ptr(string(jsonV))
	}

	setReadState(ctx, resp, rp("environment_id"), envId)
	setReadState(ctx, resp, rp("archived"), res.Archived)
	setReadState(ctx, resp, rp("created_at"), timetypes.NewRFC3339TimeValue(res.CreatedAt))
	setReadState(ctx, resp, rp("deprecated"), res.Deprecated)
	setReadState(ctx, resp, rp("description"), res.Description)
	setReadState(ctx, resp, rp("feature_flag"), res.FeatureFlag)
	setReadState(ctx, resp, rp("group_name"), res.GroupName)
	setReadState(ctx, resp, rp("name"), res.Name)
	setReadState(ctx, resp, rp("schemas"), jsontypes.NewNormalizedPointerValue(schemasJson))
	setReadState(ctx, resp, rp("updated_at"), timetypes.NewRFC3339TimeValue(res.CreatedAt))
}

func (r *EventTypeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// load state/plan
	var data EventTypeResourceModel
	var envId string
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("environment_id"), &envId)...)

	// create svix client
	svx, err := r.state.ClientWithEnvId(envId)
	if err != nil {
		resp.Diagnostics.AddError(UNABLE_TO_CREATE_SVIX_CLIENT, err.Error())
		return
	}

	// create EventTypeUpdate model
	schemas := stringToMapStringT[any](&resp.Diagnostics, data.Schemas.ValueStringPointer())
	if resp.Diagnostics.HasError() {
		return
	}
	eventType := models.EventTypeUpdate{
		Archived:    boolOrNil(data.Archived),
		Deprecated:  boolOrNil(data.Deprecated),
		Description: data.Description.ValueString(),
		FeatureFlag: strOrNil(data.FeatureFlag),
		GroupName:   strOrNil(data.GroupName),
		Schemas:     schemas,
	}

	// call api
	res, err := svx.EventType.Update(ctx, data.Name.ValueString(), eventType)
	if err != nil {
		logSvixError(&resp.Diagnostics, err, "Failed to update event type")
		return
	}

	// save state
	var schemasJson *string
	if res.Schemas != nil {
		jsonV, err := json.Marshal(res.Schemas)
		if err != nil {
			resp.Diagnostics.AddAttributeError(
				path.Root("schemas"),
				"Failed to marshal a map[string]any to a string",
				err.Error(),
			)
			return
		}
		schemasJson = ptr(string(jsonV))
	}

	setUpdateState(ctx, resp, rp("schemas"), jsontypes.NewNormalizedPointerValue(schemasJson))
	setUpdateState(ctx, resp, rp("archived"), res.Archived)
	setUpdateState(ctx, resp, rp("deprecated"), res.Deprecated)
	setUpdateState(ctx, resp, rp("description"), res.Description)
	setUpdateState(ctx, resp, rp("feature_flag"), res.FeatureFlag)
	setUpdateState(ctx, resp, rp("group_name"), res.GroupName)

}

func (r *EventTypeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// load state/plan
	var envId, name string
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("name"), &name)...)
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("environment_id"), &envId)...)

	// create svix client
	svx, err := r.state.ClientWithEnvId(envId)
	if err != nil {
		resp.Diagnostics.AddError(UNABLE_TO_CREATE_SVIX_CLIENT, err.Error())
		return
	}

	err = svx.EventType.Delete(ctx, name, &svix.EventTypeDeleteOptions{
		Expunge: ptr(false),
	})

	if err != nil {
		logSvixError(&resp.Diagnostics, err, "Failed to delete event type")
		return
	}
}
