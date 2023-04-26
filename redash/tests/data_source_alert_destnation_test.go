package test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSoureceAlertDestnation_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testAccAlertSestinationConfigBasic,
			},
			{
				Config: TestAccDataSoureceAlertDestnationConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.redash_alert_destination.my_dest", "name", "my-dest"),
					resource.TestCheckResourceAttr("data.redash_alert_destination.my_dest", "type", "email"),
					resource.TestCheckResourceAttr("data.redash_alert_destination.my_dest", "options", `{"addresses":"foo@example.com"}`),
				),
			},
		},
	})
}

const TestAccDataSoureceAlertDestnationConfigBasic = testAccAlertSestinationConfigBasic + `
data "redash_alert_destination" "my_dest" {
  name = "my-dest"
}
`
