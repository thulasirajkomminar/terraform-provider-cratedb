package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/komminarlabs/cratedb"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &OrganizationResource{}
	_ resource.ResourceWithImportState = &OrganizationResource{}
	_ resource.ResourceWithImportState = &OrganizationResource{}
)

// NewOrganizationResource is a helper function to simplify the provider implementation.
func NewOrganizationResource() resource.Resource {
	return &OrganizationResource{}
}

// OrganizationResource defines the resource implementation.
type OrganizationResource struct {
	client *cratedb.ClientWithResponses
}

// Metadata returns the resource type name.
func (r *OrganizationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organization"
}

// Schema defines the schema for the resource.
func (r *OrganizationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Creates and manages an organization.",

		Attributes: map[string]schema.Attribute{
			"email": schema.StringAttribute{
				Computed:    true,
				Description: "The notification email used in the organization.",
			},
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The id of the organization.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
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
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *OrganizationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan OrganizationModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	createOrganizationRequest := cratedb.Organization{
		Name:                 plan.Name.ValueString(),
		NotificationsEnabled: plan.NotificationsEnabled.ValueBoolPointer(),
	}

	createOrganizationResponse, err := r.client.PostApiV2OrganizationsWithResponse(ctx, createOrganizationRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating organization",
			"Could not create organization, unexpected error: "+err.Error(),
		)
		return
	}

	if createOrganizationResponse.StatusCode() != 201 {
		resp.Diagnostics.AddError(
			"Error getting organization",
			fmt.Sprintf("HTTP Status Code: %d\nStatus: %v", createOrganizationResponse.StatusCode(), createOrganizationResponse.Status()),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	organizationPlan, err := getOrganizationModel(ctx, *createOrganizationResponse.JSON201)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting organization model",
			err.Error(),
		)
		return
	}
	plan = *organizationPlan

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *OrganizationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state OrganizationModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed organization value from InfluxDB
	readOrganizationResponse, err := r.client.GetApiV2OrganizationsOrganizationIdWithResponse(ctx, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting organization",
			err.Error(),
		)
		return
	}

	if readOrganizationResponse.StatusCode() != 200 {
		resp.Diagnostics.AddError(
			"Error getting organization",
			fmt.Sprintf("HTTP Status Code: %d\nStatus: %v", readOrganizationResponse.StatusCode(), readOrganizationResponse.Status()),
		)
		return
	}

	// Map response body to model
	organizationState, err := getOrganizationModel(ctx, *readOrganizationResponse.JSON200)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting organization model",
			err.Error(),
		)
		return
	}
	// Overwrite items with refreshed state
	state = *organizationState

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *OrganizationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan OrganizationModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	updateOrganizationRequest := cratedb.OrganizationEdit{
		Id:   plan.Id.ValueStringPointer(),
		Name: plan.Name.ValueStringPointer(),
	}

	// Update existing organization
	updateOrganizationResponse, err := r.client.PutApiV2OrganizationsOrganizationIdWithResponse(ctx, plan.Id.ValueString(), updateOrganizationRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating organization",
			"Could not update organization, unexpected error: "+err.Error(),
		)
		return
	}

	if updateOrganizationResponse.StatusCode() != 200 {
		resp.Diagnostics.AddError(
			"Error updating organization",
			fmt.Sprintf("HTTP Status Code: %d\nStatus: %v", updateOrganizationResponse.StatusCode(), updateOrganizationResponse.Status()),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	organizationPlan, err := getOrganizationModel(ctx, *updateOrganizationResponse.JSON200)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting organization model",
			err.Error(),
		)
		return
	}
	plan = *organizationPlan

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *OrganizationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state OrganizationModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing organization
	deleteOrganizationsResponse, err := r.client.DeleteApiV2OrganizationsOrganizationIdWithResponse(ctx, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting organization",
			"Could not delete organization, unexpected error: "+err.Error(),
		)
		return
	}

	if deleteOrganizationsResponse.StatusCode() != 204 {
		resp.Diagnostics.AddError(
			"Error deleting organization",
			fmt.Sprintf("HTTP Status Code: %d\nStatus: %v", deleteOrganizationsResponse.StatusCode(), deleteOrganizationsResponse.Status()),
		)
		return
	}
}

// Configure adds the provider configured client to the resource.
func (r *OrganizationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
	r.client = client
}

func (r *OrganizationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
