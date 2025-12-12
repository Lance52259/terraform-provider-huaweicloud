---
subcategory: "Cloud Container Engine Autopilot (CCE Autopilot)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cce_autopilot_chart"
description: |- 
  Use this data source to get  the list of CCE Autopilot charts within huaweicloud.
---

# huaweicloud_cce_autopilot_chart

获取模板列表

## Example Usage

```hcl
data "huaweicloud_cce_autopilot_chart" "test" {
  region = "cn-north-4"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource. If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.
