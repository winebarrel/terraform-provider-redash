package redash_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlertDestinations_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testAccAlertSestinationConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("redash_alert_destination.my_dest", "name", "my-dest"),
					resource.TestCheckResourceAttr("redash_alert_destination.my_dest", "type", "email"),
					resource.TestCheckResourceAttr("redash_alert_destination.my_dest", "options", `{"addresses":"foo@example.com"}`),
				),
			},
			{
				Config: testAccAlertSestinationConfigBasic2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("redash_alert_destination.my_dest", "name", "my-dest2"),
					resource.TestCheckResourceAttr("redash_alert_destination.my_dest", "type", "email"),
					resource.TestCheckResourceAttr("redash_alert_destination.my_dest", "options", `{"addresses":"bar@example.com"}`),
				),
			},
		},
	})
}

const testAccAlertSestinationConfigBasic = `
resource "redash_alert_destination" "my_dest" {
  name = "my-dest"
  type = "email"
  options = jsonencode({
    addresses = "foo@example.com"
  })
}
`

const testAccAlertSestinationConfigBasic2 = `
resource "redash_alert_destination" "my_dest" {
  name = "my-dest2"
  type = "email"
  options = jsonencode({
    addresses = "bar@example.com"
  })
}
`
