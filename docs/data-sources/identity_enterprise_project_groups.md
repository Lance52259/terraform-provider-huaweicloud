---
subcategory: "Identity and Access Management (IAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identity_group_enterprise_projects"
description: "Use this data source to get enterprise projects of the specified IAM user group."
---

# huaweicloud_identity_group_enterprise_projects

Use this data source to get enterprise projects of the specified IAM user group.

## Example Usage

```hcl
variable "enterprise_project_id" {}

data "huaweicloud_identity_enterprise_project_groups" "groups" {
  enterprise_project_id = var.enterprise_project_id
}
```

## Argument Reference

* `enterprise_project_id` - (Optional, String) Specifies the ID of the enterprise project.

## Attribute Reference

* `groups` - Indicates the users the group contains.
  The [groups](#IdentityEnterpriseProjects_Groups) structure is documented below.

<a name="IdentityEnterpriseProjects_Groups"></a>
The `groups` block contains:

* `create_time` - Indicates the time of the IAM user groups creat.

* `description` - Indicates the description of the IAM user groups.

* `domain_id` - Indicates the domain ID of the IAM user groups.

* `id` - Indicates the ID of the IAM user groups.

* `name` - Indicates the name of the IAM user groups.


