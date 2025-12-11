---
subcategory: "Cloud Container Engine (CCE)"
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
  cluster_id = "cluster-1234"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource. If omitted, the provider-level region will be used.

* `chart_id` - (Optional, String) 模板ID

* `cluster_id` - (Required, String) 集群ID,获取方式请参见[如何获取接口URI中参数](cce_02_0271.xml)。

* `namespace` - (Optional, String) 模板对应的命名空间

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.
