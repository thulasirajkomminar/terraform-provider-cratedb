package provider

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccOrganizationsDataSource(t *testing.T) {
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

data "cratedb_organizations" "test" {
  depends_on = [cratedb_organization.test]
}
`, name),
				Check: resource.TestCheckResourceAttrWith("data.cratedb_organizations.test", "organizations.#", func(value string) error {
					n, err := strconv.Atoi(value)
					if err != nil {
						return err
					}
					if n < 1 {
						return fmt.Errorf("expected at least one organization, got %d", n)
					}
					return nil
				}),
			},
		},
	})
}
