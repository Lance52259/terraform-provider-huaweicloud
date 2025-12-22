---
subcategory: "Cloud Container Engine (CCE)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cce_flavor_specifications"
description: |-
  Use this data source to get the list of CCE Cluster Master flavor specifications within HuaweiCloud.
---

# huaweicloud_cce_flavor_specifications

Use this data source to get the list of CCE Cluster Master flavor specifications within HuaweiCloud.

## Example Usage

```hcl
variable "cluster_type" {}

data "huaweicloud_cce_flavor_specifications" "test" {
  cluster_type = var.cluster_type
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource. If omitted, the provider-level region will be used.

* `cluster_type` - (Required, String) Specifies the type of cluster.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `name` - The cluster flavor name.

* `node_capacity` - The cluster node capacity.

* `is_sold_out` - The flavor is soldout or not.

* `is_support_multi_az` - The support multi az data.

* `available_master_flavors` - The available_master_flavors data of cce cluster.

  The [available_master_flavors](#available_master_flavors_struct) structure is documented below.

<a name="available_master_flavors_struct"></a>
The `available_master_flavors` block supports:

* `name` - The cluster flavor name.

* `azs` - The azs of master flvor.

* `az_fault_domains` - The fault domain az data.
