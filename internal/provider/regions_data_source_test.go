package provider

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccRegionsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProviderConfig + `
data "cratedb_regions" "test" {}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrWith("data.cratedb_regions.test", "regions.#", func(value string) error {
						n, err := strconv.Atoi(value)
						if err != nil {
							return err
						}
						if n < 1 {
							return fmt.Errorf("expected at least one region, got %d", n)
						}
						return nil
					}),
					resource.TestCheckResourceAttrSet("data.cratedb_regions.test", "regions.0.name"),
				),
			},
		},
	})
}
