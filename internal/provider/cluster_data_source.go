package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/thulasirajkomminar/cratedb"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &ClusterDataSource{}
	_ datasource.DataSourceWithConfigure = &ClusterDataSource{}
)

// NewClusterDataSource is a helper function to simplify the provider implementation.
func NewClusterDataSource() datasource.DataSource {
	return &ClusterDataSource{}
}

// ClusterDataSource is the data source implementation.
type ClusterDataSource struct {
	client *cratedb.ClientWithResponses
}

// Metadata returns the data source type name.
func (d *ClusterDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cluster"
}

// Schema defines the schema for the data source.
func (d *ClusterDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "To retrieve a cluster.",

		Attributes: map[string]schema.Attribute{
			"allow_custom_storage": schema.BoolAttribute{
				Computed:    true,
				Description: "The allow custom storage flag.",
			},
			"allow_suspend": schema.BoolAttribute{
				Computed:    true,
				Description: "The allow suspend flag.",
			},
			"backup_schedule": schema.StringAttribute{
				Computed:    true,
				Description: "The backup schedule.",
			},
			"channel": schema.StringAttribute{
				Computed:    true,
				Description: "The channel of the cluster.",
			},
			"crate_version": schema.StringAttribute{
				Computed:    true,
				Description: "The CrateDB version of the cluster.",
			},
			"dc": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "The DublinCore of the cluster.",
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
			"deletion_protected": schema.BoolAttribute{
				Computed:    true,
				Description: "The deletion protected flag.",
			},
			"external_ip": schema.StringAttribute{
				Computed:    true,
				Description: "The external IP address.",
			},
			"fqdn": schema.StringAttribute{
				Computed:    true,
				Description: "The Fully Qualified Domain Name.",
			},
			"gc_available": schema.BoolAttribute{
				Computed:    true,
				Description: "The garbage collection available flag.",
			},
			"hardware_specs": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "The hardware specs of the cluster.",
				Attributes: map[string]schema.Attribute{
					"cpus_per_node": schema.Int32Attribute{
						Computed:    true,
						Description: "The cpus per node.",
					},
					"disk_size_per_node_bytes": schema.Int64Attribute{
						Computed:    true,
						Description: "The disk size per node in bytes.",
					},
					"disk_type": schema.StringAttribute{
						Computed:    true,
						Description: "The disk type.",
					},
					"disks_per_node": schema.Int32Attribute{
						Computed:    true,
						Description: "The disks per node.",
					},
					"heap_size_bytes": schema.Int64Attribute{
						Computed:    true,
						Description: "The heap size in bytes.",
					},
					"memory_per_node_bytes": schema.Int64Attribute{
						Computed:    true,
						Description: "The memory per node in bytes.",
					},
				},
			},
			"health": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "The health of the cluster.",
				Attributes: map[string]schema.Attribute{
					"last_seen": schema.StringAttribute{
						Computed:    true,
						Description: "The last seen time.",
					},
					"running_operation": schema.StringAttribute{
						Computed:    true,
						Description: "The type of the currently running operation. Returns an empty string if there is no operation in progress.",
					},
					"status": schema.StringAttribute{
						Computed:    true,
						Description: "The health status of the cluster.",
					},
				},
			},
			"id": schema.StringAttribute{
				Required:    true,
				Description: "The id of the cluster.",
			},
			"ip_whitelist": schema.ListNestedAttribute{
				Computed:    true,
				Description: "The IP whitelist of the cluster.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"cidr": schema.StringAttribute{
							Computed:    true,
							Description: "The CIDR.",
						},
						"description": schema.StringAttribute{
							Computed:    true,
							Description: "The description.",
						},
					},
				},
			},
			"last_async_operation": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "The last async operation of the cluster.",
				Attributes: map[string]schema.Attribute{
					"dc": schema.SingleNestedAttribute{
						Computed:    true,
						Description: "The DublinCore of the cluster.",
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
						Computed:    true,
						Description: "The id of the last async operation.",
					},
					"status": schema.StringAttribute{
						Computed:    true,
						Description: "The status of the last async operation.",
					},
					"type": schema.StringAttribute{
						Computed:    true,
						Description: "The type of the last async operation.",
					},
				},
			},
			"name": schema.StringAttribute{
				Computed:    true,
				Description: "The name of the cluster.",
			},
			"num_nodes": schema.Int32Attribute{
				Computed:    true,
				Description: "The number of nodes in the cluster.",
			},
			"origin": schema.StringAttribute{
				Computed:    true,
				Description: "The origin of the cluster.",
			},
			"product_name": schema.StringAttribute{
				Computed:    true,
				Description: "The product name of the cluster.",
			},
			"product_tier": schema.StringAttribute{
				Computed:    true,
				Description: "The product tier of the cluster.",
			},
			"product_unit": schema.Int32Attribute{
				Computed:    true,
				Description: "The product unit of the cluster.",
			},
			"project_id": schema.StringAttribute{
				Computed:    true,
				Description: "The project id of the cluster.",
			},
			"subscription_id": schema.StringAttribute{
				Computed:    true,
				Description: "The subscription id of the cluster.",
			},
			"suspended": schema.BoolAttribute{
				Computed:    true,
				Description: "The suspended flag.",
			},
			"url": schema.StringAttribute{
				Computed:    true,
				Description: "The URL of the cluster.",
			},
			"username": schema.StringAttribute{
				Computed:    true,
				Description: "The username of the cluster.",
			},
			"password": schema.StringAttribute{
				Computed:    true,
				Sensitive:   true,
				Description: "The password of the cluster.",
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *ClusterDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *ClusterDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state ClusterModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	readClusterResponse, err := d.client.GetApiV2ClustersClusterIdWithResponse(ctx, state.Id.ValueString())
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
	state = *clusterState

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
