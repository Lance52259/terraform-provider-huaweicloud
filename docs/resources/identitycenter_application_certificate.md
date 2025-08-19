---
subcategory: "IAM Identity Center"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identitycenter_application_certificate"
description: |-
  Manages an Identity Center application certificate resource within HuaweiCloud.
---

# huaweicloud_identitycenter_application_certificate

Manages an Identity Center application certificate resource within HuaweiCloud.

## Example Usage

```hcl
variable "application_instance_id" {}
variable "instance_id" {}

resource "huaweicloud_identitycenter_application_certificate" "test"{
  application_instance_id = var.application_instance_id
  instance_id =  var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) Specifies the ID of the IAM Identity Center instance.
  Changing this parameter will create a new resource.

* `application_instance_id` - (Required, String, ForceNew) Specifies the ID of the application instance.
  Changing this parameter will create a new resource.

* `application_certificate_status` - (Optional, String) Whether if you need to create an active certificate. Value options:
  + **active**

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `application_instance_certificate` - The application instance certificate.The [application_instance_certificate](#application_instance_certificate_struct) structure is documented below.

<a name="application_instance_certificate_struct"></a>
The `application_instance_certificate` block supports:

* `algorithm` - The algorithm of the application instance certificate.

* `certificate` - The certificate of the application instance.

* `certificate_id` - The ID of the application instance certificate.

* `expiry_date` - The expiry date of the application instance certificate.

* `status` - The status of the application instance certificate.

* `key_size` - The key size of the application instance certificate.

* `issue_date` - The issue date of the application instance certificate.

## Import

The IdentityCenter application certificate can be imported using the `instance_id`, `application_instance_id`
and `certificate_id` separated by slashes, e.g.

```bash
$ terraform import huaweicloud_identitycenter_application_certificate.test <instance_id>/<application_instance_id>/<certificate_id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `application_certificate_status`. It is generally
recommended running `terraform plan` after importing an IdentityCenter user. You can then decide if changes should be
applied to the IdentityCenter user, or the resource definition should be updated to align with the instance. Also, you
can ignore changes as below.

```hcl
resource "huaweicloud_identitycenter_application_certificate" "test" {
  ...

  lifecycle {
    ignore_changes = [
      application_certificate_status,
    ]
  }
}
```
