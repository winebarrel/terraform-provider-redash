package test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataAlert_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAlertBasic0,
			},
			{
				Config: testAccDataSourceAlertBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.redash_alert.my_alert", "name", "my-alert"),
					resource.TestCheckResourceAttr("data.redash_alert.my_alert", "rearm", "300"),
					resource.TestCheckResourceAttr("data.redash_alert.my_alert", "muted", "true"),
					resource.TestCheckTypeSetElemNestedAttrs("data.redash_alert.my_alert", "options.*", map[string]string{
						"column":         "?column?",
						"op":             "<",
						"value":          "1.5",
						"selector":       "first",
						"custom_subject": "subject",
						"custom_body":    "body",
						"template":       "",
					}),
				),
			},
		},
	})
}

const testAccDataSourceAlertBasic0 = testAccQueryConfigBasic + `
resource "redash_alert" "my_alert" {
  name     = "my-alert"
  query_id = redash_query.my_query.id
  options {
    column         = "?column?"
    op             = "<"
    value          = 1.5
		custom_subject = "subject"
		custom_body    = "body"
  }
	rearm = 300
	muted = true
}
`

const testAccDataSourceAlertBasic = testAccDataSourceAlertBasic0 + `
data "redash_alert" "my_alert" {
  name = "my-alert"
}
`
