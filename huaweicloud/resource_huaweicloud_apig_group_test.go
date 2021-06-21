package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/huaweicloud/golangsdk/openstack/apigw/v2/apigroups"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func TestAccApigGroupV2_basic(t *testing.T) {
	var (
		// Only letters, digits and underscores (_) are allowed in the name.
		rName        = fmt.Sprintf("tf_acc_test_%s", acctest.RandString(5))
		resourceName = "huaweicloud_apig_group.test"
		group        apigroups.Group
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckApigGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccApigGroup_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApigGroupExists(resourceName, &group),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "Created by script"),
				),
			},
			{
				// update name, description and app_code.
				Config: testAccApigGroup_update(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApigGroupExists(resourceName, &group),
					resource.TestCheckResourceAttr(resourceName, "name", rName+"_update"),
					resource.TestCheckResourceAttr(resourceName, "description", "Updated by script"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccApigInstanceSubResourceImportStateIdFunc(resourceName),
			},
		},
	})
}

func TestAccApigGroupV2_variables(t *testing.T) {
	var (
		// Only letters, digits and underscores (_) are allowed in the name.
		rName        = fmt.Sprintf("tf_acc_test_%s", acctest.RandString(5))
		resourceName = "huaweicloud_apig_group.test"
		group        apigroups.Group
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckApigGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccApigGroup_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApigGroupExists(resourceName, &group),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
				),
			},
			{
				// update name, description and app_code.
				Config: testAccApigGroup_variablesBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApigGroupExists(resourceName, &group),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "environments.#", "1"),
				),
			},
			{
				// update name, description and app_code.
				Config: testAccApigGroup_variablesUpdate(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApigGroupExists(resourceName, &group),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "environments.#", "1"),
				),
			},
			{
				Config: testAccApigGroup_basic(rName), // Remove all custom environment variables.
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApigGroupExists(resourceName, &group),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccApigInstanceSubResourceImportStateIdFunc(resourceName),
			},
		},
	})
}

func TestAccApigGroupV2_responses(t *testing.T) {
	var (
		// Only letters, digits and underscores (_) are allowed in the name.
		rName        = fmt.Sprintf("tf_acc_test_%s", acctest.RandString(5))
		resourceName = "huaweicloud_apig_group.test"
		group        apigroups.Group
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckApigGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccApigGroup_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApigGroupExists(resourceName, &group),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
				),
			},
			{
				// update name, description and app_code.
				Config: testAccApigGroup_responsesBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApigGroupExists(resourceName, &group),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "custom_responses.#", "1"),
				),
			},
			{
				Config: testAccApigGroup_responsesUpdate(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApigGroupExists(resourceName, &group),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "custom_responses.#", "1"),
				),
			},
			{
				Config: testAccApigGroup_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApigGroupExists(resourceName, &group),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccApigInstanceSubResourceImportStateIdFunc(resourceName),
			},
		},
	})
}

func testAccCheckApigGroupDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*config.Config)
	client, err := config.ApigV2Client(HW_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud APIG v2 client: %s", err)
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_apig_group" {
			continue
		}
		_, err := apigroups.Get(client, rs.Primary.Attributes["instance_id"], rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("APIG v2 API group (%s) is still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckApigGroupExists(groupName string, app *apigroups.Group) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[groupName]
		if !ok {
			return fmt.Errorf("Resource %s not found", groupName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No APIG V2 API group Id")
		}

		config := testAccProvider.Meta().(*config.Config)
		client, err := config.ApigV2Client(HW_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating HuaweiCloud APIG v2 client: %s", err)
		}
		found, err := apigroups.Get(client, rs.Primary.Attributes["instance_id"], rs.Primary.ID).Extract()
		if err != nil {
			return fmt.Errorf("APIG v2 API group not exist: %s", err)
		}
		*app = *found
		return nil
	}
}

func testAccApigGroup_base(rName string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_vpc" "test" {
  name = "%s"
  cidr = "192.168.0.0/16"
}

resource "huaweicloud_vpc_subnet" "test" {
  name       = "%s"
  vpc_id     = huaweicloud_vpc.test.id
  gateway_ip = "192.168.0.1"
  cidr       = "192.168.0.0/24"
}

resource "huaweicloud_networking_secgroup" "test" {
  name = "%s"
}

resource "huaweicloud_apig_instance" "test" {
  name                  = "%s"
  edition               = "BASIC"
  vpc_id                = huaweicloud_vpc.test.id
  subnet_id             = huaweicloud_vpc_subnet.test.id
  security_group_id     = huaweicloud_networking_secgroup.test.id
  enterprise_project_id = "%s"

  available_zone_ids = [
    data.huaweicloud_availability_zones.test.names[0],
  ]
}
`, rName, rName, rName, rName, HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccApigGroup_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_apig_group" "test" {
  name        = "%s"
  instance_id = huaweicloud_apig_instance.test.id
  description = "Created by script"
}
`, testAccApigGroup_base(rName), rName)
}

func testAccApigGroup_update(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_apig_group" "test" {
  name        = "%s_update"
  instance_id = huaweicloud_apig_instance.test.id
  description = "Updated by script"
}
`, testAccApigGroup_base(rName), rName)
}

func testAccApigGroup_variablesBase(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_apig_environment" "test" {
  name        = "%s"
  instance_id = huaweicloud_apig_instance.test.id
}
`, testAccApigGroup_base(rName), rName)
}

func testAccApigGroup_variablesBasic(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_apig_group" "test" {
  name        = "%s"
  instance_id = huaweicloud_apig_instance.test.id
  description = "Created by script"

  environments {
    environment_id = huaweicloud_apig_environment.test.id

    variables {
      name = "%s"
      value = "/terraform/testPath"
    }
  }
}
`, testAccApigGroup_variablesBase(rName), rName, rName)
}

func testAccApigGroup_variablesUpdate(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_apig_group" "test" {
  name        = "%s"
  instance_id = huaweicloud_apig_instance.test.id
  description = "Created by script"

  environments {
    environment_id = huaweicloud_apig_environment.test.id

    variables {
      name = "%s"
      value = "/terraform/newTestPath"
    }
  }
}
`, testAccApigGroup_variablesBase(rName), rName, rName)
}

func testAccApigGroup_responsesBasic(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_apig_group" "test" {
  name        = "%s"
  instance_id = huaweicloud_apig_instance.test.id
  description = "Created by script"

  custom_responses {
    name = "%s"

    responses {
      error_type  = "ACCESS_DENIED"
      body        = "{\"error_code\":\"$context.error.code\",\"request_id\":\"$context.requestId\"}"
      status_code = 402
    }
  }
}
`, testAccApigGroup_base(rName), rName, rName)
}

func testAccApigGroup_responsesUpdate(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_apig_group" "test" {
  name        = "%s"
  instance_id = huaweicloud_apig_instance.test.id
  description = "Created by script"

  custom_responses {
    name = "%s"

    responses {
      error_type  = "ACCESS_DENIED"
      body        = "{\"error_code\":\"$context.error.code\",\"error_msg\":\"$context.error.message\"}"
      status_code = 412
    }
  }
}
`, testAccApigGroup_base(rName), rName, rName)
}
