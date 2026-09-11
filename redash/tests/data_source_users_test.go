package test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceUsers_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceUsersConfigNameFilter,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.redash_users.users", "ids.#", "1"),
					resource.TestCheckResourceAttr("data.redash_users.users", "ids.0", "1"),
				),
			},
			{
				Config: testAccDataSourceUsersConfigNameFilter2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.redash_users.users", "ids.#", "1"),
					resource.TestCheckResourceAttr("data.redash_users.users", "ids.0", "1"),
				),
			},
			{
				Config: testAccDataSourceUsersConfigNameFilterNG,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.redash_users.users", "ids.#", "0"),
				),
			},
			{
				Config: testAccDataSourceUsersConfigEmailFilter,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.redash_users.users", "ids.#", "1"),
					resource.TestCheckResourceAttr("data.redash_users.users", "ids.0", "1"),
				),
			},
			{
				Config: testAccDataSourceUsersConfigEmailFilter2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.redash_users.users", "ids.#", "1"),
					resource.TestCheckResourceAttr("data.redash_users.users", "ids.0", "1"),
				),
			},
			{
				Config: testAccDataSourceUsersConfigEmailFilterNG,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.redash_users.users", "ids.#", "0"),
				),
			},
			{
				Config: testAccDataSourceUsersConfigComplexFilter,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.redash_users.users", "ids.#", "1"),
					resource.TestCheckResourceAttr("data.redash_users.users", "ids.0", "1"),
				),
			},
			{
				Config: testAccDataSourceUsersConfigComplexFilterNG,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.redash_users.users", "ids.#", "0"),
				),
			},
			{
				Config: testAccDataSourceUsersConfigComplexFilterNG2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.redash_users.users", "ids.#", "0"),
				),
			},
		},
	})
}

const testAccDataSourceUsersConfigNameFilter = `
data "redash_users" "users" {
  name = "a*n"
}
`

const testAccDataSourceUsersConfigNameFilter2 = `
data "redash_users" "users" {
  name = "ad?in"
}
`

const testAccDataSourceUsersConfigNameFilterNG = `
data "redash_users" "users" {
  name = "a*d?min"
}
`

const testAccDataSourceUsersConfigEmailFilter = `
data "redash_users" "users" {
  email = "*@example.com"
}
`

const testAccDataSourceUsersConfigEmailFilter2 = `
data "redash_users" "users" {
  email = "ad?in@example.com"
}
`

const testAccDataSourceUsersConfigEmailFilterNG = `
data "redash_users" "users" {
  email = "a*d?min@example.com"
}
`

const testAccDataSourceUsersConfigComplexFilter = `
data "redash_users" "users" {
	name  = "a*d?in"
	email = "a*d?in@example.com"
}
`

const testAccDataSourceUsersConfigComplexFilterNG = `
data "redash_users" "users" {
	name  = "a*d?min"
	email = "a*d?in@example.com"
}
`

const testAccDataSourceUsersConfigComplexFilterNG2 = `
data "redash_users" "users" {
	name  = "a*d?in"
	email = "a*d?min@example.com"
}
`
