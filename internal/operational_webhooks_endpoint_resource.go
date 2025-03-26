package internal

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	svix "github.com/svix/svix-webhooks/go"
	"github.com/svix/svix-webhooks/go/models"
)

var _ resource.Resource = &OperationalWebhooksEndpointResource{}

func NewOperationalWebhooksEndpoint() resource.Resource {
	return &OperationalWebhooksEndpointResource{}
}

type OperationalWebhooksEndpointResource struct {
	svx *svix.Svix
}

type OperationalWebhooksEndpointModel struct {
	Description types.String         `tfsdk:"description"`
	Disabled    types.Bool           `tfsdk:"disabled"`
	FilterTypes types.List           `tfsdk:"filter_types"`
	Metadata    jsontypes.Normalized `tfsdk:"metadata"`
	RateLimit   types.Int32          `tfsdk:"rate_limit"`
	Secret      types.String         `tfsdk:"secret"`
	Uid         types.String         `tfsdk:"uid"`
	Url         types.String         `tfsdk:"url"`
}

func (r *OperationalWebhooksEndpointResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "svix_operational_webhooks_endpoint"

}
func (r *OperationalWebhooksEndpointResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"description":  schema.StringAttribute{Optional: true},
			"disabled":     schema.BoolAttribute{Optional: true},
			"filter_types": schema.ListAttribute{ElementType: types.StringType, Optional: true},
			"metadata":     schema.StringAttribute{CustomType: jsontypes.NormalizedType{}, Optional: true},
			"rate_limit":   schema.Int32Attribute{Optional: true},
			"secret":       schema.StringAttribute{Optional: true, Sensitive: true},
			"uid":          schema.StringAttribute{Optional: true},
			"url":          schema.StringAttribute{Required: true},
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
	log.Println("in create")
	var data OperationalWebhooksEndpointModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	metadata := stringToMapStringT[string](&resp.Diagnostics, data.Metadata.ValueString())

	if metadata == nil {
		return
	}
	opWebhookIn := models.OperationalWebhookEndpointIn{
		Description: data.Description.ValueStringPointer(),
		Disabled:    data.Disabled.ValueBoolPointer(),
		// FilterTypes:data.FilterTypes.
		Metadata:  metadata,
		RateLimit: ptr(uint16(data.RateLimit.ValueInt32())),
		Secret:    data.Secret.ValueStringPointer(),
		Uid:       data.Uid.ValueStringPointer(),
		Url:       data.Url.ValueString(),
	}

	spw(opWebhookIn)
}
func (r *OperationalWebhooksEndpointResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
}
func (r *OperationalWebhooksEndpointResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}
func (r *OperationalWebhooksEndpointResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}
