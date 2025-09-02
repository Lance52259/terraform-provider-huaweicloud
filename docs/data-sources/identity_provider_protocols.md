---
subcategory: "Identity and Access Management (IAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identity_provider_protocols"
description: |-
  Use this data source to get the list of IAM identity provider protocols.
---

# huaweicloud_identity_provider_protocols

Use this data source to get the list of IAM identity provider protocols.

## Example Usage

```hcl
data "huaweicloud_identity_provider_protocols" "test1" {
  provider_id = "Your identity provider id"
}

data "huaweicloud_identity_provider_protocols" "test2" {
  provider_id = "Your identity provider id"
  protocol_id = "saml"
}
```

## Argument Reference

The following arguments are supported:

* `provider_id` - (Required, String) Name of an identity provider.

* `protocol_id` - (Optional, String) Specifies the protocol id. 
  The valid values are **saml** and **oidc**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `protocols` - Protocol Information List.

The `protocols` block supports:

* `id` - The protocol id.

* `mapping_id` - The mapping id of the protocol.

* `links` - The links of protocol.

The `links` block supports:

* `self` - The protocol resource link.

* `identity_provider` - The identity provider resource link.
