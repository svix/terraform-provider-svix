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
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

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
	state appState
}

type OperationalWebhooksEndpointResourceModel struct {
	EnvironmentId types.String         `tfsdk:"environment_id"`
	CreatedAt     timetypes.RFC3339    `tfsdk:"created_at"`
	Description   types.String         `tfsdk:"description"`
	Disabled      types.Bool           `tfsdk:"disabled"`
	FilterTypes   types.List           `tfsdk:"filter_types"`
	Id            types.String         `tfsdk:"id"`
	Metadata      jsontypes.Normalized `tfsdk:"metadata"`
	RateLimit     types.Int32          `tfsdk:"rate_limit"`
	Secret        types.String         `tfsdk:"secret"`
	Uid           types.String         `tfsdk:"uid"`
	UpdatedAt     timetypes.RFC3339    `tfsdk:"updated_at"`
	Url           types.String         `tfsdk:"url"`
}

func (r *OperationalWebhooksEndpointResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "svix_operational_webhooks_endpoint"
}

func (r *OperationalWebhooksEndpointResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"environment_id": schema.StringAttribute{
				Required:    true,
				Description: ENV_ID_DESC,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"created_at":  schema.StringAttribute{Computed: true, CustomType: timetypes.RFC3339Type{}},
			"description": schema.StringAttribute{Computed: true, Optional: true, Default: stringdefault.StaticString("")},
			"disabled":    schema.BoolAttribute{Computed: true, Optional: true, Default: booldefault.StaticBool(false)},
			"filter_types": schema.ListAttribute{ElementType: types.StringType, Required: true, Validators: []validator.List{
				listvalidator.SizeAtLeast(1),
				listvalidator.UniqueValues(),
				listvalidator.ValueStringsAre(stringvalidator.OneOf(opWebhookTypes...)),
			}},
			"id": schema.StringAttribute{Computed: true},
			"metadata": schema.StringAttribute{
				Computed:            true,
				CustomType:          jsontypes.NormalizedType{},
				Default:             stringdefault.StaticString("{}"),
				MarkdownDescription: "JSON object encoded as a string, use `jsonencode` to create this field",
				Optional:            true,
			},
			"rate_limit": schema.Int32Attribute{Optional: true, Validators: []validator.Int32{
				// uint16
				int32validator.AtLeast(1),
				int32validator.AtMost(65535),
			}},
			"secret": schema.StringAttribute{
				Sensitive:           true,
				Computed:            true,
				MarkdownDescription: "The endpoint's verification secret.\n" + "Format: base64 encoded random bytes prefixed with whsec_. the server generates the secret.",
			},
			"uid": schema.StringAttribute{Optional: true, Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
				stringvalidator.LengthAtMost(256),
				stringvalidator.RegexMatches(saneStringRegex(), "String must match against `^[a-zA-Z0-9\\-_.]+$`"),
			}},
			"updated_at": schema.StringAttribute{Computed: true, CustomType: timetypes.RFC3339Type{}},
			"url":        schema.StringAttribute{Required: true},
		},
	}
}

func (r *OperationalWebhooksEndpointResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *OperationalWebhooksEndpointResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// load state/plan
	var data OperationalWebhooksEndpointResourceModel
	var envId string
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("environment_id"), &envId)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// create svix client
	svx, err := r.state.clientWithEnvId(envId)
	if err != nil {
		resp.Diagnostics.AddError(UNABLE_TO_CREATE_SVIX_CLIENT, err.Error())
		return
	}

	// create the opWebhookIn struct
	var opWebhookIn models.OperationalWebhookEndpointIn
	{
		metadata := stringToMapStringT[string](&resp.Diagnostics, data.Metadata.ValueStringPointer())
		if resp.Diagnostics.HasError() {
			return
		}

		var filterTypes []string
		resp.Diagnostics.Append(data.FilterTypes.ElementsAs(ctx, &filterTypes, false)...)
		if resp.Diagnostics.HasError() {
			return
		}

		var rateLimit *uint16
		if !data.RateLimit.IsUnknown() && !data.RateLimit.IsNull() {
			rateLimit = ptr(uint16(data.RateLimit.ValueInt32()))
		}

		opWebhookIn = models.OperationalWebhookEndpointIn{
			Description: strOrNil(data.Description),
			Disabled:    boolOrNil(data.Disabled),
			FilterTypes: filterTypes,
			Metadata:    metadata,
			RateLimit:   rateLimit,
			Uid:         strOrNil(data.Uid),
			Url:         data.Url.ValueString(),
		}
	}

	// call api
	res, err := svx.OperationalWebhookEndpoint.Create(ctx, opWebhookIn, &svix.OperationalWebhookEndpointCreateOptions{
		IdempotencyKey: randStr32(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Unable to create operational webhooks endpoint", err.Error())
		return
	}
	secretRes, err := svx.OperationalWebhook.Endpoint.GetSecret(ctx, res.Id)
	if err != nil {
		resp.Diagnostics.AddError("Failed to get op webhook endpoint secret", err.Error())
		return
	}

	// save state
	filterTypesOut, diags := types.ListValueFrom(ctx, types.StringType, res.FilterTypes)
	resp.Diagnostics.Append(diags...)
	metadataOut, err := json.Marshal(res.Metadata)
	if err != nil {
		resp.Diagnostics.AddAttributeError(path.Root("metadata"), "Unable to marshal metadata to a string", err.Error())
	}
	var rateLimitOut *int32
	if res.RateLimit != nil {
		rateLimitOut = ptr(int32(*res.RateLimit))
	}

	setCreateState(ctx, resp, "environment_id", envId)
	setCreateState(ctx, resp, "created_at", timetypes.NewRFC3339TimeValue(res.CreatedAt))
	setCreateState(ctx, resp, "description", types.StringValue(res.Description))
	setCreateState(ctx, resp, "disabled", types.BoolPointerValue(res.Disabled))
	setCreateState(ctx, resp, "filter_types", filterTypesOut)
	setCreateState(ctx, resp, "id", types.StringValue(res.Id))
	setCreateState(ctx, resp, "metadata", jsontypes.NewNormalizedValue(string(metadataOut)))
	setCreateState(ctx, resp, "rate_limit", types.Int32PointerValue(rateLimitOut))
	setCreateState(ctx, resp, "secret", types.StringValue(secretRes.Key))
	setCreateState(ctx, resp, "uid", types.StringPointerValue(res.Uid))
	setCreateState(ctx, resp, "updated_at", timetypes.NewRFC3339TimeValue(res.UpdatedAt))
	setCreateState(ctx, resp, "url", types.StringValue(res.Url))

}

func (r *OperationalWebhooksEndpointResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// load state/plan
	var data OperationalWebhooksEndpointResourceModel
	var envId string
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("environment_id"), &envId)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// create svix client
	svx, err := r.state.clientWithEnvId(envId)
	if err != nil {
		resp.Diagnostics.AddError(UNABLE_TO_CREATE_SVIX_CLIENT, err.Error())
		return
	}

	// call api
	res, err := svx.OperationalWebhookEndpoint.Get(ctx, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to get op webhook endpoint", err.Error())
		return
	}
	secretRes, err := svx.OperationalWebhook.Endpoint.GetSecret(ctx, res.Id)
	if err != nil {
		resp.Diagnostics.AddError("Failed to get op webhook endpoint secret", err.Error())
		return
	}
	// save state
	filterTypesOut, diags := types.ListValueFrom(ctx, types.StringType, res.FilterTypes)
	resp.Diagnostics.Append(diags...)
	metadataOut, err := json.Marshal(res.Metadata)
	if err != nil {
		resp.Diagnostics.AddAttributeError(path.Root("metadata"), "Unable to marshal metadata to a string", err.Error())
	}
	var rateLimitOut *int32
	if res.RateLimit != nil {
		rateLimitOut = ptr(int32(*res.RateLimit))
	}

	setReadState(ctx, resp, "environment_id", envId)
	setReadState(ctx, resp, "created_at", timetypes.NewRFC3339TimeValue(res.CreatedAt))
	setReadState(ctx, resp, "description", types.StringValue(res.Description))
	setReadState(ctx, resp, "disabled", types.BoolPointerValue(res.Disabled))
	setReadState(ctx, resp, "filter_types", filterTypesOut)
	setReadState(ctx, resp, "id", types.StringValue(res.Id))
	setReadState(ctx, resp, "metadata", jsontypes.NewNormalizedValue(string(metadataOut)))
	setReadState(ctx, resp, "rate_limit", types.Int32PointerValue(rateLimitOut))
	setReadState(ctx, resp, "secret", types.StringValue(secretRes.Key))
	setReadState(ctx, resp, "uid", types.StringPointerValue(res.Uid))
	setReadState(ctx, resp, "updated_at", timetypes.NewRFC3339TimeValue(res.UpdatedAt))
	setReadState(ctx, resp, "url", types.StringValue(res.Url))
}

func (r *OperationalWebhooksEndpointResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// load state/plan
	var data OperationalWebhooksEndpointResourceModel
	var envId, epId string
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("environment_id"), &envId)...)
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &epId)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// create svix client
	svx, err := r.state.clientWithEnvId(envId)
	if err != nil {
		resp.Diagnostics.AddError(UNABLE_TO_CREATE_SVIX_CLIENT, err.Error())
		return
	}

	// call api
	metadata := stringToMapStringT[string](&resp.Diagnostics, data.Metadata.ValueStringPointer())
	if resp.Diagnostics.HasError() {
		return
	}

	var filterTypes []string
	resp.Diagnostics.Append(data.FilterTypes.ElementsAs(ctx, &filterTypes, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var rateLimit *uint16
	if !data.RateLimit.IsUnknown() && !data.RateLimit.IsNull() {
		rateLimit = ptr(uint16(data.RateLimit.ValueInt32()))
	}

	opWebhook := models.OperationalWebhookEndpointUpdate{
		Description: strOrNil(data.Description),
		Disabled:    boolOrNil(data.Disabled),
		FilterTypes: filterTypes,
		Metadata:    metadata,
		RateLimit:   rateLimit,
		Uid:         strOrNil(data.Uid),
		Url:         data.Url.ValueString(),
	}
	res, err := svx.OperationalWebhookEndpoint.Update(ctx, epId, opWebhook)
	if err != nil {
		resp.Diagnostics.AddError("Error while updating operational webhook endpoint", err.Error())
		return
	}

	// save state
	outMetadata := mapStringTToString(&resp.Diagnostics, &res.Metadata)
	if resp.Diagnostics.HasError() {
		return
	}
	setUpdateState(ctx, resp, "description", res.Description)
	setUpdateState(ctx, resp, "disabled", res.Disabled)
	setUpdateState(ctx, resp, "filter_types", res.FilterTypes)
	setUpdateState(ctx, resp, "metadata", outMetadata)
	setUpdateState(ctx, resp, "rate_limit", res.RateLimit)
	setUpdateState(ctx, resp, "uid", res.Uid)
	setUpdateState(ctx, resp, "url", res.Url)

}

func (r *OperationalWebhooksEndpointResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var envId, epId string
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("environment_id"), &envId)...)
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &epId)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// create svix client
	svx, err := r.state.clientWithEnvId(envId)
	if err != nil {
		resp.Diagnostics.AddError(UNABLE_TO_CREATE_SVIX_CLIENT, err.Error())
		return
	}

	err = svx.OperationalWebhookEndpoint.Delete(ctx, epId)
	if err != nil {
		resp.Diagnostics.AddError("Failed to delete operational webhooks endpoint", err.Error())
	}

}
