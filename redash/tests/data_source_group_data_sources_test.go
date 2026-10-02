package test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceGroupDataSources_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGroupDataSourcesConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.redash_group_data_sources.my_group", "data_sources.#", "1"),
					resource.TestCheckResourceAttrPair(
						"data.redash_group_data_sources.my_group", "data_sources.0.id",
						"redash_data_source.my_data_source", "id",
					),
					resource.TestCheckResourceAttr("data.redash_group_data_sources.my_group", "data_sources.0.name", "my-data-source"),
					resource.TestCheckResourceAttr("data.redash_group_data_sources.my_group", "data_sources.0.view_only", "false"),
				),
			},
		},
	})
}

const testAccDataSourceGroupDataSourcesConfig = testAccGroupSubscriptionConfigBasic + `
data "redash_group_data_sources" "my_group" {
  group_id = redash_group.my_group.id
}
`
