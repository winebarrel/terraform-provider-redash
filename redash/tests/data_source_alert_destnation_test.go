package test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceAlertDestnation_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testAccAlertSestinationConfigBasic,
			},
			{
				Config: TestAccDataSourceAlertDestnationConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.redash_alert_destination.my_dest", "name", "my-dest"),
					resource.TestCheckResourceAttr("data.redash_alert_destination.my_dest", "type", "email"),
					resource.TestCheckResourceAttr("data.redash_alert_destination.my_dest", "options", `{"addresses":"foo@example.com"}`),
				),
			},
		},
	})
}

const TestAccDataSourceAlertDestnationConfigBasic = testAccAlertSestinationConfigBasic + `
data "redash_alert_destination" "my_dest" {
  name = "my-dest"
}
`
