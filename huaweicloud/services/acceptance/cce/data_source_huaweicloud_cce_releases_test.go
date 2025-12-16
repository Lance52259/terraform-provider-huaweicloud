package cce

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCCEReleases_basic(t *testing.T) {
	dataSource := "data.huaweicloud_cce_releases.test"
	rName := "xxxxxxdeploy"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceCCEListReleases_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "id"),
					resource.TestCheckResourceAttrSet(dataSource, "region"),
					resource.TestCheckResourceAttr(dataSource, "releases.0.name", rName),
				),
			},
		},
	})
}

func testDataSourceCCEListReleases_basic() string {
	return fmt.Sprintf(`

data "huaweicloud_cce_releases" "test" {
  cluster_id = "3f44db15-d996-11f0-b740-0255ac1001b5"
}
`)
}
