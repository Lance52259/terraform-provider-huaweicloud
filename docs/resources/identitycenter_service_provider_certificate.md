---
subcategory: "IAM Identity Center"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identitycenter_service_provider_certificate"
description: |
  Manages an Identity Center service provider certificate resource within HuaweiCloud.
---

# huaweicloud_identitycenter_service_provider_certificate

Manages an Identity Center service provider certificate resource within HuaweiCloud.

## Example Usage

```hcl
variable "identity_store_id" {}
 
resource "huaweicloud_identitycenter_service_provider_certificate" "test"{
  identity_store_id = var.identity_store_id
}
```

## Argument Reference

The following arguments are supported:
* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `identity_store_id` - (Required, String, NonUpdatable) Specifies the ID of the identity store.

* `sp_certificate_status` - (Optional, String) Specifies the status of the service provider certificate. Value options: 
   **active**


## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `certificate` - The certificate of the service provider.
  The [certificate](#certificate_struct) structure is documented below.

<a name="certificate_struct"></a>
The `certificate` block supports:

* `state` - The state of the service provider certificate.

* `expiry_date` - The expiry date of the service provider certificate.

* `algorithm` - The algorithm of the service provider certificate.

* `x509certificate` - The certificate of the service provider.

* `certificate_id` - The ID of the service provider certificate.

## Import

The IdentityCenter service provider certificate can be imported using the `identity_store_id` and `certificate_id` separated by a slash, e.g.

```bash
$ terraform import huaweicloud_identitycenter_service_provider_certificate.test <identity_store_id>/<certificate_id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `sp_certificate_status`. It is generally
recommended running `terraform plan` after importing an IdentityCenter service provider certificate. You can then decide if changes should be
applied to the IdentityCenter service provider certificate, or the resource definition should be updated to align with the service provider certificate. Also, you
can ignore changes as below.

```hcl
resource "huaweicloud_identitycenter_service_provider_certificate" "test" {
  ...
 
  lifecycle {
    ignore_changes = [
      sp_certificate_status,
    ]
  }
}
```