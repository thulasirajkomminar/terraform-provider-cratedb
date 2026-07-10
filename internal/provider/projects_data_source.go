package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/thulasirajkomminar/cratedb-cloud-go"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &ProjectsDataSource{}
	_ datasource.DataSourceWithConfigure = &ProjectsDataSource{}
)

// NewProjectsDataSource is a helper function to simplify the provider implementation.
func NewProjectsDataSource() datasource.DataSource {
	return &ProjectsDataSource{}
}

// ProjectsDataSource is the data source implementation.
type ProjectsDataSource struct {
	client *cratedb.ClientWithResponses
}

// ProjectsDataSourceModel describes the data source data model.
type ProjectsDataSourceModel struct {
	Projects []ProjectModel `tfsdk:"projects"`
}

// Metadata returns the data source type name.
func (d *ProjectsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_projects"
}

// Schema defines the schema for the data source.
func (d *ProjectsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "To retrieve all projects.",

		Attributes: map[string]schema.Attribute{
			"projects": schema.ListNestedAttribute{
				Computed:    true,
				Description: "The list of projects.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"dc": schema.SingleNestedAttribute{
							Computed:    true,
							Description: "The DublinCore of the project.",
							Attributes: map[string]schema.Attribute{
								"created": schema.StringAttribute{
									CustomType:  timetypes.RFC3339Type{},
									Computed:    true,
									Description: "The created time.",
								},
								"modified": schema.StringAttribute{
									CustomType:  timetypes.RFC3339Type{},
									Computed:    true,
									Description: "The modified time.",
								},
							},
						},
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "The id of the project.",
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "The name of the project.",
						},
						"organization_id": schema.StringAttribute{
							Computed:    true,
							Description: "The organization id of the project.",
						},
						"region": schema.StringAttribute{
							Computed:    true,
							Description: "The region of the project.",
						},
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *ProjectsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if client := clientFromProviderData(req.ProviderData, "Data Source", &resp.Diagnostics); client != nil {
		d.client = client
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *ProjectsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state ProjectsDataSourceModel

	readProjectsResponse, err := d.client.GetApiV2ProjectsWithResponse(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting projects",
			err.Error(),
		)
		return
	}

	if readProjectsResponse.StatusCode() != 200 || readProjectsResponse.JSON200 == nil {
		resp.Diagnostics.AddError(
			"Error getting projects",
			apiErrorDetail(readProjectsResponse.HTTPResponse, readProjectsResponse.Body),
		)
		return
	}

	for _, project := range *readProjectsResponse.JSON200 {
		projectState, err := getProjectModel(ctx, project)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error getting project model",
				err.Error(),
			)
			return
		}
		state.Projects = append(state.Projects, *projectState)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
