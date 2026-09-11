package test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceGroup_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testAccGroupConfigBasic,
			},
			{
				Config: testAccDataSourceGroupConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.redash_group.my_group", "name", "my-group"),
				),
			},
		},
	})
}

const testAccDataSourceGroupConfigBasic = testAccGroupConfigBasic + `
data "redash_group" "my_group" {
  name = "my-group"
}
`
