package test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccQuery_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testAccQueryConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("redash_query.my_query", "name", "my-query"),
					resource.TestCheckResourceAttr("redash_query.my_query", "description", "my-query desc"),
					resource.TestCheckResourceAttr("redash_query.my_query", "query", "select 1"),
					resource.TestCheckNoResourceAttr("redash_query.my_query", "interval"),
					resource.TestCheckNoResourceAttr("redash_query.my_query", "tags"),
					resource.TestCheckResourceAttr("redash_query.my_query", "published", "false"),
				),
			},
			{
				Config: testAccQueryConfigBasic2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("redash_query.my_query", "name", "my-query2"),
					resource.TestCheckResourceAttr("redash_query.my_query", "description", "my-query desc2"),
					resource.TestCheckResourceAttr("redash_query.my_query", "query", "select 2"),
					resource.TestCheckTypeSetElemNestedAttrs("redash_query.my_query", "schedule.*", map[string]string{
						"interval": "600",
					}),
					resource.TestCheckNoResourceAttr("redash_query.my_query", "tags"),
					resource.TestCheckResourceAttr("redash_query.my_query", "published", "false"),
				),
			},
			{
				Config: testAccQueryConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("redash_query.my_query", "name", "my-query"),
					resource.TestCheckResourceAttr("redash_query.my_query", "description", "my-query desc"),
					resource.TestCheckResourceAttr("redash_query.my_query", "query", "select 1"),
					resource.TestCheckNoResourceAttr("redash_query.my_query", "interval"),
					resource.TestCheckNoResourceAttr("redash_query.my_query", "tags"),
				),
			},
			{
				Config: testAccQueryConfigWithTags,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("redash_query.my_query", "name", "my-query"),
					resource.TestCheckResourceAttr("redash_query.my_query", "description", "my-query desc"),
					resource.TestCheckResourceAttr("redash_query.my_query", "query", "select 1"),
					resource.TestCheckNoResourceAttr("redash_query.my_query", "interval"),
					resource.TestCheckResourceAttr("redash_query.my_query", "tags.#", "2"),
					resource.TestCheckResourceAttr("redash_query.my_query", "tags.0", "foo"),
					resource.TestCheckResourceAttr("redash_query.my_query", "tags.1", "bar"),
					resource.TestCheckResourceAttr("redash_query.my_query2", "tags.#", "3"),
					resource.TestCheckResourceAttr("redash_query.my_query2", "tags.0", "bar"),
					resource.TestCheckResourceAttr("redash_query.my_query2", "tags.1", "zoo"),
					resource.TestCheckResourceAttr("redash_query.my_query2", "tags.2", "baz"),
					resource.TestCheckResourceAttr("redash_query.my_query", "published", "false"),
				),
			},
			{
				Config: testAccQueryConfigWithTags2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("redash_query.my_query", "name", "my-query"),
					resource.TestCheckResourceAttr("redash_query.my_query", "description", "my-query desc"),
					resource.TestCheckResourceAttr("redash_query.my_query", "query", "select 1"),
					resource.TestCheckNoResourceAttr("redash_query.my_query", "interval"),
					resource.TestCheckResourceAttr("redash_query.my_query", "tags.#", "1"),
					resource.TestCheckResourceAttr("redash_query.my_query", "tags.0", "bar"),
					resource.TestCheckResourceAttr("redash_query.my_query2", "tags.#", "0"),
					resource.TestCheckResourceAttr("redash_query.my_query", "published", "false"),
				),
			},
			{
				Config: testAccQueryConfigWithPublish,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("redash_query.my_query", "name", "my-query"),
					resource.TestCheckResourceAttr("redash_query.my_query", "description", "my-query desc"),
					resource.TestCheckResourceAttr("redash_query.my_query", "query", "select 1"),
					resource.TestCheckNoResourceAttr("redash_query.my_query", "interval"),
					resource.TestCheckResourceAttr("redash_query.my_query", "tags.#", "1"),
					resource.TestCheckResourceAttr("redash_query.my_query", "tags.0", "bar"),
					resource.TestCheckResourceAttr("redash_query.my_query", "published", "true"),
				),
			},
			{
				Config: testAccQueryConfigWithUnpublish,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("redash_query.my_query", "name", "my-query"),
					resource.TestCheckResourceAttr("redash_query.my_query", "description", "my-query desc"),
					resource.TestCheckResourceAttr("redash_query.my_query", "query", "select 1"),
					resource.TestCheckNoResourceAttr("redash_query.my_query", "interval"),
					resource.TestCheckResourceAttr("redash_query.my_query", "tags.#", "1"),
					resource.TestCheckResourceAttr("redash_query.my_query", "tags.0", "bar"),
					resource.TestCheckResourceAttr("redash_query.my_query", "published", "false"),
				),
			},
		},
	})
}

const testAccQueryConfigBasic = testAccDataSourceConfigBasicPg + `
resource "redash_query" "my_query" {
	data_source_id = redash_data_source.my_data_source.id
	name           = "my-query"
	description    = "my-query desc"
	query          = "select 1"
}
`

const testAccQueryConfigBasic2 = testAccDataSourceConfigBasicPg + `
resource "redash_query" "my_query" {
	data_source_id = redash_data_source.my_data_source.id
	name           = "my-query2"
	description    = "my-query desc2"
	query          = "select 2"
	schedule {
		interval = 600
	}
}
`

const testAccQueryConfigWithTags = testAccDataSourceConfigBasicPg + `
resource "redash_query" "my_query" {
	data_source_id = redash_data_source.my_data_source.id
	name           = "my-query"
	description    = "my-query desc"
	query          = "select 1"
	tags           = ["foo", "bar"]
}

resource "redash_query" "my_query2" {
	data_source_id = redash_data_source.my_data_source.id
	name           = "my-query2"
	description    = "my-query2 desc"
	query          = "select 2"
	tags           = ["bar", "zoo", "baz"]
}
`

const testAccQueryConfigWithTags2 = testAccDataSourceConfigBasicPg + `
resource "redash_query" "my_query" {
	data_source_id = redash_data_source.my_data_source.id
	name           = "my-query"
	description    = "my-query desc"
	query          = "select 1"
	tags           = ["bar"]
}

resource "redash_query" "my_query2" {
	data_source_id = redash_data_source.my_data_source.id
	name           = "my-query2"
	description    = "my-query2 desc"
	query          = "select 2"
}
`

const testAccQueryConfigWithPublish = testAccDataSourceConfigBasicPg + `
resource "redash_query" "my_query" {
	data_source_id = redash_data_source.my_data_source.id
	name           = "my-query"
	description    = "my-query desc"
	query          = "select 1"
	tags           = ["bar"]
	published      = true
}
`

const testAccQueryConfigWithUnpublish = testAccDataSourceConfigBasicPg + `
resource "redash_query" "my_query" {
	data_source_id = redash_data_source.my_data_source.id
	name           = "my-query"
	description    = "my-query desc"
	query          = "select 1"
	tags           = ["bar"]
	published      = false
}
`
