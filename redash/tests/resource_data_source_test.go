package test

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	redashgo "github.com/winebarrel/redash-go/v2"
)

func TestAccDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceConfigBasicPg,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("redash_data_source.my_data_source", "name", "my-data-source"),
					resource.TestCheckResourceAttr("redash_data_source.my_data_source", "type", "pg"),
					resource.TestCheckResourceAttr("redash_data_source.my_data_source", "options", `{"dbname":"postgres","host":"postgres","port":5432,"user":"postgres"}`),
				),
			},
			{
				Config: testAccDataSourceConfigBasicPg2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("redash_data_source.my_data_source", "name", "my-data-source2"),
					resource.TestCheckResourceAttr("redash_data_source.my_data_source", "type", "pg"),
					resource.TestCheckResourceAttr("redash_data_source.my_data_source", "options", `{"dbname":"postgres2","host":"postgres2","port":5433,"user":"postgres2"}`),
				),
			},
		},
	})
}

func TestAccDataSource_secret(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceConfigSecretPg,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("redash_data_source.my_data_source_with_secret", "name", "my-data-source-with-secret"),
					resource.TestCheckResourceAttr("redash_data_source.my_data_source_with_secret", "type", "pg"),
					// The state must keep the real password, not the placeholder returned by the API.
					// cf. https://github.com/winebarrel/terraform-provider-redash/issues/177
					resource.TestCheckResourceAttr("redash_data_source.my_data_source_with_secret", "options", `{"dbname":"postgres","host":"postgres","password":"supersecret","port":5432,"user":"postgres"}`),
				),
			},
			{
				// Applying the same config again must yield an empty plan.
				Config:   testAccDataSourceConfigSecretPg,
				PlanOnly: true,
			},
			{
				Config: testAccDataSourceConfigSecretPg2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("redash_data_source.my_data_source_with_secret", "options", `{"dbname":"postgres","host":"postgres","password":"supersecret2","port":5432,"user":"postgres"}`),
				),
			},
		},
	})
}

const testAccDataSourceConfigBasicPg = `
resource "redash_data_source" "my_data_source" {
  name = "my-data-source"
  type = "pg"
  options = jsonencode({
    dbname = "postgres"
    host   = "postgres"
    port   = 5432
    user   = "postgres"
  })
}
`

const testAccDataSourceConfigBasicPg2 = `
resource "redash_data_source" "my_data_source" {
  name = "my-data-source2"
  type = "pg"
  options = jsonencode({
    dbname = "postgres2"
    host   = "postgres2"
    port   = 5433
    user   = "postgres2"
  })
}
`

const testAccDataSourceConfigSecretPg = `
resource "redash_data_source" "my_data_source_with_secret" {
  name = "my-data-source-with-secret"
  type = "pg"
  options = jsonencode({
    dbname   = "postgres"
    host     = "postgres"
    port     = 5432
    user     = "postgres"
    password = "supersecret"
  })
}
`

func TestAccDataSource_optionsWriteOnly(t *testing.T) {
	const addr = "redash_data_source.my_data_source_write_only"

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceConfigWriteOnlyPg,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(addr, "name", "my-data-source-write-only"),
					resource.TestCheckResourceAttr(addr, "type", "pg"),
					resource.TestCheckResourceAttr(addr, "options_wo_version", "1"),
					testAccCheckOptionsNotInState(addr, "supersecret"),
					testAccCheckAPIDataSource(addr, map[string]string{
						"dbname":   "postgres",
						"host":     "postgres",
						"user":     "postgres",
						"password": "--------",
					}),
				),
			},
			{
				Config:   testAccDataSourceConfigWriteOnlyPg,
				PlanOnly: true,
			},
			{
				// A new options_wo value with the same version must not plan an update.
				Config:   testAccDataSourceConfigWriteOnlyPgUnchangedVersion,
				PlanOnly: true,
			},
			{
				Config: testAccDataSourceConfigWriteOnlyPg2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(addr, "options_wo_version", "2"),
					testAccCheckOptionsNotInState(addr, "supersecret", "supersecret2"),
					testAccCheckAPIDataSource(addr, map[string]string{
						"dbname":   "postgres2",
						"host":     "postgres",
						"user":     "postgres",
						"password": "--------",
					}),
				),
			},
		},
	})
}

func testAccCheckOptionsNotInState(addr string, secrets ...string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[addr]
		if !ok {
			return fmt.Errorf("not found: %s", addr)
		}

		if options := rs.Primary.Attributes["options"]; options != "" {
			return fmt.Errorf("options = %q, want empty", options)
		}
		if wo := rs.Primary.Attributes["options_wo"]; wo != "" {
			return fmt.Errorf("options_wo stored in state")
		}

		for k, v := range rs.Primary.Attributes {
			for _, secret := range secrets {
				if strings.Contains(v, secret) {
					return fmt.Errorf("state attribute %s contains %q", k, secret)
				}
			}
		}

		return nil
	}
}

func testAccCheckAPIDataSource(addr string, want map[string]string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[addr]
		if !ok {
			return fmt.Errorf("not found: %s", addr)
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return err
		}

		client, err := redashgo.NewClient(testAccRedashURL, testAccRedashAPIKey)
		if err != nil {
			return err
		}

		ds, err := client.GetDataSource(context.Background(), id)
		if err != nil {
			return err
		}

		for key, wantValue := range want {
			got, ok := ds.Options[key]
			if !ok {
				return fmt.Errorf("options.%s missing in %#v", key, ds.Options)
			}
			if fmt.Sprint(got) != wantValue {
				return fmt.Errorf("options.%s = %#v, want %q", key, got, wantValue)
			}
		}

		return nil
	}
}

const testAccDataSourceConfigWriteOnlyPg = `
resource "redash_data_source" "my_data_source_write_only" {
  name = "my-data-source-write-only"
  type = "pg"
  options_wo = jsonencode({
    dbname   = "postgres"
    host     = "postgres"
    port     = 5432
    user     = "postgres"
    password = "supersecret"
  })
  options_wo_version = 1
}
`

const testAccDataSourceConfigWriteOnlyPgUnchangedVersion = `
resource "redash_data_source" "my_data_source_write_only" {
  name = "my-data-source-write-only"
  type = "pg"
  options_wo = jsonencode({
    dbname   = "postgres2"
    host     = "postgres"
    port     = 5432
    user     = "postgres"
    password = "supersecret2"
  })
  options_wo_version = 1
}
`

const testAccDataSourceConfigWriteOnlyPg2 = `
resource "redash_data_source" "my_data_source_write_only" {
  name = "my-data-source-write-only"
  type = "pg"
  options_wo = jsonencode({
    dbname   = "postgres2"
    host     = "postgres"
    port     = 5432
    user     = "postgres"
    password = "supersecret2"
  })
  options_wo_version = 2
}
`

const testAccDataSourceConfigSecretPg2 = `
resource "redash_data_source" "my_data_source_with_secret" {
  name = "my-data-source-with-secret"
  type = "pg"
  options = jsonencode({
    dbname   = "postgres"
    host     = "postgres"
    port     = 5432
    user     = "postgres"
    password = "supersecret2"
  })
}
`
