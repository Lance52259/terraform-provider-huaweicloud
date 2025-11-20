---
subcategory: "Cloud Connect (CC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cc_bandwidth_package_region_specifications"
description: |-
  Use this data source to get the list of the multi-city bandwidth package specifications.
---

# huaweicloud_cc_bandwidth_package_region_specifications

Use this data source to get the list of the multi-city bandwidth package specifications.

## Example Usage

```hcl
variable "local_region_id" {}
variable "remote_region_id" {}

data "huaweicloud_cc_bandwidth_package_region_specifications" "test" {}

data "huaweicloud_cc_bandwidth_package_region_specifications" "test1" {
  local_region_id = var.local_region_id
}

data "huaweicloud_cc_bandwidth_package_region_specifications" "test2" {
  remote_region_id = var.remote_region_id
}

data "huaweicloud_cc_bandwidth_package_region_specifications" "test3" {
  local_region_id  = var.local_region_id
  remote_region_id = var.remote_region_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `local_region_id` - (Optional, String) Specifies ID of a region for querying the list of multi-city bandwidth package configurations.

* `remote_region_id` - (Optional, String) Specifies ID of another region for querying the list of multi-city bandwidth package configurations.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `region_specifications` - The list of the multi-city bandwidth package specifications.

  The [region_specifications](#region_specifications_struct) structure is documented below.

<a name="region_specifications_struct"></a>
The `region_specifications` block supports:

* `id` - Specification ID of the bandwidth package for communications between regions.

* `name` - Specification chinese name of the bandwidth package for communications between regions.

* `en_name` - Specification english name of the bandwidth package for communications between regions.

* `local_region_id` - ID of a region of the multi-city bandwidth package configuration.

* `remote_region_id` - ID of another region of the multi-city bandwidth package configuration.

* `spec_codes` - List of the bandwidth package specifications.
  The [spec_codes](#spec_codes_struct) structure is documented below.

<a name="spec_codes_struct"></a>
The `spec_codes` block supports:

* `spec_code` - Specification code of the bandwidth package.

* `billing_mode` - Billing mode of the bandwidth package on the Chinese mainland website or the International website.

* `max_bandwidth` - Maximum bandwidth.

* `mim_bandwidth` - Minimum bandwidth.
