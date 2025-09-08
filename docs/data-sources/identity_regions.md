---
subcategory: "Identity and Access Management (IAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identity_regions"
description: ""
---

# huaweicloud_identity_regions

Use this data source to query the region list within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_identity_regions" "test1" {
}

data "huaweicloud_identity_regions" "test2" {
  region_id = "cn-north-4"
}
```

## Argument Reference

* `region_id` - (Optional, String) Specifies the region id.

## Attribute Reference

* `regions` - The region info list.

The `regions` block contains:

* `id` - Region Id.

* `type` - Region Type.

* `description` - Region Description.

* `link` - Resource link.

* `locales` - The Map of Region Name.