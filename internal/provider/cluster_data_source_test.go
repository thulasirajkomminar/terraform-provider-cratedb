package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccClusterDataSource reads an existing cluster, since deploying one just
// for the data source would be slow and costly.
func TestAccClusterDataSource(t *testing.T) {
	clusterID := envOrSkip(t, "CRATEDB_CLUSTER_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProviderConfig + fmt.Sprintf(`
data "cratedb_cluster" "test" {
  id = %q
}
`, clusterID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.cratedb_cluster.test", "id", clusterID),
					resource.TestCheckResourceAttrSet("data.cratedb_cluster.test", "name"),
					resource.TestCheckResourceAttrSet("data.cratedb_cluster.test", "crate_version"),
					resource.TestCheckResourceAttrSet("data.cratedb_cluster.test", "project_id"),
					resource.TestCheckResourceAttrSet("data.cratedb_cluster.test", "num_nodes"),
					resource.TestCheckResourceAttrSet("data.cratedb_cluster.test", "health.status"),
					resource.TestCheckResourceAttrSet("data.cratedb_cluster.test", "dc.created"),
				),
			},
		},
	})
}
