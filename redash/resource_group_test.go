package redash_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccGroup_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testAccGroupConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("redash_group.my_group", "name", "my-group"),
				),
			},
			{
				Config: testAccGroupConfigBasic2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("redash_group.my_group", "name", "my-group2"),
				),
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
