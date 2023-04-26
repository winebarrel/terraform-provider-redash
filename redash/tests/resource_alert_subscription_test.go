package test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAlertSubscription_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testAccAlertConfigBasic,
			},
			{
				Config: testAccAlertSubscriptionConfigBasic,
				Check:  testAccCheckAlertSubscription("redash_alert_subscription.my_subs"),
			},
		},
	})
}

const testAccAlertSubscriptionConfigBasic = testAccAlertConfigBasic + `
resource "redash_alert_destination" "my_dest" {
  name = "my-dest"
  type = "email"
  options = jsonencode({
    addresses = "foo@example.com"
  })
}

resource "redash_alert_subscription" "my_subs" {
	alert_id             = redash_alert.my_alert.id
	alert_destination_id = redash_alert_destination.my_dest.id
}
`

func testAccCheckAlertSubscription(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]

		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Alert Subscription (%s) ID is not set", resourceName)
		}

		alertId := rs.Primary.Attributes["alert_id"]

		if !regexp.MustCompile(`^\d+$`).MatchString(alertId) {
			return fmt.Errorf("alert_id must be number, got: %s", alertId)
		}

		destId := rs.Primary.Attributes["alert_destination_id"]

		if !regexp.MustCompile(`^\d+$`).MatchString(destId) {
			return fmt.Errorf("alert_destination_id must be number, got: %s", destId)
		}

		return nil
	}
}
