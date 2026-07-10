package provider

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/thulasirajkomminar/cratedb-cloud-go"
)

// defaultURL is the CrateDB Cloud API endpoint used when neither the `url`
// provider attribute nor the CRATEDB_URL environment variable is set.
const defaultURL = "https://console.cratedb.cloud"

// Ensure the implementation satisfies the expected interfaces.
var _ provider.Provider = &CrateDBProvider{}

// CrateDBProvider defines the provider implementation.
type CrateDBProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// CrateDBProviderModel maps provider schema data to a Go type.
type CrateDBProviderModel struct {
	APIKey    types.String `tfsdk:"api_key"`
	APISecret types.String `tfsdk:"api_secret"`
	URL       types.String `tfsdk:"url"`
}

// Metadata returns the provider type name.
func (p *CrateDBProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "cratedb"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *CrateDBProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "CrateDB provider to deploy and manage resources supported by CrateDB.",

		Attributes: map[string]schema.Attribute{
			"api_key": schema.StringAttribute{
				Description: "The API key. Can also be set with the `CRATEDB_API_KEY` environment variable.",
				Optional:    true,
				Sensitive:   true,
			},
			"api_secret": schema.StringAttribute{
				Description: "The API secret. Can also be set with the `CRATEDB_API_SECRET` environment variable.",
				Optional:    true,
				Sensitive:   true,
			},
			"url": schema.StringAttribute{
				Description: "The CrateDB Cloud URL. Can also be set with the `CRATEDB_URL` environment variable. Defaults to `" + defaultURL + "`.",
				Optional:    true,
			},
		},
	}
}

// Configure prepares a CrateDB API client for data sources and resources.
func (p *CrateDBProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	// Retrieve provider data from configuration
	var config CrateDBProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If practitioner provided a configuration value for any of the
	// attributes, it must be a known value.

	if config.APIKey.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_key"),
			"Unknown CrateDB API Key",
			"The provider cannot create the CrateDB client as there is an unknown configuration value for the CrateDB API Key. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the CRATEDB_API_KEY environment variable.",
		)
	}

	if config.APISecret.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_secret"),
			"Unknown CrateDB API Secret",
			"The provider cannot create the CrateDB client as there is an unknown configuration value for the CrateDB API Secret. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the CRATEDB_API_SECRET environment variable.",
		)
	}

	if config.URL.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("url"),
			"Unknown CrateDB Cloud URL",
			"The provider cannot create the CrateDB client as there is an unknown configuration value for the CrateDB Cloud URL. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the CRATEDB_URL environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.

	apiKey := os.Getenv("CRATEDB_API_KEY")
	apiSecret := os.Getenv("CRATEDB_API_SECRET")
	url := os.Getenv("CRATEDB_URL")

	if !config.APIKey.IsNull() {
		apiKey = config.APIKey.ValueString()
	}

	if !config.APISecret.IsNull() {
		apiSecret = config.APISecret.ValueString()
	}

	if !config.URL.IsNull() {
		url = config.URL.ValueString()
	}

	if url == "" {
		url = defaultURL
	}
	url = strings.TrimRight(url, "/")

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.

	if apiKey == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_key"),
			"Missing CrateDB API Key",
			"The provider cannot create the CrateDB client as there is a missing or empty value for the CrateDB API Key. "+
				"Set the api_key value in the configuration or use the CRATEDB_API_KEY environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if apiSecret == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_secret"),
			"Missing CrateDB API Secret",
			"The provider cannot create the CrateDB client as there is a missing or empty value for the CrateDB API Secret. "+
				"Set the api_secret value in the configuration or use the CRATEDB_API_SECRET environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "cratedb_url", url)
	tflog.Debug(ctx, "Creating CrateDB client")

	// Create a new CrateDB client using a retrying HTTP client whose
	// transport chain logs every request/response with TF_LOG=DEBUG while
	// keeping credentials out of the logs.
	client, err := cratedb.NewClientWithResponses(
		url,
		cratedb.WithRequestEditorFn(func(ctx context.Context, req *http.Request) error {
			req.Header.Set("Accept", "application/json")
			return nil
		}),
		cratedb.WithHTTPClient(newAPIHTTPClient(apiKey, apiSecret)),
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create CrateDB Client",
			"An unexpected error occurred when creating the CrateDB client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"CrateDB Client Error: "+err.Error(),
		)
		return
	}

	// Make the CrateDB client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = client
	resp.ResourceData = client
	tflog.Info(ctx, "Configured CrateDB client", map[string]any{"success": true})
}

// Resources defines the resources implemented in the provider.
func (p *CrateDBProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewClusterResource,
		NewOrganizationResource,
		NewProjectResource,
	}
}

// DataSources defines the data sources implemented in the provider.
func (p *CrateDBProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewClusterDataSource,
		NewOrganizationDataSource,
		NewOrganizationsDataSource,
		NewProjectDataSource,
		NewProjectsDataSource,
		NewRegionsDataSource,
	}
}

// New is a helper function to simplify provider server and testing implementation.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &CrateDBProvider{
			version: version,
		}
	}
}
