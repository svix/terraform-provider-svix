package internal

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	svix "github.com/svix/svix-webhooks/go"
	"github.com/svix/svix-webhooks/go/models"
)

var _ resource.Resource = &OperationalWebhooksEndpointResource{}

var opWebhookTypes = []string{
	"background_task.finished",
	"endpoint.created",
	"endpoint.deleted",
	"endpoint.disabled",
	"endpoint.enabled",
	"endpoint.updated",
	"message.attempt.exhausted",
	"message.attempt.failing",
	"message.attempt.recovered",
}

func NewOperationalWebhooksEndpoint() resource.Resource {
	return &OperationalWebhooksEndpointResource{}
}

type OperationalWebhooksEndpointResource struct {
	svx *svix.Svix
}

type OperationalWebhooksEndpointResourceModel struct {
	CreatedAt   timetypes.RFC3339    `tfsdk:"created_at"`
	Description types.String         `tfsdk:"description"`
	Disabled    types.Bool           `tfsdk:"disabled"`
	FilterTypes types.List           `tfsdk:"filter_types"`
	Id          types.String         `tfsdk:"id"`
	Metadata    jsontypes.Normalized `tfsdk:"metadata"`
	RateLimit   types.Int32          `tfsdk:"rate_limit"`
	Secret      types.String         `tfsdk:"secret"`
	Uid         types.String         `tfsdk:"uid"`
	UpdatedAt   timetypes.RFC3339    `tfsdk:"updated_at"`
	Url         types.String         `tfsdk:"url"`
}

func (r *OperationalWebhooksEndpointResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "svix_operational_webhooks_endpoint"
}

func (r *OperationalWebhooksEndpointResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"created_at":  schema.StringAttribute{Computed: true, CustomType: timetypes.RFC3339Type{}},
			"description": schema.StringAttribute{Computed: true, Optional: true, Default: stringdefault.StaticString("")},
			"disabled":    schema.BoolAttribute{Computed: true, Optional: true, Default: booldefault.StaticBool(false)},
			"filter_types": schema.ListAttribute{ElementType: types.StringType, Required: true, Validators: []validator.List{
				listvalidator.SizeAtLeast(1),
				listvalidator.UniqueValues(),
				listvalidator.ValueStringsAre(stringvalidator.OneOf(opWebhookTypes...)),
			}},
			"id":       schema.StringAttribute{Computed: true},
			"metadata": schema.StringAttribute{Computed: true, CustomType: jsontypes.NormalizedType{}, Optional: true, Default: stringdefault.StaticString("{}")},
			"rate_limit": schema.Int32Attribute{Optional: true, Validators: []validator.Int32{
				int32validator.AtLeast(1),
			}},
			"secret":     schema.StringAttribute{Optional: true, Sensitive: true},
			"uid":        schema.StringAttribute{Optional: true, Computed: true},
			"updated_at": schema.StringAttribute{Computed: true, CustomType: timetypes.RFC3339Type{}},
			"url":        schema.StringAttribute{Required: true},
		},
	}
}

func (r *OperationalWebhooksEndpointResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *OperationalWebhooksEndpointResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data OperationalWebhooksEndpointResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	metadata := stringToMapStringT[string](&resp.Diagnostics, data.Metadata.ValueString())
	if metadata == nil {
		return
	}
	var filterTypes []string
	resp.Diagnostics.Append(data.FilterTypes.ElementsAs(ctx, &filterTypes, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	opWebhookIn := models.OperationalWebhookEndpointIn{
		Description: data.Description.ValueStringPointer(),
		Disabled:    data.Disabled.ValueBoolPointer(),
		FilterTypes: filterTypes,
		Metadata:    metadata,
		RateLimit:   ptr(uint16(data.RateLimit.ValueInt32())),
		Secret:      data.Secret.ValueStringPointer(),
		Uid:         strOrNil(data.Uid),
		Url:         data.Url.ValueString(),
	}
	opts := svix.OperationalWebhookEndpointCreateOptions{
		IdempotencyKey: randStr32(),
	}
	res, err := r.svx.OperationalWebhookEndpoint.Create(ctx, opWebhookIn, &opts)
	if err != nil {
		resp.Diagnostics.AddError("Unable to create operational webhooks endpoint", err.Error())
		return
	}
	out := operationalWebhookEndpointOutToModel(ctx, &resp.Diagnostics, *res)
	if out != nil {
		resp.Diagnostics.Append(resp.State.Set(ctx, &out)...)
	}
}

func (r *OperationalWebhooksEndpointResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data OperationalWebhooksEndpointResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	res, err := r.svx.OperationalWebhookEndpoint.Get(ctx, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to get op webhook endpoint", err.Error())
		return
	}
	out := operationalWebhookEndpointOutToModel(ctx, &resp.Diagnostics, *res)
	if out != nil {
		resp.Diagnostics.Append(resp.State.Set(ctx, &out)...)
	}

}

func (r *OperationalWebhooksEndpointResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data OperationalWebhooksEndpointResourceModel
	var ep_id string
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &ep_id)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	metadata := stringToMapStringT[string](&resp.Diagnostics, data.Metadata.ValueString())
	if metadata == nil {
		return
	}
	var filterTypes []string
	resp.Diagnostics.Append(data.FilterTypes.ElementsAs(ctx, &filterTypes, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	opWebhook := models.OperationalWebhookEndpointUpdate{
		Description: data.Description.ValueStringPointer(),
		Disabled:    data.Disabled.ValueBoolPointer(),
		FilterTypes: filterTypes,
		Metadata:    metadata,
		RateLimit:   ptr(uint16(data.RateLimit.ValueInt32())),
		Uid:         strOrNil(data.Uid),
		Url:         data.Url.ValueString(),
	}
	res, err := r.svx.OperationalWebhookEndpoint.Update(ctx, ep_id, opWebhook)
	if err != nil {
		resp.Diagnostics.AddError("Error while updating operational webhook endpoint", err.Error())
		return
	}
	out_metadata := mapStringTToString(&resp.Diagnostics, res.Metadata)
	if out_metadata == nil {
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("description"), res.Description)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("disabled"), res.Disabled)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("filter_types"), res.FilterTypes)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("metadata"), out_metadata)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("rate_limit"), res.RateLimit)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("uid"), res.Uid)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("url"), res.Url)...)

}

func (r *OperationalWebhooksEndpointResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var ep_id string
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &ep_id)...)
	if resp.Diagnostics.HasError() {
		return
	}
	err := r.svx.OperationalWebhookEndpoint.Delete(ctx, ep_id)
	if err != nil {
		resp.Diagnostics.AddError("Failed to delete operational webhooks endpoint", err.Error())
	}

}

func operationalWebhookEndpointOutToModel(ctx context.Context, d *diag.Diagnostics, v models.OperationalWebhookEndpointOut) *OperationalWebhooksEndpointResourceModel {
	filterTypes, diags := types.ListValueFrom(ctx, types.StringType, v.FilterTypes)
	d.Append(diags...)
	if d.HasError() {
		return nil
	}
	metadata, err := json.Marshal(v.Metadata)
	if err != nil {
		d.AddAttributeError(path.Root("metadata"), "Unable to marshal metadata to a string", err.Error())
		return nil
	}
	ret := OperationalWebhooksEndpointResourceModel{
		CreatedAt:   timetypes.NewRFC3339TimeValue(v.CreatedAt),
		Description: types.StringValue(v.Description),
		Disabled:    types.BoolPointerValue(v.Disabled),
		FilterTypes: filterTypes,
		Id:          types.StringValue(v.Id),
		Metadata:    jsontypes.NewNormalizedValue(string(metadata)),
		RateLimit:   types.Int32Value(int32(*v.RateLimit)),
		Uid:         types.StringPointerValue(v.Uid),
		UpdatedAt:   timetypes.NewRFC3339TimeValue(v.UpdatedAt),
		Url:         types.StringValue(v.Url),
	}
	return &ret
}
