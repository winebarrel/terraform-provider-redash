package test

import (
	"errors"
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccParamQuery_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testAccParamQueryConfigTextNum1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("redash_query.my_query", "name", "my-query"),
					resource.TestCheckResourceAttr("redash_query.my_query", "description", "my-query desc"),
					resource.TestCheckResourceAttr("redash_query.my_query", "query", "select '{{ txt }}','{{ num }}'"),
					resource.TestCheckResourceAttr("redash_query.my_query", "options.0.parameter.#", "2"),
					resource.TestCheckResourceAttr("redash_query.my_query", "options.0.parameter.0.title", "tnum"),
					resource.TestCheckResourceAttr("redash_query.my_query", "options.0.parameter.0.name", "num"),
					resource.TestCheckResourceAttr("redash_query.my_query", "options.0.parameter.0.type", "number"),
					resource.TestCheckResourceAttr("redash_query.my_query", "options.0.parameter.0.value", "100"),
					resource.TestCheckResourceAttr("redash_query.my_query", "options.0.parameter.1.title", "ttext"),
					resource.TestCheckResourceAttr("redash_query.my_query", "options.0.parameter.1.name", "txt"),
					resource.TestCheckResourceAttr("redash_query.my_query", "options.0.parameter.1.type", "text"),
					resource.TestCheckResourceAttr("redash_query.my_query", "options.0.parameter.1.value", "hello"),
				),
			},
			{
				Config: testAccParamQueryConfigTextNum2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("redash_query.my_query", "name", "my-query"),
					resource.TestCheckResourceAttr("redash_query.my_query", "description", "my-query desc"),
					resource.TestCheckResourceAttr("redash_query.my_query", "query", "select '{{ txt }}','{{ txt2 }}','{{ num2 }}'"),
					resource.TestCheckResourceAttr("redash_query.my_query", "options.0.parameter.#", "3"),
					resource.TestCheckResourceAttr("redash_query.my_query", "options.0.parameter.0.title", "tnum2"),
					resource.TestCheckResourceAttr("redash_query.my_query", "options.0.parameter.0.name", "num2"),
					resource.TestCheckResourceAttr("redash_query.my_query", "options.0.parameter.0.type", "number"),
					resource.TestCheckResourceAttr("redash_query.my_query", "options.0.parameter.0.value", "1002"),
					resource.TestCheckResourceAttr("redash_query.my_query", "options.0.parameter.1.title", "ttext"),
					resource.TestCheckResourceAttr("redash_query.my_query", "options.0.parameter.1.name", "txt"),
					resource.TestCheckResourceAttr("redash_query.my_query", "options.0.parameter.1.type", "text"),
					resource.TestCheckResourceAttr("redash_query.my_query", "options.0.parameter.1.value", "hello"),
					resource.TestCheckResourceAttr("redash_query.my_query", "options.0.parameter.2.title", ""),
					resource.TestCheckResourceAttr("redash_query.my_query", "options.0.parameter.2.name", "txt2"),
					resource.TestCheckResourceAttr("redash_query.my_query", "options.0.parameter.2.type", "text"),
					resource.TestCheckResourceAttr("redash_query.my_query", "options.0.parameter.2.value", "hello2"),
				),
			},
			{
				Config: testAccParamQueryConfigTextNum3,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("redash_query.my_query", "name", "my-query"),
					resource.TestCheckResourceAttr("redash_query.my_query", "description", "my-query desc"),
					resource.TestCheckResourceAttr("redash_query.my_query", "query", "select '{{ txt2 }}','{{ num2 }}'"),
					resource.TestCheckResourceAttr("redash_query.my_query", "options.0.parameter.#", "2"),
					resource.TestCheckResourceAttr("redash_query.my_query", "options.0.parameter.0.title", "tnum2"),
					resource.TestCheckResourceAttr("redash_query.my_query", "options.0.parameter.0.name", "num2"),
					resource.TestCheckResourceAttr("redash_query.my_query", "options.0.parameter.0.type", "number"),
					resource.TestCheckResourceAttr("redash_query.my_query", "options.0.parameter.0.value", "1002"),
					resource.TestCheckResourceAttr("redash_query.my_query", "options.0.parameter.1.name", "txt2"),
					resource.TestCheckResourceAttr("redash_query.my_query", "options.0.parameter.1.type", "text"),
					resource.TestCheckResourceAttr("redash_query.my_query", "options.0.parameter.1.value", "hello2"),
				),
			},
			{
				Config: testAccParamQueryConfigRegex,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("redash_query.my_query", "name", "my-query"),
					resource.TestCheckResourceAttr("redash_query.my_query", "description", "my-query desc"),
					resource.TestCheckResourceAttr("redash_query.my_query", "query", "select '{{ rgx }}'"),
					resource.TestCheckResourceAttr("redash_query.my_query", "options.0.parameter.#", "1"),
					resource.TestCheckResourceAttr("redash_query.my_query", "options.0.parameter.0.title", "tregex"),
					resource.TestCheckResourceAttr("redash_query.my_query", "options.0.parameter.0.name", "rgx"),
					resource.TestCheckResourceAttr("redash_query.my_query", "options.0.parameter.0.type", "text-pattern"),
					resource.TestCheckResourceAttr("redash_query.my_query", "options.0.parameter.0.regex", "ab+c"),
				),
			},
			{
				Config: testAccParamQueryConfigEnum,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("redash_query.my_query", "name", "my-query"),
					resource.TestCheckResourceAttr("redash_query.my_query", "description", "my-query desc"),
					resource.TestCheckResourceAttr("redash_query.my_query", "query", "select '{{ enm }}'"),
					resource.TestCheckResourceAttr("redash_query.my_query", "options.0.parameter.#", "1"),
					resource.TestCheckResourceAttr("redash_query.my_query", "options.0.parameter.0.title", "tenum"),
					resource.TestCheckResourceAttr("redash_query.my_query", "options.0.parameter.0.name", "enm"),
					resource.TestCheckResourceAttr("redash_query.my_query", "options.0.parameter.0.type", "enum"),
					resource.TestCheckResourceAttr("redash_query.my_query", "options.0.parameter.0.enum.0.options.#", "3"),
					resource.TestCheckResourceAttr("redash_query.my_query", "options.0.parameter.0.enum.0.options.0", "aaa"),
					resource.TestCheckResourceAttr("redash_query.my_query", "options.0.parameter.0.enum.0.options.1", "bbb"),
					resource.TestCheckResourceAttr("redash_query.my_query", "options.0.parameter.0.enum.0.options.2", "ccc"),
				),
			},
			{
				Config: testAccParamQueryConfigEnumMultiValues1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("redash_query.my_query", "name", "my-query"),
					resource.TestCheckResourceAttr("redash_query.my_query", "description", "my-query desc"),
					resource.TestCheckResourceAttr("redash_query.my_query", "query", "select '{{ enm }}'"),
					resource.TestCheckResourceAttr("redash_query.my_query", "options.0.parameter.#", "1"),
					resource.TestCheckResourceAttr("redash_query.my_query", "options.0.parameter.0.title", "tenum"),
					resource.TestCheckResourceAttr("redash_query.my_query", "options.0.parameter.0.name", "enm"),
					resource.TestCheckResourceAttr("redash_query.my_query", "options.0.parameter.0.type", "enum"),
					resource.TestCheckResourceAttr("redash_query.my_query", "options.0.parameter.0.enum.0.options.#", "3"),
					resource.TestCheckResourceAttr("redash_query.my_query", "options.0.parameter.0.enum.0.options.0", "aaa"),
					resource.TestCheckResourceAttr("redash_query.my_query", "options.0.parameter.0.enum.0.options.1", "bbb"),
					resource.TestCheckResourceAttr("redash_query.my_query", "options.0.parameter.0.enum.0.options.2", "ccc"),
					resource.TestCheckResourceAttr("redash_query.my_query", "options.0.parameter.0.enum.0.multi_values.#", "1"),
					resource.TestCheckResourceAttr("redash_query.my_query", "options.0.parameter.0.enum.0.multi_values.0.quotation", ""),
					resource.TestCheckResourceAttr("redash_query.my_query", "options.0.parameter.0.enum.0.multi_values.0.separator", ","),
				),
			},
			{
				Config: testAccParamQueryConfigEnumMultiValues2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("redash_query.my_query", "name", "my-query"),
					resource.TestCheckResourceAttr("redash_query.my_query", "description", "my-query desc"),
					resource.TestCheckResourceAttr("redash_query.my_query", "query", "select '{{ enm }}'"),
					resource.TestCheckResourceAttr("redash_query.my_query", "options.0.parameter.#", "1"),
					resource.TestCheckResourceAttr("redash_query.my_query", "options.0.parameter.0.title", "tenum"),
					resource.TestCheckResourceAttr("redash_query.my_query", "options.0.parameter.0.name", "enm"),
					resource.TestCheckResourceAttr("redash_query.my_query", "options.0.parameter.0.type", "enum"),
					resource.TestCheckResourceAttr("redash_query.my_query", "options.0.parameter.0.enum.0.options.#", "3"),
					resource.TestCheckResourceAttr("redash_query.my_query", "options.0.parameter.0.enum.0.options.0", "aaa"),
					resource.TestCheckResourceAttr("redash_query.my_query", "options.0.parameter.0.enum.0.options.1", "bbb"),
					resource.TestCheckResourceAttr("redash_query.my_query", "options.0.parameter.0.enum.0.options.2", "ccc"),
					resource.TestCheckResourceAttr("redash_query.my_query", "options.0.parameter.0.enum.0.multi_values.#", "1"),
					resource.TestCheckResourceAttr("redash_query.my_query", "options.0.parameter.0.enum.0.multi_values.0.quotation", "'"),
					resource.TestCheckResourceAttr("redash_query.my_query", "options.0.parameter.0.enum.0.multi_values.0.separator", ";"),
				),
			},
			{
				Config: testAccParamQueryConfigQuery,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("redash_query.my_query", "name", "my-query"),
					resource.TestCheckResourceAttr("redash_query.my_query", "description", "my-query desc"),
					resource.TestCheckResourceAttr("redash_query.my_query", "query", "select '{{ q }}'"),
					resource.TestCheckResourceAttr("redash_query.my_query", "options.0.parameter.#", "1"),
					resource.TestCheckResourceAttr("redash_query.my_query", "options.0.parameter.0.title", "tenum"),
					resource.TestCheckResourceAttr("redash_query.my_query", "options.0.parameter.0.name", "q"),
					resource.TestCheckResourceAttr("redash_query.my_query", "options.0.parameter.0.type", "query"),
					resource.TestCheckResourceAttrWith("redash_query.my_query", "options.0.parameter.0.query_id", func(value string) error {
						i, err := strconv.Atoi(value)
						if err != nil {
							return fmt.Errorf("options.0.parameter.0.query_id must be number: %w", err)
						}
						if i <= 0 {
							return errors.New("options.0.parameter.0.query_id must be > 0")
						}
						return nil
					}),
				),
			},
		},
	})
}

func TestAccParamQuery_date(t *testing.T) {
	types := []string{
		"date",
		"datetime-local",
		"datetime-with-seconds",
		"date-range",
		"datetime-range",
		"datetime-range-with-seconds",
	}

	steps := []resource.TestStep{}

	for _, t := range types {
		steps = append(steps, resource.TestStep{
			Config: fmt.Sprintf(testAccParamQueryConfigDateTmpl, t, t, t, t),
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr("redash_query.my_query", "name", "my-query"),
				resource.TestCheckResourceAttr("redash_query.my_query", "description", "my-query desc"),
				resource.TestCheckResourceAttr("redash_query.my_query", "query", fmt.Sprintf("select '{{ n%s }}'", t)),
				resource.TestCheckResourceAttr("redash_query.my_query", "options.0.parameter.#", "1"),
				resource.TestCheckResourceAttr("redash_query.my_query", "options.0.parameter.0.title", "t"+t),
				resource.TestCheckResourceAttr("redash_query.my_query", "options.0.parameter.0.name", "n"+t),
				resource.TestCheckResourceAttr("redash_query.my_query", "options.0.parameter.0.type", t),
			),
		})
	}

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  func() { testAccPreCheck(t) },
		Steps:     steps,
	})
}

const testAccParamQueryConfigTextNum1 = testAccDataSourceConfigBasicPg + `
resource "redash_query" "my_query" {
  data_source_id = redash_data_source.my_data_source.id
  name           = "my-query"
  description    = "my-query desc"
  query          = "select '{{ txt }}','{{ num }}'"
  options {
    parameter {
      title = "ttext"
      name  = "txt"
      type  = "text"
      value = "hello"
    }
    parameter {
      title = "tnum"
      name  = "num"
      type  = "number"
      value = 100
    }
  }
}
`

const testAccParamQueryConfigTextNum2 = testAccDataSourceConfigBasicPg + `
resource "redash_query" "my_query" {
  data_source_id = redash_data_source.my_data_source.id
  name           = "my-query"
  description    = "my-query desc"
  query          = "select '{{ txt }}','{{ txt2 }}','{{ num2 }}'"
  options {
    parameter {
      title = "ttext"
      name  = "txt"
      type  = "text"
      value = "hello"
    }
    parameter {
      name  = "txt2"
      type  = "text"
      value = "hello2"
    }
    parameter {
      title = "tnum2"
      name  = "num2"
      type  = "number"
      value = 1002
    }
  }
}
`

const testAccParamQueryConfigTextNum3 = testAccDataSourceConfigBasicPg + `
resource "redash_query" "my_query" {
  data_source_id = redash_data_source.my_data_source.id
  name           = "my-query"
  description    = "my-query desc"
  query          = "select '{{ txt2 }}','{{ num2 }}'"
  options {
    parameter {
      title = "tnum2"
      name  = "num2"
      type  = "number"
      value = 1002
    }
    parameter {
      name  = "txt2"
      type  = "text"
      value = "hello2"
    }
  }
}
`

const testAccParamQueryConfigRegex = testAccDataSourceConfigBasicPg + `
resource "redash_query" "my_query" {
  data_source_id = redash_data_source.my_data_source.id
  name           = "my-query"
  description    = "my-query desc"
  query          = "select '{{ rgx }}'"
  options {
    parameter {
      title = "tregex"
      name  = "rgx"
      type  = "text-pattern"
      regex = "ab+c"
    }
  }
}
`

const testAccParamQueryConfigEnum = testAccDataSourceConfigBasicPg + `
resource "redash_query" "my_query" {
  data_source_id = redash_data_source.my_data_source.id
  name           = "my-query"
  description    = "my-query desc"
  query          = "select '{{ enm }}'"
  options {
    parameter {
      title = "tenum"
      name  = "enm"
      type  = "enum"

			enum {
				options = ["aaa", "bbb", "ccc"]
			}
    }
  }
}
`

const testAccParamQueryConfigEnumMultiValues1 = testAccDataSourceConfigBasicPg + `
resource "redash_query" "my_query" {
  data_source_id = redash_data_source.my_data_source.id
  name           = "my-query"
  description    = "my-query desc"
  query          = "select '{{ enm }}'"
  options {
    parameter {
      title = "tenum"
      name  = "enm"
      type  = "enum"

			enum {
				options = ["aaa", "bbb", "ccc"]

				multi_values {}
			}
    }
  }
}
`

const testAccParamQueryConfigEnumMultiValues2 = testAccDataSourceConfigBasicPg + `
resource "redash_query" "my_query" {
  data_source_id = redash_data_source.my_data_source.id
  name           = "my-query"
  description    = "my-query desc"
  query          = "select '{{ enm }}'"
  options {
    parameter {
      title = "tenum"
      name  = "enm"
      type  = "enum"

			enum {
				options = ["aaa", "bbb", "ccc"]

				multi_values {
					quotation = "'"
					separator = ";"
				}
			}
    }
  }
}
`

const testAccParamQueryConfigQuery = testAccDataSourceConfigBasicPg + `
resource "redash_query" "select_array" {
  data_source_id = redash_data_source.my_data_source.id
  name           = "select array"
  query          = "select unnest(array[1,2,3])"
  published      = true
}

resource "redash_query" "my_query" {
  data_source_id = redash_data_source.my_data_source.id
  name           = "my-query"
  description    = "my-query desc"
  query          = "select '{{ q }}'"
  options {
    parameter {
      title    = "tenum"
      name     = "q"
      type     = "query"
			query_id = redash_query.select_array.id
    }
  }
}
`

const testAccParamQueryConfigDateTmpl = testAccDataSourceConfigBasicPg + `
resource "redash_query" "my_query" {
  data_source_id = redash_data_source.my_data_source.id
  name           = "my-query"
  description    = "my-query desc"
  query          = "select '{{ n%s }}'"
  options {
    parameter {
      title = "t%s"
      name  = "n%s"
      type  = "%s"
    }
  }
}
`
