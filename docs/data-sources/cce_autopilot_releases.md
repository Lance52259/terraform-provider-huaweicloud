---
subcategory: "Cloud Container Engine Autopilot (CCE Autopilot)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cce_autopilot_releases"
description: |- 
  Manages the CCE autopilot releases data source in HuaweiCloud.
---

# `huaweicloud_cce_autopilot_releases`

This data source provides details about the CCE autopilot releases in HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_cce_autopilot_releases" "test" {
  cluster_id = var.cluster_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource. If omitted, the provider-level region will be used.

* `chart_id` - (Optional, String) Chart ID

* `cluster_id` - (Required, String) Clsuter ID

* `namespace` - (Optional, String) Namespace corresponding to the template

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.
