---
subcategory: "Cloud Connect (CC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cc_connection_capabilities"
description: |-
  Use this data source to get the list of cloud connection capabilities.
---

# huaweicloud_cc_connection_capabilities

Use this data source to get the list of cloud connection capabilities.

## Example Usage

```hcl
variable "resource_type" {}

data "huaweicloud_cc_connection_capabilities" "test" {}

data "huaweicloud_cc_connection_capabilities" "test1" {
  resource_type = var.resource_type
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `quota_type` - (Optional, String) Specifies the Resource type. 

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `capabilities` - The capability list.

  The [capabilities](#capabilities_struct) structure is documented below.

<a name="capabilities_struct"></a>
The `capabilities` block supports:

* `id` - Capability ID.

* `domain_id` - Domain ID.

* `resource_type` - Resource type.

* `description` - Description of the capability.

* `is_support` - Whether cloud connection capabilities are supported.

* `support_regions` - List of supported regions.

* `support_sites` - List of supported sites.

* `created_at` - Creation time of the capability.

* `updated_at` - Update time of the capability.

* `bandwidth` - The bandwidth capabilities.

  The [bandwidth](#bandwidth_struct) structure is documented below.

<a name="bandwidth_struct"></a>
The `bandwidth` block supports:

* `min` - Minimum bandwidth.

* `max` - Maximum bandwidth.
