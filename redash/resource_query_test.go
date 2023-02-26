package redash_test

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
				),
			},
			{
				Config: testAccQueryConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("redash_query.my_query", "name", "my-query"),
					resource.TestCheckResourceAttr("redash_query.my_query", "description", "my-query desc"),
					resource.TestCheckResourceAttr("redash_query.my_query", "query", "select 1"),
					resource.TestCheckNoResourceAttr("redash_query.my_query", "interval"),
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
