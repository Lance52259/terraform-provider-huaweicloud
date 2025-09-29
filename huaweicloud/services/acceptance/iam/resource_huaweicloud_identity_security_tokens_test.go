package iam

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccIdentitySecurityTokens_basic(t *testing.T) {
	resourceName := "huaweicloud_identity_security_tokens.test"
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
				Config: AccIdentitySecurityTokens_basic(userName, initPassword),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(resourceName, "credential.#"),
				),
			},
		},
	})
}

func AccIdentitySecurityTokens_basic(userName string, initPassword string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_identity_security_tokens" "test" {
  token = huaweicloud_identity_user_token.test.token
  methods = "token"
  policy = <<POLICY
{
  "Version": "1.1",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": ["obs:object:GetObject"],
      "Resource": ["OBS:*:*:object:*"],
      "Condition": {
        "StringEquals": {
          "g:DomainName": ["DomainNameExample"]
        }
      }
    }
  ]
}
POLICY
}

}
`, testAccIdentityUserToken_basic(userName, initPassword))
}
