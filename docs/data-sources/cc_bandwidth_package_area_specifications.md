---
subcategory: "Cloud Connect (CC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cc_bandwidth_package_area_specifications"
description: |-
  Use this data source to get the list of bandwidth package specifications in a geographic region.
---

# huaweicloud_cc_bandwidth_package_area_specifications

Use this data source to get the list of bandwidth package specifications in a geographic region.

## Example Usage

```hcl
variable "local_area_id" {}
variable "remote_area_id" {}

data "huaweicloud_cc_bandwidth_package_area_specifications" "test" {
  
}

data "huaweicloud_cc_bandwidth_package_area_specifications" "test1" {
  local_area_id = var.local_area_id
}

data "huaweicloud_cc_bandwidth_package_area_specifications" "test2" {
  remote_area_id = var.remote_area_id
}

data "huaweicloud_cc_bandwidth_package_area_specifications" "test3" {
  local_area_id = var.local_area_id
  remote_area_id = var.remote_area_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `local_area_id` - (Optional, String) Specifies ID of a geographic region for querying the list of bandwidth package specifications.

* `remote_area_id` - (Optional, String) Specifies ID of another geographic region for querying the list of bandwidth package specifications.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `area_specifications` - The list of bandwidth package specifications in a geographic region.

  The [area_specifications](#area_specifications_struct) structure is documented below.

<a name="area_specifications_struct"></a>
The `area_specifications` block supports:

* `id` - Specification ID of the bandwidth package for communications between regions.

* `local_area_id` - ID of a geographic region for querying the list of bandwidth package specifications.

* `remote_area_id` - ID of another geographic region for querying the list of bandwidth package specifications.

* `spec_codes` - List of the bandwidth package specifications.
  The [spec_codes](#spec_codes_struct) structure is documented below.

<a name="spec_codes_struct"></a>
The `spec_codes` block supports:

* `spec_code` - Specification code of the bandwidth package.

* `billing_mode` - Billing mode of the bandwidth package on the Chinese mainland website or the International website.

* `max_bandwidth` - Maximum bandwidth.

* `mim_bandwidth` - Minimum bandwidth.
