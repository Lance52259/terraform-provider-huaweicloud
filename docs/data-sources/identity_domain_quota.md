---
subcategory: "Identity and Access Management (IAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identity_domain_quota"
description: "Use this data source to get domain(account) quota."
---

# huaweicloud_identity_domain_quota

Use this data source to get domain(account) quota.

## Example Usage

```hcl
data "huaweicloud_identity_domain_quota" "test" {}

data "huaweicloud_identity_domain_quota" "test" {
  type = "user"
}
```

## Argument Reference

* `type` - (Optional, String) Specifies the type of quota.
  The valid values are user, group, idp, agency, policy, assigment_group_mp, assigment_agency_mp, assigment_group_ep, assigment_user_ep and mapping.

## Attribute Reference

* `resources` - Indicates the resource info list.
  The [resources](#IdentityDomainQuota_Resources) structure is documented below.

<a name="IdentityDomainQuota_Resources"></a>
The `resources` block contains:

* `type` - Indicates the type of resource

* `max` - Indicates the maximum Quota

* `min` - Indicates the minimum Quota

* `quota` - Indicates the current quota

* `used` - Indicates the used quota
