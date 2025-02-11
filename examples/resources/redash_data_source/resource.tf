resource "redash_data_source" "postgres" {
  name = "postgres"
  type = "pg"
  # see https://github.com/getredash/redash/blob/v25.1/redash/query_runner/pg.py#L149-L153
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
}
