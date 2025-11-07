package cc

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCcBandwidthPackageRegionSpecifications_basic(t *testing.T) {
	dataSource := "data.huaweicloud_cc_bandwidth_package_region_specifications.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceCcBandwidthPackageRegionSpecifications_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "region_specifications.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "region_specifications.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "region_specifications.0.en_name"),
					resource.TestCheckResourceAttrSet(dataSource, "region_specifications.0.local_region_id"),
					resource.TestCheckResourceAttrSet(dataSource, "region_specifications.0.remote_region_id"),
					resource.TestCheckResourceAttrSet(dataSource, "region_specifications.0.spec_codes.0.spec_code"),
					resource.TestCheckResourceAttrSet(dataSource, "region_specifications.0.spec_codes.0.billing_mode"),
					resource.TestCheckResourceAttrSet(dataSource, "region_specifications.0.spec_codes.0.max_bandwidth"),
					resource.TestCheckResourceAttrSet(dataSource, "region_specifications.0.spec_codes.0.mim_bandwidth"),
				),
			},
			{
				Config: testDataSourceCcBandwidthPackageRegionSpecifications_localRegion(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSource, "region_specifications.0.local_region_id", "cn-southwest-242"),
					resource.TestCheckResourceAttrSet(dataSource, "region_specifications.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "region_specifications.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "region_specifications.0.en_name"),
					resource.TestCheckResourceAttrSet(dataSource, "region_specifications.0.remote_region_id"),
					resource.TestCheckResourceAttrSet(dataSource, "region_specifications.0.spec_codes.0.spec_code"),
					resource.TestCheckResourceAttrSet(dataSource, "region_specifications.0.spec_codes.0.billing_mode"),
					resource.TestCheckResourceAttrSet(dataSource, "region_specifications.0.spec_codes.0.max_bandwidth"),
					resource.TestCheckResourceAttrSet(dataSource, "region_specifications.0.spec_codes.0.mim_bandwidth"),
				),
			},
			{
				Config: testDataSourceCcBandwidthPackageRegionSpecifications_remoteRegion(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSource, "region_specifications.0.remote_region_id", "cn-southwest-248"),
					resource.TestCheckResourceAttrSet(dataSource, "region_specifications.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "region_specifications.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "region_specifications.0.en_name"),
					resource.TestCheckResourceAttrSet(dataSource, "region_specifications.0.local_region_id"),
					resource.TestCheckResourceAttrSet(dataSource, "region_specifications.0.spec_codes.0.spec_code"),
					resource.TestCheckResourceAttrSet(dataSource, "region_specifications.0.spec_codes.0.billing_mode"),
					resource.TestCheckResourceAttrSet(dataSource, "region_specifications.0.spec_codes.0.max_bandwidth"),
					resource.TestCheckResourceAttrSet(dataSource, "region_specifications.0.spec_codes.0.mim_bandwidth"),
				),
			},
			{
				Config: testDataSourceCcBandwidthPackageRegionSpecifications_MultiRegion(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSource, "region_specifications.0.local_region_id", "cn-southwest-242"),
					resource.TestCheckResourceAttr(dataSource, "region_specifications.0.remote_region_id", "cn-southwest-248"),
					resource.TestCheckResourceAttrSet(dataSource, "region_specifications.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "region_specifications.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "region_specifications.0.en_name"),
					resource.TestCheckResourceAttrSet(dataSource, "region_specifications.0.spec_codes.0.spec_code"),
					resource.TestCheckResourceAttrSet(dataSource, "region_specifications.0.spec_codes.0.billing_mode"),
					resource.TestCheckResourceAttrSet(dataSource, "region_specifications.0.spec_codes.0.max_bandwidth"),
					resource.TestCheckResourceAttrSet(dataSource, "region_specifications.0.spec_codes.0.mim_bandwidth"),
				),
			},
		},
	})
}

func testDataSourceCcBandwidthPackageRegionSpecifications_basic() string {
	return `data "huaweicloud_cc_bandwidth_package_region_specifications" "test" {}`
}

func testDataSourceCcBandwidthPackageRegionSpecifications_localRegion() string {
	return `
data "huaweicloud_cc_bandwidth_package_region_specifications" "test" {
  local_region_id = "cn-southwest-242"	
}`
}

func testDataSourceCcBandwidthPackageRegionSpecifications_remoteRegion() string {
	return `
data "huaweicloud_cc_bandwidth_package_region_specifications" "test" {
  remote_region_id = "cn-southwest-248"
}`
}

func testDataSourceCcBandwidthPackageRegionSpecifications_MultiRegion() string {
	return `
data "huaweicloud_cc_bandwidth_package_region_specifications" "test" {
  local_region_id = "cn-southwest-242"
  remote_region_id = "cn-southwest-248"
}`
}
