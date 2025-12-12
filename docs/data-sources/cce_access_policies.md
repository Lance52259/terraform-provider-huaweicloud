---
subcategory: "Cloud Container Engine (CCE)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cce_access_policies"
description: |-
  Use this data source to get the access policies.
---

# huaweicloud_cce_access_policies

Use this data source to get the access policies.

## Example Usage

```hcl

data "huaweicloud_cce_access_policies" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `kind` - access policies type.

* `apiVersion` - The access policies api version.

* `accessPolicyList` - data in accessPolicyList.yaml

  The [accessPolicyList](#values_struct) structure is documented below.

<a name="accessPolicyList_struct"></a>
The `accessPolicyList` block supports:

* `id` - access policy id