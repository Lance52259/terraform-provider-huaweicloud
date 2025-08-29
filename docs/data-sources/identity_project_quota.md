---
subcategory: "Identity and Access Management (IAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identity_project_quota"
description: ""
---

# huaweicloud_identity_project_quota

Use this data source to get project quota.

-> **NOTE:** You *must* have IAM read privileges to use this data source.

## Example Usage

```hcl
data "huaweicloud_identity_project_quota" "test" {
  project_id = var.project_id
}
```

## Argument Reference

* `project_id` - (Required, String) Specifies the project id.

## Attribute Reference

The `resources` block contains:

* `type` - The type of resource, must be **project**

* `max` - Maximum Quota

* `min` - Minimum Quota

* `quota` - Current quota

* `used` - Used quota
