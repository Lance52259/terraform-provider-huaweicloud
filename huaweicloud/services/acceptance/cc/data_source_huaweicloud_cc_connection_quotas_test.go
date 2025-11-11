package cc

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCcCloudConnectionQuotas_basic(t *testing.T) {
	dataSource := "data.huaweicloud_cc_connection_quotas.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceCcCloudConnection_QuotaType_CloudConnection(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "quotas.#"),
					resource.TestCheckResourceAttr(dataSource, "quotas.0.quota_type", "cloud_connection"),
					resource.TestCheckResourceAttrSet(dataSource, "quotas.0.domain_id"),
					resource.TestCheckResourceAttrSet(dataSource, "quotas.0.quota_number"),
					resource.TestCheckResourceAttrSet(dataSource, "quotas.0.quota_used"),
				),
			},
			{
				Config: testDataSourceCcCloudConnection_QuotaType_CloudConnectionRegion(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "quotas.#"),
					resource.TestCheckResourceAttr(dataSource, "quotas.0.quota_type", "cloud_connection_region"),
					resource.TestCheckResourceAttr(dataSource, "quotas.0.cloud_connection_id", "a4da6488b3154b22ad5a93f858dc6a9d"),
					resource.TestCheckResourceAttrSet(dataSource, "quotas.0.domain_id"),
					resource.TestCheckResourceAttrSet(dataSource, "quotas.0.quota_number"),
					resource.TestCheckResourceAttrSet(dataSource, "quotas.0.quota_used"),
				),
			},
			{
				Config: testDataSourceCcCloudConnection_QuotaType_CloudConnectionRoute(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "quotas.#"),
					resource.TestCheckResourceAttr(dataSource, "quotas.0.quota_type", "cloud_connection_route"),
					resource.TestCheckResourceAttr(dataSource, "quotas.0.cloud_connection_id", "a4da6488b3154b22ad5a93f858dc6a9d"),
					resource.TestCheckResourceAttrSet(dataSource, "quotas.0.domain_id"),
					resource.TestCheckResourceAttrSet(dataSource, "quotas.0.quota_number"),
					resource.TestCheckResourceAttrSet(dataSource, "quotas.0.quota_used"),
				),
			},
			{
				Config: testDataSourceCcCloudConnection_QuotaType_NetworkInstance(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "quotas.#"),
					resource.TestCheckResourceAttr(dataSource, "quotas.0.quota_type", "region_network_instance"),
					resource.TestCheckResourceAttr(dataSource, "quotas.0.cloud_connection_id", "a4da6488b3154b22ad5a93f858dc6a9d"),
					resource.TestCheckResourceAttr(dataSource, "quotas.0.region_id", "cn-southwest-242"),
					resource.TestCheckResourceAttrSet(dataSource, "quotas.0.domain_id"),
					resource.TestCheckResourceAttrSet(dataSource, "quotas.0.quota_number"),
					resource.TestCheckResourceAttrSet(dataSource, "quotas.0.quota_used"),
				),
			},
		},
	})
}

func testDataSourceCcCloudConnection_QuotaType_CloudConnection() string {
	return `
data "huaweicloud_cc_connection_quotas" "test" {
  quota_type = "cloud_connection"
}`
}

func testDataSourceCcCloudConnection_QuotaType_CloudConnectionRegion() string {
	return `
data "huaweicloud_cc_connection_quotas" "test" {
  quota_type = "cloud_connection_region"
  cloud_connection_id = "a4da6488b3154b22ad5a93f858dc6a9d"
}`
}

func testDataSourceCcCloudConnection_QuotaType_CloudConnectionRoute() string {
	return `
data "huaweicloud_cc_connection_quotas" "test" {
  quota_type = "cloud_connection_route"
  cloud_connection_id = "a4da6488b3154b22ad5a93f858dc6a9d"
}`
}

func testDataSourceCcCloudConnection_QuotaType_NetworkInstance() string {
	return `
data "huaweicloud_cc_connection_quotas" "test" {
  quota_type = "region_network_instance"
  cloud_connection_id = "a4da6488b3154b22ad5a93f858dc6a9d"
  region_id = "cn-southwest-242"
}`
}
