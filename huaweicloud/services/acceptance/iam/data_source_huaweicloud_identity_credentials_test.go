package iam

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"testing"
)

func TestAccIdentityCredentialsDataSource_basic(t *testing.T) {
	dataSource := "data.huaweicloud_identity_credentials.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAdminOnly(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: AccIdentityCredentialsDataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "credentials.#"),
				),
			},
			{
				Config: AccIdentityCredentialsDataSourceByUserId_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "credentials.#"),
				),
			},
		},
	})
}

const AccIdentityCredentialsDataSource_basic = `
data "huaweicloud_identity_credentials" "test" {
}`

const AccIdentityCredentialsDataSourceByUserId_basic = `
data "huaweicloud_identity_credentials" "test" {
	user_id = "d9a1d990bc4e4dfeb30c25ea4020e007"
}`
