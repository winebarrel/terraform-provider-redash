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

resource "redash_data_source" "events" {
  name = "events"
  type = "pg"
  options = jsonencode({
    dbname = "events"
    host   = "postgres"
    port   = 5432
    user   = "postgres"
  })
}

resource "redash_group_data_sources" "my_group" {
  group_id = redash_group.my_group.id

  dynamic "data_source" {
    for_each = [
      redash_data_source.postgres,
      redash_data_source.events,
    ]

    content {
      data_source_id = data_source.value.id
      name           = data_source.value.name
    }
  }
}
