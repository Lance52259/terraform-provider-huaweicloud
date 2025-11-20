---
subcategory: "Cloud Connect (CC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cc_connection_supported_regions"
description: |-
  Use this data source to get the list of region that the cloud connection supported.
---

# huaweicloud_cc_connection_supported_regions

Use this data source to get the list of region that the cloud connection supported.

## Example Usage

```hcl
data "huaweicloud_cc_connection_supported_regions" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `regions` - The list of region that the cloud connection supported.

  The [regions](#regions_struct) structure is documented below.

<a name="regions_struct"></a>
The `regions` block supports:

* `id` - Region id.

* `name` - Region name.

* `area_id` - Geographic region id.

* `area_name` - Geographic region name.

* `used_scenes` - Cloud Connection application scenarios.
