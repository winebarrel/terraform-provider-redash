data "redash_data_source" "postgres" {
  name = "postgres"
}

resource "redash_query" "select_one" {
  data_source_id = data.redash_data_source.postgres.id
  name           = "select one"
  description    = "Select one."
  query          = "select 1"
  schedule {
    interval = 60
  }
}
