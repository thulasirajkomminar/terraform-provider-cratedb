package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/thulasirajkomminar/cratedb-cloud-go"
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
	dcObjectValue, err := getDCObjectValue(ctx, project.Dc)
	if err != nil {
		return nil, fmt.Errorf("error getting project DC value: %w", err)
	}

	return &ProjectModel{
		Dc:             dcObjectValue,
		Id:             types.StringPointerValue(project.Id),
		Name:           types.StringValue(project.Name),
		OrganizationId: types.StringValue(project.OrganizationId),
		Region:         types.StringPointerValue(project.Region),
	}, nil
}
