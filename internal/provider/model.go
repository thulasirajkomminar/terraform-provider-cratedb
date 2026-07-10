package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/thulasirajkomminar/cratedb-cloud-go"
)

// DCModel maps CrateDB DC schema data.
type DCModel struct {
	Created  timetypes.RFC3339 `tfsdk:"created"`
	Modified timetypes.RFC3339 `tfsdk:"modified"`
}

func (o DCModel) GetAttrType() map[string]attr.Type {
	return map[string]attr.Type{
		"created":  timetypes.RFC3339Type{},
		"modified": timetypes.RFC3339Type{},
	}
}

// getDCObjectValue maps the API DublinCore value to a Terraform object,
// assigning null unconditionally when the API returns nothing so cleared
// values do not linger in state.
func getDCObjectValue(ctx context.Context, dc *cratedb.DublinCore) (types.Object, error) {
	dcValue := DCModel{
		Created:  timetypes.NewRFC3339Null(),
		Modified: timetypes.NewRFC3339Null(),
	}
	if dc != nil {
		dcValue.Created = timetypes.NewRFC3339TimePointerValue(dc.Created)
		dcValue.Modified = timetypes.NewRFC3339TimePointerValue(dc.Modified)
	}

	dcObjectValue, diags := types.ObjectValueFrom(ctx, dcValue.GetAttrType(), dcValue)
	if diags.HasError() {
		return types.ObjectNull(dcValue.GetAttrType()), fmt.Errorf("error converting DublinCore value: %v", diags.Errors())
	}
	return dcObjectValue, nil
}

// intPointerToInt32Value converts an optional API integer to an Int32 value,
// preserving null when the API omits the field.
func intPointerToInt32Value(v *int) types.Int32 {
	if v == nil {
		return types.Int32Null()
	}
	return types.Int32Value(int32(*v))
}

// intPointerToInt64Value converts an optional API integer to an Int64 value,
// preserving null when the API omits the field.
func intPointerToInt64Value(v *int) types.Int64 {
	if v == nil {
		return types.Int64Null()
	}
	return types.Int64Value(int64(*v))
}
