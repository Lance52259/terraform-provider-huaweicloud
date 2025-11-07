---
subcategory: "Cloud Connect (CC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cc_bandwidth_package_lines"
description: |-
  Use this data source to get the list of the bandwidth package line.
---

# huaweicloud_cc_bandwidth_package_lines

Use this data source to get the list of the bandwidth package line.

## Example Usage

```hcl
data "huaweicloud_cc_bandwidth_package_lines" "test" {
  
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `bandwidth_package_lines` - The list of the bandwidth package line.

  The [bandwidth_package_lines](#bandwidth_package_lines_struct) structure is documented below.

<a name="bandwidth_package_lines_struct"></a>
The `bandwidth_package_lines` block supports:

* `local_region_id` - ID of a region of the bandwidth package line.

* `remote_region_id` - ID of another region of the bandwidth package line.

* `local_site_code` - Code of a site of the bandwidth package line.

* `remote_site_code` - Code of another site of the bandwidth package line.

* `support_levels` - List of supported classes.

* `spec_codes` - List of offering code.

  The [spec_codes_lines](#spec_codes_lines_struct) structure is documented below.

<a name="spec_codes_lines_struct"></a>
The `spec_codes_lines` block supports:

* `name_cn` - Specification code name of the bandwidth package line in chinese.

* `name_en` - Specification code name of the bandwidth package line in english.

* `level` - Bandwidth package class.

* `spec_code` - Offering code.

* `max_bandwidth` - Maximum bandwidth.

* `mim_bandwidth` - Minimum bandwidth.

* `support_billing_modes` - Supported offering billing mode.
