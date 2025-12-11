package cceautopilot

import (
    "fmt"
    "testing"

    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

    "github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCceAutopilotClusterUpgradeInfo_basic(t *testing.T) {
    dataSource := "data.huaweicloud_cce_autopilot_cluster_upgrade_info.test"
    rName := acceptance.RandomAccResourceName()
    dc := acceptance.InitDataSourceCheck(dataSource)

    resource.ParallelTest(t, resource.TestCase{
        PreCheck: func() {
            acceptance.TestAccPreCheck(t)
        },
        ProviderFactories: acceptance.TestAccProviderFactories,
        Steps: []resource.TestStep{
            {
                Config: testDataSourceCceAutopilotClusterUpgradeInfo_basic(rName),
                Check: resource.ComposeTestCheckFunc(
                    dc.CheckResourceExists(),
                ),
            },
        },
    })
}

func testDataSourceCceAutopilotClusterUpgradeInfo_basic(name string) string {
    return fmt.Sprintf(`
data "huaweicloud_cce_autopilot_cluster_upgrade_info" "test" {
  cluster_id = "83a085fe-d4f1-11f0-80d0-0255ac10178d"
  project_id = "0abdd2dce980d4162f8ac006608ee02d"
}
`)
}
