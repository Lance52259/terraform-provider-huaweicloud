package iam

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccIdentityRoleDataSource_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_identity_role.role_1"
	dataSourceName2 := "huaweicloud_identity_permission_detail.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)
	dc2 := acceptance.InitDataSourceCheck(dataSourceName2)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAdminOnly(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityRoleDataSourceByName,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "name", "system_all_64"),
					resource.TestCheckResourceAttr(dataSourceName, "display_name", "OBS ReadOnlyAccess"),
					resource.TestCheckResourceAttr(dataSourceName, "role_id", "0af84c1502f447fa9c2fa18083fbb87e"),
				),
			},
			{
				Config: testAccIdentityRoleDataSourceByDisplayName,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "name", "kms_adm"),
					resource.TestCheckResourceAttr(dataSourceName, "display_name", "KMS Administrator"),
					resource.TestCheckResourceAttr(dataSourceName, "role_id", "0af84c1502f447fa9c2fa18083fbb87e"),
				),
			},
			{
				Config: testAccIdentityRoleDataSourceByRoleId,
				Check: resource.ComposeTestCheckFunc(
					dc2.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "name", "kms_adm"),
					resource.TestCheckResourceAttr(dataSourceName, "display_name", "KMS Administrator"),
					resource.TestCheckResourceAttr(dataSourceName, "role_id", "0af84c1502f447fa9c2fa18083fbb87e"),
				),
			},
		},
	})
}

const testAccIdentityRoleDataSourceByName = `
data "huaweicloud_identity_role" "role_1" {
  # OBS ReadOnlyAccess
  name = "system_all_64"
}
`
const testAccIdentityRoleDataSourceByDisplayName = `
data "huaweicloud_identity_role" "role_1" {
  display_name = "KMS Administrator"
}
`
const testAccIdentityRoleDataSourceByRoleId = `
data "huaweicloud_identity_permission_detail" "test" {
role_id =  "0af84c1502f447fa9c2fa18083fbb87e"
}
`
