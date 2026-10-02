package test

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	redashgo "github.com/winebarrel/redash-go/v2"
)

func TestAccGroupDataSources_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testAccGroupDataSourcesOne,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("redash_group_data_sources.my_group", "data_source.#", "1"),
					resource.TestCheckTypeSetElemAttrPair(
						"redash_group_data_sources.my_group", "data_source.*.data_source_id",
						"redash_data_source.my_data_source", "id",
					),
					resource.TestCheckTypeSetElemNestedAttrs("redash_group_data_sources.my_group", "data_source.*", map[string]string{
						"view_only": "false",
					}),
				),
			},
			{
				Config: testAccGroupDataSourcesTwo,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("redash_group_data_sources.my_group", "data_source.#", "2"),
					resource.TestCheckTypeSetElemAttrPair(
						"redash_group_data_sources.my_group", "data_source.*.data_source_id",
						"redash_data_source.my_data_source", "id",
					),
					resource.TestCheckTypeSetElemAttrPair(
						"redash_group_data_sources.my_group", "data_source.*.data_source_id",
						"redash_data_source.extra", "id",
					),
				),
			},
			{
				Config: testAccGroupDataSourcesOne,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("redash_group_data_sources.my_group", "data_source.#", "1"),
					resource.TestCheckTypeSetElemAttrPair(
						"redash_group_data_sources.my_group", "data_source.*.data_source_id",
						"redash_data_source.my_data_source", "id",
					),
				),
			},
			{
				Config: testAccGroupDataSourcesViewOnly,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("redash_group_data_sources.my_group", "data_source.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs("redash_group_data_sources.my_group", "data_source.*", map[string]string{
						"view_only": "true",
					}),
				),
			},
			{
				ResourceName:      "redash_group_data_sources.my_group",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccGroupDataSourcesEmpty,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("redash_group_data_sources.my_group", "data_source.#", "0"),
				),
			},
		},
	})
}

func TestAccGroupDataSources_removesExtra(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testAccGroupDataSourcesOne,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("redash_group_data_sources.my_group", "data_source.#", "1"),
					testAccAddExtraGroupDataSource,
				),
				// The grant added in Check shows up as a diff.
				ExpectNonEmptyPlan: true,
			},
			{
				Config: testAccGroupDataSourcesOne,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("redash_group_data_sources.my_group", "data_source.#", "1"),
					resource.TestCheckTypeSetElemAttrPair(
						"redash_group_data_sources.my_group", "data_source.*.data_source_id",
						"redash_data_source.my_data_source", "id",
					),
				),
			},
		},
	})
}

func TestAccGroupDataSources_createWithNewDataSources(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testAccGroupDataSourcesTwo,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("redash_group_data_sources.my_group", "data_source.#", "2"),
					resource.TestCheckTypeSetElemAttrPair(
						"redash_group_data_sources.my_group", "data_source.*.data_source_id",
						"redash_data_source.my_data_source", "id",
					),
					resource.TestCheckTypeSetElemAttrPair(
						"redash_group_data_sources.my_group", "data_source.*.data_source_id",
						"redash_data_source.extra", "id",
					),
				),
			},
		},
	})
}

func TestAccGroupDataSources_createRemovesExisting(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testAccGroupDataSourcesBase,
				Check:  testAccAddExtraGroupDataSource,
			},
			{
				Config: testAccGroupDataSourcesOne,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("redash_group_data_sources.my_group", "data_source.#", "1"),
					resource.TestCheckTypeSetElemAttrPair(
						"redash_group_data_sources.my_group", "data_source.*.data_source_id",
						"redash_data_source.my_data_source", "id",
					),
				),
			},
		},
	})
}

func TestAccGroupDataSources_destroyRemovesAll(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testAccGroupDataSourcesOne,
				Check:  testAccAddExtraGroupDataSource,
				// The grant added in Check shows up as a diff.
				ExpectNonEmptyPlan: true,
			},
			{
				// Destroy only redash_group_data_sources.
				Config: testAccGroupDataSourcesBase,
				Check:  testAccCheckGroupHasNoDataSources,
			},
		},
	})
}

func testAccCheckGroupHasNoDataSources(s *terraform.State) error {
	groupId, err := strconv.Atoi(s.RootModule().Resources["redash_group.my_group"].Primary.ID)
	if err != nil {
		return err
	}

	client := testAccProvider.Meta().(*redashgo.Client)
	dsList, err := client.ListGroupDataSources(context.Background(), groupId)
	if err != nil {
		return err
	}

	if len(dsList) != 0 {
		return fmt.Errorf("group %d still has %d data sources", groupId, len(dsList))
	}

	return nil
}

func testAccAddExtraGroupDataSource(s *terraform.State) error {
	groupId, err := strconv.Atoi(s.RootModule().Resources["redash_group.my_group"].Primary.ID)
	if err != nil {
		return err
	}

	dsId, err := strconv.Atoi(s.RootModule().Resources["redash_data_source.extra"].Primary.ID)
	if err != nil {
		return err
	}

	client := testAccProvider.Meta().(*redashgo.Client)
	_, err = client.AddGroupDataSource(context.Background(), groupId, dsId)
	return err
}

const testAccGroupDataSourcesBase = testAccGroupConfigBasic + testAccDataSourceConfigBasicPg + `
resource "redash_data_source" "extra" {
	name = "extra-data-source"
	type = "pg"
	options = jsonencode({
		dbname = "postgres"
		host   = "postgres"
		port   = 5432
		user   = "postgres"
	})
}
`

const testAccGroupDataSourcesOne = testAccGroupDataSourcesBase + `
resource "redash_group_data_sources" "my_group" {
	group_id = redash_group.my_group.id

	data_source {
		data_source_id = redash_data_source.my_data_source.id
	}
}
`

const testAccGroupDataSourcesTwo = testAccGroupDataSourcesBase + `
resource "redash_group_data_sources" "my_group" {
	group_id = redash_group.my_group.id

	data_source {
		data_source_id = redash_data_source.my_data_source.id
	}

	data_source {
		data_source_id = redash_data_source.extra.id
	}
}
`

const testAccGroupDataSourcesViewOnly = testAccGroupDataSourcesBase + `
resource "redash_group_data_sources" "my_group" {
	group_id = redash_group.my_group.id

	data_source {
		data_source_id = redash_data_source.my_data_source.id
		view_only      = true
	}
}
`

const testAccGroupDataSourcesEmpty = testAccGroupDataSourcesBase + `
resource "redash_group_data_sources" "my_group" {
	group_id = redash_group.my_group.id
}
`
