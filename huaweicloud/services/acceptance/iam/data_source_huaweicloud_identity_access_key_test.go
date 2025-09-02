package iam

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"testing"
)

func TestAccIdentityAccessKeyDataSource_basic(t *testing.T) {
	dataSource := "data.huaweicloud_identity_access_key.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: AccIdentityAccessKeyDataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSource, "credentials.#", "1"),
					resource.TestCheckResourceAttr(dataSource, "credentials.0.status", "active"),
				),
			},
			{
				Config: AccIdentityAccessKeyDataSourceByAccessKey_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSource, "credential.#", "1"),
					resource.TestCheckResourceAttr(dataSource, "credential.0.status", "active"),
				),
			},
			{
				Config: AccIdentityAccessKeyDataSourceWithoutParams_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSource, "credentials.#", "1"),
					resource.TestCheckResourceAttr(dataSource, "credential.0.status", "active"),
				),
			},
		},
	})
}

const AccIdentityAccessKeyDataSource_basic = `
resource "huaweicloud_identity_user" "test_user" {
  name        = "test_huaweicloud"
  password    = "chenyanC0212"
  enabled     = true
  description = "tested by terraform"
}


resource "huaweicloud_identity_access_key" "key_1" {
  user_id     = huaweicloud_identity_user.test_user.id
  description = "access key by terraform"
  secret_file = "./credentials.csv"
}

data "huaweicloud_identity_access_key" "test" {
  user_id = huaweicloud_identity_access_key.key_1.user_id
}
`

const AccIdentityAccessKeyDataSourceByAccessKey_basic = `
resource "huaweicloud_identity_user" "test_user" {
  name        = "test_huaweicloud"
  password    = "chenyanC0212"
  enabled     = true
  description = "tested by terraform"
}


resource "huaweicloud_identity_access_key" "key_1" {
  user_id     = huaweicloud_identity_user.test_user.id
  description = "access key by terraform"
  secret_file = "./credentials.csv"
}

data "huaweicloud_identity_access_key" "test" {
  access_key = huaweicloud_identity_access_key.key_1.id
}
`

const AccIdentityAccessKeyDataSourceWithoutParams_basic = `
resource "huaweicloud_identity_user" "test_user" {
  name        = "test_huaweicloud"
  password    = "chenyanC0212"
  enabled     = true
  description = "tested by terraform"
}


resource "huaweicloud_identity_access_key" "key_1" {
  user_id     = huaweicloud_identity_user.test_user.id
  description = "access key by terraform"
  secret_file = "./credentials.csv"
}

data "huaweicloud_identity_access_key" "test" {
}
`
