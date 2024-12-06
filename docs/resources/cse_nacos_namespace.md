---
subcategory: "Cloud Service Engine (CSE)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cse_nacos_namespace"
description: |-
  Manages a namespace resource under CSE Nacos microservice engine within HuaweiCloud.
---

# huaweicloud_cse_nacos_namespace

Manages a namespace resource under CSE Nacos microservice engine within HuaweiCloud.

## Example Usage

```hcl
variable "nacos_engine_id" {}
variable "namespace_name" {}

resource "huaweicloud_cse_nacos_namespace" "test" {
  engine_id = var.nacos_engine_id
  name      = var.namespace_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Required, String, ForceNew) Specifies the region where the namespace is located.  
  Changing this will create a new resource.

* `engine_id` - (Required, String, ForceNew) Specifies the engine ID to which the namespace belongs.  
  Changing this will create a new resource.

* `name` - (Required, String) Specifies the name of the dedicated microservice.
  The name can contain `1` to `128` characters, special characters `@#$%^&*` are not allowed.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

## Import

Nacos namespace can be imported using related `engine_id` and their `id`, separated by a slash (/), e.g.

```bash
$ terraform import huaweicloud_cse_microservice.test <engine_id>/<id>
```
