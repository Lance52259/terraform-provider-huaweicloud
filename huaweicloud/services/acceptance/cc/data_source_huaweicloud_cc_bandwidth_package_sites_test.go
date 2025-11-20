package cc

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCcBandwidthPackageSites_basic(t *testing.T) {
	dataSource := "data.huaweicloud_cc_bandwidth_package_sites.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceCcBandwidthPackageSites_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "bandwidth_package_sites.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "bandwidth_package_sites.0.site_code"),
					resource.TestCheckResourceAttrSet(dataSource, "bandwidth_package_sites.0.region_id"),
					resource.TestCheckResourceAttrSet(dataSource, "bandwidth_package_sites.0.site_type"),
					resource.TestCheckResourceAttrSet(dataSource, "bandwidth_package_sites.0.name_cn"),
					resource.TestCheckResourceAttrSet(dataSource, "bandwidth_package_sites.0.name_en"),
					resource.TestCheckResourceAttrSet(dataSource, "bandwidth_package_sites.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSource, "bandwidth_package_sites.0.updated_at"),
				),
			},
			{
				Config: testDataSourceCcBandwidthPackageSites_site_code(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSource, "bandwidth_package_sites.0.site_code", "cn-south-guangzhou"),
					resource.TestCheckResourceAttrSet(dataSource, "bandwidth_package_sites.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "bandwidth_package_sites.0.region_id"),
					resource.TestCheckResourceAttrSet(dataSource, "bandwidth_package_sites.0.site_type"),
					resource.TestCheckResourceAttrSet(dataSource, "bandwidth_package_sites.0.name_cn"),
					resource.TestCheckResourceAttrSet(dataSource, "bandwidth_package_sites.0.name_en"),
					resource.TestCheckResourceAttrSet(dataSource, "bandwidth_package_sites.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSource, "bandwidth_package_sites.0.updated_at"),
				),
			},
			{
				Config: testDataSourceCcBandwidthPackageSites_region_id(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSource, "bandwidth_package_sites.0.region_id", "cn-south-1"),
					resource.TestCheckResourceAttrSet(dataSource, "bandwidth_package_sites.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "bandwidth_package_sites.0.site_code"),
					resource.TestCheckResourceAttrSet(dataSource, "bandwidth_package_sites.0.site_type"),
					resource.TestCheckResourceAttrSet(dataSource, "bandwidth_package_sites.0.name_cn"),
					resource.TestCheckResourceAttrSet(dataSource, "bandwidth_package_sites.0.name_en"),
					resource.TestCheckResourceAttrSet(dataSource, "bandwidth_package_sites.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSource, "bandwidth_package_sites.0.updated_at"),
				),
			},
			{
				Config: testDataSourceCcBandwidthPackageSites_name(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSource, "bandwidth_package_sites.0.name_en", "CN-South-Guangzhou"),
					resource.TestCheckResourceAttrSet(dataSource, "bandwidth_package_sites.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "bandwidth_package_sites.0.site_code"),
					resource.TestCheckResourceAttrSet(dataSource, "bandwidth_package_sites.0.region_id"),
					resource.TestCheckResourceAttrSet(dataSource, "bandwidth_package_sites.0.site_type"),
					resource.TestCheckResourceAttrSet(dataSource, "bandwidth_package_sites.0.name_cn"),
					resource.TestCheckResourceAttrSet(dataSource, "bandwidth_package_sites.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSource, "bandwidth_package_sites.0.updated_at"),
				),
			},
		},
	})
}

func testDataSourceCcBandwidthPackageSites_basic() string {
	return `data "huaweicloud_cc_bandwidth_package_sites" "test" {}`
}

func testDataSourceCcBandwidthPackageSites_site_code() string {
	return `
data "huaweicloud_cc_bandwidth_package_sites" "test" {
  site_code = "cn-south-guangzhou"	
}`
}

func testDataSourceCcBandwidthPackageSites_region_id() string {
	return `
data "huaweicloud_cc_bandwidth_package_sites" "test" {
  region_id = "cn-south-1"
}`
}

func testDataSourceCcBandwidthPackageSites_name() string {
	return `
data "huaweicloud_cc_bandwidth_package_sites" "test" {
  name = "CN-South-Guangzhou"
}`
}
