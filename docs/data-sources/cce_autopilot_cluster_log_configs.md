---
subcategory: "Cloud Container Engine Autopilot (CCE Autopilot)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cce_autopilot_cluster_log_configs"
description: |-
  Use this data source to get  the list of CCE Autopilot clusters log configs within huaweicloud.
---

# huaweicloud_cce_autopilot_cluster_log_configs

## Example Usage

```hcl
data "huaweicloud_cce_autopilot_cluster_log_configs" "test" {
  cluster_id = var.cluster_id
  project_id = var.project_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource. If omitted, the provider-level region will be used.

* `cluster_id` - (Required, String) Cluster ID, which can be obtained from the parameter in the interface URI.

* `project_id` - (Required, String) Project ID, which can be obtained from the parameter in the interface URI.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `ttl_in_days` - The Storage duration.

* `log_configs` - Log config information.

  The [log_configs](#log_configs_struct) structure is documented below.

<a name="log_configs_struct"></a>
The `log_configs` block supports:

* `name` - Last config server name.

* `enable` - Collected or Not.
