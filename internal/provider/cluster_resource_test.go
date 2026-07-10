package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccClusterResource deploys a real cluster and therefore needs a project
// and subscription to deploy into. Set CRATEDB_PRODUCT_NAME/CRATEDB_PRODUCT_TIER
// to override the free-tier defaults on paid accounts.
func TestAccClusterResource(t *testing.T) {
	organizationID := envOrSkip(t, "CRATEDB_ORGANIZATION_ID")
	projectID := envOrSkip(t, "CRATEDB_PROJECT_ID")
	subscriptionID := envOrSkip(t, "CRATEDB_SUBSCRIPTION_ID")
	crateVersion := envOrSkip(t, "CRATEDB_CRATE_VERSION")
	productName := envOrDefault("CRATEDB_PRODUCT_NAME", "crfree")
	productTier := envOrDefault("CRATEDB_PRODUCT_TIER", "default")

	name := acctest.RandomWithPrefix("tf-acc-test")
	firstPassword := acctest.RandomWithPrefix("tf-acc-password")
	secondPassword := acctest.RandomWithPrefix("tf-acc-password")

	clusterConfig := func(password string) string {
		return testAccProviderConfig + fmt.Sprintf(`
resource "cratedb_cluster" "test" {
  organization_id = %q
  crate_version   = %q
  name            = %q
  product_name    = %q
  product_tier    = %q
  project_id      = %q
  subscription_id = %q
  username        = "admin"
  password        = %q
}
`, organizationID, crateVersion, name, productName, productTier, projectID, subscriptionID, password)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: clusterConfig(firstPassword),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("cratedb_cluster.test", "name", name),
					resource.TestCheckResourceAttr("cratedb_cluster.test", "channel", "stable"),
					resource.TestCheckResourceAttr("cratedb_cluster.test", "product_name", productName),
					resource.TestCheckResourceAttr("cratedb_cluster.test", "project_id", projectID),
					resource.TestCheckResourceAttrSet("cratedb_cluster.test", "id"),
					resource.TestCheckResourceAttrSet("cratedb_cluster.test", "num_nodes"),
					resource.TestCheckResourceAttrSet("cratedb_cluster.test", "fqdn"),
					resource.TestCheckResourceAttrSet("cratedb_cluster.test", "dc.created"),
				),
			},
			// ImportState testing. The API never returns the password and the
			// organization id is not part of the cluster representation, so
			// both cannot be verified after import.
			{
				ResourceName:            "cratedb_cluster.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password", "organization_id"},
			},
			// Update (password change) and Read testing
			{
				Config: clusterConfig(secondPassword),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("cratedb_cluster.test", "password", secondPassword),
					resource.TestCheckResourceAttr("cratedb_cluster.test", "name", name),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
