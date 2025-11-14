---
subcategory: "Cloud Connect (CC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cc_bandwidth_package_supported_areas"
description: |-
  Use this data source to get the list of geographic region that the bandwidth package supported.
---

# huaweicloud_cc_bandwidth_package_supported_areas

Use this data source to get the list of geographic region that the bandwidth package supported.

## Example Usage

```hcl
variable "station" {}

data "huaweicloud_cc_bandwidth_package_supported_areas" "test" {}

data "huaweicloud_cc_bandwidth_package_supported_areas" "test1" {
  station = var.station
}

```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `station` - (Optional, String) Specifies Site.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `areas` - The list of geographic region that the bandwidth package supported.

  The [areas](#areas_struct) structure is documented below.

<a name="areas_struct"></a>
The `areas` block supports:

* `id` - Geographic region id.

* `name` - Geographic region name in Chinese.

* `en_name` - Geographic region name in English.

* `station` - Site.
