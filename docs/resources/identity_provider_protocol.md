---
subcategory: "Identity and Access Management (IAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identity_provider_protocol"
description: ""
---

# huaweicloud_identity_provider_protocol

Manage the protocol of identity provider within HuaweiCloud IAM service.

When you create identity provider through console or terraform resource `huaweicloud_identity_provider`, it will create 
a default protocol with a default mapping. If the mapping of the existed protocol is different from the mapping you 
specified, it will try to update the protocol when using `huaweicloud_identity_provider_protocol`.

## Example Usage

```hcl
resource "huaweicloud_identity_provider_protocol" "protocol" {
  provider_id = "Your identity provider id"
  protocol_id = "saml"
  mapping_id = "Your mapping id"
}
```

## Argument Reference

The following arguments are supported:

* `provider_id` - (Required, String) Specifies the ID of the identity provider used to manage the protocol.

* `protocol_id` - (Required, String) Specifies the identity protocol of the identity provider. The content of this field is `saml` or `oidc`.

* `mapping_id` - (Optional, String) Specifies the mapping_id for the protocol. When the identity provider type is `iam_user_sso`, there is no need to bind a mapping ID, and this field does not need to be passed; otherwise, this field is mandatory.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `links` - The links of protocol.

The `links` block supports:

* `self` - The resource link.

* `identity_provider` - The identity provider resource link.

## Import

Identity provider protocol is imported using the `{provider_id}:{protocol_id}`. For example, if you have provider_id `provider_test` and protocol `saml`, you should use `provider_test:saml` to import.

```bash
$ terraform import huaweicloud_identity_provider_protocol.protocol example_com_provider_saml
```
