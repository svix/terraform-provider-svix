package internal

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	svix "github.com/svix/svix-webhooks/go"
	"github.com/svix/svix-webhooks/go/models"
)

const (
	autoConfigTokenPrefixV1 = "auto_v1_"
	autoConfigMetadataKey   = "_svix_auto_config"
)

var _ resource.Resource = &AutoConfigResource{}

// autoConfigMetadataModifier ensures planned metadata includes the server-injected key
// so apply does not fail with an inconsistent result.
type autoConfigMetadataModifier struct{}

func (m autoConfigMetadataModifier) Description(_ context.Context) string {
	return "Merges _svix_auto_config=true into metadata (set by the AutoConfig API)."
}

func (m autoConfigMetadataModifier) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (m autoConfigMetadataModifier) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	if req.PlanValue.IsNull() || req.PlanValue.IsUnknown() {
		return
	}

	var meta map[string]string
	if err := json.Unmarshal([]byte(req.PlanValue.ValueString()), &meta); err != nil {
		return
	}
	if meta == nil {
		meta = map[string]string{}
	}
	if meta[autoConfigMetadataKey] == "true" {
		return
	}
	meta[autoConfigMetadataKey] = "true"
	out, err := json.Marshal(meta)
	if err != nil {
		return
	}
	resp.PlanValue = types.StringValue(string(out))
}

func withAutoConfigMetadata(meta *map[string]string) *map[string]string {
	if meta == nil {
		m := map[string]string{autoConfigMetadataKey: "true"}
		return &m
	}
	(*meta)[autoConfigMetadataKey] = "true"
	return meta
}

func NewAutoConfigResource() resource.Resource {
	return &AutoConfigResource{}
}

type AutoConfigResource struct{}

type AutoConfigResourceModel struct {
	Token        types.String         `tfsdk:"token"`
	Url          types.String         `tfsdk:"url"`
	FilterTypes  types.List           `tfsdk:"filter_types"`
	Channels     types.List           `tfsdk:"channels"`
	Description  types.String         `tfsdk:"description"`
	Disabled     types.Bool           `tfsdk:"disabled"`
	Headers      jsontypes.Normalized `tfsdk:"headers"`
	Metadata     jsontypes.Normalized `tfsdk:"metadata"`
	RateLimit    types.Int32          `tfsdk:"rate_limit"`
	ThrottleRate types.Int32          `tfsdk:"throttle_rate"`
	Uid          types.String         `tfsdk:"uid"`
	Version      types.Int32          `tfsdk:"version"`
	Id           types.String         `tfsdk:"id"`
	Secret       types.String         `tfsdk:"secret"`
	CreatedAt    timetypes.RFC3339    `tfsdk:"created_at"`
	UpdatedAt    timetypes.RFC3339    `tfsdk:"updated_at"`
}

type autoConfigTokenContentV1 struct {
	AppID          string `json:"aid"`
	EndpointID     string `json:"eid"`
	ServerURL      string `json:"surl"`
	EndpointSecret string `json:"esec"`
	TokenPlaintext string `json:"tok"`
}

func (r *AutoConfigResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "svix_autoconfig"
}

func (r *AutoConfigResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Configures a webhook endpoint via [Webhooks AutoConfig](https://docs.svix.com/receiving/webhooks-autoconfig). " +
			"Does not require provider `token` or `server_url`; credentials come from the AutoConfig token. " +
			"Create and update call `subscribe()`. Destroy only removes the resource from Terraform state " +
			"(the portal-created endpoint is left unchanged).",
		Attributes: map[string]schema.Attribute{
			"token": schema.StringAttribute{
				Required:            true,
				Sensitive:           true,
				MarkdownDescription: "AutoConfig token from the Application Portal (prefix `auto_v1_`).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"url": schema.StringAttribute{
				Required: true,
			},
			"filter_types": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Validators: []validator.List{
					listvalidator.UniqueValues(),
				},
			},
			"channels": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Validators: []validator.List{
					listvalidator.UniqueValues(),
				},
			},
			"description": schema.StringAttribute{
				Computed: true,
				Optional: true,
				Default:  stringdefault.StaticString(""),
			},
			"disabled": schema.BoolAttribute{
				Computed: true,
				Optional: true,
				Default:  booldefault.StaticBool(false),
			},
			"headers": schema.StringAttribute{
				Optional:            true,
				CustomType:          jsontypes.NormalizedType{},
				MarkdownDescription: "JSON object of HTTP headers encoded as a string, use `jsonencode` to create this field",
			},
			"metadata": schema.StringAttribute{
				Computed:   true,
				Optional:   true,
				CustomType: jsontypes.NormalizedType{},
				// Server always injects _svix_auto_config=true on subscribe.
				Default: stringdefault.StaticString(`{"_svix_auto_config":"true"}`),
				MarkdownDescription: "JSON object encoded as a string, use `jsonencode` to create this field. " +
					"The server always sets `_svix_auto_config` to `\"true\"`; the provider merges that key automatically.",
				PlanModifiers: []planmodifier.String{
					autoConfigMetadataModifier{},
				},
			},
			"rate_limit": schema.Int32Attribute{
				Optional:            true,
				MarkdownDescription: "Deprecated; prefer `throttle_rate`.",
				Validators: []validator.Int32{
					int32validator.AtLeast(1),
					int32validator.AtMost(65535),
				},
			},
			"throttle_rate": schema.Int32Attribute{
				Optional: true,
				Validators: []validator.Int32{
					int32validator.AtLeast(1),
					int32validator.AtMost(65535),
				},
			},
			"uid": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					stringvalidator.LengthAtMost(256),
					stringvalidator.RegexMatches(saneStringRegex(), "String must match against `^[a-zA-Z0-9\\-_.]+$`"),
				},
			},
			"version": schema.Int32Attribute{
				Computed: true,
				Optional: true,
				Validators: []validator.Int32{
					int32validator.AtLeast(1),
					int32validator.AtMost(65535),
				},
			},
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"secret": schema.StringAttribute{
				Computed:  true,
				Sensitive: true,
				MarkdownDescription: "The endpoint's verification secret from the AutoConfig token. " +
					"Format: base64 encoded random bytes prefixed with whsec_.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"created_at": schema.StringAttribute{
				Computed:   true,
				CustomType: timetypes.RFC3339Type{},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"updated_at": schema.StringAttribute{
				Computed:   true,
				CustomType: timetypes.RFC3339Type{},
			},
		},
	}
}

func (r *AutoConfigResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// AutoConfig does not use provider credentials.
}

func (r *AutoConfigResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data AutoConfigResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpointIn, ok := endpointInFromAutoConfigModel(ctx, &data, &resp.Diagnostics)
	if !ok {
		return
	}

	out, secret, err := subscribeAutoConfig(ctx, data.Token.ValueString(), endpointIn)
	if err != nil {
		logSvixError(&resp.Diagnostics, err, "Failed to subscribe AutoConfig endpoint")
		return
	}

	setCreateState(ctx, resp, rp("token"), data.Token)
	setCreateState(ctx, resp, rp("headers"), data.Headers)
	setAutoConfigCreateState(ctx, resp, out, secret)
}

func (r *AutoConfigResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// AutoConfig tokens are scoped to subscribe (update) only and cannot GET the endpoint.
	// Keep state as-is; create/update re-apply desired config via subscribe().
	var data AutoConfigResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AutoConfigResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data AutoConfigResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpointIn, ok := endpointInFromAutoConfigModel(ctx, &data, &resp.Diagnostics)
	if !ok {
		return
	}

	out, secret, err := subscribeAutoConfig(ctx, data.Token.ValueString(), endpointIn)
	if err != nil {
		logSvixError(&resp.Diagnostics, err, "Failed to subscribe AutoConfig endpoint")
		return
	}

	setUpdateState(ctx, resp, rp("token"), data.Token)
	setUpdateState(ctx, resp, rp("headers"), data.Headers)
	setAutoConfigUpdateState(ctx, resp, out, secret)
}

func (r *AutoConfigResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// AutoConfig configures a portal-created endpoint; destroy only drops Terraform state.
}

func endpointInFromAutoConfigModel(ctx context.Context, data *AutoConfigResourceModel, d *diag.Diagnostics) (models.EndpointIn, bool) {
	var filterTypes []string
	if !data.FilterTypes.IsNull() && !data.FilterTypes.IsUnknown() {
		d.Append(data.FilterTypes.ElementsAs(ctx, &filterTypes, false)...)
		if d.HasError() {
			return models.EndpointIn{}, false
		}
	}

	var channels []string
	if !data.Channels.IsNull() && !data.Channels.IsUnknown() {
		d.Append(data.Channels.ElementsAs(ctx, &channels, false)...)
		if d.HasError() {
			return models.EndpointIn{}, false
		}
	}

	metadata := stringToMapStringT[string](d, data.Metadata.ValueStringPointer())
	if d.HasError() {
		return models.EndpointIn{}, false
	}
	metadata = withAutoConfigMetadata(metadata)

	headers := stringToMapStringT[string](d, data.Headers.ValueStringPointer())
	if d.HasError() {
		return models.EndpointIn{}, false
	}

	var rateLimit *uint16
	if !data.RateLimit.IsUnknown() && !data.RateLimit.IsNull() {
		rateLimit = ptr(uint16(data.RateLimit.ValueInt32()))
	}

	var throttleRate *uint16
	if !data.ThrottleRate.IsUnknown() && !data.ThrottleRate.IsNull() {
		throttleRate = ptr(uint16(data.ThrottleRate.ValueInt32()))
	}

	var version *uint16
	if !data.Version.IsUnknown() && !data.Version.IsNull() {
		version = ptr(uint16(data.Version.ValueInt32()))
	}

	return models.EndpointIn{
		Url:          data.Url.ValueString(),
		FilterTypes:  filterTypes,
		Channels:     channels,
		Description:  strOrNil(data.Description),
		Disabled:     boolOrNil(data.Disabled),
		Headers:      headers,
		Metadata:     metadata,
		RateLimit:    rateLimit,
		ThrottleRate: throttleRate,
		Uid:          strOrNil(data.Uid),
		Version:      version,
	}, true
}

func subscribeAutoConfig(ctx context.Context, token string, endpointIn models.EndpointIn) (*models.EndpointOut, string, error) {
	content, err := decodeAutoConfigTokenV1(token)
	if err != nil {
		return nil, "", err
	}

	ac, err := svix.NewAutoConfig(token, endpointIn)
	if err != nil {
		return nil, "", err
	}

	out, err := ac.Subscribe(ctx)
	if err != nil {
		return nil, "", err
	}

	return out, content.EndpointSecret, nil
}

func decodeAutoConfigTokenV1(token string) (*autoConfigTokenContentV1, error) {
	token, ok := strings.CutPrefix(token, autoConfigTokenPrefixV1)
	if !ok {
		return nil, fmt.Errorf("invalid AutoConfig token: missing %s prefix", autoConfigTokenPrefixV1)
	}

	decoded, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return nil, fmt.Errorf("invalid AutoConfig token: %w", err)
	}

	var content autoConfigTokenContentV1
	if err := json.Unmarshal(decoded, &content); err != nil {
		return nil, fmt.Errorf("invalid AutoConfig token: %w", err)
	}
	return &content, nil
}

func setAutoConfigCreateState(ctx context.Context, resp *resource.CreateResponse, out *models.EndpointOut, secret string) {
	filterTypesOut := optionalStringList(ctx, &resp.Diagnostics, out.FilterTypes)
	channelsOut := optionalStringList(ctx, &resp.Diagnostics, out.Channels)
	if resp.Diagnostics.HasError() {
		return
	}

	metadataOut, err := json.Marshal(out.Metadata)
	if err != nil {
		resp.Diagnostics.AddAttributeError(path.Root("metadata"), "Unable to marshal metadata to a string", err.Error())
		return
	}

	var rateLimitOut *int32
	if out.RateLimit != nil {
		rateLimitOut = ptr(int32(*out.RateLimit))
	}
	var throttleRateOut *int32
	if out.ThrottleRate != nil {
		throttleRateOut = ptr(int32(*out.ThrottleRate))
	}

	setCreateState(ctx, resp, rp("url"), types.StringValue(out.Url))
	setCreateState(ctx, resp, rp("filter_types"), filterTypesOut)
	setCreateState(ctx, resp, rp("channels"), channelsOut)
	setCreateState(ctx, resp, rp("description"), types.StringValue(out.Description))
	setCreateState(ctx, resp, rp("disabled"), types.BoolPointerValue(out.Disabled))
	setCreateState(ctx, resp, rp("metadata"), jsontypes.NewNormalizedValue(string(metadataOut)))
	setCreateState(ctx, resp, rp("rate_limit"), types.Int32PointerValue(rateLimitOut))
	setCreateState(ctx, resp, rp("throttle_rate"), types.Int32PointerValue(throttleRateOut))
	setCreateState(ctx, resp, rp("uid"), types.StringPointerValue(out.Uid))
	setCreateState(ctx, resp, rp("version"), types.Int32Value(out.Version))
	setCreateState(ctx, resp, rp("id"), types.StringValue(out.Id))
	setCreateState(ctx, resp, rp("secret"), types.StringValue(secret))
	setCreateState(ctx, resp, rp("created_at"), timetypes.NewRFC3339TimeValue(out.CreatedAt))
	setCreateState(ctx, resp, rp("updated_at"), timetypes.NewRFC3339TimeValue(out.UpdatedAt))
}

func setAutoConfigUpdateState(ctx context.Context, resp *resource.UpdateResponse, out *models.EndpointOut, secret string) {
	filterTypesOut := optionalStringList(ctx, &resp.Diagnostics, out.FilterTypes)
	channelsOut := optionalStringList(ctx, &resp.Diagnostics, out.Channels)
	if resp.Diagnostics.HasError() {
		return
	}

	metadataOut, err := json.Marshal(out.Metadata)
	if err != nil {
		resp.Diagnostics.AddAttributeError(path.Root("metadata"), "Unable to marshal metadata to a string", err.Error())
		return
	}

	var rateLimitOut *int32
	if out.RateLimit != nil {
		rateLimitOut = ptr(int32(*out.RateLimit))
	}
	var throttleRateOut *int32
	if out.ThrottleRate != nil {
		throttleRateOut = ptr(int32(*out.ThrottleRate))
	}

	setUpdateState(ctx, resp, rp("url"), types.StringValue(out.Url))
	setUpdateState(ctx, resp, rp("filter_types"), filterTypesOut)
	setUpdateState(ctx, resp, rp("channels"), channelsOut)
	setUpdateState(ctx, resp, rp("description"), types.StringValue(out.Description))
	setUpdateState(ctx, resp, rp("disabled"), types.BoolPointerValue(out.Disabled))
	setUpdateState(ctx, resp, rp("metadata"), jsontypes.NewNormalizedValue(string(metadataOut)))
	setUpdateState(ctx, resp, rp("rate_limit"), types.Int32PointerValue(rateLimitOut))
	setUpdateState(ctx, resp, rp("throttle_rate"), types.Int32PointerValue(throttleRateOut))
	setUpdateState(ctx, resp, rp("uid"), types.StringPointerValue(out.Uid))
	setUpdateState(ctx, resp, rp("version"), types.Int32Value(out.Version))
	setUpdateState(ctx, resp, rp("id"), types.StringValue(out.Id))
	setUpdateState(ctx, resp, rp("secret"), types.StringValue(secret))
	setUpdateState(ctx, resp, rp("created_at"), timetypes.NewRFC3339TimeValue(out.CreatedAt))
	setUpdateState(ctx, resp, rp("updated_at"), timetypes.NewRFC3339TimeValue(out.UpdatedAt))
}

func optionalStringList(ctx context.Context, d *diag.Diagnostics, values []string) types.List {
	if len(values) == 0 {
		return types.ListNull(types.StringType)
	}
	list, diags := types.ListValueFrom(ctx, types.StringType, values)
	d.Append(diags...)
	return list
}
