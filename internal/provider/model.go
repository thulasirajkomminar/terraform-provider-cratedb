package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// DCModel maps CrateDB DC schema data.
type DCModel struct {
	Created  types.String `tfsdk:"created"`
	Modified types.String `tfsdk:"modified"`
}

func (o DCModel) GetAttrType() map[string]attr.Type {
	return map[string]attr.Type{
		"created":  types.StringType,
		"modified": types.StringType,
	}
}
