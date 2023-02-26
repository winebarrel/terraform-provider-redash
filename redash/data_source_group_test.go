package redash_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSoureceGroup_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testAccGroupConfigBasic,
			},
			{
				Config: testAccDataSoureceGroupConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.redash_group.my_group", "name", "my-group"),
				),
			},
		},
	})
}

const testAccDataSoureceGroupConfigBasic = testAccGroupConfigBasic + `
data "redash_group" "my_group" {
  name = "my-group"
}
`
