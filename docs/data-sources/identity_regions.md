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

* `regions` - Indicates the region info list.
  The [regions](#IdentityRegions_Regions) structure is documented below.

<a name="IdentityRegions_Regions"></a>
The `regions` block contains:

* `id` - Indicates the Region id.

* `type` - Indicates the Region Type.

* `description` - Indicates the Region Description.

* `link` - Indicates the Resource link.

* `locales` - Indicates the Map of Region Name.