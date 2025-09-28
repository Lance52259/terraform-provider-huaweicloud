---
subcategory: "Identity and Access Management (IAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identity_user_enterprise_projects"
description: "Use this data source to get enterprise projects of the specified IAM user id."
---

# huaweicloud_identity_user_enterprise_projects

Use this data source to get enterprise projects of the specified IAM user id

## Example Usage

```hcl
variable "user_id" {}

data "huaweicloud_identity_user_enterprise_project" "user" {
  user_id = var.user_id
}
```

## Argument Reference

* `user_id` - (Optional, String) Specifies the ID of the identity user.

## Attribute Reference

* `enterprise_projects` - Indicates the users the group contains.
  The [enterprise_projects](#IdentityEnterpriseProjects_projects) structure is documented below.

<a name="IdentityEnterpriseProjects_projects"></a>
The `enterprise_projects` block contains:

* `project_id` - Indicates the ID of the IAM user.


