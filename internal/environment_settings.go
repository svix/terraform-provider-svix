package internal

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
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

var fontFamilyEnum = []string{
	"Helvetica",
	"Roboto",
	"Open Sans",
	"Lato",
	"Source Sans Pro",
	"Raleway",
	"Ubuntu",
	"Manrope",
	"DM Sans",
	"Poppins",
	"Lexend Deca",
	"Rubik",
	"Custom",
}

var colorPaletteSchema = schema.SingleNestedAttribute{
	Optional:      true,
	PlanModifiers: []planmodifier.Object{objectplanmodifier.UseStateForUnknown()},
	Attributes: map[string]schema.Attribute{
		"primary": schema.StringAttribute{
			PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			Optional:      true,
			Description:   "Primary color",
		},
		"background": schema.StringAttribute{
			PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			Optional:      true,
			Description:   "Background",
		},
		"surface_background": schema.StringAttribute{
			PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			Optional:            true,
			Description:         "Surface Background",
			MarkdownDescription: "Background for cards, tables and other surfaces",
		},
		"surface_hover": schema.StringAttribute{
			PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			Optional:            true,
			Description:         "Surface Header",
			MarkdownDescription: "Background for card headers and table headers",
		},
		"interactive_accent": schema.StringAttribute{
			PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			Optional:            true,
			Description:         "Interactive Accent",
			MarkdownDescription: "For secondary buttons, links, and other interactive elements",
		},
		"navigation_accent": schema.StringAttribute{
			PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			Optional:            true,
			Description:         "Navigation Accent",
			MarkdownDescription: "For the top-level navigation items",
		},
		"button_primary": schema.StringAttribute{
			PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			Optional:            true,
			Description:         "Button Primary",
			MarkdownDescription: "For the main action buttons",
		},
		"text_primary": schema.StringAttribute{
			PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			Optional:      true,
			Description:   "Text Primary",
		},
		"text_danger": schema.StringAttribute{
			PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			Optional:            true,
			Description:         "Text Danger",
			MarkdownDescription: "For error messages and other warnings",
		},
	},
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
			"whitelabel_settings": schema.SingleNestedAttribute{
				PlanModifiers:       []planmodifier.Object{objectplanmodifier.UseStateForUnknown()},
				Optional:            true,
				MarkdownDescription: "Customize how the [Consumer App Portal](https://docs.svix.com/management-ui) will look for your users in this environment.",
				Attributes: map[string]schema.Attribute{
					"display_name": schema.StringAttribute{
						Optional:            true,
						Computed:            true,
						PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
						Description:         "Display Name",
						MarkdownDescription: "The name of your company or service. Visible to users in the App Portal and the [Event Catalog](https://docs.svix.com/event-types#publishing-your-event-catalog).",
					},
					"base_font_size": schema.Int64Attribute{
						PlanModifiers: []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
						Validators: []validator.Int64{
							// this is the limit we have in the front-end
							int64validator.AtMost(23),
							int64validator.AtLeast(8),
						},
						Optional:            true,
						MarkdownDescription: "This affects all text size on the screen relative to the size of the text in the main body of the page. Default: 16px",
						Description:         "Base Font Size (in pixels)",
					},
					"font_family": schema.StringAttribute{
						PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
						Validators: []validator.String{
							stringvalidator.OneOf(fontFamilyEnum...),
							customFontFamilyValidator{},
						},
						Optional:    true,
						Description: "Font Family",
						MarkdownDescription: "Can be one of `" + strings.Join(fontFamilyEnum[:len(fontFamilyEnum)-1], "`, `") + "` and `Custom`\n\n" +
							"You can also set a custom font by providing a URL to a font file. \n\n" +
							"If you chose to use the `font_family_url` make sure to set this to `Custom`\n",
					},
					"font_family_url": schema.StringAttribute{
						PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
						Validators:    []validator.String{customFontURLValidator{}},
						Optional:      true,
						Description:   "Font Family URL",
						MarkdownDescription: "URL of a woff2 font file (e.g. https://fonts.gstatic.com/s/librebaskerville.woff2)\n\n" +
							"Make sure to set `font_family` to `Custom`",
					},
					"logo_url": schema.StringAttribute{
						PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
						Optional:            true,
						Description:         "Icon URL",
						MarkdownDescription: "Used in the standalone App Portal experience. Not visible in the [embedded App Portal](https://docs.svix.com/management-ui).",
					},

					"color_palette_dark":  colorPaletteSchema,
					"color_palette_light": colorPaletteSchema,

					"border_radius": schema.SingleNestedAttribute{
						Optional:    true,
						Description: "Borders",
						Attributes: map[string]schema.Attribute{
							"button": schema.StringAttribute{
								PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
								Optional:      true,
								Validators: []validator.String{
									stringvalidator.OneOf(borderRadiusEnum...),
								},
								Description:         "Button corners",
								MarkdownDescription: "Use `none` for a square border, `lg` for large rounded `md` for medium rounded, `sm` for small rounded and `full` for Pill-shaped",
							},
							"card": schema.StringAttribute{
								PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
								Optional:      true,
								Validators: []validator.String{
									stringvalidator.OneOf(borderRadiusEnum...),
								},
								Description:         "Card corners",
								MarkdownDescription: "Use `none` for a square border, `lg` for large rounded `md` for medium rounded, `sm` for small rounded and `full` for Pill-shaped",
							},
							"input": schema.StringAttribute{
								PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
								Optional:      true,
								Validators: []validator.String{
									stringvalidator.OneOf(borderRadiusEnum...),
								},
								Description:         "Input corners",
								MarkdownDescription: "Use `none` for a square border, `lg` for large rounded `md` for medium rounded, `sm` for small rounded and `full` for Pill-shaped",
							},
						},
					},
					"channels_strings_override": schema.SingleNestedAttribute{
						PlanModifiers: []planmodifier.Object{objectplanmodifier.UseStateForUnknown()},
						Optional:      true,
						Description:   "Rename 'channels' in the App Portal, depending on the usage you give them in your application.",
						Attributes: map[string]schema.Attribute{
							"channels_help": schema.StringAttribute{
								PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
								Optional:      true,
								Description:   "Channels help text.",
							},
							"channels_many": schema.StringAttribute{
								PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
								Optional:      true,
								Description:   "Plural form.",
							},
							"channels_one": schema.StringAttribute{
								PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
								Optional:      true,
								Description:   "Singular form.",
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
			"enable_advanced_endpoint_types": schema.BoolAttribute{
				PlanModifiers: []planmodifier.Bool{boolplanmodifier.UseStateForUnknown()},
				Optional:      true,
				Computed:      true,
				Description:   "Advanced endpoint types",
				MarkdownDescription: REQUIRES_PRO_OR_ENTERPRISE_PLAN + `Allows users to configure Polling Endpoints and FIFO endpoints to get
messages. Read more about them in the [docs](https://docs.svix.com/advanced-endpoints/intro).`,
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
				MarkdownDescription: "Enable this to make your Event Catalog public. " +
					"You can find the link to the published Event Catalog at https://dashboard.svix.com/settings/organization/catalog",
			},
			"read_only": schema.BoolAttribute{
				PlanModifiers:       []planmodifier.Bool{boolplanmodifier.UseStateForUnknown()},
				Optional:            true,
				Computed:            true,
				Description:         "Read Only mode",
				MarkdownDescription: `Sets your Consumer App Portal to read only so your customers can view but not modify their data`,
			},
			"require_channel_filtering": schema.BoolAttribute{
				PlanModifiers:       []planmodifier.Bool{boolplanmodifier.UseStateForUnknown()},
				Optional:            true,
				Computed:            true,
				Description:         "Require channel filters for endpoints",
				MarkdownDescription: "If enabled, all new Endpoints must filter on at least one channel.",
			},
			"whitelabel_headers": schema.BoolAttribute{
				PlanModifiers: []planmodifier.Bool{boolplanmodifier.UseStateForUnknown()},
				Optional:      true,
				Computed:      true,
				Description:   "White label headers",
				MarkdownDescription: REQUIRES_PRO_OR_ENTERPRISE_PLAN +
					"Changes the prefix of the webhook HTTP headers to use the" +
					"`webhook-` prefix. <strong>Changing this setting can break existing integrations</strong>",
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

	Spw(data)

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
	Spw(req.State.Raw)
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
	Spw(res)
	outModel := internalSettingsOutToTF(ctx, &resp.Diagnostics, *res, envId)

	resp.Diagnostics.Append(resp.State.Set(ctx, outModel)...)
	Spw(resp.State.Raw)
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
	out := generated.EnvironmentSettingsResourceModel{
		WhitelabelSettings:          basetypes.NewObjectNull(generated.WhitelabelSettings_TF_AttributeTypes()),
		EnvironmentId:               types.StringValue(envId),
		DisableEndpointOnFailure:    types.BoolPointerValue(v.DisableEndpointOnFailure),
		EnableChannels:              types.BoolPointerValue(v.EnableChannels),
		EnableEndpointMtlsConfig:    types.BoolPointerValue(v.EnableEndpointMtlsConfig),
		EnableEndpointOauthConfig:   types.BoolPointerValue(v.EnableEndpointOauthConfig),
		EnableIntegrationManagement: types.BoolPointerValue(v.EnableIntegrationManagement),
		EnableMessageStream:         types.BoolPointerValue(v.EnableMessageStream),
		EnableTransformations:       types.BoolPointerValue(v.EnableTransformations),
		EnforceHttps:                types.BoolPointerValue(v.EnforceHttps),
		EventCatalogPublished:       types.BoolPointerValue(v.EventCatalogPublished),
		ReadOnly:                    types.BoolPointerValue(v.ReadOnly),
		RequireEndpointChannel:      types.BoolPointerValue(v.RequireEndpointChannel),
		WhitelabelHeaders:           types.BoolPointerValue(v.WhitelabelHeaders),
		WipeSuccessfulPayload:       types.BoolPointerValue(v.WipeSuccessfulPayload),
	}

	{

		whitelabelSettingsTf := generated.WhitelabelSettings{
			CustomStringsOverride: basetypes.NewObjectNull(generated.CustomStringsOverride_TF_AttributeTypes()),
			BorderRadius:          basetypes.NewObjectNull(generated.BorderRadius_AttributeTypes()),
			ColorPaletteDark:      basetypes.NewObjectNull(generated.CustomColorPalette_TF_AttributeTypes()),
			ColorPaletteLight:     basetypes.NewObjectNull(generated.CustomColorPalette_TF_AttributeTypes()),
			DisplayName:           types.StringPointerValue(v.DisplayName),
			CustomBaseFontSize:    types.Int64PointerValue(v.CustomBaseFontSize),
			CustomFontFamily:      types.StringPointerValue(v.CustomFontFamily),
			CustomFontFamilyUrl:   types.StringPointerValue(v.CustomFontFamilyUrl),
			CustomLogoUrl:         types.StringPointerValue(v.CustomLogoUrl),
		}
		if v.CustomThemeOverride != nil {
			if v.CustomThemeOverride.BorderRadius != nil {
				existingBorderRadius := v.CustomThemeOverride.BorderRadius
				borderRadiusTf := generated.BorderRadius{
					Button: types.StringPointerValue((*string)(existingBorderRadius.Button)),
					Card:   types.StringPointerValue((*string)(existingBorderRadius.Card)),
					Input:  types.StringPointerValue((*string)(existingBorderRadius.Input)),
				}
				borderRadius, diags := types.ObjectValueFrom(ctx, borderRadiusTf.AttributeTypes(), borderRadiusTf)
				d.Append(diags...)
				whitelabelSettingsTf.BorderRadius = borderRadius
			}
		}
		if v.ColorPaletteDark != nil {
			colorPaletteDarkTf := customColorPaletteToTF(*v.ColorPaletteDark)
			colorPaletteDark, diags := types.ObjectValueFrom(ctx, colorPaletteDarkTf.AttributeTypes(), colorPaletteDarkTf)
			whitelabelSettingsTf.ColorPaletteDark = colorPaletteDark
			d.Append(diags...)
		}
		if v.ColorPaletteLight != nil {
			colorPaletteLightTf := customColorPaletteToTF(*v.ColorPaletteLight)
			colorPaletteLight, diags := types.ObjectValueFrom(ctx, colorPaletteLightTf.AttributeTypes(), colorPaletteLightTf)
			whitelabelSettingsTf.ColorPaletteLight = colorPaletteLight
			d.Append(diags...)
		}

		if v.CustomStringsOverride != nil {
			customStringsOverrideTF := generated.CustomStringsOverride_TF{
				ChannelsHelp: types.StringPointerValue(v.CustomStringsOverride.ChannelsHelp),
				ChannelsMany: types.StringPointerValue(v.CustomStringsOverride.ChannelsMany),
				ChannelsOne:  types.StringPointerValue(v.CustomStringsOverride.ChannelsOne),
			}
			customStringsOverride, diags := types.ObjectValueFrom(ctx, customStringsOverrideTF.AttributeTypes(), customStringsOverrideTF)
			whitelabelSettingsTf.CustomStringsOverride = customStringsOverride
			d.Append(diags...)
		}

		whitelabelSettings, diags := types.ObjectValueFrom(ctx, whitelabelSettingsTf.AttributeTypes(), whitelabelSettingsTf)
		d.Append(diags...)
		out.WhitelabelSettings = whitelabelSettings
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

type customFontURLValidator struct{}

// Description returns a description of the validator
func (v customFontURLValidator) Description(ctx context.Context) string {
	return "When a font URL is provided, font_family must be set to 'Custom'"
}

// MarkdownDescription returns a markdown description of the validator
func (v customFontURLValidator) MarkdownDescription(ctx context.Context) string {
	return "When a font URL is provided, font_family must be set to 'Custom'"
}

// ValidateString performs the validation
func (v customFontURLValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	// Skip validation if value is unknown or null
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	// Get the font URL value
	fontURL := req.ConfigValue.ValueString()

	// Skip validation if URL is empty
	if fontURL == "" {
		return
	}

	// Get the font_family value to check
	var fontFamily types.String
	diags := req.Config.GetAttribute(ctx, path.Root("whitelabel_settings").AtName("font_family"), &fontFamily)

	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	// If font_family is unknown or null, add error
	if fontFamily.IsNull() || fontFamily.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid Font Configuration",
			"When providing a font_family_url, the font_family attribute must be set to 'Custom'",
		)
		return
	}

	// Verify that font_family is set to "Custom"
	if fontFamily.ValueString() != "Custom" {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid Font Configuration",
			"When providing a font_family_url, the font_family attribute must be set to 'Custom'",
		)
	}
}

// customFontFamilyValidator validates that when font_family is set to "Custom",
// a font_family_url must be provided
type customFontFamilyValidator struct{}

// Description returns a description of the validator
func (v customFontFamilyValidator) Description(ctx context.Context) string {
	return "When font_family is set to 'Custom', a font_family_url must be provided"
}

// MarkdownDescription returns a markdown description of the validator
func (v customFontFamilyValidator) MarkdownDescription(ctx context.Context) string {
	return "When font_family is set to 'Custom', a font_family_url must be provided"
}

// ValidateString performs the validation
func (v customFontFamilyValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	// Skip validation if value is unknown or null
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	// Get the font family value
	fontFamily := req.ConfigValue.ValueString()

	// Skip validation if not set to "Custom"
	if fontFamily != "Custom" {
		return
	}

	// Get the font_family_url value to check
	var fontURL types.String
	diags := req.Config.GetAttribute(ctx, path.Root("whitelabel_settings").AtName("font_family_url"), &fontURL)

	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	// If font_family_url is null, unknown, or empty, add error
	if fontURL.IsNull() || fontURL.IsUnknown() || fontURL.ValueString() == "" {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Missing Font URL",
			"When font_family is set to 'Custom', a valid font_family_url must be provided",
		)
	}
}
