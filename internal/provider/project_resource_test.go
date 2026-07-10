package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccProjectResource(t *testing.T) {
	organizationID := envOrSkip(t, "CRATEDB_ORGANIZATION_ID")
	region := discoverRegion(t)
	name := acctest.RandomWithPrefix("tf-acc-test")

	projectConfig := func(name string) string {
		return testAccProviderConfig + fmt.Sprintf(`
resource "cratedb_project" "test" {
  name            = %q
  organization_id = %q
  region          = %q
}
`, name, organizationID, region)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: projectConfig(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("cratedb_project.test", "name", name),
					resource.TestCheckResourceAttr("cratedb_project.test", "organization_id", organizationID),
					resource.TestCheckResourceAttr("cratedb_project.test", "region", region),
					resource.TestCheckResourceAttrSet("cratedb_project.test", "id"),
					resource.TestCheckResourceAttrSet("cratedb_project.test", "dc.created"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "cratedb_project.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: projectConfig(name + "-renamed"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("cratedb_project.test", "name", name+"-renamed"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
