package internal

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/svix/svix-webhooks/go/models"
	"github.com/svix/terraform-provider-svix/internal/generated"
)

var _ resource.Resource = &EnvironmentSettingsResource{}

func NewEnvironmentSettingsResource() resource.Resource {
	return &EnvironmentSettingsResource{}
}

type EnvironmentSettingsResource struct {
	state appState
}

var borderRadiusEnum = []string{
	"none",
	"lg",
	"md",
	"sm",
	"full",
}

// Metadata implements resource.Resource.
func (r *EnvironmentSettingsResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "svix_environment_settings"
}

func (r *EnvironmentSettingsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"environment_id": schema.StringAttribute{
				Required:    true,
				Description: ENV_ID_DESC,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"color_palette_dark": schema.SingleNestedAttribute{
				Optional:      true,
				Computed:      true,
				PlanModifiers: []planmodifier.Object{objectplanmodifier.UseStateForUnknown()},
				Attributes: map[string]schema.Attribute{
					"background_hover": schema.StringAttribute{
						PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
						Optional:      true,
						Computed:      true,
					},
					"background_primary": schema.StringAttribute{
						PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
						Optional:      true,
						Computed:      true,
					},
					"background_secondary": schema.StringAttribute{
						PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
						Optional:      true,
						Computed:      true,
					},
					"button_primary": schema.StringAttribute{
						PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
						Optional:      true,
						Computed:      true,
					},
					"interactive_accent": schema.StringAttribute{
						PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
						Optional:      true,
						Computed:      true,
					},
					"navigation_accent": schema.StringAttribute{
						PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
						Optional:      true,
						Computed:      true,
					},
					"primary": schema.StringAttribute{
						PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
						Optional:      true,
						Computed:      true,
					},
					"text_danger": schema.StringAttribute{
						PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
						Optional:      true,
						Computed:      true,
					},
					"text_primary": schema.StringAttribute{
						PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
						Optional:      true,
						Computed:      true,
					},
				},
			},
			"color_palette_light": schema.SingleNestedAttribute{
				Optional:      true,
				Computed:      true,
				PlanModifiers: []planmodifier.Object{objectplanmodifier.UseStateForUnknown()},
				Attributes: map[string]schema.Attribute{
					"background_hover": schema.StringAttribute{
						PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
						Optional:      true,
						Computed:      true,
					},
					"background_primary": schema.StringAttribute{
						PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
						Optional:      true,
						Computed:      true,
					},
					"background_secondary": schema.StringAttribute{
						PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
						Optional:      true,
						Computed:      true,
					},
					"button_primary": schema.StringAttribute{
						PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
						Optional:      true,
						Computed:      true,
					},
					"interactive_accent": schema.StringAttribute{
						PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
						Optional:      true,
						Computed:      true,
					},
					"navigation_accent": schema.StringAttribute{
						PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
						Optional:      true,
						Computed:      true,
					},
					"primary": schema.StringAttribute{
						PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
						Optional:      true,
						Computed:      true,
					},
					"text_danger": schema.StringAttribute{
						PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
						Optional:      true,
						Computed:      true,
					},
					"text_primary": schema.StringAttribute{
						PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
						Optional:      true,
						Computed:      true,
					},
				},
			},
			"custom_base_font_size": schema.Int64Attribute{
				PlanModifiers: []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
				Optional:      true,
				Computed:      true,
			},
			"custom_color": schema.StringAttribute{
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
				Optional:      true,
				Computed:      true,
			},
			"custom_font_family": schema.StringAttribute{
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
				Optional:      true,
				Computed:      true,
			},
			"custom_font_family_url": schema.StringAttribute{
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
				Optional:      true,
				Computed:      true,
			},
			"custom_logo_url": schema.StringAttribute{
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
				Optional:      true,
				Computed:      true,
			},
			"custom_strings_override": schema.SingleNestedAttribute{
				PlanModifiers: []planmodifier.Object{objectplanmodifier.UseStateForUnknown()},
				Optional:      true,
				Computed:      true,
				Attributes: map[string]schema.Attribute{
					"channels_help": schema.StringAttribute{
						PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
						Optional:      true,
						Computed:      true,
					},
					"channels_many": schema.StringAttribute{
						PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
						Optional:      true,
						Computed:      true,
					},
					"channels_one": schema.StringAttribute{
						PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
						Optional:      true,
						Computed:      true,
					},
				},
			},
			"custom_theme_override": schema.SingleNestedAttribute{
				PlanModifiers: []planmodifier.Object{objectplanmodifier.UseStateForUnknown()},
				Optional:      true,
				Computed:      true,
				Attributes: map[string]schema.Attribute{
					"border_radius": schema.SingleNestedAttribute{
						PlanModifiers: []planmodifier.Object{objectplanmodifier.UseStateForUnknown()},
						Optional:      true,
						Computed:      true,
						Attributes: map[string]schema.Attribute{
							"button": schema.StringAttribute{
								PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
								Optional:      true,
								Computed:      true,
								Validators: []validator.String{
									stringvalidator.OneOf(borderRadiusEnum...),
								},
							},
							"card": schema.StringAttribute{
								PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
								Optional:      true,
								Computed:      true,
								Validators: []validator.String{
									stringvalidator.OneOf(borderRadiusEnum...),
								},
							},
							"input": schema.StringAttribute{
								PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
								Optional:      true,
								Computed:      true,
								Validators: []validator.String{
									stringvalidator.OneOf(borderRadiusEnum...),
								},
							},
						},
					},
					"font_size": schema.SingleNestedAttribute{
						PlanModifiers: []planmodifier.Object{objectplanmodifier.UseStateForUnknown()},
						Optional:      true,
						Computed:      true,
						Attributes: map[string]schema.Attribute{
							"base": schema.Int64Attribute{
								PlanModifiers: []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
								Optional:      true,
								Computed:      true,
							},
						},
					},
				},
			},
			"disable_endpoint_on_failure": schema.BoolAttribute{
				PlanModifiers: []planmodifier.Bool{boolplanmodifier.UseStateForUnknown()},
				Optional:      true,
				Computed:      true,
				Description:   "Disable endpoint on failure",
				MarkdownDescription: `If messages to a particular endpoint have been consistently failing for
some time, we will automatically disable the endpoint and let 
you know [via webhook](https://docs.svix.com/incoming-webhooks). Read 
more about it [in the docs](https://docs.svix.com/retries#disabling-failing-endpoints).`,
			},
			"display_name": schema.StringAttribute{
				Optional:      true,
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"enable_channels": schema.BoolAttribute{
				PlanModifiers: []planmodifier.Bool{boolplanmodifier.UseStateForUnknown()},
				Optional:      true,
				Computed:      true,
				Description:   "Enable Channels",
				MarkdownDescription: `Controls whether or not your users can configure
<strong>channels</strong> from the Consumer App Portal.`,
			},
			"enable_endpoint_mtls_config": schema.BoolAttribute{
				PlanModifiers:       []planmodifier.Bool{boolplanmodifier.UseStateForUnknown()},
				Optional:            true,
				Computed:            true,
				Description:         "Enable mTLS configuration",
				MarkdownDescription: REQUIRES_ENTERPRISE_PLAN + "Allows users to configure mutual TLS (mTLS) for their endpoints.",
			},
			"enable_endpoint_oauth_config": schema.BoolAttribute{
				PlanModifiers:       []planmodifier.Bool{boolplanmodifier.UseStateForUnknown()},
				Optional:            true,
				Computed:            true,
				Description:         "Enable OAuth configuration",
				MarkdownDescription: REQUIRES_ENTERPRISE_PLAN + "Allows users to configure OAuth for their endpoints.",
			},
			"enable_integration_management": schema.BoolAttribute{
				PlanModifiers: []planmodifier.Bool{boolplanmodifier.UseStateForUnknown()},
				Optional:      true,
				Computed:      true,
				Description:   "Allow users to manage integrations",
				MarkdownDescription: `Controls whether or not your users can manage integrations from the
Consumer App Portal. We recommend disabling this if you manage
integrations on your users' behalf.`,
			},
			"enable_message_stream": schema.BoolAttribute{
				PlanModifiers: []planmodifier.Bool{boolplanmodifier.UseStateForUnknown()},
				Optional:      true,
				Computed:      true,
				Description:   "Advanced endpoint types",
				MarkdownDescription: REQUIRES_PRO_OR_ENTERPRISE_PLAN + `Allows users to configure Polling Endpoints and FIFO endpoints to get
messages. Read more about them in the [docs](https://docs.svix.com/advanced-endpoints/intro).`,
			},
			"enable_msg_atmpt_log": schema.BoolAttribute{
				PlanModifiers: []planmodifier.Bool{boolplanmodifier.UseStateForUnknown()},
				Optional:      true,
				Computed:      true,
			},
			"enable_otlp": schema.BoolAttribute{
				PlanModifiers: []planmodifier.Bool{boolplanmodifier.UseStateForUnknown()},
				Optional:      true,
				Computed:      true,
			},
			"enable_transformations": schema.BoolAttribute{
				PlanModifiers: []planmodifier.Bool{boolplanmodifier.UseStateForUnknown()},
				Optional:      true,
				Computed:      true,
				Description:   "Transformations",
				MarkdownDescription: `Controls whether or not your users can add transformations to their
endpoints. Transformations are code that can change a message's HTTP
method, destination URL, and payload body in-flight.`,
			},
			"enforce_https": schema.BoolAttribute{
				PlanModifiers:       []planmodifier.Bool{boolplanmodifier.UseStateForUnknown()},
				Optional:            true,
				Computed:            true,
				Description:         "HTTPS Only Endpoints",
				MarkdownDescription: "Enforces HTTPS on all endpoints of this environment",
			},
			"event_catalog_published": schema.BoolAttribute{
				PlanModifiers: []planmodifier.Bool{boolplanmodifier.UseStateForUnknown()},
				Optional:      true,
				Computed:      true,
			},
			"read_only": schema.BoolAttribute{
				PlanModifiers:       []planmodifier.Bool{boolplanmodifier.UseStateForUnknown()},
				Optional:            true,
				Computed:            true,
				Description:         "Read Only mode",
				MarkdownDescription: `Sets your Consumer App Portal to read only so your customers can view but not modify their data`,
			},
			"require_endpoint_channel": schema.BoolAttribute{
				PlanModifiers:       []planmodifier.Bool{boolplanmodifier.UseStateForUnknown()},
				Optional:            true,
				Computed:            true,
				Description:         "Require channel filters for endpoints",
				MarkdownDescription: "If enabled, all new Endpoints must filter on at least one channel.",
			},
			"retry_policy": schema.ListAttribute{
				PlanModifiers: []planmodifier.List{listplanmodifier.UseStateForUnknown()},
				Optional:      true,
				Computed:      true,
				ElementType:   types.StringType,
			},
			"show_use_svix_play": schema.BoolAttribute{
				PlanModifiers: []planmodifier.Bool{boolplanmodifier.UseStateForUnknown()},
				Optional:      true,
				Computed:      true,
			},
			"whitelabel_headers": schema.BoolAttribute{
				PlanModifiers: []planmodifier.Bool{boolplanmodifier.UseStateForUnknown()},
				Optional:      true,
				Computed:      true,
				Description:   "White label headers",
				MarkdownDescription: REQUIRES_PRO_OR_ENTERPRISE_PLAN +
					"Changes the prefix of the webhook HTTP headers to use the" +
					"`webhook-` prefix. <string>Changing this setting can break existing integrations<strong/>",
			},
			"wipe_successful_payload": schema.BoolAttribute{
				PlanModifiers: []planmodifier.Bool{boolplanmodifier.UseStateForUnknown()},
				Optional:      true,
				Computed:      true,
				Description:   "Delete successful message payloads",
				MarkdownDescription: REQUIRES_PRO_OR_ENTERPRISE_PLAN + `Delete message payloads from Svix after they are successfully
delivered to the endpoint. Only affects messages sent after this
setting is enabled.`,
			},
		},
	}

}

func (r *EnvironmentSettingsResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// we won't be created the settings here. rather for any defined field, we will run an update query
func (r *EnvironmentSettingsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// load state/plan
	var data generated.EnvironmentSettingsResourceModel
	var envId string
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("environment_id"), &envId)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// create svix client
	svx, err := r.state.InternalClientWithEnvId(envId)
	if err != nil {
		resp.Diagnostics.AddError(UNABLE_TO_CREATE_SVIX_CLIENT, err.Error())
		return
	}

	currentSettings, err := svx.Management.EnvironmentSettings.Get(ctx)
	if err != nil {
		logSvixError(&resp.Diagnostics, err, "Failed to get environment settings")
		return
	}
	settingsIn := generated.PatchSettingsInternalInWithPlan(ctx, &resp.Diagnostics, currentSettings, data)
	if resp.Diagnostics.HasError() {
		return
	}

	// call api
	res, err := svx.Management.EnvironmentSettings.Update(ctx, settingsIn)
	if err != nil {
		logSvixError(&resp.Diagnostics, err, "Failed to update environment settings")
		return
	}

	outModel := internalSettingsOutToTF(ctx, &resp.Diagnostics, *res, envId)
	resp.Diagnostics.Append(resp.State.Set(ctx, outModel)...)
}

func (r *EnvironmentSettingsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// load state/plan
	var envId string
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("environment_id"), &envId)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// create svix client
	svx, err := r.state.InternalClientWithEnvId(envId)
	if err != nil {
		resp.Diagnostics.AddError(UNABLE_TO_CREATE_SVIX_CLIENT, err.Error())
		return
	}

	res, err := svx.Management.EnvironmentSettings.Get(ctx)
	if err != nil {
		logSvixError(&resp.Diagnostics, err, "Failed to get environment settings")
		return
	}

	outModel := internalSettingsOutToTF(ctx, &resp.Diagnostics, *res, envId)
	resp.Diagnostics.Append(resp.State.Set(ctx, outModel)...)
}

func (r *EnvironmentSettingsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// load state/plan
	var data generated.EnvironmentSettingsResourceModel
	var envId string
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("environment_id"), &envId)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// create svix client
	svx, err := r.state.InternalClientWithEnvId(envId)
	if err != nil {
		resp.Diagnostics.AddError(UNABLE_TO_CREATE_SVIX_CLIENT, err.Error())
		return
	}

	currentSettings, err := svx.Management.EnvironmentSettings.Get(ctx)
	if err != nil {
		logSvixError(&resp.Diagnostics, err, "Failed to get environment settings")
		return
	}
	settingsIn := generated.PatchSettingsInternalInWithPlan(ctx, &resp.Diagnostics, currentSettings, data)
	if resp.Diagnostics.HasError() {
		return
	}

	// call api
	res, err := svx.Management.EnvironmentSettings.Update(ctx, settingsIn)
	if err != nil {
		logSvixError(&resp.Diagnostics, err, "Failed to update environment settings")
		return
	}

	outModel := internalSettingsOutToTF(ctx, &resp.Diagnostics, *res, envId)
	resp.Diagnostics.Append(resp.State.Set(ctx, outModel)...)
}

func (r *EnvironmentSettingsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// we don't delete the env settings
	// env delete will delete the env settings
}

func customColorPaletteToTF(v models.CustomColorPalette) generated.CustomColorPalette_TF {
	return generated.CustomColorPalette_TF{
		BackgroundHover:     types.StringPointerValue(v.BackgroundHover),
		BackgroundPrimary:   types.StringPointerValue(v.BackgroundPrimary),
		BackgroundSecondary: types.StringPointerValue(v.BackgroundSecondary),
		ButtonPrimary:       types.StringPointerValue(v.ButtonPrimary),
		InteractiveAccent:   types.StringPointerValue(v.InteractiveAccent),
		NavigationAccent:    types.StringPointerValue(v.NavigationAccent),
		Primary:             types.StringPointerValue(v.Primary),
		TextDanger:          types.StringPointerValue(v.TextDanger),
		TextPrimary:         types.StringPointerValue(v.TextPrimary),
	}

}
func internalSettingsOutToTF(ctx context.Context, d *diag.Diagnostics, v models.SettingsInternalOut, envId string) generated.EnvironmentSettingsResourceModel {
	retryPolicy, diags := types.ListValueFrom(ctx, types.StringType, v.RetryPolicy)
	d.Append(diags...)

	out := generated.EnvironmentSettingsResourceModel{
		ColorPaletteDark:            basetypes.NewObjectNull(generated.CustomColorPalette_TF_AttributeTypes()),
		ColorPaletteLight:           basetypes.NewObjectNull(generated.CustomColorPalette_TF_AttributeTypes()),
		CustomStringsOverride:       basetypes.NewObjectNull(generated.CustomStringsOverride_TF_AttributeTypes()),
		CustomThemeOverride:         basetypes.NewObjectNull(generated.CustomThemeOverride_TF_AttributeTypes()),
		EnvironmentId:               types.StringValue(envId),
		CustomBaseFontSize:          types.Int64PointerValue(v.CustomBaseFontSize),
		CustomColor:                 types.StringPointerValue(v.CustomColor),
		CustomFontFamily:            types.StringPointerValue(v.CustomFontFamily),
		CustomFontFamilyUrl:         types.StringPointerValue(v.CustomFontFamilyUrl),
		CustomLogoUrl:               types.StringPointerValue(v.CustomLogoUrl),
		DisableEndpointOnFailure:    types.BoolPointerValue(v.DisableEndpointOnFailure),
		DisplayName:                 types.StringPointerValue(v.DisplayName),
		EnableChannels:              types.BoolPointerValue(v.EnableChannels),
		EnableEndpointMtlsConfig:    types.BoolPointerValue(v.EnableEndpointMtlsConfig),
		EnableEndpointOauthConfig:   types.BoolPointerValue(v.EnableEndpointOauthConfig),
		EnableIntegrationManagement: types.BoolPointerValue(v.EnableIntegrationManagement),
		EnableMessageStream:         types.BoolPointerValue(v.EnableMessageStream),
		EnableMsgAtmptLog:           types.BoolPointerValue(v.EnableMsgAtmptLog),
		EnableOtlp:                  types.BoolPointerValue(v.EnableOtlp),
		EnableTransformations:       types.BoolPointerValue(v.EnableTransformations),
		EnforceHttps:                types.BoolPointerValue(v.EnforceHttps),
		EventCatalogPublished:       types.BoolPointerValue(v.EventCatalogPublished),
		ReadOnly:                    types.BoolPointerValue(v.ReadOnly),
		RequireEndpointChannel:      types.BoolPointerValue(v.RequireEndpointChannel),
		RetryPolicy:                 retryPolicy,
		ShowUseSvixPlay:             types.BoolPointerValue(v.ShowUseSvixPlay),
		WhitelabelHeaders:           types.BoolPointerValue(v.WhitelabelHeaders),
		WipeSuccessfulPayload:       types.BoolPointerValue(v.WipeSuccessfulPayload),
	}
	if v.ColorPaletteDark != nil {
		colorPaletteDarkTf := customColorPaletteToTF(*v.ColorPaletteDark)
		colorPaletteDark, diags := types.ObjectValueFrom(ctx, colorPaletteDarkTf.AttributeTypes(), colorPaletteDarkTf)
		out.ColorPaletteDark = colorPaletteDark
		d.Append(diags...)
	}
	if v.ColorPaletteLight != nil {
		Spw(v.ColorPaletteLight)
		colorPaletteLightTf := customColorPaletteToTF(*v.ColorPaletteLight)
		Spw(colorPaletteLightTf)
		colorPaletteLight, diags := types.ObjectValueFrom(ctx, colorPaletteLightTf.AttributeTypes(), colorPaletteLightTf)
		Spw(colorPaletteLight)
		out.ColorPaletteLight = colorPaletteLight
		d.Append(diags...)
	}

	if v.CustomStringsOverride != nil {
		customStringsOverrideTF := generated.CustomStringsOverride_TF{
			ChannelsHelp: types.StringPointerValue(v.CustomStringsOverride.ChannelsHelp),
			ChannelsMany: types.StringPointerValue(v.CustomStringsOverride.ChannelsMany),
			ChannelsOne:  types.StringPointerValue(v.CustomStringsOverride.ChannelsOne),
		}
		customStringsOverride, diags := types.ObjectValueFrom(ctx, customStringsOverrideTF.AttributeTypes(), customStringsOverrideTF)
		out.CustomStringsOverride = customStringsOverride
		d.Append(diags...)
	}

	if v.CustomThemeOverride != nil {
		customThemeOverrideTF := generated.CustomThemeOverride_TF{
			FontSize:     basetypes.NewObjectNull(generated.FontSizeConfig_TF_AttributeTypes()),
			BorderRadius: basetypes.NewObjectNull(generated.BorderRadiusConfig_TF_AttributeTypes()),
		}
		if v.CustomThemeOverride.BorderRadius != nil {
			borderRadiusTF := generated.BorderRadiusConfig_TF{
				Button: BorderRadiusEnumStringValue(v.CustomThemeOverride.BorderRadius.Button),
				Card:   BorderRadiusEnumStringValue(v.CustomThemeOverride.BorderRadius.Card),
				Input:  BorderRadiusEnumStringValue(v.CustomThemeOverride.BorderRadius.Input),
			}
			borderRadius, diags := types.ObjectValueFrom(ctx, borderRadiusTF.AttributeTypes(), borderRadiusTF)
			customThemeOverrideTF.BorderRadius = borderRadius
			d.Append(diags...)
		}

		if v.CustomThemeOverride.FontSize != nil {
			var base types.Int64
			if v.CustomThemeOverride.FontSize.Base != nil {
				base = basetypes.NewInt64Value(int64(*v.CustomThemeOverride.FontSize.Base))
			} else {
				base = basetypes.NewInt64Null()
			}
			fontSizeTF := generated.FontSizeConfig_TF{Base: base}
			fontSize, diags := types.ObjectValueFrom(ctx, fontSizeTF.AttributeTypes(), fontSizeTF)
			customThemeOverrideTF.FontSize = fontSize
			d.Append(diags...)
		}
		customThemeOverride, diags := types.ObjectValueFrom(ctx, customThemeOverrideTF.AttributeTypes(), customThemeOverrideTF)
		out.CustomThemeOverride = customThemeOverride
		d.Append(diags...)
	}

	return out
}

// golang is f***ing dumb ):
func BorderRadiusEnumStringValue(v *models.BorderRadiusEnum) basetypes.StringValue {
	if v == nil {
		return types.StringPointerValue(nil)
	}
	switch *v {
	case models.BORDERRADIUSENUM_NONE:
		return types.StringValue("none")
	case models.BORDERRADIUSENUM_LG:
		return types.StringValue("lg")
	case models.BORDERRADIUSENUM_MD:
		return types.StringValue("md")
	case models.BORDERRADIUSENUM_SM:
		return types.StringValue("sm")
	case models.BORDERRADIUSENUM_FULL:
		return types.StringValue("full")
	}
	return types.StringPointerValue(nil)
}
