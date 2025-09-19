package iam

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccIdentityUnscopedTokenWithIdToken(t *testing.T) {
	idpId := "ShibbolethIdP_Terraform_OIDC"
	protocolId := "oidc"
	idToken := "eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICJqVnBkWDU1WWNjYUpHelRFQnFid2tNcEFqQXJmdnRnRkVjeTVidEYzcE9JIn0.eyJleHAiOjE3NTgyNzEwMDUsImlhdCI6MTc1ODI3MDEwNSwiYXV0aF90aW1lIjoxNzU4MjUyMDg1LCJqdGkiOiI4NGFjYTY3Ni03OTQ1LTQwOWYtYjZlMC01ZmI4ZDAyMmIzZTEiLCJpc3MiOiJodHRwczovLzEwLjIzNS4yMzEuMTg5OjU0My9hdXRoL3JlYWxtcy9tYXN0ZXIiLCJhdWQiOiJjeHctX09JREMiLCJzdWIiOiJkNGMwZjY1Ni1lZWZiLTQwYjMtODhkYS04MWI5NjRlYjExMTAiLCJ0eXAiOiJJRCIsImF6cCI6ImN4dy1fT0lEQyIsIm5vbmNlIjoiSEFzRGtLWlMzcGNoSmJiRU94cEY3TlVjNWJ3ZzZDckliSGwiLCJzZXNzaW9uX3N0YXRlIjoiMGMyMzE3OTUtYTc4Yy00NDI3LThiNjEtZWI3ODZhOGUxNDdkIiwiYWNyIjoiMCIsInNfaGFzaCI6IlBMcGlEdDRndDJvZVFrMHhiRXlYeFEiLCJHcm91cCI6W10sImVtYWlsX3ZlcmlmaWVkIjpmYWxzZSwiVXNlck5hbWUiOiIwMDg3OTI4NiIsIkdyb3VwcyI6W10sInByZWZlcnJlZF91c2VybmFtZSI6IjAwODc5Mjg2In0.CIcd7ISk0Pg1Z9v11qG1iKGVunAiIHIWacXVYz0aE5LO1Vv1C-2q2WBNEf6i4sf3ulZmokRPyRP9s26tRwNb0VK1FM0Cq2j2j43WKimRzusxxU81oReGywUTtzVGT3MosofDNLhWVstp9EhT6tzsOMW-mUA_uLdzmRJx5xkJ_HX96fJZ_24CUQcUVa-bz0CShPBsKtUlGW_0cKUfmps5F2KY2lP1Ku3LjLCQW6-Qk837Vrb51EjIS7VtHQr5FjqigZuyePZjuhDkkc07eAOaDlB6Q31V2fouqz31n13sskF17jeI383eJG8gmSSt6f7M1-HyZMx0I6eMfR9930DMXQ"
	resourceName := "huaweicloud_identity_unscoped_token_with_id_token.test"

	// Avoid CheckDestroy because the token can not be destroyed.
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityUnscopedTokenWithIdToken(idpId, protocolId, idToken),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "token"),
					resource.TestCheckResourceAttrSet(resourceName, "username"),
					resource.TestCheckResourceAttrSet(resourceName, "groups.#"),
					resource.TestCheckResourceAttrSet(resourceName, "expires_at"),
				),
			},
		},
	})
}

func testAccIdentityUnscopedTokenWithIdToken(idpId, protocolId, idToken string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identity_unscoped_token_with_id_token" "test" {
  idp_id = "%[1]s"
  protocol_id = "%[2]s"
  id_token = "%[3]s"
}
`, idpId, protocolId, idToken)
}
