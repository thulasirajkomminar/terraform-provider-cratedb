package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/thulasirajkomminar/cratedb-cloud-go"
)

// apiTime tolerates timestamps without a timezone offset (e.g.
// "2026-07-10T10:41:02.983000"), which the regions endpoint returns and the
// generated client cannot parse as time.Time. Offset-less timestamps are
// interpreted as UTC.
type apiTime struct {
	time.Time
}

func (t *apiTime) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	for _, layout := range []string{time.RFC3339Nano, "2006-01-02T15:04:05.999999999"} {
		ts, err := time.Parse(layout, s)
		if err == nil {
			t.Time = ts
			return nil
		}
	}
	return fmt.Errorf("unsupported timestamp format %q", s)
}

// apiRegion mirrors cratedb.Region with tolerant timestamp decoding.
type apiRegion struct {
	Dc *struct {
		Created  *apiTime `json:"created"`
		Modified *apiTime `json:"modified"`
	} `json:"dc"`
	Deprecated       *bool    `json:"deprecated"`
	Description      *string  `json:"description"`
	IsEdgeRegion     *bool    `json:"is_edge_region"`
	LastSeen         *apiTime `json:"last_seen"`
	Name             *string  `json:"name"`
	OrganizationId   *string  `json:"organization_id"`
	Status           *string  `json:"status"`
	UpgradeAvailable *bool    `json:"upgrade_available"`
}

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &RegionsDataSource{}
	_ datasource.DataSourceWithConfigure = &RegionsDataSource{}
)

// NewRegionsDataSource is a helper function to simplify the provider implementation.
func NewRegionsDataSource() datasource.DataSource {
	return &RegionsDataSource{}
}

// RegionsDataSource is the data source implementation.
type RegionsDataSource struct {
	client *cratedb.ClientWithResponses
}

// RegionsDataSourceModel describes the data source data model.
type RegionsDataSourceModel struct {
	OrganizationId types.String  `tfsdk:"organization_id"`
	Regions        []RegionModel `tfsdk:"regions"`
}

// RegionModel maps CrateDB region schema data.
type RegionModel struct {
	Dc               types.Object      `tfsdk:"dc"`
	Deprecated       types.Bool        `tfsdk:"deprecated"`
	Description      types.String      `tfsdk:"description"`
	IsEdgeRegion     types.Bool        `tfsdk:"is_edge_region"`
	LastSeen         timetypes.RFC3339 `tfsdk:"last_seen"`
	Name             types.String      `tfsdk:"name"`
	OrganizationId   types.String      `tfsdk:"organization_id"`
	Status           types.String      `tfsdk:"status"`
	UpgradeAvailable types.Bool        `tfsdk:"upgrade_available"`
}

// Metadata returns the data source type name.
func (d *RegionsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_regions"
}

// Schema defines the schema for the data source.
func (d *RegionsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "To retrieve all regions available to the current user.",

		Attributes: map[string]schema.Attribute{
			"organization_id": schema.StringAttribute{
				Optional:    true,
				Description: "Filter regions by organization id, to include the organization's edge regions.",
			},
			"regions": schema.ListNestedAttribute{
				Computed:    true,
				Description: "The list of regions.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"dc": schema.SingleNestedAttribute{
							Computed:    true,
							Description: "The DublinCore of the region.",
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
						"deprecated": schema.BoolAttribute{
							Computed:    true,
							Description: "Whether the region is deprecated.",
						},
						"description": schema.StringAttribute{
							Computed:    true,
							Description: "The description of the region.",
						},
						"is_edge_region": schema.BoolAttribute{
							Computed:    true,
							Description: "Whether the region is an edge region.",
						},
						"last_seen": schema.StringAttribute{
							CustomType:  timetypes.RFC3339Type{},
							Computed:    true,
							Description: "The last seen time of the region.",
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "The name of the region.",
						},
						"organization_id": schema.StringAttribute{
							Computed:    true,
							Description: "The organization id of the region, for organization-specific edge regions.",
						},
						"status": schema.StringAttribute{
							Computed:    true,
							Description: "The status of the region.",
						},
						"upgrade_available": schema.BoolAttribute{
							Computed:    true,
							Description: "Whether an upgrade is available for the region.",
						},
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *RegionsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if client := clientFromProviderData(req.ProviderData, "Data Source", &resp.Diagnostics); client != nil {
		d.client = client
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *RegionsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state RegionsDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := &cratedb.GetApiV2RegionsParams{}
	if !state.OrganizationId.IsNull() {
		organizationId := state.OrganizationId.ValueString()
		params.OrganizationId = &organizationId
	}

	// Use the raw client and decode the body with tolerant timestamp types:
	// the generated typed client fails on this endpoint's offset-less
	// timestamps.
	readRegionsResponse, err := d.client.GetApiV2Regions(ctx, params)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting regions",
			err.Error(),
		)
		return
	}
	defer func() { _ = readRegionsResponse.Body.Close() }()

	body, err := io.ReadAll(readRegionsResponse.Body)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting regions",
			"Could not read regions response: "+err.Error(),
		)
		return
	}

	if readRegionsResponse.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(
			"Error getting regions",
			apiErrorDetail(readRegionsResponse, body),
		)
		return
	}

	var apiRegions []apiRegion
	if err := json.Unmarshal(body, &apiRegions); err != nil {
		resp.Diagnostics.AddError(
			"Error getting regions",
			"Could not parse regions response: "+err.Error()+"\n"+apiErrorDetail(readRegionsResponse, body),
		)
		return
	}

	for _, region := range apiRegions {
		regionState, err := getRegionModel(ctx, region)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error getting region model",
				err.Error(),
			)
			return
		}
		state.Regions = append(state.Regions, *regionState)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func getRegionModel(ctx context.Context, region apiRegion) (*RegionModel, error) {
	var dc *cratedb.DublinCore
	if region.Dc != nil {
		dc = &cratedb.DublinCore{}
		if region.Dc.Created != nil {
			dc.Created = &region.Dc.Created.Time
		}
		if region.Dc.Modified != nil {
			dc.Modified = &region.Dc.Modified.Time
		}
	}

	dcObjectValue, err := getDCObjectValue(ctx, dc)
	if err != nil {
		return nil, fmt.Errorf("error getting region DC value: %w", err)
	}

	lastSeen := timetypes.NewRFC3339Null()
	if region.LastSeen != nil {
		lastSeen = timetypes.NewRFC3339TimeValue(region.LastSeen.Time)
	}

	return &RegionModel{
		Dc:               dcObjectValue,
		Deprecated:       types.BoolPointerValue(region.Deprecated),
		Description:      types.StringPointerValue(region.Description),
		IsEdgeRegion:     types.BoolPointerValue(region.IsEdgeRegion),
		LastSeen:         lastSeen,
		Name:             types.StringPointerValue(region.Name),
		OrganizationId:   types.StringPointerValue(region.OrganizationId),
		Status:           types.StringPointerValue(region.Status),
		UpgradeAvailable: types.BoolPointerValue(region.UpgradeAvailable),
	}, nil
}
