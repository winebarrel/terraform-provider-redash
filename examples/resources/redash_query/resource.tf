resource "redash_data_source" "postgres" {
  name = "postgres"
  type = "pg"
  options = jsonencode({
    dbname = "postgres"
    host   = "postgres"
    port   = 5432
    user   = "postgres"
  })
}

resource "redash_query" "select_one" {
  data_source_id = redash_data_source.postgres.id
  name           = "select one"
  description    = "Select one."
  query          = "select 1"
  schedule {
    interval = 60
  }
  tags = ["foo", "bar"]
}
