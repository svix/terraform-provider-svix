package internal

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	svix "github.com/svix/svix-webhooks/go"
	"github.com/svix/svix-webhooks/go/models"
)

var _ resource.Resource = &EventTypeResource{}

func NewEventTypeResource() resource.Resource {
	return &EventTypeResource{}
}

type EventTypeResource struct {
	svx *svix.Svix
}

type EventTypeResourceModel struct {
	Archived    types.Bool           `tfsdk:"archived"`
	CreatedAt   timetypes.RFC3339    `tfsdk:"created_at"`
	Deprecated  types.Bool           `tfsdk:"deprecated"`
	Description types.String         `tfsdk:"description"`
	FeatureFlag types.String         `tfsdk:"feature_flag"`
	GroupName   types.String         `tfsdk:"group_name"`
	Name        types.String         `tfsdk:"name"`
	Schemas     jsontypes.Normalized `tfsdk:"schemas"`
	UpdatedAt   timetypes.RFC3339    `tfsdk:"updated_at"`
}

func (r *EventTypeResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "svix_event_type"
}

func (r *EventTypeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"archived":    schema.BoolAttribute{Computed: true, Optional: true, Default: booldefault.StaticBool(false)},
			"created_at":  schema.StringAttribute{Computed: true, CustomType: timetypes.RFC3339Type{}},
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
			"name": schema.StringAttribute{Required: true, Validators: []validator.String{
				stringvalidator.LengthAtMost(256),
				stringvalidator.RegexMatches(saneStringRegex(), "String must match against `^[a-zA-Z0-9\\-_.]+$`"),
			}},
			"schemas":    schema.StringAttribute{Optional: true, CustomType: jsontypes.NormalizedType{}},
			"updated_at": schema.StringAttribute{Computed: true, CustomType: timetypes.RFC3339Type{}},
		},
	}
}

func (r *EventTypeResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*svix.Svix)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *svix.Svix, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}
	r.svx = client
}

func (r *EventTypeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data EventTypeResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Converting schemas into a `map[string]any`")
	var schema *map[string]any
	if data.Schemas.IsNull() {
		schema = nil
	} else {
		resp.Diagnostics.Append(data.Schemas.Unmarshal(&schema)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}
	tflog.Debug(ctx, "Creating EventTypeIn struct")
	eventTypeIn := models.EventTypeIn{
		Name:        data.Name.ValueString(),
		Description: data.Description.ValueString(),
		Archived:    boolOrNil(data.Archived),
		Deprecated:  boolOrNil(data.Deprecated),
		Schemas:     schema,
		FeatureFlag: strOrNil(data.FeatureFlag),
		GroupName:   strOrNil(data.GroupName),
	}
	tflog.Debug(ctx, "Sending `EventType.Create` request")
	reqOpts := svix.EventTypeCreateOptions{
		IdempotencyKey: randStr32(),
	}
	res, err := r.svx.EventType.Create(ctx, eventTypeIn, &reqOpts)
	if err != nil {
		resp.Diagnostics.AddError("Unable to create event type", err.Error())
		return
	}

	out := eventTypeOutToTFModel(ctx, &resp.Diagnostics, *res)
	if out != nil {
		resp.Diagnostics.Append(resp.State.Set(ctx, &out)...)
	}
}
func (r *EventTypeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data EventTypeResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	res, err := r.svx.EventType.Get(ctx, data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error while fetching state", err.Error())
		return
	}
	out := eventTypeOutToTFModel(ctx, &resp.Diagnostics, *res)
	if out != nil {
		resp.Diagnostics.Append(resp.State.Set(ctx, &out)...)
	}
}

func (r *EventTypeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data EventTypeResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

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

	res, err := r.svx.EventType.Update(ctx, data.Name.ValueString(), eventType)
	if err != nil {
		resp.Diagnostics.AddError("Error while updating event type", err.Error())
		return
	}

	schemasStr := mapStringTToString(&resp.Diagnostics, res.Schemas)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("schemas"), jsontypes.NewNormalizedPointerValue(schemasStr))...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("archived"), res.Archived)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("deprecated"), res.Deprecated)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("description"), res.Description)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("feature_flag"), res.FeatureFlag)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("group_name"), res.GroupName)...)

}
func (r *EventTypeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data EventTypeResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	opts := svix.EventTypeDeleteOptions{
		Expunge: ptr(true),
	}
	err := r.svx.EventType.Delete(ctx, data.Name.ValueString(), &opts)
	if err != nil {
		resp.Diagnostics.AddError("Failed to delete event type", err.Error())
	}
}

func eventTypeOutToTFModel(ctx context.Context, d *diag.Diagnostics, v models.EventTypeOut) *EventTypeResourceModel {
	tflog.Debug(ctx, "Converting schemas back into a `NormalizedValue`")
	var schemas jsontypes.Normalized
	if v.Schemas == nil {
		schemas = jsontypes.NewNormalizedNull()
	} else {
		schemasString, err := json.Marshal(v.Schemas)
		if err != nil {
			d.AddError("Unable to marshal `schemas` back into model (after response)", err.Error())
			return nil
		}
		schemas = jsontypes.NewNormalizedValue(string(schemasString))
	}
	out := EventTypeResourceModel{
		Archived:    types.BoolPointerValue(v.Archived),
		CreatedAt:   timetypes.NewRFC3339TimeValue(v.CreatedAt),
		Deprecated:  types.BoolValue(v.Deprecated),
		Description: types.StringValue(v.Description),
		FeatureFlag: types.StringPointerValue(v.FeatureFlag),
		GroupName:   types.StringPointerValue(v.GroupName),
		Name:        types.StringValue(v.Name),
		Schemas:     schemas,
		UpdatedAt:   timetypes.NewRFC3339TimeValue(v.UpdatedAt),
	}
	return &out
}
