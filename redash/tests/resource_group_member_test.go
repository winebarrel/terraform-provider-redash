package test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccGroupMember_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testAccGroupConfigBasic + testAccDataSoureceUserConfigName,
			},
			{
				Config: testAccGroupMemberBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGroupMember("redash_group_member.my_member"),
				),
			},
		},
	})
}

const testAccGroupMemberBasic = testAccGroupConfigBasic + testAccDataSoureceUserConfigName + `
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
