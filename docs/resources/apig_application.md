---
subcategory: "API Gateway (APIG)"
---

# huaweicloud_apig_application

Manages an APIG application resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "app_name" {}
variable "app_code" {}

resource "huaweicloud_apig_application" "test" {
  instance_id = var.instance_id
  name        = var.app_name
  description = "Created by script"

  app_codes {
    variable = var.app_code
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the APIG application resource.
  If omitted, the provider-level region will be used.
  Changing this will create a new APIG application resource.

* `name` - (Required, String) Specifies the name of the API application.
  The API group name consists of 3 to 64 characters, starting with a letter.
  Only letters, digits and underscores (_) are allowed.
  Chinese characters must be in UTF-8 or Unicode format.

* `description` - (Optional, String) Specifies the description about the APIG application.
  The description contain a maximum of 255 characters and the angle brackets (< and >) are not allowed.
  Chinese characters must be in UTF-8 or Unicode format.

* `instance_id` - (Required, String, ForceNew) Specifies an ID of the APIG dedicated instance to which the APIG
  application belongs to.
  Changing this will create a new APIG application resource.

* `app_codes` - (Required, List) Specifies an array of one or more application codes which the APIG application
  belongs to.
  Up to five application codes can be created.
  The `app_codes` object structure is documented below.

The `app_codes` block supports:

* `code` - (Optional, String) Specifies the application code.
  The code consists of 64 to 180 characters, starting with a letter, digit, plus sign (+) or slash (/).
  Only letters, digits and following special special characters are allowed: !@#$%+-_/=

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the APIG application.
* `app_key` - The APIG application key.
* `app_secret` - The APIG application secret.
* `registraion_time` - Registration time, in RFC-3339 format.
* `update_time` - Time when the API group was last modified, in RFC-3339 format.
* `app_codes/id` - ID of the APIG application code.

## Import

APIG Applications can be imported using their `id` and the ID of the APIG instance to which the application belongs,
separated by a slash, e.g.
```
$ terraform import huaweicloud_apig_application.test <instance id>/<id>
```
