package cceautopilot

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCceAutopilotChart_basic(t *testing.T) {
	dataSource := "data.huaweicloud_cce_autopilot_chart.test"
	resourceName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceCceAutopilotChart_basic(resourceName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSource, "region", acceptance.HW_REGION_NAME),
				),
			},
		},
	})
}

func testDataSourceCceAutopilotChart_basic(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_cce_autopilot_chart" "test" {
  region = "%[1]s"
}
`, acceptance.HW_REGION_NAME)
}
