---
subcategory: "Identity and Access Management (IAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identity_group_enterprise_projects"
description: |-
  Use this data source to get enterprise projects of the specified IAM user group id.
---

# huaweicloud_identity_group_enterprise_projects

Use this data source to get enterprise projects of the specified IAM user group id.

## Example Usage

```hcl
variable "group_id" {}

data "huaweicloud_identity_group_enterprise_project" "group" {
   group_id = var.group_id
}
```

## Argument Reference

* `group_id` - (Required, String) Specifies the ID of the identity group.

## Attribute Reference

* `enterprise_projects` - Indicates the projects contains. 
  The [enterprise_projects](#IdentityEnterpriseProjects_projects) structure is documented below.

<a name="IdentityEnterpriseProjects_projects"></a>
The `enterprise_projects` block contains:

* `project_id` - Indicates the ID of enterprise projects.


