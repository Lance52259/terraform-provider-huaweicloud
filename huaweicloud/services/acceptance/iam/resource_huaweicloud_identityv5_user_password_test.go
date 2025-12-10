package iam

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccIdentityV5UserPassword(t *testing.T) {
	password := "TestPassword#$%111"
	originalPassword := "TestPassword#$%"
	resourceName := "huaweicloud_identityv5_user_password.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityV5UserPassword_basic(password, originalPassword),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "new_password", password),
					resource.TestCheckResourceAttr(resourceName, "old_password", originalPassword),
				),
			},
			{
				Config: testAccIdentityV5UserPassword_update(password, originalPassword),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "new_password", originalPassword),
					resource.TestCheckResourceAttr(resourceName, "old_password", password),
				),
			},
		},
	})
}

func testAccIdentityV5UserPassword_basic(newPassword, oldPassword string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identityv5_user_password" "test" {
  new_password = "%[1]s"
  old_password = "%[2]s"
}
`, newPassword, oldPassword)
}

func testAccIdentityV5UserPassword_update(newPassword, oldPassword string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identityv5_user_password" "test" {
  new_password = "%[1]s"
  old_password = "%[2]s"
  enable_force_new = true
}
`, oldPassword, newPassword)
}
