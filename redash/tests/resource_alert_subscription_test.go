package test

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	redashgo "github.com/winebarrel/redash-go/v2"
)

func TestAccAlertSubscription_basic(t *testing.T) {
	var id string

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testAccAlertConfigBasic,
			},
			{
				Config: testAccAlertSubscriptionConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlertSubscription("redash_alert_subscription.my_subs"),
					testAccCaptureID("redash_alert_subscription.my_subs", &id),
				),
			},
			{
				ResourceName:      "redash_alert_subscription.my_subs",
				ImportState:       true,
				ImportStateIdFunc: testAccAlertSubscriptionImportID("redash_alert_subscription.my_subs"),
				ImportStateVerify: true,
			},
			{
				PreConfig:          testAccRemoveAlertSubscription(t, &id),
				Config:             testAccAlertSubscriptionConfigBasic,
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
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

func testAccAlertSubscriptionImportID(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]

		if !ok {
			return "", fmt.Errorf("Not found: %s", resourceName)
		}

		return rs.Primary.Attributes["alert_id"] + "/" + rs.Primary.Attributes["alert_destination_id"], nil
	}
}

func testAccRemoveAlertSubscription(t *testing.T, id *string) func() {
	return func() {
		alertIdStr, subsIdStr, _ := strings.Cut(*id, "/")
		alertId, _ := strconv.Atoi(alertIdStr)
		subsId, _ := strconv.Atoi(subsIdStr)
		client := testAccProvider.Meta().(*redashgo.Client)

		err := client.RemoveAlertSubscription(t.Context(), alertId, subsId)
		if err != nil {
			t.Fatal(err)
		}
	}
}
