package test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSoureceQuery_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testAccQueryConfigBasic2,
			},
			{
				Config: testAccDataSoureceQueryConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.redash_query.my_query", "name", "my-query2"),
					resource.TestCheckResourceAttr("data.redash_query.my_query", "description", "my-query desc2"),
					resource.TestCheckResourceAttr("data.redash_query.my_query", "query", "select 2"),
					resource.TestCheckTypeSetElemNestedAttrs("data.redash_query.my_query", "schedule.*", map[string]string{
						"interval": "600",
					}),
					resource.TestCheckNoResourceAttr("data.redash_query.my_query", "tags"),
				),
			},
			{
				Config: testAccQueryConfigWithTags,
			},
			{
				Config: testAccDataSoureceQueryConfigWithTags,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.redash_query.my_query", "name", "my-query2"),
					resource.TestCheckResourceAttr("data.redash_query.my_query", "description", "my-query2 desc"),
					resource.TestCheckResourceAttr("data.redash_query.my_query", "query", "select 2"),
					resource.TestCheckNoResourceAttr("data.redash_query.my_query", "schedule"),
					resource.TestCheckResourceAttr("data.redash_query.my_query", "tags.#", "3"),
					resource.TestCheckResourceAttr("data.redash_query.my_query", "tags.0", "bar"),
					resource.TestCheckResourceAttr("data.redash_query.my_query", "tags.1", "zoo"),
					resource.TestCheckResourceAttr("data.redash_query.my_query", "tags.2", "baz"),
				),
			},
		},
	})
}

const testAccDataSoureceQueryConfigBasic = testAccQueryConfigBasic2 + `
data "redash_query" "my_query" {
  name = "my-query2"
}
`

const testAccDataSoureceQueryConfigWithTags = testAccQueryConfigWithTags + `
data "redash_query" "my_query" {
  name = "my-query2"
}
`
