---
subcategory: "Cloud Container Engine (CCE)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cce_flavor_specifications"
description: |-
  Use this data source to get the list of CCE Cluster Master flavor specifications within HuaweiCloud.
---

# Resource HuaweiCloud Cce Show Flavor Specifications

## Example Usage

```hcl
data "huaweicloud_cce_flavor_specifications" "test" {
  clusterType = ""
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource. If omitted, the provider-level region will be used.

* `clusterType` - (Required, String) 集群类型

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `name` - The data source ID.
* `nodeCapacity` - The data source ID.
* `isSoldOut` - The data source ID.
* `isSupportMultiAZ` - The data source ID.

* `availableMasterFlavors` - availableMasterFlavors.yaml中的数据,数据结构以具体的模板为准

  The [availableMasterFlavors](#availableMasterFlavors_struct) structure is documented below.

<a name="availableMasterFlavors_struct"></a>
The `availableMasterFlavors` block supports:

* `name` -

* `azs` -

* `azFaultDomains` -
