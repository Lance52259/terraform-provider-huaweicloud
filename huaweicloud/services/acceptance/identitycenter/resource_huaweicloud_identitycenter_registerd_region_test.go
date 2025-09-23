package identitycenter

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getRegisteredRegionResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	getRegisteredRegionClient, err := cfg.NewServiceClient("identitycenter", region)
	if err != nil {
		return nil, fmt.Errorf("error creating Identity Center client: %s", err)
	}

	getRegisteredRegionHttpUrl := "v1/registered-regions"
	getRegisteredRegionPath := getRegisteredRegionClient.Endpoint + getRegisteredRegionHttpUrl

	getRegisteredRegionOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getRegisteredRegionResp, err := getRegisteredRegionClient.Request("GET", getRegisteredRegionPath, &getRegisteredRegionOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving registered region: %s", err)
	}

	return utils.FlattenResponse(getRegisteredRegionResp)
}

func TestAccRegisteredRegion_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.HW_REGION_NAME
	rName := "huaweicloud_identitycenter_registered_region.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getRegisteredRegionResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckMultiAccount(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testRegisteredRegion_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "region_id", acceptance.HW_REGION_NAME),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testRegisteredRegion_basic(name string) string {
	return fmt.Sprintf(`

resource "huaweicloud_identitycenter_registered_region" "test" {
  region_id = "%s"
}
`, name)
}
