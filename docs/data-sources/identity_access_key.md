---
subcategory: "Identity and Access Management (IAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identity_access_key"
description: ""
---

# huaweicloud_identity_access_key

Query all permanent access key and specified access key data source within Huaweicloud IAM service.

## Example Usage

```hcl
resource "huaweicloud_identity_user" "test_user" {
  name        = "my_test_huaweicloud"
  password    = "password@123!"
  enabled     = true
  description = "tested by terraform"
}
resource "huaweicloud_identity_access_key" "key_1" {
  user_id     = huaweicloud_identity_user.test_user.id
}
data "huaweicloud_identity_access_key" "test" {
  user_id = huaweicloud_identity_access_key.key_1.user_id
}

```
```hcl
resource "huaweicloud_identity_user" "test_user" {
  name        = "my_test_huaweicloud"
  password    = "password@123!"
  enabled     = true
  description = "tested by terraform"
}
resource "huaweicloud_identity_access_key" "key_1" {
  user_id  = huaweicloud_identity_user.test_user.id
}
data "huaweicloud_identity_access_key" "test" {
  access_key = huaweicloud_identity_access_key.key_1.id
}
```
```hcl
data "huaweicloud_identity_access_key" "test3" {
}
```

## Argument Reference

* `user_id` - (Optional, String) The IAM user ID associated with the access_key.

* `access_key` - (Optional, String) access_key associated with the identity user.

## Attribute Reference
* `credentials` -  The credentials carrying user information.

* `user_id` -  The IAM user ID to which the access_key belongs.

* `access` - The Access Key (AK) that was queried.

* `status` - The status of the access_key. Valid values: **active**, The access_key is enabled. **inactive**: The access_key is disabled.

* `create_time` - The creation time of the access_key. The value is a UTC time in the YYYY-MM-DDTHH:mm:ss.ssssssZ format.

* `description` - The description of the access_key.

* `last_use_time` - The field is only present when you are certified by access_key

