package test

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	redashgo "github.com/winebarrel/redash-go/v2"
)

func TestAccGroupMember_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testAccGroupConfigBasic + testAccDataSourceUserConfigName,
			},
			{
				Config: testAccGroupMemberBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGroupMember("redash_group_member.my_member"),
				),
			},
			{
				ResourceName:      "redash_group_member.my_member",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// A user that is not in the group cannot be imported.
				ResourceName: "redash_group_member.my_member",
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources["redash_group_member.my_member"]
					if !ok {
						return "", fmt.Errorf("Not found: redash_group_member.my_member")
					}

					return rs.Primary.Attributes["group_id"] + "/999999", nil
				},
				ExpectError: regexp.MustCompile(`Cannot import non-existent remote object`),
			},
		},
	})
}

func TestAccGroupMember_removedOutOfBand(t *testing.T) {
	const addr = "redash_group_member.my_member"

	var groupId, userId int

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testAccGroupMemberBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGroupMember(addr),
					testAccCaptureGroupMember(addr, &groupId, &userId),
				),
			},
			{
				PreConfig: func() {
					if err := testAccNewClient(t).RemoveGroupMember(context.Background(), groupId, userId); err != nil {
						t.Fatal(err)
					}
				},
				// The membership is gone, so the plan must show a create.
				Config:             testAccGroupMemberBasic,
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
			{
				Config: testAccGroupMemberBasic,
				Check:  testAccCheckGroupMember(addr),
			},
		},
	})
}

func TestAccGroupMember_groupDeletedOutOfBand(t *testing.T) {
	const addr = "redash_group_member.my_member"

	// Redash answers GET /api/groups/:id with 500 for a missing group, so the
	// group is not managed by Terraform here.
	client := testAccNewClient(t)
	group, err := client.CreateGroup(context.Background(), &redashgo.CreateGroupInput{Name: "my-group-deleted"})
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() { client.DeleteGroup(context.Background(), group.ID) }) //nolint:errcheck

	config := testAccDataSourceUserConfigName + fmt.Sprintf(`
resource "redash_group_member" "my_member" {
	group_id = %d
	user_id  = data.redash_user.admin.id
}
`, group.ID)

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: config,
				Check:  testAccCheckGroupMember(addr),
			},
			{
				PreConfig: func() {
					if err := client.DeleteGroup(context.Background(), group.ID); err != nil {
						t.Fatal(err)
					}
				},
				// The group is gone, so the plan must show a create, not an error.
				Config:             config,
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccNewClient(t *testing.T) *redashgo.Client {
	t.Helper()

	client, err := redashgo.NewClient(testAccRedashURL, testAccRedashAPIKey)
	if err != nil {
		t.Fatal(err)
	}

	return client
}

// testAccCaptureGroupMember saves the group and user IDs of the member so that
// a later PreConfig can change Redash behind Terraform's back.
func testAccCaptureGroupMember(addr string, groupId, userId *int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[addr]
		if !ok {
			return fmt.Errorf("Not found: %s", addr)
		}

		var err error

		if *groupId, err = strconv.Atoi(rs.Primary.Attributes["group_id"]); err != nil {
			return err
		}

		*userId, err = strconv.Atoi(rs.Primary.Attributes["user_id"])

		return err
	}
}

const testAccGroupMemberBasic = testAccGroupConfigBasic + testAccDataSourceUserConfigName + `
resource "redash_group_member" "my_member" {
	group_id = redash_group.my_group.id
	user_id  = data.redash_user.admin.id
}
`

func testAccCheckGroupMember(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]

		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Group Member (%s) ID is not set", resourceName)
		}

		groupId := rs.Primary.Attributes["group_id"]

		if !regexp.MustCompile(`^\d+$`).MatchString(groupId) {
			return fmt.Errorf("group_id must be number, got: %s", groupId)
		}

		userId := rs.Primary.Attributes["user_id"]

		if !regexp.MustCompile(`^\d+$`).MatchString(userId) {
			return fmt.Errorf("user_id must be number, got: %s", userId)
		}

		return nil
	}
}
