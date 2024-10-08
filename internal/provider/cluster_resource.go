package provider

import (
	"context"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/komminarlabs/cratedb"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &ClusterResource{}
	_ resource.ResourceWithImportState = &ClusterResource{}
	_ resource.ResourceWithImportState = &ClusterResource{}
)

// NewClusterResource is a helper function to simplify the provider implementation.
func NewClusterResource() resource.Resource {
	return &ClusterResource{}
}

// ClusterResource defines the resource implementation.
type ClusterResource struct {
	client *cratedb.ClientWithResponses
}

// Metadata returns the resource type name.
func (r *ClusterResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cluster"
}

// Schema defines the schema for the resource.
func (r *ClusterResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Creates and manages a cluster.",

		Attributes: map[string]schema.Attribute{
			"organization_id": schema.StringAttribute{
				Required:    true,
				Description: "The organization id of the cluster.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"allow_custom_storage": schema.BoolAttribute{
				Computed:    true,
				Description: "The allow custom storage flag.",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"allow_suspend": schema.BoolAttribute{
				Computed:    true,
				Description: "The allow suspend flag.",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"backup_schedule": schema.StringAttribute{
				Computed:    true,
				Description: "The backup schedule.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"channel": schema.StringAttribute{
				Computed:    true,
				Optional:    true,
				Default:     stringdefault.StaticString("stable"),
				Description: "The channel of the cluster. Default is 'stable'.",
			},
			"crate_version": schema.StringAttribute{
				Required:    true,
				Description: "The CrateDB version of the cluster.",
			},
			"dc": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "The DublinCore of the cluster.",
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{
					"created": schema.StringAttribute{
						Computed:    true,
						Description: "The created time.",
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"modified": schema.StringAttribute{
						Computed:    true,
						Description: "The modified time.",
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
				},
			},
			"deletion_protected": schema.BoolAttribute{
				Computed:    true,
				Description: "The deletion protected flag.",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"external_ip": schema.StringAttribute{
				Computed:    true,
				Description: "The external IP address.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"fqdn": schema.StringAttribute{
				Computed:    true,
				Description: "The Fully Qualified Domain Name.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"gc_available": schema.BoolAttribute{
				Computed:    true,
				Description: "The garbage collection available flag.",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"hardware_specs": schema.SingleNestedAttribute{
				Computed:    true,
				Optional:    true,
				Description: "The hardware specs of the cluster.",
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{
					"cpus_per_node": schema.Int32Attribute{
						Computed:    true,
						Optional:    true,
						Description: "The cpus per node.",
					},
					"disk_size_per_node_bytes": schema.Int64Attribute{
						Computed:    true,
						Optional:    true,
						Description: "The disk size per node in bytes.",
					},
					"disk_type": schema.StringAttribute{
						Computed:    true,
						Optional:    true,
						Description: "The disk type.",
					},
					"disks_per_node": schema.Int32Attribute{
						Computed:    true,
						Optional:    true,
						Description: "The disks per node.",
					},
					"heap_size_bytes": schema.Int64Attribute{
						Computed:    true,
						Optional:    true,
						Description: "The heap size in bytes.",
					},
					"memory_per_node_bytes": schema.Int64Attribute{
						Computed:    true,
						Optional:    true,
						Description: "The memory per node in bytes.",
					},
				},
			},
			"health": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "The health of the cluster.",
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{
					"status": schema.StringAttribute{
						Computed:    true,
						Description: "The health status of the cluster.",
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
				},
			},
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The id of the cluster.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"ip_whitelist": schema.SetNestedAttribute{
				Computed:    true,
				Description: "The IP whitelist of the cluster.",
				Default:     setdefault.StaticValue(types.SetNull(ClusterIpWhitelistModel{}.GetAttrType())),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"cidr": schema.StringAttribute{
							Computed:    true,
							Description: "The CIDR.",
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
						},
						"description": schema.StringAttribute{
							Computed:    true,
							Description: "The description.",
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
						},
					},
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the cluster.",
			},
			"num_nodes": schema.Int32Attribute{
				Computed:    true,
				Description: "The number of nodes in the cluster.",
				PlanModifiers: []planmodifier.Int32{
					int32planmodifier.UseStateForUnknown(),
				},
			},
			"origin": schema.StringAttribute{
				Computed:    true,
				Description: "The origin of the cluster.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"product_name": schema.StringAttribute{
				Required:    true,
				Description: "The product name of the cluster.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					stringvalidator.LengthAtMost(512),
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^\w[\w\-\. ]*$`),
						"Product name must start with a letter and contain only letters, numbers, hyphens, underscores, and periods.",
					),
				},
			},
			"product_tier": schema.StringAttribute{
				Required:    true,
				Description: "The product tier of the cluster.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					stringvalidator.LengthAtMost(512),
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^\w[\w\-\. ]*$`),
						"Product name must start with a letter and contain only letters, numbers, hyphens, underscores, and periods.",
					),
				},
			},
			"product_unit": schema.Int32Attribute{
				Computed:    true,
				Optional:    true,
				Default:     int32default.StaticInt32(0),
				Description: "The product unit of the cluster. Default is `0`.",
			},
			"project_id": schema.StringAttribute{
				Required:    true,
				Description: "The project id of the cluster.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(36),
					stringvalidator.LengthAtMost(36),
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`),
						"Project ID must be a valid UUID.",
					),
				},
			},
			"subscription_id": schema.StringAttribute{
				Required:    true,
				Description: "The subscription id of the cluster.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					stringvalidator.LengthAtMost(512),
				},
			},
			"suspended": schema.BoolAttribute{
				Computed:    true,
				Description: "The suspended flag.",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"url": schema.StringAttribute{
				Computed:    true,
				Description: "The URL of the cluster.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"username": schema.StringAttribute{
				Required:    true,
				Description: "The username of the cluster.",
			},
			"password": schema.StringAttribute{
				Required:    true,
				Sensitive:   true,
				Description: "The password of the cluster.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(24),
				},
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *ClusterResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ClusterModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	productUnit := int(plan.ProductUnit.ValueInt32())
	password := plan.Password
	organizationId := plan.OrganizationId
	createPartialClusterRequest := cratedb.PartialCluster{
		Channel:      plan.Channel.ValueStringPointer(),
		CrateVersion: plan.CrateVersion.ValueString(),
		Name:         plan.Name.ValueString(),
		ProductName:  plan.ProductName.ValueString(),
		ProductTier:  plan.ProductTier.ValueString(),
		ProductUnit:  &productUnit,
		Username:     plan.Username.ValueString(),
		Password:     password.ValueStringPointer(),
	}

	createClusterRequest := cratedb.ClusterProvision{
		Cluster:        createPartialClusterRequest,
		ProjectId:      plan.ProjectId.ValueStringPointer(),
		SubscriptionId: plan.SubscriptionId.ValueString(),
	}

	createClusterResponse, err := r.client.PostApiV2OrganizationsOrganizationIdClustersWithResponse(ctx, organizationId.ValueString(), createClusterRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating cluster",
			"Could not create cluster, unexpected error: "+err.Error(),
		)
		return
	}

	if createClusterResponse.StatusCode() != 201 {
		resp.Diagnostics.AddError(
			"Error creating cluster",
			fmt.Sprintf("HTTP Status Code: %d\nStatus: %v", createClusterResponse.StatusCode(), createClusterResponse.Status()),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	clusterPlan, err := getClusterModel(ctx, *createClusterResponse.JSON201)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting cluster model",
			err.Error(),
		)
		return
	}
	plan = *clusterPlan
	plan.OrganizationId = organizationId
	plan.Password = password

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *ClusterResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state ClusterModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed cluster value from InfluxDB
	password := state.Password
	organizationId := state.OrganizationId
	readClusterResponse, err := r.client.GetApiV2ClustersClusterIdWithResponse(ctx, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting cluster",
			err.Error(),
		)
		return
	}

	if readClusterResponse.StatusCode() != 200 {
		resp.Diagnostics.AddError(
			"Error getting cluster",
			fmt.Sprintf("HTTP Status Code: %d\nStatus: %v", readClusterResponse.StatusCode(), readClusterResponse.Status()),
		)
		return
	}

	// Map response body to model
	clusterState, err := getClusterModel(ctx, *readClusterResponse.JSON200)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting cluster model",
			err.Error(),
		)
		return
	}
	// Overwrite items with refreshed state
	state = *clusterState
	state.OrganizationId = organizationId
	state.Password = password

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *ClusterResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ClusterModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	password := plan.Password
	organizationId := plan.OrganizationId
	updateClusterRequest := cratedb.ClusterEdit{
		Password: plan.Password.ValueStringPointer(),
	}

	// Update existing cluster
	updateClusterResponse, err := r.client.PatchApiV2ClustersClusterIdWithResponse(ctx, plan.Id.ValueString(), updateClusterRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating cluster",
			"Could not update cluster, unexpected error: "+err.Error(),
		)
		return
	}

	if updateClusterResponse.StatusCode() != 200 {
		resp.Diagnostics.AddError(
			"Error updating cluster",
			fmt.Sprintf("HTTP Status Code: %d\nStatus: %v", updateClusterResponse.StatusCode(), updateClusterResponse.Status()),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	clusterPlan, err := getClusterModel(ctx, *updateClusterResponse.JSON200)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting cluster model",
			err.Error(),
		)
		return
	}
	plan = *clusterPlan
	plan.OrganizationId = organizationId
	plan.Password = password

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *ClusterResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ClusterModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing cluster
	deleteClustersResponse, err := r.client.DeleteApiV2ClustersClusterIdWithResponse(ctx, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting cluster",
			"Could not delete cluster, unexpected error: "+err.Error(),
		)
		return
	}

	if deleteClustersResponse.StatusCode() != 204 {
		resp.Diagnostics.AddError(
			"Error deleting cluster",
			fmt.Sprintf("HTTP Status Code: %d\nStatus: %v", deleteClustersResponse.StatusCode(), deleteClustersResponse.Status()),
		)
		return
	}
}

// Configure adds the provider configured client to the resource.
func (r *ClusterResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ClusterResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
