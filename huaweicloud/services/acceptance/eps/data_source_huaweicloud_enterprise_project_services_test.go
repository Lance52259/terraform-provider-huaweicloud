package eps

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataEnterpriseProjectServices_basic(t *testing.T) {
	all := "data.huaweicloud_enterprise_project_services.test"
	dc := acceptance.InitDataSourceCheck(all)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataEnterpriseProjectServices_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(all, "services.#"),
					resource.TestCheckOutput("huaweicloud_enterprise_project_services", "true"),
				),
			},
		},
	})
}

const testAccDataEnterpriseProjectServices_basic = `
data "huaweicloud_enterprise_project_services" "test" {}

output "huaweicloud_enterprise_project_services" {
  value = length(data.huaweicloud_enterprise_project_services.test.services) > 0
}
`
