package cc

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCcCloudConnectionCapabilities_basic(t *testing.T) {
	dataSource := "data.huaweicloud_cc_connection_capabilities.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceCcCloudConnectionCapabilities_Basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "capabilities.#"),
					resource.TestCheckResourceAttrSet(dataSource, "capabilities.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "capabilities.0.domain_id"),
					resource.TestCheckResourceAttrSet(dataSource, "capabilities.0.resource_type"),
				),
			},
			{
				Config: testDataSourceCcCloudConnectionCapabilities_ResourceType_Ipv6_Support_Regions(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "capabilities.#"),
					resource.TestCheckResourceAttr(dataSource, "capabilities.0.resource_type", "ipv6_support_regions"),
					resource.TestCheckResourceAttrSet(dataSource, "capabilities.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "capabilities.0.domain_id"),
					resource.TestCheckResourceAttrSet(dataSource, "capabilities.0.support_regions.#"),
					resource.TestCheckResourceAttrSet(dataSource, "capabilities.0.support_regions.0"),
				),
			},
			{
				Config: testDataSourceCcCloudConnectionCapabilities_ResourceType_DomainBandwidth(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "capabilities.#"),
					resource.TestCheckResourceAttr(dataSource, "capabilities.0.resource_type", "domain_bandwidth"),
					resource.TestCheckResourceAttrSet(dataSource, "capabilities.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "capabilities.0.domain_id"),
					resource.TestCheckResourceAttrSet(dataSource, "capabilities.0.bandwidth.#"),
					resource.TestCheckResourceAttrSet(dataSource, "capabilities.0.bandwidth.0.min"),
					resource.TestCheckResourceAttrSet(dataSource, "capabilities.0.bandwidth.0.max"),
				),
			},
		},
	})
}

func testDataSourceCcCloudConnectionCapabilities_Basic() string {
	return `
data "huaweicloud_cc_connection_capabilities" "test" {
  
}`
}

func testDataSourceCcCloudConnectionCapabilities_ResourceType_Ipv6_Support_Regions() string {
	return `
data "huaweicloud_cc_connection_capabilities" "test" {
  resource_type = "ipv6_support_regions"
}`
}

func testDataSourceCcCloudConnectionCapabilities_ResourceType_DomainBandwidth() string {
	return `
data "huaweicloud_cc_connection_capabilities" "test" {
  resource_type = "domain_bandwidth"
}`
}
