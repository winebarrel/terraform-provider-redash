package test

import (
	"errors"
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSoureceParamQuery_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testAccParamQueryConfigTextNum1,
			},
			{
				Config: testAccDataSoureceParamQueryConfigTextNum,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.redash_query.my_query", "name", "my-query"),
					resource.TestCheckResourceAttr("data.redash_query.my_query", "description", "my-query desc"),
					resource.TestCheckResourceAttr("data.redash_query.my_query", "query", "select '{{ txt }}','{{ num }}'"),
					resource.TestCheckResourceAttr("data.redash_query.my_query", "options.0.parameter.#", "2"),
					resource.TestCheckResourceAttr("data.redash_query.my_query", "options.0.parameter.0.title", "tnum"),
					resource.TestCheckResourceAttr("data.redash_query.my_query", "options.0.parameter.0.name", "num"),
					resource.TestCheckResourceAttr("data.redash_query.my_query", "options.0.parameter.0.type", "number"),
					resource.TestCheckResourceAttr("data.redash_query.my_query", "options.0.parameter.0.value", "100"),
					resource.TestCheckResourceAttr("data.redash_query.my_query", "options.0.parameter.1.title", "ttext"),
					resource.TestCheckResourceAttr("data.redash_query.my_query", "options.0.parameter.1.name", "txt"),
					resource.TestCheckResourceAttr("data.redash_query.my_query", "options.0.parameter.1.type", "text"),
					resource.TestCheckResourceAttr("data.redash_query.my_query", "options.0.parameter.1.value", "hello"),
				),
			},
			{
				Config: testAccParamQueryConfigRegex,
			},
			{
				Config: testAccDataSoureceParamQueryConfigRegex,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.redash_query.my_query", "name", "my-query"),
					resource.TestCheckResourceAttr("data.redash_query.my_query", "description", "my-query desc"),
					resource.TestCheckResourceAttr("data.redash_query.my_query", "query", "select '{{ rgx }}'"),
					resource.TestCheckResourceAttr("data.redash_query.my_query", "options.0.parameter.#", "1"),
					resource.TestCheckResourceAttr("data.redash_query.my_query", "options.0.parameter.0.title", "tregex"),
					resource.TestCheckResourceAttr("data.redash_query.my_query", "options.0.parameter.0.name", "rgx"),
					resource.TestCheckResourceAttr("data.redash_query.my_query", "options.0.parameter.0.type", "text-pattern"),
					resource.TestCheckResourceAttr("data.redash_query.my_query", "options.0.parameter.0.regex", "ab+c"),
				),
			},
			{
				Config: testAccParamQueryConfigEnum,
			},
			{
				Config: testAccDataSoureceParamQueryConfigEnum,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.redash_query.my_query", "name", "my-query"),
					resource.TestCheckResourceAttr("data.redash_query.my_query", "description", "my-query desc"),
					resource.TestCheckResourceAttr("data.redash_query.my_query", "query", "select '{{ enm }}'"),
					resource.TestCheckResourceAttr("data.redash_query.my_query", "options.0.parameter.#", "1"),
					resource.TestCheckResourceAttr("data.redash_query.my_query", "options.0.parameter.0.title", "tenum"),
					resource.TestCheckResourceAttr("data.redash_query.my_query", "options.0.parameter.0.name", "enm"),
					resource.TestCheckResourceAttr("data.redash_query.my_query", "options.0.parameter.0.type", "enum"),
					resource.TestCheckResourceAttr("data.redash_query.my_query", "options.0.parameter.0.enum.0.options.#", "3"),
					resource.TestCheckResourceAttr("data.redash_query.my_query", "options.0.parameter.0.enum.0.options.0", "aaa"),
					resource.TestCheckResourceAttr("data.redash_query.my_query", "options.0.parameter.0.enum.0.options.1", "bbb"),
					resource.TestCheckResourceAttr("data.redash_query.my_query", "options.0.parameter.0.enum.0.options.2", "ccc"),
				),
			},
			{
				Config: testAccParamQueryConfigEnumMultiValues2,
			},
			{
				Config: testAccDataSoureceParamQueryConfigEnumMultiValues,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.redash_query.my_query", "name", "my-query"),
					resource.TestCheckResourceAttr("data.redash_query.my_query", "description", "my-query desc"),
					resource.TestCheckResourceAttr("data.redash_query.my_query", "query", "select '{{ enm }}'"),
					resource.TestCheckResourceAttr("data.redash_query.my_query", "options.0.parameter.#", "1"),
					resource.TestCheckResourceAttr("data.redash_query.my_query", "options.0.parameter.0.title", "tenum"),
					resource.TestCheckResourceAttr("data.redash_query.my_query", "options.0.parameter.0.name", "enm"),
					resource.TestCheckResourceAttr("data.redash_query.my_query", "options.0.parameter.0.type", "enum"),
					resource.TestCheckResourceAttr("data.redash_query.my_query", "options.0.parameter.0.enum.0.options.#", "3"),
					resource.TestCheckResourceAttr("data.redash_query.my_query", "options.0.parameter.0.enum.0.options.0", "aaa"),
					resource.TestCheckResourceAttr("data.redash_query.my_query", "options.0.parameter.0.enum.0.options.1", "bbb"),
					resource.TestCheckResourceAttr("data.redash_query.my_query", "options.0.parameter.0.enum.0.options.2", "ccc"),
					resource.TestCheckResourceAttr("data.redash_query.my_query", "options.0.parameter.0.enum.0.multi_values.#", "1"),
					resource.TestCheckResourceAttr("data.redash_query.my_query", "options.0.parameter.0.enum.0.multi_values.0.quotation", "'"),
					resource.TestCheckResourceAttr("data.redash_query.my_query", "options.0.parameter.0.enum.0.multi_values.0.separator", ";"),
				),
			},
			{
				Config: testAccParamQueryConfigQuery,
			},
			{
				Config: testAccDataSoureceParamQueryConfigQuery,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.redash_query.my_query", "name", "my-query"),
					resource.TestCheckResourceAttr("data.redash_query.my_query", "description", "my-query desc"),
					resource.TestCheckResourceAttr("data.redash_query.my_query", "query", "select '{{ q }}'"),
					resource.TestCheckResourceAttr("data.redash_query.my_query", "options.0.parameter.#", "1"),
					resource.TestCheckResourceAttr("data.redash_query.my_query", "options.0.parameter.0.title", "tenum"),
					resource.TestCheckResourceAttr("data.redash_query.my_query", "options.0.parameter.0.name", "q"),
					resource.TestCheckResourceAttr("data.redash_query.my_query", "options.0.parameter.0.type", "query"),
					resource.TestCheckResourceAttrWith("data.redash_query.my_query", "options.0.parameter.0.query_id", func(value string) error {
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

func TestAccDataSoureceParamQuery_date(t *testing.T) {
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
		})

		steps = append(steps, resource.TestStep{
			Config: fmt.Sprintf(testAccDataSoureceParamQueryConfigDateTmpl, t, t, t, t),
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr("data.redash_query.my_query", "name", "my-query"),
				resource.TestCheckResourceAttr("data.redash_query.my_query", "description", "my-query desc"),
				resource.TestCheckResourceAttr("data.redash_query.my_query", "query", fmt.Sprintf("select '{{ n%s }}'", t)),
				resource.TestCheckResourceAttr("data.redash_query.my_query", "options.0.parameter.#", "1"),
				resource.TestCheckResourceAttr("data.redash_query.my_query", "options.0.parameter.0.title", "t"+t),
				resource.TestCheckResourceAttr("data.redash_query.my_query", "options.0.parameter.0.name", "n"+t),
				resource.TestCheckResourceAttr("data.redash_query.my_query", "options.0.parameter.0.type", t),
			),
		})
	}

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  func() { testAccPreCheck(t) },
		Steps:     steps,
	})
}

const testAccDataSoureceParamQueryConfigTextNum = testAccParamQueryConfigTextNum1 + `
data "redash_query" "my_query" {
	query_id = redash_query.my_query.id
	name     = "my-query"
}
`

const testAccDataSoureceParamQueryConfigRegex = testAccParamQueryConfigRegex + `
data "redash_query" "my_query" {
	query_id = redash_query.my_query.id
	name     = "my-query"
}
`

const testAccDataSoureceParamQueryConfigEnum = testAccParamQueryConfigEnum + `
data "redash_query" "my_query" {
	query_id = redash_query.my_query.id
	name     = "my-query"
}
`

const testAccDataSoureceParamQueryConfigEnumMultiValues = testAccParamQueryConfigEnumMultiValues2 + `
data "redash_query" "my_query" {
	query_id = redash_query.my_query.id
	name     = "my-query"
}
`

const testAccDataSoureceParamQueryConfigQuery = testAccParamQueryConfigQuery + `
data "redash_query" "my_query" {
	query_id = redash_query.my_query.id
	name     = "my-query"
}
`

const testAccDataSoureceParamQueryConfigDateTmpl = testAccParamQueryConfigDateTmpl + `
data "redash_query" "my_query" {
	query_id = redash_query.my_query.id
	name     = "my-query"
}
`
