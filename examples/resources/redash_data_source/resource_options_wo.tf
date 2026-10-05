resource "redash_data_source" "example_database" {
  name = "example-database"
  type = "pg"
  # options_wo is not stored in Terraform state. Requires Terraform 1.11 or later.
  # see https://github.com/getredash/redash/blob/v25.1/redash/query_runner/pg.py#L149-L153
  options_wo = jsonencode({
    dbname   = "example-database"
    host     = "postgres"
    port     = 5432
    user     = "example-database"
    password = "example-password"
  })
  # Stored in state. Increment this to apply a new options_wo value.
  options_wo_version = 1
}
