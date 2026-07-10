package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/thulasirajkomminar/cratedb-cloud-go"
)

// OrganizationModel maps CrateDB organization schema data.
type OrganizationModel struct {
	Dc                   types.Object `tfsdk:"dc"`
	Email                types.String `tfsdk:"email"`
	Id                   types.String `tfsdk:"id"`
	Name                 types.String `tfsdk:"name"`
	NotificationsEnabled types.Bool   `tfsdk:"notifications_enabled"`
	PlanType             types.Int32  `tfsdk:"plan_type"`
	ProjectCount         types.Int32  `tfsdk:"project_count"`
	RoleFQN              types.String `tfsdk:"role_fqn"`
}

func getOrganizationModel(ctx context.Context, organization cratedb.Organization) (*OrganizationModel, error) {
	dcObjectValue, err := getDCObjectValue(ctx, organization.Dc)
	if err != nil {
		return nil, fmt.Errorf("error getting organization DC value: %w", err)
	}

	email := types.StringNull()
	if organization.Email != nil {
		email = types.StringValue(string(*organization.Email))
	}

	planType := types.Int32Null()
	if organization.PlanType != nil {
		planType = types.Int32Value(int32(*organization.PlanType))
	}

	roleFQN := types.StringNull()
	if organization.RoleFqn != nil {
		roleFQN = types.StringValue(string(*organization.RoleFqn))
	}

	return &OrganizationModel{
		Dc:                   dcObjectValue,
		Email:                email,
		Id:                   types.StringPointerValue(organization.Id),
		Name:                 types.StringValue(organization.Name),
		NotificationsEnabled: types.BoolPointerValue(organization.NotificationsEnabled),
		PlanType:             planType,
		ProjectCount:         intPointerToInt32Value(organization.ProjectCount),
		RoleFQN:              roleFQN,
	}, nil
}
