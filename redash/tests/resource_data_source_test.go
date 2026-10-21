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
			{
				ResourceName:      "redash_data_source.my_data_source",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccDataSource_inDefaultGroup(t *testing.T) {
	const addr = "redash_data_source.my_data_source_private"

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				// Redash adds a new data source to the default group.
				Config: testAccDataSourceConfigPrivate("true"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(addr, "in_default_group", "true"),
					testAccCheckDataSourceInDefaultGroup(addr, true),
				),
			},
			{
				// Changing it to false removes the existing data source.
				Config: testAccDataSourceConfigPrivate("false"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(addr, "in_default_group", "false"),
					testAccCheckDataSourceInDefaultGroup(addr, false),
				),
			},
			{
				Config:   testAccDataSourceConfigPrivate("false"),
				PlanOnly: true,
			},
			{
				ResourceName:      addr,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Changing it back to true adds the data source again.
				Config: testAccDataSourceConfigPrivate("true"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(addr, "in_default_group", "true"),
					testAccCheckDataSourceInDefaultGroup(addr, true),
				),
			},
			{
				Config:  testAccDataSourceConfigPrivate("true"),
				Destroy: true,
			},
			{
				// Creating with false removes the new data source.
				Config: testAccDataSourceConfigPrivate("false"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(addr, "in_default_group", "false"),
					testAccCheckDataSourceInDefaultGroup(addr, false),
				),
			},
			{
				Config:   testAccDataSourceConfigPrivate("false"),
				PlanOnly: true,
			},
		},
	})
}

func testAccCheckDataSourceInDefaultGroup(addr string, want bool) resource.TestCheckFunc {
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

		groups, err := client.ListGroups(context.Background())
		if err != nil {
			return err
		}

		for _, g := range groups {
			if g.Type != "builtin" || g.Name != "default" {
				continue
			}

			dsList, err := client.ListGroupDataSources(context.Background(), g.ID)
			if err != nil {
				return err
			}

			got := false

			for _, ds := range dsList {
				if ds.ID == id {
					got = true
				}
			}

			if got != want {
				return fmt.Errorf("data source %d in default group = %t, want %t", id, got, want)
			}

			return nil
		}

		return fmt.Errorf("default group not found")
	}
}

func testAccDataSourceConfigPrivate(inDefaultGroup string) string {
	return fmt.Sprintf(`
resource "redash_data_source" "my_data_source_private" {
  name = "my-data-source-private"
  type = "pg"
  options = jsonencode({
    dbname = "postgres"
    host   = "postgres"
    port   = 5432
    user   = "postgres"
  })
  in_default_group = %s
}
`, inDefaultGroup)
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
				// Redash does not return options_wo, and import records the API
				// options, so skip the option attributes.
				ResourceName:            addr,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"options", "options_wo", "options_wo_version"},
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

func TestAccDataSource_optionsWriteOnlyRename(t *testing.T) {
	const addr = "redash_data_source.my_data_source_write_only_rename"

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceConfigWriteOnlyRename("my-data-source-write-only-rename"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(addr, "name", "my-data-source-write-only-rename"),
					testAccCheckOptionsNotInState(addr, "supersecret"),
				),
			},
			{
				// Changing only the name must keep the options in Redash.
				Config: testAccDataSourceConfigWriteOnlyRename("my-data-source-write-only-rename2"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(addr, "name", "my-data-source-write-only-rename2"),
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
		},
	})
}

func TestAccDataSource_optionsToWriteOnly(t *testing.T) {
	const addr = "redash_data_source.my_data_source_migrate"

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceConfigMigrateOptions("supersecret"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(addr, "options", `{"dbname":"postgres","host":"postgres","password":"supersecret","port":5432,"user":"postgres"}`),
				),
			},
			{
				// Switching from options to options_wo removes options from the state.
				Config: testAccDataSourceConfigMigrateWriteOnly,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(addr, "options_wo_version", "1"),
					testAccCheckOptionsNotInState(addr, "supersecret", "supersecret2"),
					testAccCheckAPIDataSource(addr, map[string]string{
						"dbname":   "postgres2",
						"host":     "postgres",
						"user":     "postgres",
						"password": "--------",
					}),
				),
			},
			{
				Config:   testAccDataSourceConfigMigrateWriteOnly,
				PlanOnly: true,
			},
			{
				// Switching back to options stores options in the state again.
				Config: testAccDataSourceConfigMigrateOptions("supersecret3"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(addr, "options", `{"dbname":"postgres","host":"postgres","password":"supersecret3","port":5432,"user":"postgres"}`),
				),
			},
			{
				Config:   testAccDataSourceConfigMigrateOptions("supersecret3"),
				PlanOnly: true,
			},
		},
	})
}

func testAccDataSourceConfigWriteOnlyRename(name string) string {
	return fmt.Sprintf(`
resource "redash_data_source" "my_data_source_write_only_rename" {
  name = %q
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
`, name)
}

func testAccDataSourceConfigMigrateOptions(password string) string {
	return fmt.Sprintf(`
resource "redash_data_source" "my_data_source_migrate" {
  name = "my-data-source-migrate"
  type = "pg"
  options = jsonencode({
    dbname   = "postgres"
    host     = "postgres"
    port     = 5432
    user     = "postgres"
    password = %q
  })
}
`, password)
}

const testAccDataSourceConfigMigrateWriteOnly = `
resource "redash_data_source" "my_data_source_migrate" {
  name = "my-data-source-migrate"
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
