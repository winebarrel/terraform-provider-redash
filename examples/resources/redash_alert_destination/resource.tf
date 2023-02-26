resource "redash_alert_destination" "my_dest" {
  name = "my-dest"
  type = "email"
  options = jsonencode({
    addresses = "foo@example.com"
  })
}
