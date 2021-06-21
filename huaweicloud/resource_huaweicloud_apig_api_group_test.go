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

func TestAccApigApiGroupV2_basic(t *testing.T) {
	var (
		// Only letters, digits and underscores (_) are allowed in the name.
		rName        = fmt.Sprintf("tf_acc_test_%s", acctest.RandString(5))
		resourceName = "huaweicloud_apig_api_group.test"
		group        apigroups.Group
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckApigApiGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccApigApiGroup_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApigApiGroupExists(resourceName, &group),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "Created by script"),
				),
			},
			{
				// update name, description and app_code.
				Config: testAccApigApiGroup_update(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApigApiGroupExists(resourceName, &group),
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

func testAccCheckApigApiGroupDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*config.Config)
	client, err := config.ApigV2Client(HW_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud APIG v2 client: %s", err)
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_apig_application" {
			continue
		}
		_, err := apigroups.Get(client, rs.Primary.Attributes["instance_id"], rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("APIG v2 API group (%s) is still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckApigApiGroupExists(groupName string, app *apigroups.Group) resource.TestCheckFunc {
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

func testAccApigApiGroup_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_apig_api_group" "test" {
  name        = "%s"
  instance_id = huaweicloud_apig_instance.test.id
  description = "Created by script"
}
`, testAccApigInstance_basic(rName), rName)
}

func testAccApigApiGroup_update(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_apig_api_group" "test" {
  name        = "%s_update"
  instance_id = huaweicloud_apig_instance.test.id
  description = "Updated by script"
}
`, testAccApigInstance_basic(rName), rName)
}
