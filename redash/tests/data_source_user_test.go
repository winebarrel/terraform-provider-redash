package test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceUser_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceUserConfigName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.redash_user.admin", "name", "admin"),
					resource.TestCheckResourceAttr("data.redash_user.admin", "email", "admin@example.com"),
				),
			},
			{
				Config: testAccDataSourceUserConfigEmail,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.redash_user.admin", "name", "admin"),
					resource.TestCheckResourceAttr("data.redash_user.admin", "email", "admin@example.com"),
				),
			},
		},
	})
}

const testAccDataSourceUserConfigName = `
data "redash_user" "admin" {
  name = "admin"
}
`

const testAccDataSourceUserConfigEmail = `
data "redash_user" "admin" {
  email = "admin@example.com"
}
`
