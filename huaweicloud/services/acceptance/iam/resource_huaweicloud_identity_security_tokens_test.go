package iam

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccIdentitySecurityTokens(t *testing.T) {
	resourceName := "huaweicloud_identity_security_tokens.test"
	agencyName := acceptance.RandomAccResourceName()
	userName := acceptance.RandomAccResourceName()
	initPassword := acceptance.RandomPassword()
	dc := acceptance.InitDataSourceCheck(resourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAdminOnly(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: AccIdentitySecurityTokensByAgency(agencyName, userName, initPassword),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(resourceName, "expires_at"),
					resource.TestCheckResourceAttrSet(resourceName, "access"),
					resource.TestCheckResourceAttrSet(resourceName, "secret"),
					resource.TestCheckResourceAttrSet(resourceName, "securitytoken"),
				),
			},
			{
				Config: AccIdentitySecurityTokensByToken(userName, initPassword),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(resourceName, "expires_at"),
					resource.TestCheckResourceAttrSet(resourceName, "access"),
					resource.TestCheckResourceAttrSet(resourceName, "secret"),
					resource.TestCheckResourceAttrSet(resourceName, "securitytoken"),
				),
			},
		},
	})
}

func AccIdentitySecurityTokensByAgency(agencyName, userName, initPassword string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_identity_user_token" "test" {
  account_name = "%[2]s"
  user_name    = "%[3]s"
  password     = "%[4]s"
}

resource "huaweicloud_identity_security_tokens" "test" {
  token       = huaweicloud_identity_user_token.test.token
  methods     = "assume_role"
  agency_name = huaweicloud_identity_agency.test.name
  domain_name = huaweicloud_identity_agency.test.delegated_domain_name
  policy      = "{\"Version\":\"1.1\",\"Statement\":[{\"Effect\":\"Allow\",\"Action\":[\"obs:object:GetObject\"],\"Resource\":[\"OBS:*:*:object:*\"]}]}"
}
`, testAccIdentityAgency_domain(agencyName), acceptance.HW_DOMAIN_NAME, acceptance.HW_USER_NAME, acceptance.HW_USER_PWD)
}

func AccIdentitySecurityTokensByToken(userName, initPassword string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_identity_security_tokens" "test" {
  token            = huaweicloud_identity_user_token.test.token
  methods          = "token"
  duration_seconds = 3600
}
`, testAccIdentityUserToken_basic(userName, initPassword))
}
