 ---
subcategory: "Identity and Access Management (IAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identity_key"
description: ""
---

# huaweicloud_identity_credentials

Query all permanent access key data source within Huaweicloud

-> **NOTE:** You *must* have admin privileges to use this data source.

## Example Usage

```hcl
data "huaweicloud_identity_credentials" "test" {

}

```
```hcl
data "huaweicloud_identity_credentials" "test" {
    user_id = "d9a1d990bc4e4dfeb30c25ea4020e007"
}

```

## Argument Reference

* `user_id` - (Optional, String) The IAM user ID associated with the credentials.

## Attribute Reference

* `user_id` -  The IAM user ID to which the credentials belongs.

* `access` - The Access Key (AK) that was queried.

* `status` - The status of the credentials. Valid values: **active**, The credentials is enabled. **inactive**: The credentials is disabled.

* `create_time` - The creation time of the credentials. The value is a UTC time in the YYYY-MM-DDTHH:mm:ss.ssssssZ format.

* `description` - The description of the credentials.


