package internal

// svix_ingest_source_resource
import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	svix "github.com/svix/svix-webhooks/go"
	"github.com/svix/svix-webhooks/go/models"
)

var _ resource.Resource = &SvixIngestSourceResource{}

type SvixIngestSourceResource struct {
	state appState
}

type SvixIngestSourceResourceModel struct {
	EnvironmentId types.String         `tfsdk:"environment_id"`
	Id            types.String         `tfsdk:"id"`
	Type          types.String         `tfsdk:"type"`
	Name          types.String         `tfsdk:"name"`
	Uid           types.String         `tfsdk:"uid"`
	Config        jsontypes.Normalized `tfsdk:"config"`
	IngestUrl     types.String         `tfsdk:"ingest_url"`
	CreatedAt     timetypes.RFC3339    `tfsdk:"created_at"`
	UpdatedAt     timetypes.RFC3339    `tfsdk:"updated_at"`
}

func NewSvixIngestSourceResource() resource.Resource {
	return &SvixIngestSourceResource{}
}

func (r *SvixIngestSourceResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// before adding any new types to this list, ensure that `createIngestConfigFromCurrentValAndPlan` works correctly
// We are currently assuming all ingest configs are a flat json object
var ingestSourceInTypeKeys = []string{
	"generic-webhook",
	"cron",
	"adobe-sign",
	"beehiiv",
	"brex",
	"clerk",
	"docusign",
	"github",
	"guesty",
	"hubspot",
	"incident-io",
	"lithic",
	"nash",
	"pleo",
	"replicate",
	"resend",
	"safebase",
	"sardine",
	"segment",
	"shopify",
	"slack",
	"stripe",
	"stych",
	"svix",
	"zoom",
}

func ingestSourceInTypesForDocs() []string {
	var types []string
	types = append(types, ingestSourceInTypeKeys...)

	for i, val := range types {
		types[i] = fmt.Sprintf("`%s`", val)
	}
	return types
}

func (r *SvixIngestSourceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "svix_ingest_source"
}

func (r *SvixIngestSourceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"environment_id": schema.StringAttribute{
				Required:    true,
				Description: ENV_ID_DESC,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"type": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.OneOf(ingestSourceInTypeKeys...),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				MarkdownDescription: "Can be one of " + strings.Join(ingestSourceInTypesForDocs(), ", "),
			},
			"name": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(2),
					stringvalidator.LengthAtMost(256),
				},
			},
			"uid": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					stringvalidator.LengthAtMost(60),
				},
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"config": schema.StringAttribute{
				Optional:   true,
				CustomType: jsontypes.NormalizedType{},
				Default:    nil,
				Sensitive:  true,
				MarkdownDescription: "The config may include sensitive fields(webhook signing secret for example)\n\n" +
					"Documentation for the config can be found in the [API docs](https://api.svix.com/docs#tag/Ingest-Source/operation/v1.ingest.source.create)",
			},
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"ingest_url": schema.StringAttribute{
				Computed: true,
				Optional: true,
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

func (r *SvixIngestSourceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// load state/plan
	var data SvixIngestSourceResourceModel
	var envId string
	var currentConfig *string
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("environment_id"), &envId)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("config"), &currentConfig)...)
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

	config, err := ingestSourceInConfigFromJsonStringAndType(currentConfig, data.Type.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Unable to parse ingest source config", err.Error())
		return
	}

	ingestIn := models.IngestSourceIn{
		Type: models.IngestSourceInTypeFromString[data.Type.ValueString()],
		Name: data.Name.ValueString(),
		Uid:  strOrNil(data.Uid),
	}
	if config != nil {
		ingestIn.Config = config
	}

	res, err := svx.Ingest.Source.Create(
		ctx,
		ingestIn,
		&svix.IngestSourceCreateOptions{
			IdempotencyKey: randStr32(),
		},
	)

	if err != nil {
		logSvixError(&resp.Diagnostics, err, "Failed to create ingest source")
		return
	}
	var configOut *string
	if currentConfig != nil {
		configOut, err = createIngestConfigFromCurrentValAndPlan(*currentConfig, res.Config)
		if err != nil {
			resp.Diagnostics.AddAttributeError(rp("config"), "Unable to save field to state", err.Error())
			return
		}
	}

	setCreateState(ctx, resp, rp("environment_id"), envId)
	setCreateState(ctx, resp, rp("id"), res.Id)
	setCreateState(ctx, resp, rp("type"), string(res.Type))
	setCreateState(ctx, resp, rp("name"), res.Name)
	setCreateState(ctx, resp, rp("uid"), res.Uid)
	setCreateState(ctx, resp, rp("config"), jsontypes.NewNormalizedPointerValue(configOut))
	setCreateState(ctx, resp, rp("ingest_url"), types.StringPointerValue(res.IngestUrl))
	setCreateState(ctx, resp, rp("created_at"), timetypes.NewRFC3339TimeValue(res.CreatedAt))
	setCreateState(ctx, resp, rp("updated_at"), timetypes.NewRFC3339TimeValue(res.UpdatedAt))
}

func (r *SvixIngestSourceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// load state/plan
	var envId, srcId string
	var currentConfig *string
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("config"), &currentConfig)...)
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("environment_id"), &envId)...)
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &srcId)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// create svix client
	svx, err := r.state.ClientWithEnvId(envId)
	if err != nil {
		resp.Diagnostics.AddError(UNABLE_TO_CREATE_SVIX_CLIENT, err.Error())
		return
	}

	res, err := svx.Ingest.Source.Get(ctx, srcId)
	if err != nil {
		logSvixError(&resp.Diagnostics, err, "Failed to get ingest source")
		return
	}

	var configOut *string
	if currentConfig != nil {
		configOut, err = createIngestConfigFromCurrentValAndPlan(*currentConfig, res.Config)
		if err != nil {
			resp.Diagnostics.AddAttributeError(rp("config"), "Unable to save field to state", err.Error())
			return
		}
	}

	setReadState(ctx, resp, rp("environment_id"), envId)
	setReadState(ctx, resp, rp("id"), res.Id)
	setReadState(ctx, resp, rp("type"), string(res.Type))
	setReadState(ctx, resp, rp("name"), res.Name)
	setReadState(ctx, resp, rp("uid"), res.Uid)
	setReadState(ctx, resp, rp("config"), jsontypes.NewNormalizedPointerValue(configOut))
	setReadState(ctx, resp, rp("ingest_url"), types.StringPointerValue(res.IngestUrl))
	setReadState(ctx, resp, rp("created_at"), timetypes.NewRFC3339TimeValue(res.CreatedAt))
	setReadState(ctx, resp, rp("updated_at"), timetypes.NewRFC3339TimeValue(res.UpdatedAt))
}

func (r *SvixIngestSourceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// load state/plan
	var data SvixIngestSourceResourceModel
	var envId, srcId string
	var currentConfig *string
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &srcId)...)
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("environment_id"), &envId)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("config"), &currentConfig)...)
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

	config, err := ingestSourceInConfigFromJsonStringAndType(currentConfig, data.Type.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Unable to parse ingest source config", err.Error())
		return
	}

	ingestIn := models.IngestSourceIn{
		Type: models.IngestSourceInTypeFromString[data.Type.ValueString()],
		Name: data.Name.ValueString(),
		Uid:  strOrNil(data.Uid),
	}
	if config != nil {
		ingestIn.Config = config
	}

	res, err := svx.Ingest.Source.Update(ctx, srcId, ingestIn)

	if err != nil {
		logSvixError(&resp.Diagnostics, err, "Failed to Update ingest source")
		return
	}
	var configOut *string
	if currentConfig != nil {
		configOut, err = createIngestConfigFromCurrentValAndPlan(*currentConfig, res.Config)
		if err != nil {
			resp.Diagnostics.AddAttributeError(rp("config"), "Unable to save field to state", err.Error())
			return
		}
	}

	setUpdateState(ctx, resp, rp("environment_id"), envId)
	setUpdateState(ctx, resp, rp("id"), res.Id)
	setUpdateState(ctx, resp, rp("type"), string(res.Type))
	setUpdateState(ctx, resp, rp("name"), res.Name)
	setUpdateState(ctx, resp, rp("uid"), res.Uid)
	setUpdateState(ctx, resp, rp("config"), jsontypes.NewNormalizedPointerValue(configOut))
	setUpdateState(ctx, resp, rp("ingest_url"), types.StringPointerValue(res.IngestUrl))
	setUpdateState(ctx, resp, rp("created_at"), timetypes.NewRFC3339TimeValue(res.CreatedAt))
	setUpdateState(ctx, resp, rp("updated_at"), timetypes.NewRFC3339TimeValue(res.UpdatedAt))

}

func (r *SvixIngestSourceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// load state/plan
	var envId, srcId string
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("environment_id"), &envId)...)
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &srcId)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// create svix client
	svx, err := r.state.ClientWithEnvId(envId)
	if err != nil {
		resp.Diagnostics.AddError(UNABLE_TO_CREATE_SVIX_CLIENT, err.Error())
		return
	}

	err = svx.Ingest.Source.Delete(ctx, srcId)
	if err != nil {
		logSvixError(&resp.Diagnostics, err, "Failed to delete ingest source")
		return
	}
}

func ingestSourceInConfigFromJsonStringAndType(jsonString *string, typ string) (models.IngestSourceInConfig, error) {
	if jsonString == nil {
		return nil, nil
	}
	de := json.NewDecoder(bytes.NewReader([]byte(*jsonString)))
	de.DisallowUnknownFields()
	var err error
	var ret models.IngestSourceInConfig
	switch typ {
	case "generic-webhook":
	case "adobe-sign":
		var c models.AdobeSignConfig
		err = de.Decode(&c)
		ret = c
	case "cron":
		var c models.CronConfig
		err = de.Decode(&c)
		ret = c
	case "docusign":
		var c models.DocusignConfig
		err = de.Decode(&c)
		ret = c
	case "github":
		var c models.GithubConfig
		err = de.Decode(&c)
		ret = c
	case "hubspot":
		var c models.HubspotConfig
		err = de.Decode(&c)
		ret = c
	case "segment":
		var c models.SegmentConfig
		err = de.Decode(&c)
		ret = c
	case "shopify":
		var c models.ShopifyConfig
		err = de.Decode(&c)
		ret = c
	case "slack":
		var c models.SlackConfig
		err = de.Decode(&c)
		ret = c
	case "stripe":
		var c models.StripeConfig
		err = de.Decode(&c)
		ret = c
	case "beehiiv", "brex", "clerk", "guesty", "incident-io", "lithic", "nash", "pleo", "replicate", "resend", "safebase", "sardine", "stych", "svix":
		var c models.SvixConfig
		err = de.Decode(&c)
		ret = c
	case "zoom":
		var c models.ZoomConfig
		err = de.Decode(&c)
		ret = c
	}
	if err != nil {
		return nil, err
	}

	return ret, nil
}

// create a copy of the plan json map, and override any with any values that exist in jsonCurrent
// this is useful when the API wont return a attribute from a config (like secret)
func createIngestConfigFromCurrentValAndPlan(configPlan string, configCurrent models.IngestSourceOutConfig) (*string, error) {
	// NOTE: this only works for json objects nested 1 level deep
	var plan map[string]any
	var current map[string]any
	err := json.Unmarshal([]byte(configPlan), &plan)
	if err != nil {
		return nil, err
	}
	configCurrentJson, err := json.Marshal(configCurrent)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(configCurrentJson, &current)
	if err != nil {
		return nil, err
	}

	// create the output map
	outMap := map[string]any{}
	// initialize it with the values in the plan
	for k, v := range plan {
		outMap[k] = v
	}

	// now fill in values for existing map
	for k, v := range current {
		outMap[k] = v
	}

	ret, err := json.Marshal(&outMap)
	if err != nil {
		return nil, err
	}

	return ptr(string(ret)), nil
}
