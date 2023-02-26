data "redash_alert" "my_alert" {
  name = "my-alert"
}

resource "redash_alert_destination" "my_dest" {
  name = "my-dest"
  type = "email"
  options = jsonencode({
    addresses = "foo@example.com"
  })
}

resource "redash_alert_subscription" "my_subs" {
  alert_id             = data.redash_alert.my_alert.id
  alert_destination_id = redash_alert_destination.my_dest.id
}
