data "redash_users" "example_com" {
  email = "*@example.com"
}

output "example_com_user_ids" {
  value = data.redash_users.example_com.ids
}
