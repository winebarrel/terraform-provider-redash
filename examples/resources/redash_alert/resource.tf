resource "redash_query" "select_one" {
  data_source_id = redash_data_source.postgres.id
  name           = "select one"
  description    = "Select one."
  query          = "select 1"
  schedule {
    interval = 60
  }
}

resource "redash_alert" "my_alert" {
  name     = "my-alert"
  query_id = redash_query.select_one.id
  options {
    column         = "value"
    op             = "<"
    value          = 1
    custom_subject = "alert"
    custom_body    = "service down"
  }
}
