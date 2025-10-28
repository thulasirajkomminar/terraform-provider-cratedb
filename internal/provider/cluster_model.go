package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/thulasirajkomminar/cratedb"
)

// ClusterModel maps CrateDB cluster schema data.
type ClusterModel struct {
	OrganizationId     types.String              `tfsdk:"organization_id"`
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

// ClusterDCModel maps CrateDB cluster HardwareSpecs schema data.
type ClusterHardwareSpecsModel struct {
	CpusPerNode          types.Int32  `tfsdk:"cpus_per_node"`
	DiskSizePerNodeBytes types.Int64  `tfsdk:"disk_size_per_node_bytes"`
	DiskType             types.String `tfsdk:"disk_type"`
	DisksPerNode         types.Int32  `tfsdk:"disks_per_node"`
	HeapSizeBytes        types.Int64  `tfsdk:"heap_size_bytes"`
	MemoryPerNodeBytes   types.Int64  `tfsdk:"memory_per_node_bytes"`
}

func (c ClusterHardwareSpecsModel) GetAttrType() map[string]attr.Type {
	return map[string]attr.Type{
		"cpus_per_node":            types.Int32Type,
		"disk_size_per_node_bytes": types.Int64Type,
		"disk_type":                types.StringType,
		"disks_per_node":           types.Int32Type,
		"heap_size_bytes":          types.Int64Type,
		"memory_per_node_bytes":    types.Int64Type,
	}
}

// ClusterHealthModel maps CrateDB cluster Health schema data.
type ClusterHealthModel struct {
	Status types.String `tfsdk:"status"`
}

func (c ClusterHealthModel) GetAttrType() map[string]attr.Type {
	return map[string]attr.Type{
		"status": types.StringType,
	}
}

// ClusterIpWhitelistModel maps CrateDB cluster IpWhitelist schema data.
type ClusterIpWhitelistModel struct {
	Cidr        types.String `tfsdk:"cidr"`
	Description types.String `tfsdk:"description"`
}

func (c ClusterIpWhitelistModel) GetAttrType() attr.Type {
	return types.ObjectType{AttrTypes: map[string]attr.Type{
		"cidr":        types.StringType,
		"description": types.StringType,
	}}
}

func getClusterModel(ctx context.Context, cluster cratedb.Cluster) (*ClusterModel, error) {
	dcValue := DCModel{
		Created:  types.StringValue(cluster.Dc.Created.String()),
		Modified: types.StringValue(cluster.Dc.Modified.String()),
	}

	dcObjectValue, diags := types.ObjectValueFrom(ctx, dcValue.GetAttrType(), dcValue)
	if diags.HasError() {
		return nil, fmt.Errorf("error getting cluster DC value")
	}

	hardwareSpecsValue := ClusterHardwareSpecsModel{
		CpusPerNode:          types.Int32Value(int32(*cluster.HardwareSpecs.CpusPerNode)),
		DiskSizePerNodeBytes: types.Int64Value(int64(*cluster.HardwareSpecs.DiskSizePerNodeBytes)),
		DiskType:             types.StringPointerValue(cluster.HardwareSpecs.DiskType),
		DisksPerNode:         types.Int32Value(int32(*cluster.HardwareSpecs.DisksPerNode)),
		HeapSizeBytes:        types.Int64Value(int64(*cluster.HardwareSpecs.HeapSizeBytes)),
		MemoryPerNodeBytes:   types.Int64Value(int64(*cluster.HardwareSpecs.MemoryPerNodeBytes)),
	}

	hardwareSpecsObjectValue, diags := types.ObjectValueFrom(ctx, hardwareSpecsValue.GetAttrType(), hardwareSpecsValue)
	if diags.HasError() {
		return nil, fmt.Errorf("error getting cluster hardware specs value")
	}

	healthValue := ClusterHealthModel{
		Status: types.StringValue(string(*cluster.Health.Status)),
	}

	healthObjectValue, diags := types.ObjectValueFrom(ctx, healthValue.GetAttrType(), healthValue)
	if diags.HasError() {
		return nil, fmt.Errorf("error getting cluster health value")
	}

	var ipWhitelistValues []ClusterIpWhitelistModel
	if cluster.IpWhitelist != nil {
		for _, ipWhitelist := range *cluster.IpWhitelist {
			ipWhitelistValues = append(ipWhitelistValues, ClusterIpWhitelistModel{
				Cidr:        types.StringValue(ipWhitelist.Cidr),
				Description: types.StringValue(string(*ipWhitelist.Description)),
			})
		}
	}

	return &ClusterModel{
		Dc:                 dcObjectValue,
		HardwareSpecs:      hardwareSpecsObjectValue,
		Health:             healthObjectValue,
		Id:                 types.StringPointerValue(cluster.Id),
		IpWhitelist:        ipWhitelistValues,
		AllowCustomStorage: types.BoolPointerValue(cluster.AllowCustomStorage),
		AllowSuspend:       types.BoolPointerValue(cluster.AllowSuspend),
		BackupSchedule:     types.StringPointerValue(cluster.BackupSchedule),
		Channel:            types.StringPointerValue(cluster.Channel),
		CrateVersion:       types.StringValue(cluster.CrateVersion),
		DeletionProtected:  types.BoolPointerValue(cluster.DeletionProtected),
		ExternalIp:         types.StringPointerValue(cluster.ExternalIp),
		Fqdn:               types.StringPointerValue(cluster.Fqdn),
		GcAvailable:        types.BoolPointerValue(cluster.GcAvailable),
		Name:               types.StringValue(cluster.Name),
		NumNodes:           types.Int32Value(int32(*cluster.NumNodes)),
		Origin:             types.StringPointerValue(cluster.Origin),
		ProductName:        types.StringValue(cluster.ProductName),
		ProductTier:        types.StringValue(cluster.ProductTier),
		ProductUnit:        types.Int32Value(int32(*cluster.ProductUnit)),
		ProjectId:          types.StringValue(cluster.ProjectId),
		SubscriptionId:     types.StringPointerValue(cluster.SubscriptionId),
		Suspended:          types.BoolPointerValue(cluster.Suspended),
		Url:                types.StringPointerValue(cluster.Url),
		Username:           types.StringValue(cluster.Username),
	}, nil
}
