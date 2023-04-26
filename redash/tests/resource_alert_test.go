package test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlert_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testAccAlertConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("redash_alert.my_alert", "name", "my-alert"),
					resource.TestCheckResourceAttr("redash_alert.my_alert", "rearm", "0"),
					resource.TestCheckTypeSetElemNestedAttrs("redash_alert.my_alert", "options.*", map[string]string{
						"column":         "?column?",
						"op":             ">",
						"value":          "1",
						"custom_subject": "",
						"custom_body":    "",
						"template":       "",
					}),
					resource.TestCheckResourceAttr("redash_alert.my_alert", "muted", "false"),
				),
			},
			{
				Config: testAccAlertConfigRearm,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("redash_alert.my_alert", "name", "my-alert"),
					resource.TestCheckResourceAttr("redash_alert.my_alert", "rearm", "300"),
					resource.TestCheckTypeSetElemNestedAttrs("redash_alert.my_alert", "options.*", map[string]string{
						"column":         "?column?",
						"op":             ">",
						"value":          "1",
						"custom_subject": "",
						"custom_body":    "",
						"template":       "",
					}),
					resource.TestCheckResourceAttr("redash_alert.my_alert", "muted", "false"),
				),
			},
			{
				Config: testAccAlertConfigTemplate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("redash_alert.my_alert", "name", "my-alert"),
					resource.TestCheckResourceAttr("redash_alert.my_alert", "rearm", "0"),
					resource.TestCheckTypeSetElemNestedAttrs("redash_alert.my_alert", "options.*", map[string]string{
						"column":         "?column?",
						"op":             "<",
						"value":          "3",
						"custom_subject": "subject",
						"custom_body":    "body",
						"template":       "",
					}),
					resource.TestCheckResourceAttr("redash_alert.my_alert", "muted", "false"),
				),
			},
			{
				Config: testAccAlertConfigMuted,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("redash_alert.my_alert", "name", "my-alert"),
					resource.TestCheckResourceAttr("redash_alert.my_alert", "rearm", "0"),
					resource.TestCheckTypeSetElemNestedAttrs("redash_alert.my_alert", "options.*", map[string]string{
						"column":         "?column?",
						"op":             ">",
						"value":          "1",
						"custom_subject": "",
						"custom_body":    "",
						"template":       "",
					}),
					resource.TestCheckResourceAttr("redash_alert.my_alert", "muted", "true"),
				),
			},
			{
				Config: testAccAlertConfigUnmuted,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("redash_alert.my_alert", "name", "my-alert"),
					resource.TestCheckResourceAttr("redash_alert.my_alert", "rearm", "0"),
					resource.TestCheckTypeSetElemNestedAttrs("redash_alert.my_alert", "options.*", map[string]string{
						"column":         "?column?",
						"op":             ">",
						"value":          "1",
						"custom_subject": "",
						"custom_body":    "",
						"template":       "",
					}),
					resource.TestCheckResourceAttr("redash_alert.my_alert", "muted", "false"),
				),
			},
		},
	})
}

const testAccAlertConfigBasic = testAccQueryConfigBasic + `
resource "redash_alert" "my_alert" {
  name     = "my-alert"
  query_id = redash_query.my_query.id
  options {
    column = "?column?"
    op     = ">"
    value  = 1
  }
}
`

const testAccAlertConfigRearm = testAccQueryConfigBasic + `
resource "redash_alert" "my_alert" {
  name     = "my-alert"
  query_id = redash_query.my_query.id
  options {
    column = "?column?"
    op     = ">"
    value  = 1
  }
	rearm = 300
}
`

const testAccAlertConfigTemplate = testAccQueryConfigBasic + `
resource "redash_alert" "my_alert" {
  name     = "my-alert"
  query_id = redash_query.my_query.id
  options {
    column         = "?column?"
    op             = "<"
    value          = 3
		custom_subject = "subject"
		custom_body    = "body"
  }
}
`

const testAccAlertConfigMuted = testAccQueryConfigBasic + `
resource "redash_alert" "my_alert" {
  name     = "my-alert"
  query_id = redash_query.my_query.id
  options {
    column = "?column?"
    op     = ">"
    value  = 1
  }
	muted = true
}
`

const testAccAlertConfigUnmuted = testAccQueryConfigBasic + `
resource "redash_alert" "my_alert" {
  name     = "my-alert"
  query_id = redash_query.my_query.id
  options {
    column = "?column?"
    op     = ">"
    value  = 1
  }
	muted = false
}
`
