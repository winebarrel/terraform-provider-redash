package test

import (
	"context"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccGroup_basic(t *testing.T) {
	var id string

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testAccGroupConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("redash_group.my_group", "name", "my-group"),
					func(s *terraform.State) error {
						id = s.RootModule().Resources["redash_group.my_group"].Primary.ID
						return nil
					},
				),
			},
			{
				Config: testAccGroupConfigBasic2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("redash_group.my_group", "name", "my-group2"),
					// renamed in place
					resource.TestCheckResourceAttrPtr("redash_group.my_group", "id", &id),
				),
			},
			{
				ResourceName:      "redash_group.my_group",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccGroup_deletedOutOfBand(t *testing.T) {
	var id int

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testAccGroupConfigBasic,
				Check: func(s *terraform.State) error {
					var err error
					id, err = strconv.Atoi(s.RootModule().Resources["redash_group.my_group"].Primary.ID)
					return err
				},
			},
			{
				PreConfig: func() {
					if err := testAccNewClient(t).DeleteGroup(context.Background(), id); err != nil {
						t.Fatal(err)
					}
				},
				// The group is gone, so the plan must show a create, not an error.
				Config:             testAccGroupConfigBasic,
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
			{
				Config: testAccGroupConfigBasic,
				Check:  resource.TestCheckResourceAttr("redash_group.my_group", "name", "my-group"),
			},
		},
	})
}

const testAccGroupConfigBasic = `
resource "redash_group" "my_group" {
	name = "my-group"
}
`

const testAccGroupConfigBasic2 = `
resource "redash_group" "my_group" {
	name = "my-group2"
}
`
