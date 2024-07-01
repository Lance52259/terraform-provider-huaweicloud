package apig

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/apigw/dedicated/v2/apigroups"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func getGroupFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.ApigV2Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating APIG v2 client: %s", err)
	}
	return apigroups.Get(client, state.Primary.Attributes["instance_id"], state.Primary.ID).Extract()
}

func TestAccGroup_basic(t *testing.T) {
	var (
		group apigroups.Group

		name       = acceptance.RandomAccResourceName()
		updateName = acceptance.RandomAccResourceName()
		baseConfig = testAccGroup_base(name)

		rName = "huaweicloud_apig_group.test"
		rc    = acceptance.InitResourceCheck(rName, &group, getGroupFunc)

		rNameWithVariables = "huaweicloud_apig_group.with_variables"
		rcWithVariables    = acceptance.InitResourceCheck(rNameWithVariables, &group, getGroupFunc)

		rNameWithUrlDomain = "huaweicloud_apig_group.with_url_domain"
		rcWithUrlDomain    = acceptance.InitResourceCheck(rNameWithUrlDomain, &group, getGroupFunc)

		rNameWithDomainAccess = "huaweicloud_apig_group.with_domain_access"
		rcWithDomainAccess    = acceptance.InitResourceCheck(rNameWithDomainAccess, &group, getGroupFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				// Check whether illegal group name ​​can be intercepted normally (create phase).
				Config:      testAccGroup_basic_step1(baseConfig, name),
				ExpectError: regexp.MustCompile("Invalid parameter value"),
			},
			{
				Config: testAccGroup_basic_step2(baseConfig, name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "instance_id", "huaweicloud_apig_instance.test", "id"),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "Created by script"),
				),
			},
			{
				Config: testAccGroup_basic_step3(baseConfig, updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", updateName),
					resource.TestCheckResourceAttr(rName, "description", ""),
					resource.TestMatchResourceAttr(rName, "created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestMatchResourceAttr(rName, "updated_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccGroupImportStateFunc(rName),
			},
			{
				Config: testAccGroup_basic_step4(baseConfig, name),
				Check: resource.ComposeTestCheckFunc(
					rcWithVariables.CheckResourceExists(),
					resource.TestCheckResourceAttr(rNameWithVariables, "environment.#", "2"),
				),
			},
			{
				Config: testAccGroup_basic_step5(baseConfig, name),
				Check: resource.ComposeTestCheckFunc(
					rcWithVariables.CheckResourceExists(),
					resource.TestCheckResourceAttr(rNameWithVariables, "environment.#", "2"),
				),
			},
			{
				Config: testAccGroup_basic_step6(baseConfig, name),
				Check: resource.ComposeTestCheckFunc(
					rcWithUrlDomain.CheckResourceExists(),
					// since the order in the schema is inconsistent with the order of data obtained by the API, other parameters are not verified.
					resource.TestCheckResourceAttr(rNameWithUrlDomain, "url_domains.#", "2"),
					resource.TestCheckResourceAttrSet(rNameWithUrlDomain, "url_domains.0.min_ssl_version"),
					resource.TestCheckResourceAttr(rNameWithUrlDomain, "url_domains.0.is_http_redirect_to_https", "false"),
				),
			},
			{
				Config: testAccGroup_basic_step7(baseConfig, name),
				Check: resource.ComposeTestCheckFunc(
					rcWithUrlDomain.CheckResourceExists(),
					resource.TestCheckResourceAttr(rNameWithUrlDomain, "url_domains.#", "1"),
					resource.TestCheckResourceAttr(rNameWithUrlDomain, "url_domains.0.name", "www.terraform.test3.com"),
					resource.TestCheckResourceAttr(rNameWithUrlDomain, "url_domains.0.min_ssl_version", "TLSv1.1"),
					resource.TestCheckResourceAttr(rNameWithUrlDomain, "url_domains.0.is_http_redirect_to_https", "true"),
				),
			},
			{
				// Check whether illegal URL domain ​​can be intercepted normally (update phase).
				Config:      testAccGroup_basic_step8(baseConfig, name),
				ExpectError: regexp.MustCompile("error binding domain name to the API group"),
			},
			{
				Config: testAccGroup_basic_step9(baseConfig, name),
				Check: resource.ComposeTestCheckFunc(
					rcWithDomainAccess.CheckResourceExists(),
					resource.TestCheckResourceAttr(rNameWithDomainAccess, "domain_access_enabled", "false"),
				),
			},
			{
				Config: testAccGroup_basic_step10(baseConfig, name),
				Check: resource.ComposeTestCheckFunc(
					rcWithDomainAccess.CheckResourceExists(),
					resource.TestCheckResourceAttr(rNameWithDomainAccess, "domain_access_enabled", "true"),
				),
			},
		},
	})
}

func testAccGroupImportStateFunc(rsName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[rsName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", rsName, rs)
		}
		if rs.Primary.Attributes["instance_id"] == "" || rs.Primary.ID == "" {
			return "", fmt.Errorf("missing some attributes, want '<instance_id>/<id>', but got '%s/%s'",
				rs.Primary.Attributes["instance_id"], rs.Primary.ID)
		}
		return fmt.Sprintf("%s/%s", rs.Primary.Attributes["instance_id"], rs.Primary.ID), nil
	}
}

func testAccGroup_base(name string) string {
	return fmt.Sprintf(`
variable "variables_configuration" {
  type = list(object({
    name  = string
    value = string
  }))
  default = [
    {name="TEST_VAR_1", value="TEST_VALUE_1"},
    {name="TEST_VAR_2", value="TEST_VALUE_2"},
    {name="TEST_VAR_3", value="TEST_VALUE_3"},
    {name="TEST_VAR_2", value="TEST_VALUE_4"}, // same variable name, but value is different.
  ]
}

%[1]s

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_apig_instances" "test" {
  name = "tf_test_randx"
}

resource "huaweicloud_apig_instance" "test" {
  name                  = "%[2]s"
  edition               = "BASIC"
  vpc_id                = huaweicloud_vpc.test.id
  subnet_id             = huaweicloud_vpc_subnet.test.id
  security_group_id     = huaweicloud_networking_secgroup.test.id
  enterprise_project_id = "0"

  availability_zones = [
    data.huaweicloud_availability_zones.test.names[0],
  ]
}

resource "huaweicloud_apig_environment" "test" {
  count = 2

  instance_id = huaweicloud_apig_instance.test.id
  name        = format("%[2]s_%%d", count.index)
}
`, common.TestBaseNetwork(name), name)
}

func testAccGroup_basic_step1(baseConfig, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_apig_group" "test" {
  instance_id = huaweicloud_apig_instance.test.id
  name        = "INVALID_GROUP_NAME_WITH_SPECIAL_CHAR!"
  description = "Created by script"
}
`, baseConfig, name)
}

func testAccGroup_basic_step2(baseConfig, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_apig_group" "test" {
  instance_id = huaweicloud_apig_instance.test.id
  name        = "%[2]s"
  description = "Created by script"
}
`, baseConfig, name)
}

func testAccGroup_basic_step3(baseConfig, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_apig_group" "test" {
  instance_id = huaweicloud_apig_instance.test.id
  name        = "%[2]s"
}
`, baseConfig, name)
}

// Create two environments for the group, and add a total of three variables to the two environments.
// Each of the two environments has a variable with the same name and different value.
func testAccGroup_basic_withVariables(baseConfig, name string, offset int) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_apig_group" "with_variables" {
  instance_id = huaweicloud_apig_instance.test.id
  name        = "%[2]s_with_variables"

  environment {
    environment_id = huaweicloud_apig_environment.test[0].id

    dynamic "variable" {
      for_each = slice(var.variables_configuration, 0+%[3]d, 2+%[3]d)

      content {
        name  = variable.value.name
        value = variable.value.value
      }
    }
  }
  environment {
    environment_id = huaweicloud_apig_environment.test[1].id

    dynamic "variable" {
      for_each = slice(var.variables_configuration, 1+%[3]d, 3+%[3]d)

      content {
        name  = variable.value.name
        value = variable.value.value
      }
    }
  }
}
`, baseConfig, name, offset)
}

func testAccGroup_basic_step4(baseConfig, name string) string {
	return testAccGroup_basic_withVariables(baseConfig, name, 0)
}

func testAccGroup_basic_step5(baseConfig, name string) string {
	return testAccGroup_basic_withVariables(baseConfig, name, 1)
}

func testAccGroup_basic_step6(baseConfig, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_apig_group" "with_url_domain" {
  instance_id = huaweicloud_apig_instance.test.id
  name        = "%[2]s_with_url_domain"

  url_domains {
    name = "www.terraform.test1.com"
  }
  url_domains {
    name = "www.terraform.test2.com"
  }
}
`, baseConfig, name)
}

func testAccGroup_basic_step7(baseConfig, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_apig_group" "with_url_domain" {
  instance_id = huaweicloud_apig_instance.test.id
  name        = "%[2]s_with_url_domain"

  url_domains {
    name                      = "www.terraform.test3.com"
    min_ssl_version           = "TLSv1.1"
    is_http_redirect_to_https = true
  }
}
`, baseConfig, name)
}

func testAccGroup_basic_step8(baseConfig, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_apig_group" "with_url_domain" {
  instance_id = huaweicloud_apig_instance.test.id
  name        = "%[2]s_with_url_domain"

  url_domains {
    name                      = "INVALID_URL_DOMAIN"
    min_ssl_version           = "TLSv1.1"
    is_http_redirect_to_https = true
  }
}
`, baseConfig, name)
}

func testAccGroup_basic_step9(baseConfig, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_apig_group" "with_domain_access" {
  instance_id           = huaweicloud_apig_instance.test.id
  name                  = "%[2]s_with_domain_access"
  domain_access_enabled = false
}
`, baseConfig, name)
}

func testAccGroup_basic_step10(baseConfig, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_apig_group" "with_domain_access" {
  instance_id           = huaweicloud_apig_instance.test.id
  name                  = "%[2]s_with_domain_access"
  domain_access_enabled = true
}
`, baseConfig, name)
}
