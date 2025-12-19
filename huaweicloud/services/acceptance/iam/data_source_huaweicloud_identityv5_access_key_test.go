package iam

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceIdentityV5AccessKey_basic(t *testing.T) {
	resourceName := "data.huaweicloud_identityv5_access_key.test"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceIdentityV5AccessKey_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "user_id"),
					resource.TestCheckResourceAttrSet(resourceName, "access_key_id"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "last_used_at"),
				),
			},
		},
	})
}

func testAccDataSourceIdentityV5AccessKey_basic() string {
	return `
resource "huaweicloud_identityv5_user" "user_1" {
  name = "Test_accessKey_7"
}

resource "huaweicloud_identityv5_access_key" "key_1" {
  user_id = huaweicloud_identityv5_user.user_1.id
  status  = "inactive"
}

data "huaweicloud_identityv5_access_key" "test" {
  user_id = huaweicloud_identityv5_access_key.key_1.user_id
}
`
}
