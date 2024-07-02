package apig

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/apig"
)

func getInstanceFeatureFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.ApigV2Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating APIG client: %s", err)
	}

	return apig.GetInstanceFeature(client, state.Primary.Attributes["instance_id"], state.Primary.ID)
}

func TestAccInstanceFeature_basic(t *testing.T) {
	var (
		feature interface{}

		rName      = "huaweicloud_apig_instance_feature.test"
		name       = acceptance.RandomAccResourceName()
		baseConfig = testAccInstanceFeature_base(name)
		rc         = acceptance.InitResourceCheck(rName, &feature, getInstanceFeatureFunc)
	)

	// Avoid CheckDestroy because this resource already exists and does not need to be deleted.
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				// Check whether illegal feature name ​​can be intercepted normally (create phase).
				Config:      testAccInstanceFeature_basic_step1(baseConfig),
				ExpectError: regexp.MustCompile(`error creating instance feature`),
			},
			{
				Config: testAccInstanceFeature_basic_step2(baseConfig),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", "ratelimit"),
					resource.TestCheckResourceAttr(rName, "enabled", "true"),
					resource.TestCheckResourceAttr(rName, "config", "{\"api_limits\":200}"),
				),
			},
			{
				// Check whether illegal feature content ​​can be intercepted normally (update phase).
				Config:      testAccInstanceFeature_basic_step3(baseConfig),
				ExpectError: regexp.MustCompile(`error updating feature \(ratelimit\) under specified instance \([a-f0-9-]+\)`),
			},
			{
				Config: testAccInstanceFeature_basic_step4(baseConfig),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", "ratelimit"),
					resource.TestCheckResourceAttr(rName, "config", "{\"api_limits\":300}"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccInstanceFeatureResourceImportStateFunc(rName),
			},
		},
	})
}

func testAccInstanceFeatureResourceImportStateFunc(rName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[rName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", rName, rs)
		}
		instanceId := rs.Primary.Attributes["instance_id"]
		featureName := rs.Primary.ID
		if instanceId == "" || featureName == "" {
			return "", fmt.Errorf("missing some attributes, want '<instance_id>/<name>', but '%s/%s'",
				instanceId, featureName)
		}
		return fmt.Sprintf("%s/%s", instanceId, featureName), nil
	}
}

func testAccInstanceFeature_base(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_apig_instance" "test" {
  name                  = "%[2]s"
  edition               = "BASIC"
  vpc_id                = huaweicloud_vpc.test.id
  subnet_id             = huaweicloud_vpc_subnet.test.id
  security_group_id     = huaweicloud_networking_secgroup.test.id
  enterprise_project_id = "0"
  availability_zones    = try(slice(data.huaweicloud_availability_zones.test.names, 0, 1), null)
}
`, common.TestBaseNetwork(name), name)
}

func testAccInstanceFeature_basic_step1(baseConfig string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_apig_instance_feature" "test" {
  instance_id = huaweicloud_apig_instance.test.id
  name        = "INVALID_FEATURE_NAME"
  enabled     = true

  config = jsonencode({
    api_limits = 200
  })
}
`, baseConfig)
}

func testAccInstanceFeature_basic_step2(baseConfig string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_apig_instance_feature" "test" {
  instance_id = huaweicloud_apig_instance.test.id
  name        = "ratelimit"
  enabled     = true

  config = jsonencode({
    api_limits = 200
  })
}
`, baseConfig)
}

func testAccInstanceFeature_basic_step3(baseConfig string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_apig_instance_feature" "test" {
  instance_id = huaweicloud_apig_instance.test.id
  name        = "ratelimit"
  enabled     = true

  config = jsonencode({
    INVALID_KEY = "INVALID_VALUE"
  })
}
`, baseConfig)
}

func testAccInstanceFeature_basic_step4(baseConfig string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_apig_instance_feature" "test" {
  instance_id = huaweicloud_apig_instance.test.id
  name        = "ratelimit"
  enabled     = true

  config = jsonencode({
    api_limits = 300
  })
}
`, baseConfig)
}
