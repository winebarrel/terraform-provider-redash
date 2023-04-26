package test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourece_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceConfigBasicPg,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("redash_data_source.my_data_source", "name", "my-data-source"),
					resource.TestCheckResourceAttr("redash_data_source.my_data_source", "type", "pg"),
					resource.TestCheckResourceAttr("redash_data_source.my_data_source", "options", `{"dbname":"postgres","host":"postgres","port":5432,"user":"postgres"}`),
				),
			},
			{
				Config: testAccDataSourceConfigBasicPg2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("redash_data_source.my_data_source", "name", "my-data-source2"),
					resource.TestCheckResourceAttr("redash_data_source.my_data_source", "type", "pg"),
					resource.TestCheckResourceAttr("redash_data_source.my_data_source", "options", `{"dbname":"postgres2","host":"postgres2","port":5433,"user":"postgres2"}`),
				),
			},
		},
	})
}

const testAccDataSourceConfigBasicPg = `
resource "redash_data_source" "my_data_source" {
  name = "my-data-source"
  type = "pg"
  options = jsonencode({
    dbname = "postgres"
    host   = "postgres"
    port   = 5432
    user   = "postgres"
  })
}
`

const testAccDataSourceConfigBasicPg2 = `
resource "redash_data_source" "my_data_source" {
  name = "my-data-source2"
  type = "pg"
  options = jsonencode({
    dbname = "postgres2"
    host   = "postgres2"
    port   = 5433
    user   = "postgres2"
  })
}
`
