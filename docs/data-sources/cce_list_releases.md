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
  cluster_id = var.cluster_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource. If omitted, the provider-level region will be used.

* `chart_id` - (Optional, String) Chart ID

* `cluster_id` - (Required, String) Cluster ID

* `namespace` - (Optional, String) Namespace corresponding to the template.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.
