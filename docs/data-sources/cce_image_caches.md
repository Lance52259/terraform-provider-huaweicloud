---
subcategory: "Cloud Container Engine (CCE)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cce_image_caches"
description: |-
  Use this data source to get the image caches.
---

# huaweicloud_cce_image_caches

Use this data source to get the image caches.

## Example Usage

```hcl

data "huaweicloud_cce_image_caches" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The image cache ID.
