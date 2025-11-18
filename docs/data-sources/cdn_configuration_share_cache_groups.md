---
subcategory: "Content Delivery Network (CDN)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cdn_configuration_share_cache_groups"
description: |-
  Use this data source to get a list of share cache groups within HuaweiCloud.
---

# huaweicloud_cdn_configuration_share_cache_groups

Use this data source to get a list of share cache groups within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_cdn_configuration_share_cache_groups" "test" {
}

output "groups" {
  value = data.huaweicloud_cdn_configuration_share_cache_groups.test.share_cache_groups
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the share cache groups are located.  
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `share_cache_groups` - The list of the share cache groups.  
  The [share_cache_groups](#cdn_configuration_share_cache_groups_share_cache_groups) structure is documented below.

<a name="cdn_configuration_share_cache_groups_share_cache_groups"></a>
The `share_cache_groups` block supports:

* `id` - The ID of the share cache group.

* `group_name` - The name of the share cache group.

* `primary_domain` - The primary domain name.

* `share_cache_records` - The list of associated domain names.  
  The [share_cache_records](#cdn_configuration_share_cache_groups_share_cache_groups_share_cache_records) structure is
  documented below.

* `create_time` - The creation time of the share cache group, in RFC3339 format.

<a name="cdn_configuration_share_cache_groups_share_cache_groups_share_cache_records"></a>
The `share_cache_records` block supports:

* `domain_name` - The associated domain name.
