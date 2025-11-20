---
subcategory: "Resource Governance Center (RGC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rgc_account_enabled_control"
description: |-
  Use this data source to get enabled control information for an managed account in Resource Governance Center.
---

# huaweicloud_rgc_account_enabled_control

Use this data source to get enabled control information for an managed account in Resource Governance Center.

## Example Usage

```hcl
variable "control_id" {}
variable "managed_account_id" {}

data "huaweicloud_rgc_account_enabled_control" "test" {
  control_id         = var.control_id
  managed_account_id = var.managed_account_id
}
```

## Argument Reference

* `region` - (Optional, String) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `control_id` - (Required, String) The ID of the control policy.

* `managed_account_id` - (Required, String) The ID of the managed account.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `control` - Information about the enabled control policy.

The [control](#control) structure is documented below.

<a name="control"></a>
The `control` block supports:

* `name` - The name of the control policy.

* `description` - The description of the control policy.

* `behavior` - The type of the control policy. It can be Proactive, Detective, or Preventive.

* `owner` - The source of the managed account creation.
