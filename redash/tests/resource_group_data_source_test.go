package test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccGroupDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceConfigBasicPg + testAccGroupConfigBasic,
			},
			{
				Config: testAccGroupSubscriptionConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGroupDataSource("redash_group_data_source.my_gds"),
					resource.TestCheckResourceAttr("redash_group_data_source.my_gds", "view_only", "false"),
				),
			},
			{
				Config: testAccGroupSubscriptionConfigViewOnly,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGroupDataSource("redash_group_data_source.my_gds"),
					resource.TestCheckResourceAttr("redash_group_data_source.my_gds", "view_only", "true"),
				),
			},
		},
	})
}

const testAccGroupSubscriptionConfigBasic = testAccDataSourceConfigBasicPg + testAccGroupConfigBasic + `
resource "redash_group_data_source" "my_gds" {
	group_id       = redash_group.my_group.id
	data_source_id = redash_data_source.my_data_source.id
}
`

const testAccGroupSubscriptionConfigViewOnly = testAccDataSourceConfigBasicPg + testAccGroupConfigBasic + `
resource "redash_group_data_source" "my_gds" {
	group_id       = redash_group.my_group.id
	data_source_id = redash_data_source.my_data_source.id
	view_only      = true
}
`

func testAccCheckGroupDataSource(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]

		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Group Data Source (%s) ID is not set", resourceName)
		}

		groupId := rs.Primary.Attributes["group_id"]

		if !regexp.MustCompile(`^\d+$`).MatchString(groupId) {
			return fmt.Errorf("group_id must be number, got: %s", groupId)
		}

		dsId := rs.Primary.Attributes["data_source_id"]

		if !regexp.MustCompile(`^\d+$`).MatchString(dsId) {
			return fmt.Errorf("data_source_id must be number, got: %s", dsId)
		}

		return nil
	}
}
