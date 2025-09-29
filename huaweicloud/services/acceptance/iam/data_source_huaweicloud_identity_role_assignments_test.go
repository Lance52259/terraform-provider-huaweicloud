package iam

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccIdentityRoleAssignments_basic(t *testing.T) {
	dataSource := "data.huaweicloud_identity_role_assignments.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityRoleAssignments_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSource, "role_assignments.#", "7"),
				),
			},
			{
				Config: testAccIdentityRoleAssignments_user(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSource, "role_assignments.#", "7"),
				),
			},
			{
				Config: testAccIdentityRoleAssignments_domain(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSource, "role_assignments.#", "2"),
				),
			},
		},
	})
}

func testAccIdentityRoleAssignments_basic() string {
	return `
  data "huaweicloud_identity_role_assignments" "test" {}
`
}
func testAccIdentityRoleAssignments_user() string {
	return `
  data "huaweicloud_identity_role_assignments" "test" {
	 is_inherited = true
}
`
}
func testAccIdentityRoleAssignments_domain() string {
	return `
  data "huaweicloud_identity_role" "role_1" {
 	  name = "system_all_64"
}

  data "huaweicloud_identity_role_assignments" "test" {
	 is_inherited = true
     role_id 	  = data.huaweicloud_identity_role.role_1.role_id
}
`
}
