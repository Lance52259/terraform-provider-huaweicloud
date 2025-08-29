---
subcategory: "Identity and Access Management (IAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identity_domain_quota"
description: ""
---

# huaweicloud_identity_domain_quota

Use this data source to get domain(account) quota.

-> **NOTE:** You *must* have IAM read privileges to use this data source.

## Example Usage

```hcl
data "huaweicloud_identity_domain_quota" "test" {
}

data "huaweicloud_identity_domain_quota" "test" {
  type = "user"
}
```

## Argument Reference

* `type` - (Optional, String) Specifies the type of quota.
  The valid values are user, group, idp, agency, policy, assigment_group_mp, assigment_agency_mp, assigment_group_ep, assigment_user_ep and mapping.

## Attribute Reference

The `resources` block contains:

* `type` - The type of resource

* `max` - Maximum Quota

* `min` - Minimum Quota

* `quota` - Current quota

* `used` - Used quota
