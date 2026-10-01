package test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccGroup_basic(t *testing.T) {
	var id string

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testAccGroupConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("redash_group.my_group", "name", "my-group"),
					func(s *terraform.State) error {
						id = s.RootModule().Resources["redash_group.my_group"].Primary.ID
						return nil
					},
				),
			},
			{
				Config: testAccGroupConfigBasic2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("redash_group.my_group", "name", "my-group2"),
					// renamed in place
					resource.TestCheckResourceAttrPtr("redash_group.my_group", "id", &id),
				),
			},
			{
				ResourceName:      "redash_group.my_group",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccGroupConfigBasic = `
resource "redash_group" "my_group" {
	name = "my-group"
}
`

const testAccGroupConfigBasic2 = `
resource "redash_group" "my_group" {
	name = "my-group2"
}
`
