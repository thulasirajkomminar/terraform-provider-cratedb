package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/komminarlabs/cratedb"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &ProjectDataSource{}
	_ datasource.DataSourceWithConfigure = &ProjectDataSource{}
)

// NewProjectDataSource is a helper function to simplify the provider implementation.
func NewProjectDataSource() datasource.DataSource {
	return &ProjectDataSource{}
}

// ProjectDataSource is the data source implementation.
type ProjectDataSource struct {
	client *cratedb.ClientWithResponses
}

// Metadata returns the data source type name.
func (d *ProjectDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_project"
}

// Schema defines the schema for the data source.
func (d *ProjectDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "To retrieve a project.",

		Attributes: map[string]schema.Attribute{
			"dc": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "The DublinCore of the project.",
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
			"id": schema.StringAttribute{
				Required:    true,
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
	}
}

// Configure adds the provider configured client to the data source.
func (d *ProjectDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *ProjectDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state ProjectModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	readProjectResponse, err := d.client.GetApiV2ProjectsProjectIdWithResponse(ctx, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting project",
			err.Error(),
		)
		return
	}

	if readProjectResponse.StatusCode() != 200 {
		resp.Diagnostics.AddError(
			"Error getting project",
			fmt.Sprintf("HTTP Status Code: %d\nStatus: %v", readProjectResponse.StatusCode(), readProjectResponse.Status()),
		)
		return
	}

	// Map response body to model
	projectState, err := getProjectModel(ctx, *readProjectResponse.JSON200)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting project model",
			err.Error(),
		)
		return
	}
	state = *projectState

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
