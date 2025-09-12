package iam

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccIdentityEndpointDataSource_basic(t *testing.T) {
	dataSource := "data.huaweicloud_identity_endpoint.test"
	userName := acceptance.RandomAccResourceName()
	initPassword := acceptance.RandomPassword()

	config := testAccIdentityEndpoint_basic(userName, initPassword)
	configWithId := testAccIdentityEndpointById_basic(userName, initPassword)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: configWithId,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSource, "endpoints.#", "1"),
					resource.TestCheckResourceAttr(dataSource, "endpoints.0.enabled", "true"),
				),
			},
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSource, "endpoints.#", "1062"),
					resource.TestCheckResourceAttr(dataSource, "endpoints.0.enabled", "true"),
					resource.TestCheckResourceAttr(dataSource, "links.%", "3"),
				),
			},
		},
	})
}

func testAccIdentityEndpoint_basic(userName, initPassword string) string {
	return fmt.Sprintf(`
%[1]s
 
resource "huaweicloud_identity_user_token" "test" {
  account_name = "%[2]s"
  user_name    = huaweicloud_identity_user.user_1.name
  password     = "%[3]s"
}
data "huaweicloud_identity_endpoint" "test" {
  token = huaweicloud_identity_user_token.test.token
}
 
`, testAccIdentityUser_basic(userName, initPassword), acceptance.HW_DOMAIN_NAME, initPassword)
}
func testAccIdentityEndpointById_basic(userName, initPassword string) string {
	return fmt.Sprintf(`
%[1]s
 
resource "huaweicloud_identity_user_token" "test" {
  account_name = "%[2]s"
  user_name    = huaweicloud_identity_user.user_1.name
  password     = "%[3]s"
}
data "huaweicloud_identity_endpoint" "test" {
  token = huaweicloud_identity_user_token.test.token
  endpoint_id = "0046cca357c94165b7a10ec2c01bdf60"
}
 
`, testAccIdentityUser_basic(userName, initPassword), acceptance.HW_DOMAIN_NAME, initPassword)
}
