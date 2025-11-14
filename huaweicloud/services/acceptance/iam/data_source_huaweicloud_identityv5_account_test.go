package iam

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceIdentityV5Account_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_identityv5_account.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceIdentityV5Account_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "feature_status"),
				),
			},
		},
	})
}

func testAccDataSourceIdentityV5Account_basic() string {
	return `
data "huaweicloud_identityv5_account" "test" {
	feature_name = "v5_console"
}
`
}
