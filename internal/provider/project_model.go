package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/komminarlabs/cratedb"
)

// ProjectModel maps CrateDB project schema data.
type ProjectModel struct {
	Dc             types.Object `tfsdk:"dc"`
	Id             types.String `tfsdk:"id"`
	Name           types.String `tfsdk:"name"`
	OrganizationId types.String `tfsdk:"organization_id"`
	Region         types.String `tfsdk:"region"`
}

func getProjectModel(ctx context.Context, project cratedb.Project) (*ProjectModel, error) {
	dcValue := DCModel{
		Created:  types.StringValue(project.Dc.Created.String()),
		Modified: types.StringValue(project.Dc.Modified.String()),
	}

	dcObjectValue, diags := types.ObjectValueFrom(ctx, dcValue.GetAttrType(), dcValue)
	if diags.HasError() {
		return nil, fmt.Errorf("error getting organization DC value")
	}

	return &ProjectModel{
		Dc:             dcObjectValue,
		Id:             types.StringPointerValue(project.Id),
		Name:           types.StringValue(project.Name),
		OrganizationId: types.StringValue(project.OrganizationId),
		Region:         types.StringPointerValue(project.Region),
	}, nil
}
