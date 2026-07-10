package provider

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccProjectsDataSource(t *testing.T) {
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

data "cratedb_projects" "test" {
  depends_on = [cratedb_project.test]
}
`, name, organizationID, region),
				Check: resource.TestCheckResourceAttrWith("data.cratedb_projects.test", "projects.#", func(value string) error {
					n, err := strconv.Atoi(value)
					if err != nil {
						return err
					}
					if n < 1 {
						return fmt.Errorf("expected at least one project, got %d", n)
					}
					return nil
				}),
			},
		},
	})
}
