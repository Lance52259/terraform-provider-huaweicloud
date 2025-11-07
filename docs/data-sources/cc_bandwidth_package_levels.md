---
subcategory: "Cloud Connect (CC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cc_bandwidth_package_levels"
description: |-
  Use this data source to get the list of the bandwidth package class.
---

# huaweicloud_cc_bandwidth_package_levels

Use this data source to get the list of the bandwidth package class.

## Example Usage

```hcl
variable "level" {}
variable "name" {}

data "huaweicloud_cc_bandwidth_package_levels" "test" {
  
}

data "huaweicloud_cc_bandwidth_package_levels" "test1" {
  level = var.level
}

data "huaweicloud_cc_bandwidth_package_levels" "test2" {
  name = var.name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `level` - (Optional, String) Specifies Bandwidth package class.

* `name` - (Optional, String) Specifies Name used for fuzzy search.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `bandwidth_package_levels` - The list of the bandwidth package class.

  The [bandwidth_package_levels](#bandwidth_package_levels_struct) structure is documented below.

<a name="bandwidth_package_levels_struct"></a>
The `bandwidth_package_levels` block supports:

* `id` - ID of the bandwidth package class.

* `level` - Bandwidth package class.

* `display_priority` - Priority of the bandwidth package. A smaller value indicates a higher priority.

* `name_cn` - Bandwidth package class in chinese.

* `name_en` - Bandwidth package class in english.

* `description` - Description. Angle brackets (<>) are not allowed.

* `created_at` - Creation time of bandwidth package class.

* `updated_at` - Update time of bandwidth package class.
