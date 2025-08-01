package sms

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceSmsSourceServerErrors_basic(t *testing.T) {
	dataSource := "data.huaweicloud_sms_source_server_errors.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceSmsSourceServerErrors_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
				),
			},
		},
	})
}

func testDataSourceDataSourceSmsSourceServerErrors_basic() string {
	return `
data "huaweicloud_sms_source_server_errors" "test" {}
`
}
