resource "redash_group" "my_group" {
  name = "my-group"
}

data "redash_user" "admin" {
  name = "admin"
}

resource "redash_group_member" "my_member" {
  group_id = redash_group.my_group.id
  user_id  = data.redash_user.admin.id
}
