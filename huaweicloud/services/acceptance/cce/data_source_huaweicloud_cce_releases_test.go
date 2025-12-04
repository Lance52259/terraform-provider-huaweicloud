package cce

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCCEReleases_basic(t *testing.T) {
	dataSource := "data.huaweicloud_cce_list_releases.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceCCEListReleases_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "id"),
					resource.TestCheckResourceAttrSet(dataSource, "region"),
				),
			},
		},
	})
}

func testDataSourceCCEListReleases_basic(name string) string {
	return fmt.Sprintf(`

data "huaweicloud_cce_list_releases" "test" {
  cluster_id = "5064e9b7-cb89-11f0-99e1-0255ac1001b2"
}
`)
}
