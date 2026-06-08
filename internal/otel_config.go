package internal

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type OtelConfig_TF struct {
	Url               types.String `tfsdk:"url"`
	AdditionalHeaders types.Map    `tfsdk:"additional_headers"`
}

func OtelConfig_TF_AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"url":                types.StringType,
		"additional_headers": types.MapType{ElemType: types.StringType},
	}
}

func (v *OtelConfig_TF) AttributeTypes() map[string]attr.Type {
	return OtelConfig_TF_AttributeTypes()
}
