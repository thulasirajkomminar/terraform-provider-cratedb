package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccOrganizationResource(t *testing.T) {
	name := acctest.RandomWithPrefix("tf-acc-test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccProviderConfig + fmt.Sprintf(`
resource "cratedb_organization" "test" {
  name = %q
}
`, name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("cratedb_organization.test", "name", name),
					resource.TestCheckResourceAttrSet("cratedb_organization.test", "id"),
					resource.TestCheckResourceAttrSet("cratedb_organization.test", "dc.created"),
					resource.TestCheckResourceAttrSet("cratedb_organization.test", "project_count"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "cratedb_organization.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccProviderConfig + fmt.Sprintf(`
resource "cratedb_organization" "test" {
  name = %q
}
`, name+"-renamed"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("cratedb_organization.test", "name", name+"-renamed"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
