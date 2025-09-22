---
subcategory: "Identity and Access Management (IAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identity_catalog"
description: "You *must* have IAM read privileges and scope of token is **project** to use this data source."
---

# huaweicloud_identity_group

Use this data source to get details of the specified IAM user group.

-> **NOTE:** You *must* have IAM read privileges and scope of token is **project** to use this data source.

## Example Usage

```hcl
variable "project_token" {}

data "huaweicloud_identity_catalog" "catalog" {
  project_token = var.project_token
}
```

## Argument Reference

* `project_token` - (Require, String) Token for IAM user. No special permissions required, but the token's scope must be project.

## Attribute Reference

* `catalog` - Indicates service catalog information list
  The [catalog](#IdentityCatalog_Catalog) structure is documented below.

* `link` - Indicates resource link information

<a name="IdentityCatalog_Catalog"></a>
The `catalog` block contains:

* `endpoints` - Indicates the **endpoints** information.
  The [endpoints](#IdentityCatalog_Endpoints) structure is documented below.

* `name` - Indicates the service name.

* `id` - Indicates the service ID.

* `type` - Indicates the service type.

<a name="IdentityCatalog_Endpoints"></a>
The `endpoints` block contains:

* `id` - Indicates the ID of endpoints.

* `interface` - Indicates the endpoint plane.

* `region` - Indicates the region of endpoints.

* `region_id` - Indicates region ID of endpoints.

* `url` - Indicates the url of endpoints.


