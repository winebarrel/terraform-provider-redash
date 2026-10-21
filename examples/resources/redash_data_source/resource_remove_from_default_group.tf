resource "redash_data_source" "private_database" {
  name = "private-database"
  type = "pg"
  options = jsonencode({
    dbname = "private-database"
    host   = "postgres"
    port   = 5432
    user   = "private-database"
  })
  # Redash grants every new data source to the default group, which includes all users.
  # Remove that grant on create, then grant access to specific groups.
  remove_from_default_group = true
}
