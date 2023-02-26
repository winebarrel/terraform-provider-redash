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

data "redash_alert_destination" "my_dest" {
  name = "my-dest"
}

resource "redash_alert_subscription" "my_subs" {
  alert_id             = redash_alert.my_alert.id
  alert_destination_id = data.redash_alert_destination.my_dest.id
}
