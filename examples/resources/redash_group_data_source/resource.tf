resource "redash_group" "my_group" {
  name = "my-group"
}

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

resource "redash_group_data_source" "my_group_ds" {
  group_id       = redash_group.my_group.id
  data_source_id = redash_data_source.postgres.id
  view_only      = true
}
