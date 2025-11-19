package eps

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataEnterpriseProjectAssociatedResources_basic(t *testing.T) {
	all := "data.huaweicloud_enterprise_project_associated_resources.test"
	dc := acceptance.InitDataSourceCheck(all)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataEnterpriseProjectAssociatedResources_basic_test(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(all, "associated_resources.#"),
				),
			},
		},
	})
}

func testAccDataEnterpriseProjectAssociatedResources_basic_test() string {
	return fmt.Sprintf(`
data "huaweicloud_enterprise_project_associated_resources" "test" {
  resource_id = "9ff4bcd5-88b7-4c28-911f-6d5be757024b"
  project_id = "0bc8d0ad6980f56e2f30c00bf75b72b2"
  region_id = "cn-north-4"
  resource_type = "ecs"
}

output "huaweicloud_enterprise_project_associated_resources" {
  value = data.huaweicloud_enterprise_project_associated_resources.test
}
`)
}
