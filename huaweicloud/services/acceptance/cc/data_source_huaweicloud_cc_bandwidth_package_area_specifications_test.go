package cc

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCcBandwidthPackageAreaSpecifications_basic(t *testing.T) {
	dataSource := "data.huaweicloud_cc_bandwidth_package_area_specifications.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceCcBandwidthPackageAreaSpecifications_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "area_specifications.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "area_specifications.0.local_area_id"),
					resource.TestCheckResourceAttrSet(dataSource, "area_specifications.0.remote_area_id"),
					resource.TestCheckResourceAttrSet(dataSource, "area_specifications.0.spec_codes.0.spec_code"),
					resource.TestCheckResourceAttrSet(dataSource, "area_specifications.0.spec_codes.0.billing_mode"),
					resource.TestCheckResourceAttrSet(dataSource, "area_specifications.0.spec_codes.0.max_bandwidth"),
					resource.TestCheckResourceAttrSet(dataSource, "area_specifications.0.spec_codes.0.mim_bandwidth"),
				),
			},
			{
				Config: testDataSourceCcBandwidthPackageAreaSpecifications_localArea(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSource, "area_specifications.0.local_area_id", "Africa"),
					resource.TestCheckResourceAttrSet(dataSource, "area_specifications.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "area_specifications.0.remote_area_id"),
					resource.TestCheckResourceAttrSet(dataSource, "area_specifications.0.spec_codes.0.spec_code"),
					resource.TestCheckResourceAttrSet(dataSource, "area_specifications.0.spec_codes.0.billing_mode"),
					resource.TestCheckResourceAttrSet(dataSource, "area_specifications.0.spec_codes.0.max_bandwidth"),
					resource.TestCheckResourceAttrSet(dataSource, "area_specifications.0.spec_codes.0.mim_bandwidth"),
				),
			},
			{
				Config: testDataSourceCcBandwidthPackageAreaSpecifications_remoteArea(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSource, "area_specifications.0.remote_area_id", "Asia-Pacific"),
					resource.TestCheckResourceAttrSet(dataSource, "area_specifications.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "area_specifications.0.local_area_id"),
					resource.TestCheckResourceAttrSet(dataSource, "area_specifications.0.spec_codes.0.spec_code"),
					resource.TestCheckResourceAttrSet(dataSource, "area_specifications.0.spec_codes.0.billing_mode"),
					resource.TestCheckResourceAttrSet(dataSource, "area_specifications.0.spec_codes.0.max_bandwidth"),
					resource.TestCheckResourceAttrSet(dataSource, "area_specifications.0.spec_codes.0.mim_bandwidth"),
				),
			},
			{
				Config: testDataSourceCcBandwidthPackageAreaSpecifications_MultiArea(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSource, "area_specifications.0.local_area_id", "Africa"),
					resource.TestCheckResourceAttr(dataSource, "area_specifications.0.remote_area_id", "Asia-Pacific"),
					resource.TestCheckResourceAttrSet(dataSource, "area_specifications.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "area_specifications.0.spec_codes.0.spec_code"),
					resource.TestCheckResourceAttrSet(dataSource, "area_specifications.0.spec_codes.0.billing_mode"),
					resource.TestCheckResourceAttrSet(dataSource, "area_specifications.0.spec_codes.0.max_bandwidth"),
					resource.TestCheckResourceAttrSet(dataSource, "area_specifications.0.spec_codes.0.mim_bandwidth"),
				),
			},
		},
	})
}

func testDataSourceCcBandwidthPackageAreaSpecifications_basic() string {
	return `data "huaweicloud_cc_bandwidth_package_area_specifications" "test" {}`
}

func testDataSourceCcBandwidthPackageAreaSpecifications_localArea() string {
	return `
data "huaweicloud_cc_bandwidth_package_area_specifications" "test" {
  local_area_id = "Africa"	
}`
}

func testDataSourceCcBandwidthPackageAreaSpecifications_remoteArea() string {
	return `
data "huaweicloud_cc_bandwidth_package_area_specifications" "test" {
  remote_area_id = "Asia-Pacific"
}`
}

func testDataSourceCcBandwidthPackageAreaSpecifications_MultiArea() string {
	return `
data "huaweicloud_cc_bandwidth_package_area_specifications" "test" {
  local_area_id  = "Africa"
  remote_area_id = "Asia-Pacific"
}`
}
