package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccOrganizationDataSource(t *testing.T) {
	name := acctest.RandomWithPrefix("tf-acc-test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProviderConfig + fmt.Sprintf(`
resource "cratedb_organization" "test" {
  name = %q
}

data "cratedb_organization" "test" {
  id = cratedb_organization.test.id
}
`, name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.cratedb_organization.test", "name", name),
					resource.TestCheckResourceAttrPair("data.cratedb_organization.test", "id", "cratedb_organization.test", "id"),
					resource.TestCheckResourceAttrSet("data.cratedb_organization.test", "dc.created"),
				),
			},
		},
	})
}
