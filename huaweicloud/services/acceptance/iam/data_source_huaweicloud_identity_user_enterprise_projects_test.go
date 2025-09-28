package iam

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccIdentityUserEnterpriseProjects_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_identity_user_enterprise_projects.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityUserEnterpriseProject(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "enterprise_projects.#"),
				),
			},
		},
	})
}

func testAccIdentityUserEnterpriseProject() string {
	return fmt.Sprintf(`
data "huaweicloud_identity_user_enterprise_projects" "test" {
 	user_id = "%s"
}
`, acceptance.HW_USER_ID)
}
