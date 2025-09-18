package iam

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"testing"
)

func TestAccIdentityRoleAssignments_basic(t *testing.T) {
	dataSource := "data.huaweicloud_identity_role_assignments.test"
	config := testAccIdentityRoleAssignments_basic()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSource, "role_assignments.#", "1"),
				),
			},
		},
	})
}

func testAccIdentityRoleAssignments_basic() string {
	return `
  data "huaweicloud_identity_role_assignments" "test" {
}
`
}
