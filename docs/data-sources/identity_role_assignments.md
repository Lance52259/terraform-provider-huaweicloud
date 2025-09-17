---
subcategory: "Identity and Access Management (IAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identity_role_assignments"
description: "Use this data-source to get role assignments of an IAM user within HuaweiCloud."
---

# huaweicloud_identity_role_assignments

Use this data-source to get role assignments of an IAM user within HuaweiCloud.

## Example Usage

```hcl
variable "domain_id" {}

data "huaweicloud_identity_role_assignments" "test" {
 domain_id = var.domaind_id
}
```

## Argument Reference

* `domain_id` - (Required, String) Specifies pending query account ID

* `domain_id` - (Optional, String) Specifies strategy ID

* `subject` - (Optional, String) Specifies authorization entity, value range: **user**, **group**, **agency**. This parameter can only choose one from 
  * **subject.user_id**, 
  * **subject.group_id**, 
  * **subject.agency_id**.

* `subject_user_id` - (Optional, String) Specifies authorized IAM user ID

* `subject_group_id` - (Optional, String) Specifies authorized user group ID

* `subject_agency_id` - (Optional, String) Specifies authorized delegation ID

* `scope` - (Optional, String) Specifies Scope of authorization, value range: **project**, **domain**, **enterprise_project**. This parameter can only choose one among 
  * **scope.project_id**, 
  * **scope.domain_id**, 
  * **scope.enterprise_projects_id**.

* `scope_project_id` - (Optional, String) Specifies authorized project ID

* `scope_domain_id` - (Optional, String) Specifies pending query account ID

* `scope_enterprise_projects_id` - (Optional, String) Specifies authorized enterprise project ID

* `is_inherited` - (Optional, String) Specifies takes effect when the parameter scope=domain or scope.domain_id exists. 
  * **true**: query records based on authorization for all projects. 
  * **false**: query records based on global service authorization.
  Default: **false**

* `include_group` - (Optional, String) Specifies Whether to include records based on the permissions granted to the user groups of the IAM user.
  * **true**: Query records based on IAM user permissions and permissions granted to user groups of IAM users. 
  * **false**: Only query records based on IAM user permissions.
  Default: **true**

* `page` - (Optional, String) Specifies the number of pages of data in a paged query, with a minimum query value of 1, must exist simultaneously with per_page.

* `per_page` - (Optional, String) Specifies the number of items per page in a paginated query, with a range of [1,50], must exist simultaneously with the page.

## Attribute Reference

* `role_assignments` - Indicates authorization information.   
  The [role_assignments](#IdentityRoleAssignments_Assignments ) structure is documented below.

* `total_num` - Indicates returns the total number of authorization records.

<a name="IdentityRoleAssignments_Assignments"></a>
The `role_assignments` block contains:

* `user_id` - Indicates IAM user ID

* `role_id` - Indicates permission ID

* `group_id` - Indicates user Group ID

* `agency_id` - Indicates commission ID

* `scope` - Indicates the scope of the authorization
  The [scope](#IdentityRoleAssignments_Scope ) structure is documented below.

* `is_inherited` - Indicates  based on authorization for all projects

<a name="IdentityRoleAssignments_Scope"></a>
The `scope` block contains:

* `project_id` - Indicates information based on IAM project authorization. 

* `domain_id` - Indicates information based on global service or authorization of all items. 

* `enterprise_project_id` - Indicates information based on enterprise project authorization.