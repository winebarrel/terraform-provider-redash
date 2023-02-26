resource "redash_query" "select_one" {
  name = "select one"
}

resource "redash_alert" "my_alert" {
  name     = "my-alert"
  query_id = data.redash_query.select_one.id
  options {
    column         = "value"
    op             = "<"
    value          = 1
    custom_subject = "alert"
    custom_body    = "service down"
  }
}
