package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccProjectDataSource(t *testing.T) {
	organizationID := envOrSkip(t, "CRATEDB_ORGANIZATION_ID")
	region := discoverRegion(t)
	name := acctest.RandomWithPrefix("tf-acc-test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProviderConfig + fmt.Sprintf(`
resource "cratedb_project" "test" {
  name            = %q
  organization_id = %q
  region          = %q
}

data "cratedb_project" "test" {
  id = cratedb_project.test.id
}
`, name, organizationID, region),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.cratedb_project.test", "name", name),
					resource.TestCheckResourceAttr("data.cratedb_project.test", "organization_id", organizationID),
					resource.TestCheckResourceAttr("data.cratedb_project.test", "region", region),
					resource.TestCheckResourceAttrSet("data.cratedb_project.test", "dc.created"),
				),
			},
		},
	})
}
