---
subcategory: "Cloud Container Engine Autopilot (CCE Autopilot)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cce_autopilot_addon_instances"
description: |-
  Use this datasource to get the list of CCE Autopilot Addon Instances.
---

# huaweicloud_cce_autopilot_addon_instances

Use this datasource to get the list of CCE Autopilot Addon Instances.

## Example Usage

```hcl
data "huaweicloud_cce_autopilot_addon_instances" "test" {
  cluster_id = var.cluster_id  # Required
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `cluster_id` - (Required, String) The ID of the cluster.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `apiVersion` - API version, default value "v3"

* `items` - addon instance list
  The [items](#items_struct) structure is documented below.

* `kind` - API type, default value "Addon"

The `items` block supports:

* `apiVersion` - API version, default value "v3"

* `kind` - API type, default value "Addon"

* `metadata` - Basic information, which is the element type of the collection class, containing a set of attributes defined by different names.
  The [metadata](#items_metadata_struct) structure is documented below.

* `spec` - Indicates the element type of the collection class
  The [spec](#items_spec_struct) structure is documented below.

* `status` - addon instance status
  The [status](#items_status_struct) structure is documented below.

The `metadata` block supports:

* `alias` - addon alias

* `annotations` - addon annotation
  - install:default value is {"addon.install/type":"install"}
  - upgrade:default value is {"addon.upgrade/type":"upgrade"}

* `creationTimestamp` - creation time

* `labels` - addon labels

* `name` - addon name

* `uid` - addon uid

* `updateTimestamp` - update time

The `spec` block supports:

* `addonTemplateLabels` - addon template labels

* `addonTemplateLogo` - addon template logo

* `addonTemplateName` - addon template name,  for example coredns

* `addonTemplateType` - addon template type

* `clusterID` - cluster id

* `description` - addon description

* `values` - addon template installation parameters

* `version` - addon version, for example 1.0.0

The `status` block supports:

* `Reason` - addon installation failure cause

* `currentVersion` - addon current version
  The [currentVersion](#status_currentVersion_struct) structure is documented below.

* `isRollbackable` - indicates whether the addon version can be rolled back to the source version.

* `message` - installation error details

* `previousVersion` - version before addon upgrade or rollback

* `status` - addon Instance Status, the values are as follows:
  - running
  - abnormal
  - installing
  - installFailed
  - upgrading
  - upgradeFailed
  - deleting
  - deleteFailed
  - deleteSuccess
  - available
  - rollbacking
  - rollbackFailed
  - unknown

* `targetVersions` - This addon version is the cluster version that can be upgraded.

The `currentVersion` block supports:

* `creationTimestamp` - creation time

* `input` - addon installation parameters

* `stable` - is or not the stable version 

* `supportVersions` - supported cluster version
  The [supportVersions](#currentVersion_supportVersions_struct) structure is documented below.

* `translate` - translate info

* `updateTimestamp` - update time

* `version` - addon version

The `supportVersions` block supports:

* `category` - cluster type
  **Value Range:**
  - CCE:CCE Standard Cluster
  - Turbo:CCE Turbo Cluster
  - Autopilot:CCE Autopilot Cluster

  **Default Value**
  default value is CCE Standard,CCE Turbo Cluster

* `clusterType` - cluster type

* `clusterVersion` - cluster version
