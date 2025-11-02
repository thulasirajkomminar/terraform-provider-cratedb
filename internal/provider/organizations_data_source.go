package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/thulasirajkomminar/cratedb-cloud-go"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &OrganizationsDataSource{}
	_ datasource.DataSourceWithConfigure = &OrganizationsDataSource{}
)

// NewOrganizationsDataSource is a helper function to simplify the provider implementation.
func NewOrganizationsDataSource() datasource.DataSource {
	return &OrganizationsDataSource{}
}

// OrganizationsDataSource is the data source implementation.
type OrganizationsDataSource struct {
	client *cratedb.ClientWithResponses
}

// OrganizationsDataSourceModel describes the data source data model.
type OrganizationsDataSourceModel struct {
	Organizations []OrganizationModel `tfsdk:"organizations"`
}

// Metadata returns the data source type name.
func (d *OrganizationsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organizations"
}

// Schema defines the schema for the data source.
func (d *OrganizationsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "To retrieve all organizations.",

		Attributes: map[string]schema.Attribute{
			"organizations": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"email": schema.StringAttribute{
							Computed:    true,
							Description: "The notification email used in the organization.",
						},
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "The id of the organization.",
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "The name of the organization.",
						},
						"notifications_enabled": schema.BoolAttribute{
							Computed:    true,
							Description: "Whether notifications enabled for the organization.",
						},
						"plan_type": schema.Int32Attribute{
							Computed:    true,
							Description: "The support plan type used in the organization.",
						},
						"project_count": schema.Int32Attribute{
							Computed:    true,
							Description: "The project count in the organization.",
						},
						"role_fqn": schema.StringAttribute{
							Computed:    true,
							Description: "The role FQN.",
						},
						"dc": schema.SingleNestedAttribute{
							Computed:    true,
							Description: "The DublinCore of the organization.",
							Attributes: map[string]schema.Attribute{
								"created": schema.StringAttribute{
									Computed:    true,
									Description: "The created time.",
								},
								"modified": schema.StringAttribute{
									Computed:    true,
									Description: "The modified time.",
								},
							},
						},
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *OrganizationsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*cratedb.ClientWithResponses)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected cratedb.ClientWithResponses, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}
	d.client = client
}

// Read refreshes the Terraform state with the latest data.
func (d *OrganizationsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state OrganizationsDataSourceModel

	readOrganizationsResponse, err := d.client.GetApiV2OrganizationsWithResponse(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting organizations",
			err.Error(),
		)
		return
	}

	if readOrganizationsResponse.StatusCode() != 200 {
		resp.Diagnostics.AddError(
			"Error formatting error response",
			fmt.Sprintf("HTTP Status Code: %d\nStatus: %v", readOrganizationsResponse.StatusCode(), readOrganizationsResponse.Status()),
		)
		return
	}

	// Map response body to model
	for _, organization := range *readOrganizationsResponse.JSON200 {
		organizationState, err := getOrganizationModel(ctx, organization)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error getting organization model",
				err.Error(),
			)
			return
		}
		state.Organizations = append(state.Organizations, *organizationState)
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
