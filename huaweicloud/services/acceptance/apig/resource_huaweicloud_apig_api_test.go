package apig

import (
	"fmt"
	"testing"

	"github.com/huaweicloud/golangsdk/openstack/apigw/v2/apis"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccApigApiV2_basic(t *testing.T) {
	var (
		// The dedica letters, digits and underscores (_) are allowed in the name.
		rName        = fmt.Sprintf("tf_acc_test_%s", acctest.RandString(5))
		resourceName = "huaweicloud_apig_api.test"
		api          apis.ApiResp
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAccPreCheck(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: testAccCheckApigApiDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccApigApi_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApigApiExists(resourceName, &api),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "port", "80"),
					resource.TestCheckResourceAttr(resourceName, "member_type", "ECS"),
					resource.TestCheckResourceAttr(resourceName, "algorithm", "WRR"),
					resource.TestCheckResourceAttr(resourceName, "protocol", "HTTP"),
					resource.TestCheckResourceAttr(resourceName, "path", "/"),
					resource.TestCheckResourceAttr(resourceName, "members.#", "1"),
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

func testAccCheckApigApiDestroy(s *terraform.State) error {
	config := acceptance.TestAccProvider.Meta().(*config.Config)
	client, err := config.ApigV2Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud APIG v2 client: %s", err)
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_apig_api" {
			continue
		}
		_, err := apis.Get(client, rs.Primary.Attributes["instance_id"], rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("APIG v2 API (%s) is still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckApigApiExists(n string, app *apis.ApiResp) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource %s not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no api id")
		}

		config := acceptance.TestAccProvider.Meta().(*config.Config)
		client, err := config.ApigV2Client(acceptance.HW_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating HuaweiCloud APIG v2 client: %s", err)
		}
		found, err := apis.Get(client, rs.Primary.Attributes["instance_id"], rs.Primary.ID).Extract()
		if err != nil {
			return err
		}
		*app = *found
		return nil
	}
}

// func testAccApigApi_base(rName string) string {
// 	return fmt.Sprintf(`

// `, testAccApigApplication_base(rName), rName, rName)
// }

func testAccApigApi_basic(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_apig_api" "test" {
  instance_id             = "525c166ffbae4458b1bd8b91ac25231c"
  group_id                = "13c37e9e8ffe45159a7826eae086be7c"
  name                    = "%s"
  version                 = "v0.1.0"
  request_protocol        = "BOTH"
  request_method          = "GET"
  request_path            = "/terraform/{resource_name}"
  security_authentication = "APP"
  matching                = "Exact"
  success_response        = "Congratulations"

  request_params {
    name        = "resource_name"
    type        = "STRING"
    location    = "PATH"
    is_required = true
    maximum     = 3
    minimum     = 64
  }
  
  backend_params {
    name      = "resourceName"
    location  = "PATH"
    req_param = "resource_name"
  }

  web {
    path           = "/getResourceName/{resourceName}"
    vpc_channel_id = "b7a7e3b173a14462be2011a17f03dae4"
    request_method = "GET"
    protocol       = "HTTPS"
    timeout        = 6000
  }

  web_policy {
    name             = "web_policy_1"
    request_protocol = "HTTPS"
    method           = "GET"
    effective_mode   = "ANY"
    path             = ""
    timeout          = 6000
    vpc_channel_id   = "b7a7e3b173a14462be2011a17f03dae4"

    backend_params {
  	name      = "resourceName"
  	location  = "PATH"
  	req_param = "resource_name"
    }

    conditions {
  	source     = "param"
  	param_name = "resourceName"
  	type       = "exact"
  	value      = "noone"
    }
  }
}
`, rName)
}

func testAccApigApi_update(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_apig_api" "test" {
  instance_id             = huaweicloud_apig_instance.test.id
  group_id                = huaweicloud_apig_group.test.id
  name                    = "tf_acc_test_0712"
  version                 = "v0.1.0"
  request_protocol        = "BOTH"
  request_method          = "GET"
  request_path            = "/terraform/{resource_name}"
  security_authentication = "APP"
  matching                = "Exact"
  success_response        = "Congratulations"
  failure_response        = "Oh no"

  request_params {
    name        = "resource_name"
    type        = "STRING"
    location    = "PATH"
    is_required = true
    maximum     = 3
    minimum     = 64
  }

  backend_params {
    name      = "resourceName"
    location  = "PATH"
    req_param = "resource_name"
  }

  web {
    path           = "/getResourceName/{resourceName}"
    vpc_channel_id = huaweicloud_apig_vpc_channel.test.id
    request_method = "GET"
    protocol       = "HTTPS"
    timeout        = 6000
  }

  web_policy {
    name             = "web_policy_1"
    request_protocol = "HTTPS"
    method           = "GET"
    effective_mode   = "ANY"
    path             = ""
    timeout          = 6000
    vpc_channel_id   = huaweicloud_apig_vpc_channel.test.id

    backend_params {
      name      = "resourceName"
      location  = "PATH"
      req_param = "resource_name"
    }

    conditions {
      source     = "param"
      param_name = "resourceName"
      type       = "exact"
      value      = "noone"
    }
  }
}
`, rName)
}
