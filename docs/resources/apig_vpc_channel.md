---
subcategory: "API Gateway (APIG)"
---

# huaweicloud_apig_vpc_channel

Manages an APIG VPC channel resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "channel_name" {}
variable "ecs_1" {}
variable "ecs_2" {}

resource "huaweicloud_apig_vpc_channel" "test" {
  instance_id = var.instance_id
  name        = var.app_name
  port        = 8080
  member_type = "ecs"
  protocol    = "HTTPS"
  path        = "/"
  
  members {
    id     = var.ecs_1
    weight = 50
  }
  
  members {
    id     = var.ecs_2
    weight = 50
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the APIG application resource.
  If omitted, the provider-level region will be used.
  Changing this will create a new APIG vpc channel resource.

* `instance_id` - (Required, String, ForceNew) Specifies an ID of the APIG dedicated instance to which the APIG
  vpc channel belongs to.
  Changing this will create a new APIG vpc channel resource.

* `name` - (Required, String) Specifies the name of the VPC channel.
  The channel name consists of 3 to 64 characters, starting with a letter.
  Only letters, digits and underscores (_) are allowed.
  Chinese characters must be in UTF-8 or Unicode format.

* `port` - (Optional, Int) Specifies the host port of the VPC channel.
  The valid value is range from 1 to 65535.

* `member_type` - (Optional, String) Specifies the type of the backend service.
  The valid types are *ecs* and *ip*.

* `algorithm` - (Optional, String) Specifies the type of the backend service.
  The valid types are *WRR*, *WLC*, *SH* and *URI hashing*, default to *WRR*.

* `protocol` - (Optional, String) Specifies the protocol for performing health checks on backend servers in the VPC
  channel.
  The valid values are *TCP*, *HTTP* and *HTTPS*, default to *TCP*.

* `path` - (Optional, String) Specifies the destination path for health checks.
  Required if `protocol` is *HTTP* or *HTTPS*.

* `healthy_threshold` - (Optional, Int) Specifies the healthy threshold, which refers to the number of consecutive
  successful checks required for a backend server to be considered healthy.
  The valid value is range from 2 to 10.

* `unhealthy_threshold` - (Optional, Int) Specifies the unhealthy threshold, which refers to the number of consecutive
  failed checks required for a backend server to be considered unhealthy.
  The valid value is range from 2 to 10.

* `timeout` - (Optional, Int) Specifies the timeout for determining whether a health check fails, in second.
  The value must be less than the value of time_interval.
  The valid value is range from 2 to 30.

* `interval` - (Optional, Int) Specifies the interval between consecutive checks, in second.
  The valid value is range from 5 to 300.

* `servers` - (Optional, List) Specifies an array of one or more backend server IDs that bind the VPC channel.
  The object structure is documented below.

The `servers` block supports:

* `id` - (Optional, String, ForceNew) Specifies the backend server ID.
  Required if `member_type` is *ecs*.
  This parameter and `ip_address` are alternative.

* `ip_address` - (Optional, String, ForceNew) Specifies the backend server address.
  Required if `member_type` is *ip*.

* `weight` - (Optional, Int, ForceNew) Specifies the backend server weight.
  The valid values are range from 1 to 100, default to 1.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the APIG application.
* `create_time` - Time when the channel created, in RFC-3339 format.
* `status` - The status of VPC channel.

## Import

APIG VPC Channels can be imported using their `id` and ID of the APIG dedicated instance to which the application
belongs, separated by a slash, e.g.
```
$ terraform import huaweicloud_apig_vpc_channel.test <instance id>/<id>
```
