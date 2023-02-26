package redash_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSoureceDataSourece_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceConfigBasicPg,
			},
			{
				Config: testAccDataSoureceDataSourceConfigBasicPg,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.redash_data_source.my_data_source", "name", "my-data-source"),
					resource.TestCheckResourceAttr("data.redash_data_source.my_data_source", "type", "pg"),
					resource.TestCheckResourceAttr("data.redash_data_source.my_data_source", "options", `{"dbname":"postgres","host":"postgres","port":5432,"user":"postgres"}`),
				),
			},
		},
	})
}

const testAccDataSoureceDataSourceConfigBasicPg = testAccDataSourceConfigBasicPg + `
data "redash_data_source" "my_data_source" {
  name = "my-data-source"
}
`
