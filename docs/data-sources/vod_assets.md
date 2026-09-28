---
subcategory: "Video on Demand (VOD)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_vod_assets"
description: |-
  Use this data source to get the list of VOD assets within HuaweiCloud.
---

# huaweicloud_vod_assets

Use this data source to get the list of VOD assets within HuaweiCloud.

## Example Usage

```hcl
variable "asset_id" {}

data "huaweicloud_vod_assets" "test" {
  asset_ids = [var.asset_id]
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the VOD assets.
  If omitted, the provider-level region will be used.

* `asset_ids` - (Optional, List) Specifies the asset IDs used to query. A maximum of `10` IDs can be specified.

* `status` - (Optional, List) Specifies the asset statuses used to query.
  Multiple values can be specified at the same time. Valid values are as follows:
  + **CREATING**
  + **FAILED**
  + **CREATED**
  + **PUBLISHED**
  + **WAITING_TRANSCODE**
  + **TRANSCODING**
  + **TRANSCODE_SUCCEED**
  + **TRANSCODE_FAILED**
  + **THUMBNAILING**
  + **THUMBNAIL_SUCCEED**
  + **THUMBNAIL_FAILED**
  + **UN_REVIEW**
  + **REVIEWING**
  + **REVIEW_SUSPICIOUS**
  + **REVIEW_PASSED**
  + **REVIEW_FAILED**
  + **REVIEW_BLOCKED**

* `start_time` - (Optional, String) Specifies the start time used to query, in the format `yyyymmddhhmmss` (UTC).

* `end_time` - (Optional, String) Specifies the end time used to query, in the format `yyyymmddhhmmss` (UTC).

* `category_id` - (Optional, Int) Specifies the category ID used to query.

* `tags` - (Optional, String) Specifies the asset tags used to query.
  A single tag cannot exceed `16` bytes, and a maximum of `16` tags are supported.
  Multiple tags are separated by commas and encoded in UTF-8.

* `query_string` - (Optional, String) Specifies the fuzzy query string used to match asset title, description,
  or category name.

* `media_types` - (Optional, List) Specifies the media types used to query. A maximum of `20` types can be specified.
  Valid values include:
  + Video: **MP4**, **TS**, **MOV**, **MXF**, **MPG**, **FLV**, **WMV**, **AVI**, **M4V**, **F4V**, **MPEG**
  + Audio: **MP3**, **OGG**, **WAV**, **WMA**, **APE**, **FLAC**, **AAC**, **AC3**, **MMF**, **AMR**, **M4A**,
    **M4R**, **WV**, **MP2**

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `assets` - All VOD assets that match the filter parameters.
  The [assets](#vod_assets) structure is documented below.

<a name="vod_assets"></a>
The `assets` block supports:

* `asset_id` - The asset ID.

* `title` - The asset title.

* `description` - The asset description.

* `duration` - The asset duration, in seconds.

* `duration_ms` - The asset duration, in milliseconds.

* `size` - The asset size, in bytes.

* `original_url` - The original playback URL of the asset.

* `category` - The category name of the asset.

* `covers` - The cover URLs of the asset.

* `create_time` - The creation time of the asset, in the format `yyyymmddhhmmss` (UTC).

* `asset_status` - The status of the asset. Valid values are as follows:
  + **CREATING**
  + **FAILED**
  + **CREATED**
  + **PUBLISHED**
  + **DELETED**

* `transcode_status` - The transcoding status of the asset. Valid values are as follows:
  + **UN_TRANSCODE**
  + **WAITING_TRANSCODE**
  + **TRANSCODING**
  + **TRANSCODE_SUCCEED**
  + **TRANSCODE_FAILED**

* `thumbnail_status` - The thumbnail status of the asset. Valid values are as follows:
  + **UN_THUMBNAIL**
  + **THUMBNAILING**
  + **THUMBNAIL_SUCCEED**
  + **THUMBNAIL_FAILED**

* `review_status` - The review status of the asset. Valid values are as follows:
  + **UN_REVIEW**
  + **REVIEWING**
  + **REVIEW_SUSPICIOUS**
  + **REVIEW_PASSED**
  + **REVIEW_FAILED**
  + **REVIEW_BLOCKED**

* `exec_desc` - The execution description of the asset tasks.

* `media_type` - The media type of the asset.
