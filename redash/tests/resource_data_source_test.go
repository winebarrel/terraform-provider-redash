package test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourece_basic(t *testing.T) {
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

func TestAccDataSourece_secret(t *testing.T) {
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
