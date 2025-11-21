package cc

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCcBandwidthPackageSupportedAreas_basic(t *testing.T) {
	dataSource := "data.huaweicloud_cc_bandwidth_package_supported_areas.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceCcBandwidthPackageSupportedAreas_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "areas.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "areas.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "areas.0.en_name"),
					resource.TestCheckResourceAttrSet(dataSource, "areas.0.station"),
				),
			},
			{
				Config: testDataSourceCcBandwidthPackageSupportedAreas_station(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSource, "areas.0.station", "chinaStation"),
					resource.TestCheckResourceAttrSet(dataSource, "areas.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "areas.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "areas.0.en_name"),
				),
			},
		},
	})
}

func testDataSourceCcBandwidthPackageSupportedAreas_basic() string {
	return `data "huaweicloud_cc_bandwidth_package_supported_areas" "test" {}`
}

func testDataSourceCcBandwidthPackageSupportedAreas_station() string {
	return `
data "huaweicloud_cc_bandwidth_package_supported_areas" "test" {
  station = "chinaStation"	
}`
}
