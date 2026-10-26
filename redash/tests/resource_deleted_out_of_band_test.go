package test

import (
	"context"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	redashgo "github.com/winebarrel/redash-go/v2"
)

// testAccDeletedOutOfBand applies config, deletes the object behind address
// through the API, and expects the next plan to create it again instead of
// failing on refresh.
func testAccDeletedOutOfBand(t *testing.T, config, address string, del func(*redashgo.Client, string) error) {
	var id string

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: func(s *terraform.State) error {
					id = s.RootModule().Resources[address].Primary.ID
					return nil
				},
			},
			{
				PreConfig: func() {
					if err := del(testAccNewClient(t), id); err != nil {
						t.Fatal(err)
					}
				},
				Config:             config,
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
			{
				Config: config,
				Check:  resource.TestCheckResourceAttrSet(address, "id"),
			},
		},
	})
}

func atoi(t *testing.T, s string) int {
	t.Helper()

	n, err := strconv.Atoi(s)
	if err != nil {
		t.Fatal(err)
	}

	return n
}

func TestAccDataSource_deletedOutOfBand(t *testing.T) {
	testAccDeletedOutOfBand(t, testAccDataSourceConfigBasicPg, "redash_data_source.my_data_source",
		func(c *redashgo.Client, id string) error {
			return c.DeleteDataSource(context.Background(), atoi(t, id))
		})
}

func TestAccGroupDataSources_groupDeletedOutOfBand(t *testing.T) {
	testAccDeletedOutOfBand(t, testAccGroupDataSourcesOne, "redash_group_data_sources.my_group",
		func(c *redashgo.Client, id string) error {
			return c.DeleteGroup(context.Background(), atoi(t, id))
		})
}

func TestAccGroupDataSource_groupDeletedOutOfBand(t *testing.T) {
	testAccDeletedOutOfBand(t, testAccGroupSubscriptionConfigBasic, "redash_group_data_source.my_gds",
		func(c *redashgo.Client, id string) error {
			groupId, _, _ := strings.Cut(id, "/")
			return c.DeleteGroup(context.Background(), atoi(t, groupId))
		})
}

func TestAccUser_deletedOutOfBand(t *testing.T) {
	testAccDeletedOutOfBand(t, testAccUserConfigBasic, "redash_user.my_user",
		func(c *redashgo.Client, id string) error {
			return c.DeleteUser(context.Background(), atoi(t, id))
		})
}

func TestAccAlertDestination_deletedOutOfBand(t *testing.T) {
	testAccDeletedOutOfBand(t, testAccAlertDestinationConfigBasic, "redash_alert_destination.my_dest",
		func(c *redashgo.Client, id string) error {
			return c.DeleteDestination(context.Background(), atoi(t, id))
		})
}
