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

func TestAccGroupMember_basic(t *testing.T) {
	var id string

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
					testAccCaptureID("redash_group_member.my_member", &id),
				),
			},
			{
				ResourceName:      "redash_group_member.my_member",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				PreConfig:          testAccRemoveGroupMember(t, &id),
				Config:             testAccGroupMemberBasic,
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
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

func testAccRemoveGroupMember(t *testing.T, id *string) func() {
	return func() {
		groupIdStr, memberIdStr, _ := strings.Cut(*id, "/")
		groupId, _ := strconv.Atoi(groupIdStr)
		memberId, _ := strconv.Atoi(memberIdStr)
		client := testAccProvider.Meta().(*redashgo.Client)

		err := client.RemoveGroupMember(t.Context(), groupId, memberId)
		if err != nil {
			t.Fatal(err)
		}
	}
}
