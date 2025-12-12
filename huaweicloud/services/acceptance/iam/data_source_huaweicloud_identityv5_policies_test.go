package iam

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceIdentityV5Policies_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_identityv5_policies.test"
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceIdentityV5Policies_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "policies.#"),
				),
			},
			{
				Config: testAccDataSourceIdentityV5PoliciesWithId_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "policies.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "policies.0.policy_id", "NATReadOnlyPolicy"),
					resource.TestCheckResourceAttr(dataSourceName, "policies.0.policy_name", "NATReadOnlyPolicy"),
				),
			},
		},
	})
}

var testAccDataSourceIdentityV5Policies_basic = `
data "huaweicloud_identityv5_policies" "test" {
  language = "en-us"
}
`

var testAccDataSourceIdentityV5PoliciesWithId_basic = `
data "huaweicloud_identityv5_policies" "test" {
  policy_id = "NATReadOnlyPolicy"
}
`
