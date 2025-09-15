package iam

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccIdentityEndpointDataSource_basic(t *testing.T) {
	dataSource := "data.huaweicloud_identity_endpoints.test"

	config := testAccIdentityEndpoints_basic()
	configWithId := testAccIdentityEndpointsById_basic()
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
				),
			},
		},
	})
}

func testAccIdentityEndpoints_basic() string {
	return `
data "huaweicloud_identity_endpoints" "test" {
}
`
}

func testAccIdentityEndpointsById_basic() string {
	return `
data "huaweicloud_identity_endpoints" "test" {
  endpoint_id = "0046cca357c94165b7a10ec2c01bdf60"
}
`
}
