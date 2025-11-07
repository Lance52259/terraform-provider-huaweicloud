---
subcategory: "Cloud Connect (CC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cc_bandwidth_package_sites"
description: |-
  Use this data source to get the list of the bandwidth package site.
---

# huaweicloud_cc_bandwidth_package_sites

Use this data source to get the list of the bandwidth package site.

## Example Usage

```hcl
variable "site_code" {}
variable "region_id" {}
variable "name" {}

data "huaweicloud_cc_bandwidth_package_sites" "test" {
  
}

data "huaweicloud_cc_bandwidth_package_sites" "test1" {
  site_code = var.site_code
}

data "huaweicloud_cc_bandwidth_package_sites" "test2" {
  region_id = var.region_id
}

data "huaweicloud_cc_bandwidth_package_sites" "test3" {
  name = var.name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `site_code` - (Optional, String) Specifies Site code.

* `region_id` - (Optional, String) Specifies Region ID for query.

* `name` - (Optional, String) Specifies Name used for fuzzy search.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `bandwidth_package_sites` - The list of the bandwidth package site.

  The [bandwidth_package_sites](#bandwidth_package_sites_struct) structure is documented below.

<a name="bandwidth_package_sites_struct"></a>
The `bandwidth_package_sites` block supports:

* `id` - ID of the bandwidth package site.

* `site_code` - Code of the bandwidth package site.

* `region_id` - Region ID.

* `site_type` - Site type. The default type is region.

* `name_cn` - Site name in chinese.

* `name_en` - Site name in english.

* `description` - Description. Angle brackets (<>) are not allowed.

* `created_at` - Creation time of the site.

* `updated_at` - Update time of the site.
