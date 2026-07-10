package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/thulasirajkomminar/cratedb-cloud-go"
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

// ClusterDataSourceModel maps the data source schema data. Unlike
// ClusterModel it has no organization_id, which is not part of the API's
// cluster representation.
type ClusterDataSourceModel struct {
	AllowCustomStorage types.Bool                `tfsdk:"allow_custom_storage"`
	AllowSuspend       types.Bool                `tfsdk:"allow_suspend"`
	BackupSchedule     types.String              `tfsdk:"backup_schedule"`
	Channel            types.String              `tfsdk:"channel"`
	CrateVersion       types.String              `tfsdk:"crate_version"`
	Dc                 types.Object              `tfsdk:"dc"`
	DeletionProtected  types.Bool                `tfsdk:"deletion_protected"`
	ExternalIp         types.String              `tfsdk:"external_ip"`
	Fqdn               types.String              `tfsdk:"fqdn"`
	GcAvailable        types.Bool                `tfsdk:"gc_available"`
	HardwareSpecs      types.Object              `tfsdk:"hardware_specs"`
	Health             types.Object              `tfsdk:"health"`
	Id                 types.String              `tfsdk:"id"`
	IpWhitelist        []ClusterIpWhitelistModel `tfsdk:"ip_whitelist"`
	Name               types.String              `tfsdk:"name"`
	NumNodes           types.Int32               `tfsdk:"num_nodes"`
	Origin             types.String              `tfsdk:"origin"`
	ProductName        types.String              `tfsdk:"product_name"`
	ProductTier        types.String              `tfsdk:"product_tier"`
	ProductUnit        types.Int32               `tfsdk:"product_unit"`
	ProjectId          types.String              `tfsdk:"project_id"`
	SubscriptionId     types.String              `tfsdk:"subscription_id"`
	Suspended          types.Bool                `tfsdk:"suspended"`
	Url                types.String              `tfsdk:"url"`
	Username           types.String              `tfsdk:"username"`
	Password           types.String              `tfsdk:"password"`
}

// clusterDataSourceModelFrom converts the shared cluster model to the data
// source model.
func clusterDataSourceModelFrom(m ClusterModel) ClusterDataSourceModel {
	return ClusterDataSourceModel{
		AllowCustomStorage: m.AllowCustomStorage,
		AllowSuspend:       m.AllowSuspend,
		BackupSchedule:     m.BackupSchedule,
		Channel:            m.Channel,
		CrateVersion:       m.CrateVersion,
		Dc:                 m.Dc,
		DeletionProtected:  m.DeletionProtected,
		ExternalIp:         m.ExternalIp,
		Fqdn:               m.Fqdn,
		GcAvailable:        m.GcAvailable,
		HardwareSpecs:      m.HardwareSpecs,
		Health:             m.Health,
		Id:                 m.Id,
		IpWhitelist:        m.IpWhitelist,
		Name:               m.Name,
		NumNodes:           m.NumNodes,
		Origin:             m.Origin,
		ProductName:        m.ProductName,
		ProductTier:        m.ProductTier,
		ProductUnit:        m.ProductUnit,
		ProjectId:          m.ProjectId,
		SubscriptionId:     m.SubscriptionId,
		Suspended:          m.Suspended,
		Url:                m.Url,
		Username:           m.Username,
		Password:           m.Password,
	}
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
	if client := clientFromProviderData(req.ProviderData, "Data Source", &resp.Diagnostics); client != nil {
		d.client = client
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *ClusterDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state ClusterDataSourceModel

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

	if readClusterResponse.StatusCode() != 200 || readClusterResponse.JSON200 == nil {
		resp.Diagnostics.AddError(
			"Error getting cluster",
			apiErrorDetail(readClusterResponse.HTTPResponse, readClusterResponse.Body),
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
	state = clusterDataSourceModelFrom(*clusterState)

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
