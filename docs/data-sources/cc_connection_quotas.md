---
subcategory: "Cloud Connect (CC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cc_connection_quotas"
description: |-
  Use this data source to get the list of cloud connection quotas.
---

# huaweicloud_cc_connection_quotas

Use this data source to get the list of cloud connection quotas.

## Example Usage

```hcl
variable "cloud_connection_id" {}
variable "region_id" {}

data "huaweicloud_cc_connection_quotas" "test" {
  quota_type = "cloud_connection"
}

data "huaweicloud_cc_connection_quotas" "test1" {
  quota_type = "cloud_connection_region"
  cloud_connection_id = var.cloud_connection_id
}

data "huaweicloud_cc_connection_quotas" "test2" {
  quota_type = "cloud_connection_route"
  cloud_connection_id = var.cloud_connection_id
}

data "huaweicloud_cc_connection_quotas" "test3" {
  quota_type = "region_network_instance"
  cloud_connection_id = var.cloud_connection_id
  region_id = var.region_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `quota_type` - (Optional, String) Specifies the quota type. The value can be any one in [cloud_connection, cloud_connection_region, cloud_connection_route, region_network_instance].
  cloud_connection: the maximum number of cloud connections that can be created in an account;
  cloud_connection_region: the maximum number of regions where a cloud connection can be used;
  cloud_connection_route: the maximum number of routes that can be added to a cloud connection;
  region_network_instance: the maximum number of network instances that can be loaded to a cloud connection in a region.

* `cloud_connection_id` - (Optional, String) Specifies the Cloud connection ID.
  This parameter is mandatory when you query the value of each of the three parameters: cloud_connection_region, cloud_connection_route, and region_network_instance.

* `region_id` - (Optional, String) Specifies Region ID. This parameter is mandatory when you query the value of region_network_instance

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `quotas` - The quota list.

  The [quotas](#quotas_struct) structure is documented below.

<a name="quotas_struct"></a>
The `quotas` block supports:

* `domain_id` - Domain ID.

* `quota_type` - Quota type.

* `quota_number` - Quota number.

* `quota_used` - Quota used number.

* `cloud_connection_id` - Cloud connection ID.

* `region_id` - Region ID.
