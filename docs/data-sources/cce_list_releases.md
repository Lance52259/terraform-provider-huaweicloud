---
subcategory: "Cloud Container Engine (CCE)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cce_list_releases"
description: |-
  Use this data source to get the list of CCE releases within HuaweiCloud.
---

# `huaweicloud_cce_list_releases`

## Example Usage

```hcl
data "huaweicloud_cce_list_releases" "test" {
  cluster_id = "your_cluster_id"
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
