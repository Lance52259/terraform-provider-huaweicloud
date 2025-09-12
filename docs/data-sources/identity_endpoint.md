---
subcategory: "Identity and Access Management (IAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identity_endpoint"
description: "Use this data source to get information of endpoint within HuaweiCloud."
---

# huaweicloud_identity_endpoint

Use this data source to get information of endpoint within HuaweiCloud.

-> **NOTE:** You *must* have IAM read privileges to use this data source.

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
data "huaweicloud_identity_endpoint" "test" {
  token = huaweicloud_identity_user_token.test.token
  endpoint_id = "0046cca357c94165b7a10ec2c01bdf60"
}
```

```hcl
variable "account_name" {}
variable "user_name" {}
variable "password" {}

resource "huaweicloud_identity_user_token" "test" {
  account_name = var.account_name
  user_name    = var.user_name
  password     = var.password
}
data "huaweicloud_identity_endpoint" "test" {
  token = huaweicloud_identity_user_token.test.token
}
```
## Argument Reference

The following arguments are supported:

* `endpoint_id` - (Optional, String) Specifies the terminal node ID to be queried

* `interface` - (Optional, String) Specifies terminal node plane.

* `service_id` - (Optional, String) Specifies the ID of service.

* `token` - (Optional, String) Specifies token of the IAM user (no special permissions required)

## Attribute Reference

* `endpoint` - Indicates service endpoint information.
  The [endpoint](#IdentityEndpoint_Endpoint) structure is documented below.

* `links` - Indicates the resource URL. This will only be included in the result set if the endpoint_id parameter is omitted from the input
  The [links](#IdentityLinks_Links) structure is documented below.

<a name="IdentityEndpoint_Endpoint"></a>
The `endpoint` block contains:

* `service_id` - Indicates the service ID.

* `region_id` - Indicates the region ID.

* `links` - Indicates the resource URL information.

* `id` - Indicates the endpoint ID.

* `interface` - Indicates the terminal node plane  .

* `region` - Indicates the region where the endpoint is located.

* `url` - Indicates the URL or network address of the endpoint.

* `enabled` - Indicates whether the endpoint is accessible and operational.The value can be:
    * **true**: accessible;
    * **false**: operational;

<a name="IdentityLinks_Links"></a>
The `links` block contains:

* `self` - Indicates the resource URL.

* `next` - Indicates the URL of the next adjacent resource. It is **null** if there is no next resource.

* `previous` - Indicates the URL of the previous adjacent resource. It is **null** if there is no previous resource.