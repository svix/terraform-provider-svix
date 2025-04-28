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
	CustomStringsOverride       basetypes.ObjectValue `tfsdk:"channels_strings_override"`
	DisableEndpointOnFailure    types.Bool            `tfsdk:"disable_endpoint_on_failure"`
	EnableChannels              types.Bool            `tfsdk:"enable_channels"`
	EnableEndpointMtlsConfig    types.Bool            `tfsdk:"enable_endpoint_mtls_config"`
	EnableEndpointOauthConfig   types.Bool            `tfsdk:"enable_endpoint_oauth_config"`
	EnableIntegrationManagement types.Bool            `tfsdk:"enable_integration_management"`
	EnableMessageStream         types.Bool            `tfsdk:"enable_advanced_endpoint_types"`
	EnableTransformations       types.Bool            `tfsdk:"enable_transformations"`
	EnforceHttps                types.Bool            `tfsdk:"enforce_https"`
	EventCatalogPublished       types.Bool            `tfsdk:"event_catalog_published"`
	ReadOnly                    types.Bool            `tfsdk:"read_only"`
	RequireEndpointChannel      types.Bool            `tfsdk:"require_endpoint_channel"`
	WhitelabelHeaders           types.Bool            `tfsdk:"whitelabel_headers"`
	WipeSuccessfulPayload       types.Bool            `tfsdk:"wipe_successful_payload"`

	WhitelabelSettings basetypes.ObjectValue `tfsdk:"whitelabel_settings"`
}

type WhitelabelSettings struct {
	DisplayName         types.String `tfsdk:"display_name"`
	CustomBaseFontSize  types.Int64  `tfsdk:"base_font_size"`
	CustomFontFamily    types.String `tfsdk:"font_family"`
	CustomFontFamilyUrl types.String `tfsdk:"font_family_url"`
	CustomLogoUrl       types.String `tfsdk:"logo_url"`

	BorderRadius      basetypes.ObjectValue `tfsdk:"border_radius"`
	ColorPaletteDark  basetypes.ObjectValue `tfsdk:"color_palette_dark"`
	ColorPaletteLight basetypes.ObjectValue `tfsdk:"color_palette_light"`
}

func WhitelabelSettings_TF_AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"display_name":    types.StringType,
		"base_font_size":  types.Int64Type,
		"font_family":     types.StringType,
		"font_family_url": types.StringType,
		"logo_url":        types.StringType,

		"color_palette_dark": basetypes.ObjectType{
			AttrTypes: CustomColorPalette_TF_AttributeTypes(),
		},
		"color_palette_light": basetypes.ObjectType{
			AttrTypes: CustomColorPalette_TF_AttributeTypes(),
		},

		"border_radius": basetypes.ObjectType{
			AttrTypes: map[string]attr.Type{
				"button": types.StringType,
				"card":   types.StringType,
				"input":  types.StringType,
			},
		},
	}
}

func (v *WhitelabelSettings) AttributeTypes() map[string]attr.Type {
	return WhitelabelSettings_TF_AttributeTypes()
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

	if !planedModel.WhitelabelSettings.IsUnknown() && !planedModel.WhitelabelSettings.IsNull() {
		var planedWhitelabelSettings WhitelabelSettings
		d.Append(planedModel.WhitelabelSettings.As(ctx, &planedWhitelabelSettings, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
		if !planedWhitelabelSettings.DisplayName.IsUnknown() {
			outModel.DisplayName = planedWhitelabelSettings.DisplayName.ValueStringPointer()
		}
		if !planedWhitelabelSettings.CustomBaseFontSize.IsUnknown() {
			outModel.CustomBaseFontSize = planedWhitelabelSettings.CustomBaseFontSize.ValueInt64Pointer()
		}
		if !planedWhitelabelSettings.CustomFontFamily.IsUnknown() {
			outModel.CustomFontFamily = planedWhitelabelSettings.CustomFontFamily.ValueStringPointer()
		}
		if !planedWhitelabelSettings.CustomFontFamilyUrl.IsUnknown() {
			outModel.CustomFontFamilyUrl = planedWhitelabelSettings.CustomFontFamilyUrl.ValueStringPointer()
		}
		if !planedWhitelabelSettings.CustomLogoUrl.IsUnknown() {
			outModel.CustomLogoUrl = planedWhitelabelSettings.CustomLogoUrl.ValueStringPointer()
		}

		{
			var planedBorderRadius BorderRadius
			d.Append(planedWhitelabelSettings.BorderRadius.As(ctx, &planedBorderRadius, basetypes.ObjectAsOptions{
				UnhandledNullAsEmpty:    false,
				UnhandledUnknownAsEmpty: false,
			})...)
			if !planedBorderRadius.Button.IsUnknown() || !planedBorderRadius.Card.IsUnknown() || !planedBorderRadius.Input.IsUnknown() {
				borderRadiusOut := svixmodels.BorderRadiusConfig{
					Button: nil,
					Card:   nil,
					Input:  nil,
				}
				if !planedBorderRadius.Button.IsUnknown() && !planedBorderRadius.Button.IsNull() {
					borderRadiusOut.Button = ptr(svixmodels.BorderRadiusEnumFromString[planedBorderRadius.Button.ValueString()])
				}
				if !planedBorderRadius.Card.IsUnknown() && !planedBorderRadius.Card.IsNull() {
					borderRadiusOut.Card = ptr(svixmodels.BorderRadiusEnumFromString[planedBorderRadius.Card.ValueString()])
				}
				if !planedBorderRadius.Input.IsUnknown() && !planedBorderRadius.Input.IsNull() {
					borderRadiusOut.Input = ptr(svixmodels.BorderRadiusEnumFromString[planedBorderRadius.Input.ValueString()])
				}
				outModel.CustomThemeOverride = &svixmodels.CustomThemeOverride{
					BorderRadius: &borderRadiusOut,
				}
			}

		}

		{
			var planedColorPaletteDark CustomColorPalette_TF
			d.Append(planedWhitelabelSettings.ColorPaletteDark.As(ctx, &planedColorPaletteDark, basetypes.ObjectAsOptions{
				UnhandledNullAsEmpty:    false,
				UnhandledUnknownAsEmpty: false,
			})...)
			colorPaletteOut := patchColorPaletteWithPlan(ctx, d, planedWhitelabelSettings)
			outModel.ColorPaletteDark = colorPaletteOut
		}

		{
			var planedColorPaletteLight CustomColorPalette_TF
			d.Append(planedWhitelabelSettings.ColorPaletteLight.As(ctx, &planedColorPaletteLight, basetypes.ObjectAsOptions{
				UnhandledNullAsEmpty:    false,
				UnhandledUnknownAsEmpty: false,
			})...)
			colorPaletteOut := patchColorPaletteWithPlan(ctx, d, planedWhitelabelSettings)
			outModel.ColorPaletteLight = colorPaletteOut
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

	if !planedModel.DisableEndpointOnFailure.IsUnknown() {
		outModel.DisableEndpointOnFailure = planedModel.DisableEndpointOnFailure.ValueBoolPointer()
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
	Primary             types.String `tfsdk:"primary"`
	BackgroundPrimary   types.String `tfsdk:"background"`
	BackgroundSecondary types.String `tfsdk:"surface_background"`
	BackgroundHover     types.String `tfsdk:"surface_hover"`
	InteractiveAccent   types.String `tfsdk:"interactive_accent"`
	NavigationAccent    types.String `tfsdk:"navigation_accent"`
	ButtonPrimary       types.String `tfsdk:"button_primary"`
	TextPrimary         types.String `tfsdk:"text_primary"`
	TextDanger          types.String `tfsdk:"text_danger"`
}

func CustomColorPalette_TF_AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"primary":            types.StringType,
		"background":         types.StringType,
		"surface_background": types.StringType,
		"surface_hover":      types.StringType,
		"interactive_accent": types.StringType,
		"navigation_accent":  types.StringType,
		"button_primary":     types.StringType,
		"text_primary":       types.StringType,
		"text_danger":        types.StringType,
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

type BorderRadius struct {
	Button types.String `tfsdk:"button"`
	Card   types.String `tfsdk:"card"`
	Input  types.String `tfsdk:"input"`
}

func BorderRadius_AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"button": types.StringType,
		"card":   types.StringType,
		"input":  types.StringType,
	}
}
func (v *BorderRadius) AttributeTypes() map[string]attr.Type {
	return BorderRadius_AttributeTypes()
}

func PatchBorderRadiusConfigWithPlan(
	ctx context.Context,
	d *diag.Diagnostics,
	existingModel *svixmodels.BorderRadiusConfig,
	planedModel BorderRadius,
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

// // if unknown return nil, else return value
// func strOrNil(v types.String) *string {
// 	if v.IsUnknown() {
// 		return nil
// 	}
// 	return v.ValueStringPointer()
// }

func patchColorPaletteWithPlan(ctx context.Context, d *diag.Diagnostics, planedWhitelabelSettings WhitelabelSettings) *svixmodels.CustomColorPalette {
	var planedColorPalette CustomColorPalette_TF
	d.Append(planedWhitelabelSettings.ColorPaletteLight.As(ctx, &planedColorPalette, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    false,
		UnhandledUnknownAsEmpty: false,
	})...)
	if !planedColorPalette.Primary.IsUnknown() ||
		!planedColorPalette.BackgroundPrimary.IsUnknown() ||
		!planedColorPalette.BackgroundSecondary.IsUnknown() ||
		!planedColorPalette.BackgroundHover.IsUnknown() ||
		!planedColorPalette.InteractiveAccent.IsUnknown() ||
		!planedColorPalette.NavigationAccent.IsUnknown() ||
		!planedColorPalette.ButtonPrimary.IsUnknown() ||
		!planedColorPalette.TextPrimary.IsUnknown() ||
		!planedColorPalette.TextDanger.IsUnknown() {
		colorPaletteOut := svixmodels.CustomColorPalette{}
		if !planedColorPalette.Primary.IsUnknown() {
			colorPaletteOut.Primary = planedColorPalette.Primary.ValueStringPointer()
		}
		if !planedColorPalette.BackgroundPrimary.IsUnknown() {
			colorPaletteOut.BackgroundPrimary = planedColorPalette.BackgroundPrimary.ValueStringPointer()
		}
		if !planedColorPalette.BackgroundSecondary.IsUnknown() {
			colorPaletteOut.BackgroundSecondary = planedColorPalette.BackgroundSecondary.ValueStringPointer()
		}
		if !planedColorPalette.BackgroundHover.IsUnknown() {
			colorPaletteOut.BackgroundHover = planedColorPalette.BackgroundHover.ValueStringPointer()
		}
		if !planedColorPalette.InteractiveAccent.IsUnknown() {
			colorPaletteOut.InteractiveAccent = planedColorPalette.InteractiveAccent.ValueStringPointer()
		}
		if !planedColorPalette.NavigationAccent.IsUnknown() {
			colorPaletteOut.NavigationAccent = planedColorPalette.NavigationAccent.ValueStringPointer()
		}
		if !planedColorPalette.ButtonPrimary.IsUnknown() {
			colorPaletteOut.ButtonPrimary = planedColorPalette.ButtonPrimary.ValueStringPointer()
		}
		if !planedColorPalette.TextPrimary.IsUnknown() {
			colorPaletteOut.TextPrimary = planedColorPalette.TextPrimary.ValueStringPointer()
		}
		if !planedColorPalette.TextDanger.IsUnknown() {
			colorPaletteOut.TextDanger = planedColorPalette.TextDanger.ValueStringPointer()
		}
		return &colorPaletteOut
	} else {
		return nil
	}

}
