package test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSoureceUsers_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testAccDataSoureceUsersConfigNameFilter,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.redash_users.users", "ids.#", "1"),
					resource.TestCheckResourceAttr("data.redash_users.users", "ids.0", "1"),
				),
			},
			{
				Config: testAccDataSoureceUsersConfigNameFilter2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.redash_users.users", "ids.#", "1"),
					resource.TestCheckResourceAttr("data.redash_users.users", "ids.0", "1"),
				),
			},
			{
				Config: testAccDataSoureceUsersConfigNameFilterNG,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.redash_users.users", "ids.#", "0"),
				),
			},
			{
				Config: testAccDataSoureceUsersConfigEmailFilter,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.redash_users.users", "ids.#", "1"),
					resource.TestCheckResourceAttr("data.redash_users.users", "ids.0", "1"),
				),
			},
			{
				Config: testAccDataSoureceUsersConfigEmailFilter2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.redash_users.users", "ids.#", "1"),
					resource.TestCheckResourceAttr("data.redash_users.users", "ids.0", "1"),
				),
			},
			{
				Config: testAccDataSoureceUsersConfigEmailFilterNG,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.redash_users.users", "ids.#", "0"),
				),
			},
			{
				Config: testAccDataSoureceUsersConfigComplexFilter,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.redash_users.users", "ids.#", "1"),
					resource.TestCheckResourceAttr("data.redash_users.users", "ids.0", "1"),
				),
			},
			{
				Config: testAccDataSoureceUsersConfigComplexFilterNG,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.redash_users.users", "ids.#", "0"),
				),
			},
			{
				Config: testAccDataSoureceUsersConfigComplexFilterNG2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.redash_users.users", "ids.#", "0"),
				),
			},
		},
	})
}

const testAccDataSoureceUsersConfigNameFilter = `
data "redash_users" "users" {
  name = "a*n"
}
`

const testAccDataSoureceUsersConfigNameFilter2 = `
data "redash_users" "users" {
  name = "ad?in"
}
`

const testAccDataSoureceUsersConfigNameFilterNG = `
data "redash_users" "users" {
  name = "a*d?min"
}
`

const testAccDataSoureceUsersConfigEmailFilter = `
data "redash_users" "users" {
  email = "*@example.com"
}
`

const testAccDataSoureceUsersConfigEmailFilter2 = `
data "redash_users" "users" {
  email = "ad?in@example.com"
}
`

const testAccDataSoureceUsersConfigEmailFilterNG = `
data "redash_users" "users" {
  email = "a*d?min@example.com"
}
`

const testAccDataSoureceUsersConfigComplexFilter = `
data "redash_users" "users" {
	name  = "a*d?in"
	email = "a*d?in@example.com"
}
`

const testAccDataSoureceUsersConfigComplexFilterNG = `
data "redash_users" "users" {
	name  = "a*d?min"
	email = "a*d?in@example.com"
}
`

const testAccDataSoureceUsersConfigComplexFilterNG2 = `
data "redash_users" "users" {
	name  = "a*d?in"
	email = "a*d?min@example.com"
}
`
