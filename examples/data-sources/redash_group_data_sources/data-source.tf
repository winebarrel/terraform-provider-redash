data "redash_group_data_sources" "default" {
  group_id = 2
}

output "default_data_sources" {
  value = data.redash_group_data_sources.default.data_sources
}
