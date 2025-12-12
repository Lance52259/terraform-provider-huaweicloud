---
subcategory: "Cloud Container Engine (CCE)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cce_chart_values"
description: |-
  Use this data source to get the list of CCE chart values within HuaweiCloud.
---

# Resource HuaweiCloud Cce Show Chart Values

## Example Usage

```hcl
data "huaweicloud_cce_chart_values" "test" {
  chart_id = var.chart_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource. If omitted, the provider-level region will be used.

* `chart_id` - (Required, String) Chart ID

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `values` - data in values.yaml

  The [values](#values_struct) structure is documented below.

<a name="values_struct"></a>
The `values` block supports:

* `basic` -
