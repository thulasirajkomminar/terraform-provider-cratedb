package provider

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// testAccProtoV6ProviderFactories are used to instantiate a provider during
// acceptance testing. The factory function will be invoked for every Terraform
// CLI command executed to create a provider server to which the CLI can
// reattach.
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"cratedb": providerserver.NewProtocol6WithError(New("test")()),
}

const testAccProviderConfig = `provider "cratedb" {}
`

// testAccPreCheck validates that the credentials needed by every acceptance
// test are available. Tests with additional requirements call envOrSkip for
// the extra environment variables they need.
func testAccPreCheck(t *testing.T) {
	t.Helper()
	if os.Getenv("CRATEDB_API_KEY") == "" || os.Getenv("CRATEDB_API_SECRET") == "" {
		t.Fatal("CRATEDB_API_KEY and CRATEDB_API_SECRET must be set for acceptance tests")
	}
}

// envOrSkip returns the value of the environment variable or skips the test
// when it is not set, so the acceptance test suite degrades gracefully on
// accounts that cannot provide every fixture.
func envOrSkip(t *testing.T, name string) string {
	t.Helper()
	v := os.Getenv(name)
	if v == "" {
		t.Skipf("%s must be set for this acceptance test", name)
	}
	return v
}

// envOrDefault returns the value of the environment variable, or the default
// when it is not set.
func envOrDefault(name, def string) string {
	if v := os.Getenv(name); v != "" {
		return v
	}
	return def
}

// discoverRegion returns CRATEDB_REGION when set, and otherwise picks the
// first non-deprecated, non-edge region from the API so the project tests can
// run without manual region configuration.
func discoverRegion(t *testing.T) string {
	t.Helper()

	if v := os.Getenv("CRATEDB_REGION"); v != "" {
		return v
	}

	// resource.Test skips without TF_ACC; don't call the API before it does.
	if os.Getenv("TF_ACC") == "" {
		t.Skip("TF_ACC must be set for acceptance tests")
	}
	testAccPreCheck(t)

	url := strings.TrimRight(envOrDefault("CRATEDB_URL", defaultURL), "/") + "/api/v2/regions/"
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		t.Fatalf("building regions discovery request: %v", err)
	}
	req.SetBasicAuth(os.Getenv("CRATEDB_API_KEY"), os.Getenv("CRATEDB_API_SECRET"))
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("discovering a region: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Fatalf("discovering a region: HTTP %d (set CRATEDB_REGION to skip discovery)", resp.StatusCode)
	}

	var regions []apiRegion
	if err := json.Unmarshal(body, &regions); err != nil {
		t.Fatalf("parsing regions discovery response: %v", err)
	}

	for _, region := range regions {
		if region.Name == nil {
			continue
		}
		if region.Deprecated != nil && *region.Deprecated {
			continue
		}
		if region.IsEdgeRegion != nil && *region.IsEdgeRegion {
			continue
		}
		return *region.Name
	}

	t.Skip("no usable region found; set CRATEDB_REGION")
	return ""
}
