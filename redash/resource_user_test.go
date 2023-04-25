package redash_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccUser_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testAccUserConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("redash_user.my_user", "name", "John Smith"),
					resource.TestCheckResourceAttr("redash_user.my_user", "email", "jhons@my.example.com"),
				),
			},
			{
				Config: testAccUserConfigBasic2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("redash_user.my_user", "name", "Dave Smith"),
					resource.TestCheckResourceAttr("redash_user.my_user", "email", "dhons@my.example.com"),
				),
			},
		},
	})
}

const testAccUserConfigBasic = `
resource "redash_user" "my_user" {
	name  = "John Smith"
	email = "jhons@my.example.com"
}
`

const testAccUserConfigBasic2 = `
resource "redash_user" "my_user" {
	name  = "Dave Smith"
	email = "dhons@my.example.com"
}
`
