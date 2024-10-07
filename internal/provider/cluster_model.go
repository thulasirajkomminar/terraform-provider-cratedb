package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/komminarlabs/cratedb"
)

// ClusterModel maps CrateDB cluster schema data.
type ClusterModel struct {
	AllowCustomStorage types.Bool     `tfsdk:"allow_custom_storage"`
	AllowSuspend       types.Bool     `tfsdk:"allow_suspend"`
	BackupSchedule     types.String   `tfsdk:"backup_schedule"`
	Channel            types.String   `tfsdk:"channel"`
	CrateVersion       types.String   `tfsdk:"crate_version"`
	Dc                 types.Object   `tfsdk:"dc"`
	DeletionProtected  types.Bool     `tfsdk:"deletion_protected"`
	ExternalIp         types.String   `tfsdk:"external_ip"`
	Fqdn               types.String   `tfsdk:"fqdn"`
	GcAvailable        types.Bool     `tfsdk:"gc_available"`
	HardwareSpecs      types.Object   `tfsdk:"hardware_specs"`
	Health             types.Object   `tfsdk:"health"`
	Id                 types.String   `tfsdk:"id"`
	IpWhitelist        []types.Object `tfsdk:"ip_whitelist"`
	LastAsyncOperation types.Object   `tfsdk:"last_async_operation"`
	Name               types.String   `tfsdk:"name"`
	NumNodes           types.Int32    `tfsdk:"num_nodes"`
	Origin             types.String   `tfsdk:"origin"`
	ProductName        types.String   `tfsdk:"product_name"`
	ProductTier        types.String   `tfsdk:"product_tier"`
	ProductUnit        types.Int32    `tfsdk:"product_unit"`
	ProjectId          types.String   `tfsdk:"project_id"`
	SubscriptionId     types.String   `tfsdk:"subscription_id"`
	Suspended          types.Bool     `tfsdk:"suspended"`
	Url                types.String   `tfsdk:"url"`
	Username           types.String   `tfsdk:"username"`
	Password           types.String   `tfsdk:"password"`
}

// ClusterDCModel maps CrateDB cluster HardwareSpecs schema data.
type ClusterHardwareSpecsModel struct {
	CpusPerNode          types.Int32  `tfsdk:"cpus_per_node"`
	DiskSizePerNodeBytes types.Int32  `tfsdk:"disk_size_per_node_bytes"`
	DiskType             types.String `tfsdk:"disk_type"`
	DisksPerNode         types.Int32  `tfsdk:"disks_per_node"`
	HeapSizeBytes        types.Int32  `tfsdk:"heap_size_bytes"`
	MemoryPerNodeBytes   types.Int32  `tfsdk:"memory_per_node_bytes"`
}

func (c ClusterHardwareSpecsModel) GetAttrType() map[string]attr.Type {
	return map[string]attr.Type{
		"cpus_per_node":            types.Int32Type,
		"disk_size_per_node_bytes": types.Int32Type,
		"disk_type":                types.StringType,
		"disks_per_node":           types.Int32Type,
		"heap_size_bytes":          types.Int32Type,
		"memory_per_node_bytes":    types.Int32Type,
	}
}

// ClusterHealthModel maps CrateDB cluster Health schema data.
type ClusterHealthModel struct {
	LastSeen         types.String `tfsdk:"last_seen"`
	RunningOperation types.String `tfsdk:"running_operation"`
	Status           types.String `tfsdk:"status"`
}

func (c ClusterHealthModel) GetAttrType() map[string]attr.Type {
	return map[string]attr.Type{
		"last_seen":         types.StringType,
		"running_operation": types.StringType,
		"status":            types.StringType,
	}
}

// ClusterIpWhitelistModel maps CrateDB cluster IpWhitelist schema data.
type ClusterIpWhitelistModel struct {
	Cidr        types.String `tfsdk:"cidr"`
	Description types.String `tfsdk:"description"`
}

func (c ClusterIpWhitelistModel) GetAttrType() map[string]attr.Type {
	return map[string]attr.Type{
		"cidr":        types.StringType,
		"description": types.StringType,
	}
}

// ClusterLastAsyncOperationModel maps CrateDB cluster LastAsyncOperation schema data.
type ClusterLastAsyncOperationModel struct {
	Dc     types.Object `tfsdk:"dc"`
	Id     types.String `tfsdk:"id"`
	Status types.String `tfsdk:"status"`
	Type   types.String `tfsdk:"type"`
}

func (c ClusterLastAsyncOperationModel) GetAttrType() map[string]attr.Type {
	return map[string]attr.Type{
		"dc":     types.ObjectType{AttrTypes: DCModel{}.GetAttrType()},
		"id":     types.StringType,
		"status": types.StringType,
		"type":   types.StringType,
	}
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
		DiskSizePerNodeBytes: types.Int32Value(int32(*cluster.HardwareSpecs.DiskSizePerNodeBytes)),
		DiskType:             types.StringValue(*cluster.HardwareSpecs.DiskType),
		DisksPerNode:         types.Int32Value(int32(*cluster.HardwareSpecs.DisksPerNode)),
		HeapSizeBytes:        types.Int32Value(int32(*cluster.HardwareSpecs.HeapSizeBytes)),
		MemoryPerNodeBytes:   types.Int32Value(int32(*cluster.HardwareSpecs.MemoryPerNodeBytes)),
	}

	hardwareSpecsObjectValue, diags := types.ObjectValueFrom(ctx, hardwareSpecsValue.GetAttrType(), hardwareSpecsValue)
	if diags.HasError() {
		return nil, fmt.Errorf("error getting cluster hardware specs value")
	}

	healthValue := ClusterHealthModel{
		LastSeen:         types.StringValue(cluster.Health.LastSeen.String()),
		RunningOperation: types.StringValue(string(*cluster.Health.RunningOperation)),
		Status:           types.StringValue(string(*cluster.Health.Status)),
	}

	healthObjectValue, diags := types.ObjectValueFrom(ctx, healthValue.GetAttrType(), healthValue)
	if diags.HasError() {
		return nil, fmt.Errorf("error getting cluster health value")
	}

	ipWhitelistObjectValueList := make([]types.Object, 0)
	for _, ipWhitelist := range *cluster.IpWhitelist {
		ipWhitelistValue := ClusterIpWhitelistModel{
			Cidr:        types.StringValue(ipWhitelist.Cidr),
			Description: types.StringValue(string(*ipWhitelist.Description)),
		}

		ipWhitelistObjectValue, diags := types.ObjectValueFrom(ctx, ipWhitelistValue.GetAttrType(), ipWhitelistValue)
		if diags.HasError() {
			return nil, fmt.Errorf("error getting cluster ip whitelist value")
		}
		ipWhitelistObjectValueList = append(ipWhitelistObjectValueList, ipWhitelistObjectValue)
	}

	lastAsyncOperationValue := ClusterLastAsyncOperationModel{
		Id:     types.StringValue(*cluster.LastAsyncOperation.Id),
		Status: types.StringValue(string(*cluster.LastAsyncOperation.Status)),
		Type:   types.StringValue(string(*cluster.LastAsyncOperation.Type)),
	}
	lastAsyncOperationObjectValue, diags := types.ObjectValueFrom(ctx, lastAsyncOperationValue.GetAttrType(), lastAsyncOperationValue)
	if diags.HasError() {
		return nil, fmt.Errorf("error getting cluster last async operation value")
	}

	clusterModel := ClusterModel{
		Dc:                 dcObjectValue,
		HardwareSpecs:      hardwareSpecsObjectValue,
		Health:             healthObjectValue,
		Id:                 types.StringValue(*cluster.Id),
		IpWhitelist:        ipWhitelistObjectValueList,
		LastAsyncOperation: lastAsyncOperationObjectValue,
		AllowCustomStorage: types.BoolValue(*cluster.AllowCustomStorage),
		AllowSuspend:       types.BoolValue(*cluster.AllowSuspend),
		BackupSchedule:     types.StringValue(*cluster.BackupSchedule),
		Channel:            types.StringValue(*cluster.Channel),
		CrateVersion:       types.StringValue(cluster.CrateVersion),
		DeletionProtected:  types.BoolValue(*cluster.DeletionProtected),
		ExternalIp:         types.StringValue(*cluster.ExternalIp),
		Fqdn:               types.StringValue(*cluster.Fqdn),
		GcAvailable:        types.BoolValue(*cluster.GcAvailable),
		Name:               types.StringValue(cluster.Name),
		NumNodes:           types.Int32Value(int32(*cluster.NumNodes)),
		Origin:             types.StringValue(*cluster.Origin),
		ProductName:        types.StringValue(cluster.ProductName),
		ProductTier:        types.StringValue(cluster.ProductTier),
		ProductUnit:        types.Int32Value(int32(*cluster.ProductUnit)),
		ProjectId:          types.StringValue(cluster.ProjectId),
		SubscriptionId:     types.StringValue(*cluster.SubscriptionId),
		Suspended:          types.BoolValue(*cluster.Suspended),
		Url:                types.StringValue(*cluster.Url),
		Username:           types.StringValue(cluster.Username),
		Password:           types.StringValue(*cluster.Password),
	}
	return &clusterModel, nil
}
