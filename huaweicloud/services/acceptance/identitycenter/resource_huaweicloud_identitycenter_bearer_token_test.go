package identitycenter

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getIdentityCenterBearerTokenResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME

	var (
		listBearerTokensHttpUrl = "v1/identity-stores/{identity_store_id}/tenant/{tenant_id}/bearer-token"
		listBearerTokensProduct = "identitystore"
	)
	listBearerTokensClient, err := cfg.NewServiceClient(listBearerTokensProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating Identity Center Client: %s", err)
	}

	listBearerTokensPath := listBearerTokensClient.Endpoint + listBearerTokensHttpUrl
	listBearerTokensPath = strings.ReplaceAll(listBearerTokensPath, "{identity_store_id}", state.Primary.Attributes["identity_store_id"])
	listBearerTokensPath = strings.ReplaceAll(listBearerTokensPath, "{tenant_id}", state.Primary.Attributes["tenant_id"])

	listBearerTokensOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	tokenId := state.Primary.ID

	listBearerTokensResp, err := listBearerTokensClient.Request("GET",
		listBearerTokensPath, &listBearerTokensOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving Identity Center bearer token: %s", err)
	}

	listBearerTokensRespBody, err := utils.FlattenResponse(listBearerTokensResp)
	if err != nil {
		return nil, err
	}

	tokens := utils.PathSearch("bearer_tokens", listBearerTokensRespBody, nil)

	if tokens == nil {
		return nil, fmt.Errorf("error get Identity Center bearer token")
	}

	for _, v := range tokens.([]interface{}) {
		if tokenId == utils.PathSearch("token_id", v, "") {
			return v, nil
		}
	}

	return nil, fmt.Errorf("error get Identity Center bearer token")
}

func TestAccIdentityCenterBearerToken_basic(t *testing.T) {
	var obj interface{}

	rName := "huaweicloud_identitycenter_bearer_token.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getIdentityCenterBearerTokenResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckMultiAccount(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testIdentityCenterBearerToken_basic(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "identity_store_id",
						"data.huaweicloud_identitycenter_instance.test", "identity_store_id"),
					resource.TestCheckResourceAttrPair(rName, "tenant_id",
						"huaweicloud_identitycenter_tenant.test", "id"),
					resource.TestCheckResourceAttrSet(rName, "token.0.creation_time"),
					resource.TestCheckResourceAttrSet(rName, "token.0.expiration_time"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testIdentityCenterBearerTokenImportState(rName),
			},
		},
	})
}

func testIdentityCenterBearerToken_basic() string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_identitycenter_bearer_token" "test" {
  depends_on        = [huaweicloud_identitycenter_tenant.test]
  identity_store_id = data.huaweicloud_identitycenter_instance.test.identity_store_id
  tenant_id         = huaweicloud_identitycenter_tenant.test.id
}
`, testIdentityCenterTenant_basic())
}

func testIdentityCenterBearerTokenImportState(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", name, rs)
		}

		identityStoreId := rs.Primary.Attributes["identity_store_id"]
		if identityStoreId == "" {
			return "", fmt.Errorf("attribute (identity_store_id) of Resource (%s) not found: %s", name, rs)
		}

		tenantId := rs.Primary.Attributes["tenant_id"]
		if tenantId == "" {
			return "", fmt.Errorf("attribute (tenant_id) of Resource (%s) not found: %s", name, rs)
		}

		return identityStoreId + "/" + tenantId + "/" + rs.Primary.ID, nil
	}
}
