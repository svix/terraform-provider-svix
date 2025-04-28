// Package generated. most of the code here was generated using the codegen and a custom template
package generated

import (
	"context"
	"log"

	"github.com/davecgh/go-spew/spew"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	svixmodels "github.com/svix/svix-webhooks/go/models"
)

func ptr[T any](value T) *T {
	return &value
}
func Spw(v any) {
	log.Println(spew.Sdump(v))
}

// Terraform wrapper around `svixmodels.FontSizeConfig`
type FontSizeConfig_TF struct {
	Base types.Int64 `tfsdk:"base"`
}

func FontSizeConfig_TF_AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"base": types.Int64Type,
	}
}
func (v *FontSizeConfig_TF) AttributeTypes() map[string]attr.Type {
	return FontSizeConfig_TF_AttributeTypes()
}

func PatchFontSizeConfigWithPlan(
	ctx context.Context,
	d *diag.Diagnostics,
	existingModel *svixmodels.FontSizeConfig,
	planedModel FontSizeConfig_TF,
) svixmodels.FontSizeConfig {
	// initialize model as empty
	outModel := svixmodels.FontSizeConfig{}
	// load variables from the existing model
	if existingModel != nil {
		outModel.Base = existingModel.Base
	}
	// override fields in outModel with variables from planed model
	if !planedModel.Base.IsUnknown() {
		if planedModel.Base.IsNull() {
			outModel.Base = nil
		} else {
			outModel.Base = ptr(uint16(planedModel.Base.ValueInt64()))
		}
	}
	return outModel
}

// Terraform wrapper around `svixmodels.SettingsInternalIn`
type EnvironmentSettingsResourceModel struct {
	EnvironmentId               types.String          `tfsdk:"environment_id"`
	ColorPaletteDark            basetypes.ObjectValue `tfsdk:"color_palette_dark"`
	ColorPaletteLight           basetypes.ObjectValue `tfsdk:"color_palette_light"`
	CustomBaseFontSize          types.Int64           `tfsdk:"base_font_size"`
	CustomColor                 types.String          `tfsdk:"custom_color"`
	CustomFontFamily            types.String          `tfsdk:"font_family"`
	CustomFontFamilyUrl         types.String          `tfsdk:"font_family_url"`
	CustomLogoUrl               types.String          `tfsdk:"custom_logo_url"`
	CustomStringsOverride       basetypes.ObjectValue `tfsdk:"custom_strings_override"`
	CustomThemeOverride         basetypes.ObjectValue `tfsdk:"custom_theme_override"`
	DisableEndpointOnFailure    types.Bool            `tfsdk:"disable_endpoint_on_failure"`
	DisplayName                 types.String          `tfsdk:"display_name"`
	EnableChannels              types.Bool            `tfsdk:"enable_channels"`
	EnableEndpointMtlsConfig    types.Bool            `tfsdk:"enable_endpoint_mtls_config"`
	EnableEndpointOauthConfig   types.Bool            `tfsdk:"enable_endpoint_oauth_config"`
	EnableIntegrationManagement types.Bool            `tfsdk:"enable_integration_management"`
	EnableMessageStream         types.Bool            `tfsdk:"enable_message_stream"`
	EnableTransformations       types.Bool            `tfsdk:"enable_transformations"`
	EnforceHttps                types.Bool            `tfsdk:"enforce_https"`
	EventCatalogPublished       types.Bool            `tfsdk:"event_catalog_published"`
	ReadOnly                    types.Bool            `tfsdk:"read_only"`
	RequireEndpointChannel      types.Bool            `tfsdk:"require_endpoint_channel"`
	WhitelabelHeaders           types.Bool            `tfsdk:"whitelabel_headers"`
	WipeSuccessfulPayload       types.Bool            `tfsdk:"wipe_successful_payload"`
}

func PatchSettingsInternalInWithPlan(
	ctx context.Context,
	d *diag.Diagnostics,
	existingModel *svixmodels.SettingsInternalOut,
	planedModel EnvironmentSettingsResourceModel,
) svixmodels.SettingsInternalIn {
	// initialize model as empty
	outModel := svixmodels.SettingsInternalIn{}
	// load variables from the existing model
	if existingModel != nil {
		outModel.ColorPaletteDark = existingModel.ColorPaletteDark
		outModel.ColorPaletteLight = existingModel.ColorPaletteLight
		outModel.CustomBaseFontSize = existingModel.CustomBaseFontSize
		outModel.CustomColor = existingModel.CustomColor
		outModel.CustomFontFamily = existingModel.CustomFontFamily
		outModel.CustomFontFamilyUrl = existingModel.CustomFontFamilyUrl
		outModel.CustomLogoUrl = existingModel.CustomLogoUrl
		outModel.CustomStringsOverride = existingModel.CustomStringsOverride
		outModel.CustomThemeOverride = existingModel.CustomThemeOverride
		outModel.DisableEndpointOnFailure = existingModel.DisableEndpointOnFailure
		outModel.DisplayName = existingModel.DisplayName
		outModel.EnableChannels = existingModel.EnableChannels
		outModel.EnableEndpointMtlsConfig = existingModel.EnableEndpointMtlsConfig
		outModel.EnableEndpointOauthConfig = existingModel.EnableEndpointOauthConfig
		outModel.EnableIntegrationManagement = existingModel.EnableIntegrationManagement
		outModel.EnableMessageStream = existingModel.EnableMessageStream
		outModel.EnableMsgAtmptLog = existingModel.EnableMsgAtmptLog
		outModel.EnableOtlp = existingModel.EnableOtlp
		outModel.EnableTransformations = existingModel.EnableTransformations
		outModel.EnforceHttps = existingModel.EnforceHttps
		outModel.EventCatalogPublished = existingModel.EventCatalogPublished
		outModel.ReadOnly = existingModel.ReadOnly
		outModel.RequireEndpointChannel = existingModel.RequireEndpointChannel
		outModel.WhitelabelHeaders = existingModel.WhitelabelHeaders
		outModel.WipeSuccessfulPayload = existingModel.WipeSuccessfulPayload
	}
	// override fields in outModel with variables from planed model
	if !planedModel.ColorPaletteDark.IsUnknown() {
		if planedModel.ColorPaletteDark.IsNull() {
			outModel.ColorPaletteDark = nil
		} else {
			var existingColorPaletteDark CustomColorPalette_TF
			d.Append(planedModel.ColorPaletteDark.As(ctx, &existingColorPaletteDark, basetypes.ObjectAsOptions{
				UnhandledNullAsEmpty:    false,
				UnhandledUnknownAsEmpty: false,
			})...)
			outModel.ColorPaletteDark = ptr(PatchCustomColorPaletteWithPlan(ctx, d, existingModel.ColorPaletteDark, existingColorPaletteDark))
		}
	}
	if !planedModel.ColorPaletteLight.IsUnknown() {
		if planedModel.ColorPaletteLight.IsNull() {
			outModel.ColorPaletteLight = nil
		} else {
			var existingColorPaletteLight CustomColorPalette_TF
			d.Append(planedModel.ColorPaletteLight.As(ctx, &existingColorPaletteLight, basetypes.ObjectAsOptions{
				UnhandledNullAsEmpty:    false,
				UnhandledUnknownAsEmpty: false,
			})...)
			outModel.ColorPaletteLight = ptr(PatchCustomColorPaletteWithPlan(ctx, d, existingModel.ColorPaletteLight, existingColorPaletteLight))
		}
	}
	if !planedModel.CustomStringsOverride.IsUnknown() {
		if planedModel.CustomStringsOverride.IsNull() {
			outModel.CustomStringsOverride = nil
		} else {
			var existingCustomStringsOverride CustomStringsOverride_TF
			d.Append(planedModel.CustomStringsOverride.As(ctx, &existingCustomStringsOverride, basetypes.ObjectAsOptions{
				UnhandledNullAsEmpty:    false,
				UnhandledUnknownAsEmpty: false,
			})...)
			outModel.CustomStringsOverride = ptr(PatchCustomStringsOverrideWithPlan(ctx, d, existingModel.CustomStringsOverride, existingCustomStringsOverride))
		}
	}

	if !planedModel.CustomThemeOverride.IsUnknown() {
		if planedModel.CustomThemeOverride.IsNull() {
			outModel.CustomThemeOverride = nil
		} else {
			var existingCustomThemeOverride CustomThemeOverride_TF
			d.Append(planedModel.CustomThemeOverride.As(ctx, &existingCustomThemeOverride, basetypes.ObjectAsOptions{
				UnhandledNullAsEmpty:    false,
				UnhandledUnknownAsEmpty: false,
			})...)
			outModel.CustomThemeOverride = ptr(PatchCustomThemeOverrideWithPlan(ctx, d, existingModel.CustomThemeOverride, existingCustomThemeOverride))
		}
	}

	if !planedModel.CustomBaseFontSize.IsUnknown() {
		outModel.CustomBaseFontSize = planedModel.CustomBaseFontSize.ValueInt64Pointer()
	}
	if !planedModel.CustomColor.IsUnknown() {
		outModel.CustomColor = planedModel.CustomColor.ValueStringPointer()
	}
	if !planedModel.CustomFontFamily.IsUnknown() {
		outModel.CustomFontFamily = planedModel.CustomFontFamily.ValueStringPointer()
	}
	if !planedModel.CustomFontFamilyUrl.IsUnknown() {
		outModel.CustomFontFamilyUrl = planedModel.CustomFontFamilyUrl.ValueStringPointer()
	}
	if !planedModel.CustomLogoUrl.IsUnknown() {
		outModel.CustomLogoUrl = planedModel.CustomLogoUrl.ValueStringPointer()
	}
	if !planedModel.DisableEndpointOnFailure.IsUnknown() {
		outModel.DisableEndpointOnFailure = planedModel.DisableEndpointOnFailure.ValueBoolPointer()
	}
	if !planedModel.DisplayName.IsUnknown() {
		outModel.DisplayName = planedModel.DisplayName.ValueStringPointer()
	}
	if !planedModel.EnableChannels.IsUnknown() {
		outModel.EnableChannels = planedModel.EnableChannels.ValueBoolPointer()
	}
	if !planedModel.EnableEndpointMtlsConfig.IsUnknown() {
		outModel.EnableEndpointMtlsConfig = planedModel.EnableEndpointMtlsConfig.ValueBoolPointer()
	}
	if !planedModel.EnableEndpointOauthConfig.IsUnknown() {
		outModel.EnableEndpointOauthConfig = planedModel.EnableEndpointOauthConfig.ValueBoolPointer()
	}
	if !planedModel.EnableIntegrationManagement.IsUnknown() {
		outModel.EnableIntegrationManagement = planedModel.EnableIntegrationManagement.ValueBoolPointer()
	}
	if !planedModel.EnableMessageStream.IsUnknown() {
		outModel.EnableMessageStream = planedModel.EnableMessageStream.ValueBoolPointer()
	}
	if !planedModel.EnableTransformations.IsUnknown() {
		outModel.EnableTransformations = planedModel.EnableTransformations.ValueBoolPointer()
	}
	if !planedModel.EnforceHttps.IsUnknown() {
		outModel.EnforceHttps = planedModel.EnforceHttps.ValueBoolPointer()
	}
	if !planedModel.EventCatalogPublished.IsUnknown() {
		outModel.EventCatalogPublished = planedModel.EventCatalogPublished.ValueBoolPointer()
	}
	if !planedModel.ReadOnly.IsUnknown() {
		outModel.ReadOnly = planedModel.ReadOnly.ValueBoolPointer()
	}
	if !planedModel.RequireEndpointChannel.IsUnknown() {
		outModel.RequireEndpointChannel = planedModel.RequireEndpointChannel.ValueBoolPointer()
	}
	if !planedModel.WhitelabelHeaders.IsUnknown() {
		outModel.WhitelabelHeaders = planedModel.WhitelabelHeaders.ValueBoolPointer()
	}
	if !planedModel.WipeSuccessfulPayload.IsUnknown() {
		outModel.WipeSuccessfulPayload = planedModel.WipeSuccessfulPayload.ValueBoolPointer()
	}
	return outModel
}

// Terraform wrapper around `svixmodels.CustomColorPalette`
type CustomColorPalette_TF struct {
	BackgroundHover     types.String `tfsdk:"background_hover"`
	BackgroundPrimary   types.String `tfsdk:"background_primary"`
	BackgroundSecondary types.String `tfsdk:"background_secondary"`
	ButtonPrimary       types.String `tfsdk:"button_primary"`
	InteractiveAccent   types.String `tfsdk:"interactive_accent"`
	NavigationAccent    types.String `tfsdk:"navigation_accent"`
	Primary             types.String `tfsdk:"primary"`
	TextDanger          types.String `tfsdk:"text_danger"`
	TextPrimary         types.String `tfsdk:"text_primary"`
}

func CustomColorPalette_TF_AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"background_hover":     types.StringType,
		"background_primary":   types.StringType,
		"background_secondary": types.StringType,
		"button_primary":       types.StringType,
		"interactive_accent":   types.StringType,
		"navigation_accent":    types.StringType,
		"primary":              types.StringType,
		"text_danger":          types.StringType,
		"text_primary":         types.StringType,
	}
}

func (v *CustomColorPalette_TF) AttributeTypes() map[string]attr.Type {
	return CustomColorPalette_TF_AttributeTypes()

}

func PatchCustomColorPaletteWithPlan(
	ctx context.Context,
	d *diag.Diagnostics,
	existingModel *svixmodels.CustomColorPalette,
	planedModel CustomColorPalette_TF,
) svixmodels.CustomColorPalette {
	// initialize model as empty
	outModel := svixmodels.CustomColorPalette{}
	// load variables from the existing model
	if existingModel != nil {
		outModel.BackgroundHover = existingModel.BackgroundHover
		outModel.BackgroundPrimary = existingModel.BackgroundPrimary
		outModel.BackgroundSecondary = existingModel.BackgroundSecondary
		outModel.ButtonPrimary = existingModel.ButtonPrimary
		outModel.InteractiveAccent = existingModel.InteractiveAccent
		outModel.NavigationAccent = existingModel.NavigationAccent
		outModel.Primary = existingModel.Primary
		outModel.TextDanger = existingModel.TextDanger
		outModel.TextPrimary = existingModel.TextPrimary
	}
	// override fields in outModel with variables from planed model
	if !planedModel.BackgroundHover.IsUnknown() {
		outModel.BackgroundHover = planedModel.BackgroundHover.ValueStringPointer()
	}
	if !planedModel.BackgroundPrimary.IsUnknown() {
		outModel.BackgroundPrimary = planedModel.BackgroundPrimary.ValueStringPointer()
	}
	if !planedModel.BackgroundSecondary.IsUnknown() {
		outModel.BackgroundSecondary = planedModel.BackgroundSecondary.ValueStringPointer()
	}
	if !planedModel.ButtonPrimary.IsUnknown() {
		outModel.ButtonPrimary = planedModel.ButtonPrimary.ValueStringPointer()
	}
	if !planedModel.InteractiveAccent.IsUnknown() {
		outModel.InteractiveAccent = planedModel.InteractiveAccent.ValueStringPointer()
	}
	if !planedModel.NavigationAccent.IsUnknown() {
		outModel.NavigationAccent = planedModel.NavigationAccent.ValueStringPointer()
	}
	if !planedModel.Primary.IsUnknown() {
		outModel.Primary = planedModel.Primary.ValueStringPointer()
	}
	if !planedModel.TextDanger.IsUnknown() {
		outModel.TextDanger = planedModel.TextDanger.ValueStringPointer()
	}
	if !planedModel.TextPrimary.IsUnknown() {
		outModel.TextPrimary = planedModel.TextPrimary.ValueStringPointer()
	}
	return outModel
}

// Terraform wrapper around `svixmodels.CustomStringsOverride`
type CustomStringsOverride_TF struct {
	ChannelsHelp types.String `tfsdk:"channels_help"`
	ChannelsMany types.String `tfsdk:"channels_many"`
	ChannelsOne  types.String `tfsdk:"channels_one"`
}

func CustomStringsOverride_TF_AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"channels_help": types.StringType,
		"channels_many": types.StringType,
		"channels_one":  types.StringType,
	}
}

func (v *CustomStringsOverride_TF) AttributeTypes() map[string]attr.Type {
	return CustomStringsOverride_TF_AttributeTypes()
}

func PatchCustomStringsOverrideWithPlan(
	ctx context.Context,
	d *diag.Diagnostics,
	existingModel *svixmodels.CustomStringsOverride,
	planedModel CustomStringsOverride_TF,
) svixmodels.CustomStringsOverride {
	// initialize model as empty
	outModel := svixmodels.CustomStringsOverride{}
	// load variables from the existing model
	if existingModel != nil {
		outModel.ChannelsHelp = existingModel.ChannelsHelp
		outModel.ChannelsMany = existingModel.ChannelsMany
		outModel.ChannelsOne = existingModel.ChannelsOne
	}
	// override fields in outModel with variables from planed model
	if !planedModel.ChannelsHelp.IsUnknown() {
		outModel.ChannelsHelp = planedModel.ChannelsHelp.ValueStringPointer()
	}
	if !planedModel.ChannelsMany.IsUnknown() {
		outModel.ChannelsMany = planedModel.ChannelsMany.ValueStringPointer()
	}
	if !planedModel.ChannelsOne.IsUnknown() {
		outModel.ChannelsOne = planedModel.ChannelsOne.ValueStringPointer()
	}
	return outModel
}

// Terraform wrapper around `svixmodels.CustomThemeOverride`
type CustomThemeOverride_TF struct {
	BorderRadius basetypes.ObjectValue `tfsdk:"border_radius"`
}

func CustomThemeOverride_TF_AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"border_radius": basetypes.ObjectType{
			AttrTypes: map[string]attr.Type{
				"button": types.StringType,
				"card":   types.StringType,
				"input":  types.StringType,
			},
		},
	}
}

func (v *CustomThemeOverride_TF) AttributeTypes() map[string]attr.Type {
	return CustomThemeOverride_TF_AttributeTypes()
}

func PatchCustomThemeOverrideWithPlan(
	ctx context.Context,
	d *diag.Diagnostics,
	existingModel *svixmodels.CustomThemeOverride,
	planedModel CustomThemeOverride_TF,
) svixmodels.CustomThemeOverride {
	// initialize model as empty
	outModel := svixmodels.CustomThemeOverride{}
	// load variables from the existing model
	if existingModel != nil {
		outModel.BorderRadius = existingModel.BorderRadius
		outModel.FontSize = existingModel.FontSize
	}
	// override fields in outModel with variables from planed model
	if !planedModel.BorderRadius.IsUnknown() {
		if planedModel.BorderRadius.IsNull() {
			outModel.BorderRadius = nil
		} else {
			var existingBorderRadius BorderRadiusConfig_TF
			d.Append(planedModel.BorderRadius.As(ctx, &existingBorderRadius, basetypes.ObjectAsOptions{
				UnhandledNullAsEmpty:    false,
				UnhandledUnknownAsEmpty: false,
			})...)
			var existingBorderRadiusConfig *svixmodels.BorderRadiusConfig
			if existingModel != nil {
				if existingModel.BorderRadius != nil {
					existingBorderRadiusConfig = existingModel.BorderRadius
				}
			}

			outModel.BorderRadius = ptr(PatchBorderRadiusConfigWithPlan(ctx, d, existingBorderRadiusConfig, existingBorderRadius))
		}

	}
	return outModel
}

// Terraform wrapper around `svixmodels.BorderRadiusConfig`
type BorderRadiusConfig_TF struct {
	Button types.String `tfsdk:"button"`
	Card   types.String `tfsdk:"card"`
	Input  types.String `tfsdk:"input"`
}

func BorderRadiusConfig_TF_AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"button": types.StringType,
		"card":   types.StringType,
		"input":  types.StringType,
	}
}
func (v *BorderRadiusConfig_TF) AttributeTypes() map[string]attr.Type {
	return BorderRadiusConfig_TF_AttributeTypes()
}

func PatchBorderRadiusConfigWithPlan(
	ctx context.Context,
	d *diag.Diagnostics,
	existingModel *svixmodels.BorderRadiusConfig,
	planedModel BorderRadiusConfig_TF,
) svixmodels.BorderRadiusConfig {
	// initialize model as empty
	outModel := svixmodels.BorderRadiusConfig{}
	// load variables from the existing model
	if existingModel != nil {
		outModel.Button = existingModel.Button
		outModel.Card = existingModel.Card
		outModel.Input = existingModel.Input
	}
	// override fields in outModel with variables from planed model
	if !planedModel.Button.IsUnknown() {
		if planedModel.Button.IsNull() {
			outModel.Button = nil
		} else {
			outModel.Button = ptr(svixmodels.BorderRadiusEnum(planedModel.Button.ValueString()))
		}
	}
	if !planedModel.Card.IsUnknown() {
		if planedModel.Card.IsNull() {
			outModel.Card = nil
		} else {
			outModel.Card = ptr(svixmodels.BorderRadiusEnum(planedModel.Card.ValueString()))
		}
	}
	if !planedModel.Input.IsUnknown() {
		if planedModel.Input.IsNull() {
			outModel.Input = nil
		} else {
			outModel.Input = ptr(svixmodels.BorderRadiusEnum(planedModel.Input.ValueString()))
		}
	}
	return outModel
}
