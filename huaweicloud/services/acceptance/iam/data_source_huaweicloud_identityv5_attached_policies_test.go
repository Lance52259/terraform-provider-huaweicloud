package iam

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceIdentityV5AttachedPolicies_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_identityv5_attached_policies.test"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceIdentityV5AttachedPolicies_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "group_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "attached_policies.#"),
				),
			},
		},
	})
}

func testAccDataSourceIdentityV5AttachedPolicies_basic() string {
	return `
data "huaweicloud_identityv5_attached_policies" "test" {
  group_id = "044f2b04cea84d96841df9b4ad19d91c"
}
`
}
