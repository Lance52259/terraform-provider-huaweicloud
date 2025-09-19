package iam

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccIdentityTokenWithIdToken(t *testing.T) {
	idpId := "ShibbolethIdP_Terraform_OIDC"
	idToken := "eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICJqVnBkWDU1WWNjYUpHelRFQnFid2tNcEFqQXJmdnRnRkVjeTVidEYzcE9JIn0.eyJleHAiOjE3NTgyNzc0NDMsImlhdCI6MTc1ODI3NjU0MywiYXV0aF90aW1lIjoxNzU4MjUyMDg1LCJqdGkiOiIzMDVjNjFhNy1hMmI2LTQyZmQtOTkwZi1jYjY3NTljZWNjZGYiLCJpc3MiOiJodHRwczovLzEwLjIzNS4yMzEuMTg5OjU0My9hdXRoL3JlYWxtcy9tYXN0ZXIiLCJhdWQiOiJjeHctX09JREMiLCJzdWIiOiJkNGMwZjY1Ni1lZWZiLTQwYjMtODhkYS04MWI5NjRlYjExMTAiLCJ0eXAiOiJJRCIsImF6cCI6ImN4dy1fT0lEQyIsIm5vbmNlIjoiNFVQZU94YnI3VXNGeVk2ZkNNN3ZReGFJZWJ0QmozbnZ0RDYiLCJzZXNzaW9uX3N0YXRlIjoiMGMyMzE3OTUtYTc4Yy00NDI3LThiNjEtZWI3ODZhOGUxNDdkIiwiYWNyIjoiMCIsInNfaGFzaCI6IkRHM0V4bjg0bm1QZjc2cXlUaHdzdXciLCJHcm91cCI6W10sImVtYWlsX3ZlcmlmaWVkIjpmYWxzZSwiVXNlck5hbWUiOiIwMDg3OTI4NiIsIkdyb3VwcyI6W10sInByZWZlcnJlZF91c2VybmFtZSI6IjAwODc5Mjg2In0.UH5WI3PxzjOD7eYN_Wlx4XeeqqhfAyYCNfu1W9wHKJhQbye0FoMib5WnFnr8N-CiDUtKpczNeaXg5Di8PI7yfOUM-thsthhxgVSYzyBPtXFyYDYtC7zHemMuLZ4DD0PT5s-sJ824PVgqUKFu-uuy7DSun3r-dRPPTb-omYYYnNTYxDzN3-NUQdpGWwmU5ztNum1w0NveLWUO1LpOI2BCIHLiBbHrAB-Llu6DOA7BeWYEppLkjltSuJOtZiAgo9OT3IjhuCIVmf5Z11KZPCDB7Hw4EQMHHfTqcwlu7eoR_9JNBWuemuQmzD8D7DFTl9zAO9-3DQskdAITT7FUjE4OYA"
	domainName := "Timcross-WEI"
	projectName := "cn-north-4"
	resourceName := "huaweicloud_identity_token_with_id_token.test"

	// Avoid CheckDestroy because the token can not be destroyed.
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityTokenWithIdToken(idpId, idToken),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "token"),
					resource.TestCheckResourceAttrSet(resourceName, "username"),
					resource.TestCheckResourceAttrSet(resourceName, "groups.#"),
					resource.TestCheckResourceAttrSet(resourceName, "expires_at"),
				),
			},
			{
				Config: testAccIdentityTokenWithIdTokenWithDomainName(idpId, idToken, domainName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "token"),
					resource.TestCheckResourceAttrSet(resourceName, "username"),
					resource.TestCheckResourceAttrSet(resourceName, "groups.#"),
					resource.TestCheckResourceAttrSet(resourceName, "expires_at"),
				),
			},
			{
				Config: testAccIdentityTokenWithIdTokenWithProjectName(idpId, idToken, projectName),
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

func testAccIdentityTokenWithIdToken(idpId, idToken string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identity_token_with_id_token" "test" {
  idp_id = "%[1]s"
  id_token = "%[2]s"
}
`, idpId, idToken)
}

func testAccIdentityTokenWithIdTokenWithDomainName(idpId, idToken, domainName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identity_token_with_id_token" "test" {
  idp_id = "%[1]s"
  id_token = "%[2]s"
  domain_name = "%[3]s"
}
`, idpId, idToken, domainName)
}

func testAccIdentityTokenWithIdTokenWithProjectName(idpId, idToken, projectName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identity_token_with_id_token" "test" {
  idp_id = "%[1]s"
  id_token = "%[2]s"
  project_name = "%[3]s"
}
`, idpId, idToken, projectName)
}
