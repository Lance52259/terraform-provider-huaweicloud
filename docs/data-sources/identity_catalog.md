---
subcategory: "Identity and Access Management (IAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identity_catalog"
description: |-
  Use this data source to query the service catalog within HuaweiCloud.
---

# huaweicloud_identity_catalog

se this data source to query the service catalog within HuaweiCloud.

-> **NOTE:** You *must* have IAM read privileges and **project** scope to use this data source.

## Example Usage

```hcl
variable "account_name" {}
variable "user_name" {}
variable "password" {}

resource "huaweicloud_identity_user_token" "test" {
  account_name = var.account_name
  user_name    = var.user_name
  password     = var.password
}
data "huaweicloud_identity_catalog" "test" {
  project_token = huaweicloud_identity_user_token.test.token
}
```

## Argument Reference

The following arguments are supported:

* `project_token` - (Require, String) Specifies token for IAM user. No special permissions required, but the token's scope must be project.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `catalog` - Indicates service catalog information list.
  The [catalog](#IdentityCatalog_Catalog) structure is documented below.

<a name="IdentityCatalog_Catalog"></a>
The `catalog` block supports:

* `endpoints` - Indicates service endpoint information.
  The [endpoints](#IdentityCatalog_Enpoints) structure is documented below.

* `id` - Indicates the service id.

* `name` - Indicates the service name.

* `type` - Indicates the service type.

* <a name="IdentityCatalog_Enpoints"></a>
The `catalog` block supports:

* `id` - Indicates the endpoint id.

* `interface` - Indicates the endpoint interface.

* `region` - Indicates the endpoint region.

* `region_id` - Indicates the endpoint region id.

* `url` - Indicates the endpoint url.