---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "redash_user Data Source - redash"
subcategory: ""
description: |-
  
---

# redash_user (Data Source)



## Example Usage

```terraform
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
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `email` (String)
- `name` (String)

### Read-Only

- `id` (String) The ID of this resource.
