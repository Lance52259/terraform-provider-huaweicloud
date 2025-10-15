---
subcategory: "Identity and Access Management (IAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identity_user_v5"
description: |-
  Manages an IAM user v5 resource within HuaweiCloud.
---

# huaweicloud_identity_user_v5

## Example Usage

### Create IAM User

```hcl
variable "name" {}

resource "huaweicloud_identity_user_v5" "user" {
  name        = var.name
  enabled     = true
  description = "tested by terraform"
}
```
 
### Update IAM User

```hcl
variable "name" {}

resource "huaweicloud_identity_user_v5" "user" {
  name        = var.name
  enabled     = true
  description = "tested by terraform"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String) Specifies the name of the user. The username consists of `1` to `64` characters. It can
  contain only uppercase letters, lowercase letters, digits, spaces, and special characters (-_) and cannot start with a
  digit.

* `description` - (Optional, String) Specifies the description of the user.

* `enabled` - (Optional, Bool) Specifies whether the user is enabled or disabled. Valid values are **true** and **false**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `user_name` - Indicates the IAM username.

* `is_root_user` - Indicates whether the user is root user or user.

* `created_at` - Indicates the time when the IAM user was created.

* `urn` - Indicates the uniform resource name.

* `last_login_at` - Indicates the last login time of the IAM user. If null, it indicates that they have never logged in.

* `user_id` - Indicates the IAM user ID.

* `tags` - Indicates the Custom Tag List.
  The [tags](#IdentityUserV5_Tags) structure is documented below.

<a name="IdentityUserV5_Tags"></a>
The `tags` block contains:

* `tag_key` - Indicates the service ID.

* `tag_value` - Indicates the service ID.


