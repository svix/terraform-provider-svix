package internal

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	svix "github.com/svix/svix-webhooks/go"
	"github.com/svix/svix-webhooks/go/models"
)

var _ resource.Resource = &EventTypeOpenapiImportResource{}

type EventTypeOpenapiImportResource struct {
	state appState
}

type EventTypeOpenapiImportResourceModel struct {
	EnvironmentId     types.String `tfsdk:"environment_id"`
	ReplaceAll        types.Bool   `tfsdk:"replace_all"`
	SpecRaw           types.String `tfsdk:"spec_raw"`
	CreatedEventTypes types.List   `tfsdk:"created_event_types"`
}

func NewEventTypeOpenapiImportResource() resource.Resource {
	return &EventTypeOpenapiImportResource{}
}

func (r *EventTypeOpenapiImportResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *EventTypeOpenapiImportResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "svix_event_type_openapi_import"
}

func (r *EventTypeOpenapiImportResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Given an OpenAPI spec, create new or update existing event types. If an existing `archived` event type is updated, it will be unarchived.\n\n" +
			"The importer will convert all webhooks found in the either the `webhooks` or `x-webhooks` top-level.\n\n" +
			"Import a list of event types from webhooks defined in an OpenAPI spec.\n\n" +
			"The OpenAPI spec is specified in the `raw_spec` field a YAML or JSON string",
		Attributes: map[string]schema.Attribute{
			"environment_id": schema.StringAttribute{
				Required:    true,
				Description: ENV_ID_DESC,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"replace_all": schema.BoolAttribute{
				Computed:            true,
				Optional:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Default `false`. If `true`, all existing event types that are not in the spec will be archived.",
			},
			"spec_raw": schema.StringAttribute{
				Required: true,
				MarkdownDescription: "A string, parsed by the server as YAML or JSON.\n\n" +
					"If the spec includes event types already defined (either by using terraform, the API, or the frontend), they will be overwritten",
			},
			"created_event_types": schema.ListAttribute{
				Computed:            true,
				MarkdownDescription: "List of the created event types",
				ElementType:         types.StringType,
			},
		},
	}

}

func (r *EventTypeOpenapiImportResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// load state/plan
	var data EventTypeOpenapiImportResourceModel
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

	// call API
	res, err := svx.EventType.ImportOpenapi(
		ctx,
		models.EventTypeImportOpenApiIn{
			ReplaceAll: data.ReplaceAll.ValueBoolPointer(),
			SpecRaw:    data.SpecRaw.ValueStringPointer(),
		},
		&svix.EventTypeImportOpenapiOptions{
			IdempotencyKey: randStr32(),
		},
	)
	if err != nil {
		logSvixError(&resp.Diagnostics, err, "Failed to import event types")
		return
	}

	// save state
	createdEventTypes, diags := types.ListValueFrom(ctx, types.StringType, res.Data.Modified)
	resp.Diagnostics.Append(diags...)

	setCreateState(ctx, resp, rp("environment_id"), envId)
	setCreateState(ctx, resp, rp("replace_all"), data.ReplaceAll)
	setCreateState(ctx, resp, rp("spec_raw"), data.SpecRaw)
	setCreateState(ctx, resp, rp("created_event_types"), createdEventTypes)

}

func (r *EventTypeOpenapiImportResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// null resource, this read does nothing
	var data EventTypeOpenapiImportResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	setReadState(ctx, resp, rp("environment_id"), data.EnvironmentId)
	setReadState(ctx, resp, rp("replace_all"), data.ReplaceAll)
	setReadState(ctx, resp, rp("spec_raw"), data.SpecRaw)
	setReadState(ctx, resp, rp("created_event_types"), data.CreatedEventTypes)

}

func (r *EventTypeOpenapiImportResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// load state/plan
	var data EventTypeOpenapiImportResourceModel
	var envId string
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("environment_id"), &envId)...)
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

	// call API
	res, err := svx.EventType.ImportOpenapi(
		ctx,
		models.EventTypeImportOpenApiIn{
			ReplaceAll: data.ReplaceAll.ValueBoolPointer(),
			SpecRaw:    data.SpecRaw.ValueStringPointer(),
		},
		&svix.EventTypeImportOpenapiOptions{
			IdempotencyKey: randStr32(),
		},
	)
	if err != nil {
		logSvixError(&resp.Diagnostics, err, "Failed to update event types")
		return
	}

	// save state
	createdEventTypes, diags := types.ListValueFrom(ctx, types.StringType, res.Data.Modified)
	resp.Diagnostics.Append(diags...)

	setUpdateState(ctx, resp, rp("environment_id"), envId)
	setUpdateState(ctx, resp, rp("replace_all"), data.ReplaceAll)
	setUpdateState(ctx, resp, rp("spec_raw"), data.SpecRaw)
	setUpdateState(ctx, resp, rp("created_event_types"), createdEventTypes)
}

func (r *EventTypeOpenapiImportResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// load state/plan
	var envId string
	var eventTypesToDelete []string
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("environment_id"), &envId)...)
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("created_event_types"), &eventTypesToDelete)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// create svix client
	svx, err := r.state.ClientWithEnvId(envId)
	if err != nil {
		resp.Diagnostics.AddError(UNABLE_TO_CREATE_SVIX_CLIENT, err.Error())
		return
	}

	for _, eventTypeName := range eventTypesToDelete {
		err = svx.EventType.Delete(ctx, eventTypeName, nil)
		if err != nil {
			logSvixError(&resp.Diagnostics, err, fmt.Sprintf("Failed to delete event type %s", eventTypeName))
			return
		}

	}
}
