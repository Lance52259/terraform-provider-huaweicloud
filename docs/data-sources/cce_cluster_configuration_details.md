---
subcategory: "Cloud Container Engine (CCE)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cce_cluster_configurationdetails"
description: |-
  Use this data source to get the configuration details within HuaweiCloud.
---

# huaweicloud_cce_cluster_configurationdetails

Use this data source to get the configurations of a CCE cluster within HuaweiCloud.

## Example Usage

```hcl

data "huaweicloud_cce_cluster_configurationdetails" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to obtain the CCE cluster configurations. If omitted, the
  provider-level region will be used.


## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `configurations` - The map of configurations.
