package redash_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSoureceUser_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testAccDataSoureceUserConfigName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.redash_user.admin", "name", "admin"),
					resource.TestCheckResourceAttr("data.redash_user.admin", "email", "admin@example.com"),
				),
			},
			{
				Config: testAccDataSoureceUserConfigEmail,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.redash_user.admin", "name", "admin"),
					resource.TestCheckResourceAttr("data.redash_user.admin", "email", "admin@example.com"),
				),
			},
		},
	})
}

const testAccDataSoureceUserConfigName = `
data "redash_user" "admin" {
  name = "admin"
}
`

const testAccDataSoureceUserConfigEmail = `
data "redash_user" "admin" {
  email = "admin@example.com"
}
`
