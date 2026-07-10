package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/thulasirajkomminar/cratedb-cloud-go"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &ProjectResource{}
	_ resource.ResourceWithConfigure   = &ProjectResource{}
	_ resource.ResourceWithImportState = &ProjectResource{}
)

// NewProjectResource is a helper function to simplify the provider implementation.
func NewProjectResource() resource.Resource {
	return &ProjectResource{}
}

// ProjectResource defines the resource implementation.
type ProjectResource struct {
	client *cratedb.ClientWithResponses
}

// Metadata returns the resource type name.
func (r *ProjectResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_project"
}

// Schema defines the schema for the resource.
func (r *ProjectResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Creates and manages a project.",

		Attributes: map[string]schema.Attribute{
			"dc": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "The DublinCore of the project.",
				Attributes: map[string]schema.Attribute{
					"created": schema.StringAttribute{
						CustomType:  timetypes.RFC3339Type{},
						Computed:    true,
						Description: "The created time.",
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"modified": schema.StringAttribute{
						CustomType:  timetypes.RFC3339Type{},
						Computed:    true,
						Description: "The modified time.",
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
				},
			},
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The id of the project.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the project.",
			},
			"organization_id": schema.StringAttribute{
				Required:    true,
				Description: "The organization id of the project.",
			},
			"region": schema.StringAttribute{
				Required:    true,
				Description: "The region of the project.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *ProjectResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ProjectModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	createProjectRequest := cratedb.Project{
		Name:           plan.Name.ValueString(),
		OrganizationId: plan.OrganizationId.ValueString(),
		Region:         plan.Region.ValueStringPointer(),
	}

	createProjectResponse, err := r.client.PostApiV2ProjectsWithResponse(ctx, createProjectRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating project",
			"Could not create project, unexpected error: "+err.Error(),
		)
		return
	}

	if createProjectResponse.StatusCode() != 201 || createProjectResponse.JSON201 == nil {
		resp.Diagnostics.AddError(
			"Error creating project",
			apiErrorDetail(createProjectResponse.HTTPResponse, createProjectResponse.Body),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	projectPlan, err := getProjectModel(ctx, *createProjectResponse.JSON201)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting project model",
			err.Error(),
		)
		return
	}
	plan = *projectPlan

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *ProjectResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state ProjectModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed project value from API
	readProjectResponse, err := r.client.GetApiV2ProjectsProjectIdWithResponse(ctx, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting project",
			err.Error(),
		)
		return
	}

	// If the project no longer exists, remove it from state so Terraform
	// plans a re-create instead of failing the refresh.
	if isNotFound(readProjectResponse.HTTPResponse) {
		resp.State.RemoveResource(ctx)
		return
	}

	if readProjectResponse.StatusCode() != 200 || readProjectResponse.JSON200 == nil {
		resp.Diagnostics.AddError(
			"Error getting project",
			apiErrorDetail(readProjectResponse.HTTPResponse, readProjectResponse.Body),
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
	// Overwrite items with refreshed state
	state = *projectState

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *ProjectResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ProjectModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	updateProjectRequest := cratedb.ProjectEdit{
		Name: plan.Name.ValueString(),
	}

	// Update existing project
	updateProjectResponse, err := r.client.PatchApiV2ProjectsProjectIdWithResponse(ctx, plan.Id.ValueString(), updateProjectRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating project",
			"Could not update project, unexpected error: "+err.Error(),
		)
		return
	}

	if updateProjectResponse.StatusCode() != 200 || updateProjectResponse.JSON200 == nil {
		resp.Diagnostics.AddError(
			"Error updating project",
			apiErrorDetail(updateProjectResponse.HTTPResponse, updateProjectResponse.Body),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	projectPlan, err := getProjectModel(ctx, *updateProjectResponse.JSON200)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting project model",
			err.Error(),
		)
		return
	}
	plan = *projectPlan

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *ProjectResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ProjectModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing project
	deleteProjectsResponse, err := r.client.DeleteApiV2ProjectsProjectIdWithResponse(ctx, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting project",
			"Could not delete project, unexpected error: "+err.Error(),
		)
		return
	}

	// A project already deleted out-of-band is fine: the desired outcome
	// (no project) is achieved.
	if isNotFound(deleteProjectsResponse.HTTPResponse) {
		return
	}

	if deleteProjectsResponse.StatusCode() != 204 {
		resp.Diagnostics.AddError(
			"Error deleting project",
			apiErrorDetail(deleteProjectsResponse.HTTPResponse, deleteProjectsResponse.Body),
		)
		return
	}
}

// Configure adds the provider configured client to the resource.
func (r *ProjectResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if client := clientFromProviderData(req.ProviderData, "Resource", &resp.Diagnostics); client != nil {
		r.client = client
	}
}

func (r *ProjectResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Read refreshes the project by its id, so the import identifier is the
	// project id.
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
