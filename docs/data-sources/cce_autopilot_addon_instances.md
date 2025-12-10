---
subcategory: "CCE Autopilot"
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
  cluster_id = "cluster_id"  # Required
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

* `apiVersion` - API版本,固定值“v3”,该值不可修改。

* `items` - 插件实例列表
  The [items](#items_struct) structure is documented below.

* `kind` - API类型,固定值“Addon”,该值不可修改。

The `items` block supports:

* `apiVersion` - API版本,固定值“v3”,该值不可修改。

* `kind` - API类型,固定值“Addon”,该值不可修改。

* `metadata` - 基本信息,为集合类的元素类型,包含一组由不同名称定义的属性
  The [metadata](#items_metadata_struct) structure is documented below.

* `spec` - spec是集合类的元素类型,内容为插件实例具体信息,实例的详细描述主体部分都在spec中给出
  The [spec](#items_spec_struct) structure is documented below.

* `status` - 插件实例状态
  The [status](#items_status_struct) structure is documented below.

The `metadata` block supports:

* `alias` - 插件别名

* `annotations` - 插件注解,由key/value组成
  - 安装:固定值为{"addon.install/type":"install"}
  - 升级:固定值为{"addon.upgrade/type":"upgrade"}

* `creationTimestamp` - 创建时间

* `labels` - 插件标签,key/value对格式,接口保留字段,填写不会生效

* `name` - 插件名称

* `uid` - 唯一id标识

* `updateTimestamp` - 更新时间

The `spec` block supports:

* `addonTemplateLabels` - 插件模板所属类型

* `addonTemplateLogo` - 插件模板logo图片的地址

* `addonTemplateName` - 插件模板名称,如coredns

* `addonTemplateType` - 插件模板类型

* `clusterID` - 集群id

* `description` - 插件模板描述

* `values` - 插件模板安装参数(各插件不同),请根据具体插件模板信息填写安装参数。

* `version` - 插件模板版本号,如1.0.0

The `status` block supports:

* `Reason` - 插件安装失败原因

* `currentVersion` - 当前插件实例使用的具体插件版本信息
  The [currentVersion](#status_currentVersion_struct) structure is documented below.

* `isRollbackable` - 是否支持回滚到插件升级前的插件版本

* `message` - 安装错误详情

* `previousVersion` - 插件升级或回滚前的版本

* `status` - 插件实例状态, 取值如下
  - running:运行中,表示插件全部实例状态都在运行中,插件正常使用。
  - abnormal:不可用,表示插件状态异常,插件不可使用。可单击插件名称查看实例异常事件。
  - installing:安装中,表示插件正在安装中。
  - installFailed:安装失败,表示插件安装失败,需要卸载后重新安装。
  - upgrading:升级中,表示插件正在更新中。
  - upgradeFailed:升级失败,表示插件升级失败,可重试升级或卸载后重新安装。
  - deleting:删除中,表示插件正在删除中。
  - deleteFailed:删除失败,表示插件删除失败,可重试卸载。
  - deleteSuccess:删除成功,表示插件删除成功。
  - available:部分就绪,表示插件下只有部分实例状态为运行中,插件部分功能可用。
  - rollbacking:回滚中,表示插件正在回滚中。
  - rollbackFailed:回滚失败,表示插件回滚失败,可重试回滚或卸载后重新安装。
  - unknown:未知状态,表示插件模板实例不存在。

* `targetVersions` - 此插件版本,支持升级的集群版本

The `currentVersion` block supports:

* `creationTimestamp` - 创建时间

* `input` - 插件安装参数

* `stable` - 是否为稳定版本

* `supportVersions` - 支持集群版本号
  The [supportVersions](#currentVersion_supportVersions_struct) structure is documented below.

* `translate` - 供界面使用的翻译信息

* `updateTimestamp` - 更新时间

* `version` - 插件版本号

The `supportVersions` block supports:

* `category` - 作用的集群类型
  **取值范围:**
  - CCE:CCE Standard集群
  - Turbo:CCE Turbo集群
  - Autopilot:CCE Autopilot集群

  **默认取值**
  为空时默认为CCE Standard,CCE Turbo集群

* `clusterType` - 支持的集群类型

* `clusterVersion` - 支持的集群版本(正则表达式)
