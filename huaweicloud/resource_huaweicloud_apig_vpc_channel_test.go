package huaweicloud

import (
	"fmt"
	"github.com/huaweicloud/golangsdk/openstack/apigw/v2/channels"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func TestAccApigVpcChannelV2_basic(t *testing.T) {
	var (
		// Only letters, digits and underscores (_) are allowed in the name.
		rName        = fmt.Sprintf("tf_acc_test_%s", acctest.RandString(5))
		resourceName = "huaweicloud_apig_vpc_channel.test"
		channel      channels.VpcChannel
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckApigVpcChannelDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccApigVpcChannel_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApigVpcChannelExists(resourceName, &channel),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "Created by script"),
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

func testAccCheckApigVpcChannelDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*config.Config)
	client, err := config.ApigV2Client(HW_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud APIG v2 client: %s", err)
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_apig_vpc_channel" {
			continue
		}
		_, err := channels.Get(client, rs.Primary.Attributes["instance_id"], rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("APIG v2 Vpc Channel (%s) is still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckApigVpcChannelExists(n string, app *channels.VpcChannel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource %s not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no vpc channel id")
		}

		config := testAccProvider.Meta().(*config.Config)
		client, err := config.ApigV2Client(HW_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating HuaweiCloud APIG v2 client: %s", err)
		}
		found, err := channels.Get(client, rs.Primary.Attributes["instance_id"], rs.Primary.ID).Extract()
		if err != nil {
			return err
		}
		*app = *found
		return nil
	}
}

func testAccApigVpcChannel_base(rName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_images_image" "test" {
  name        = "Ubuntu 18.04 server 64bit"
  most_recent = true
}

data "huaweicloud_compute_flavors" "test" {
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  performance_type  = "normal"
  cpu_core_count    = 2
  memory_size       = 4
}

resource "huaweicloud_networking_secgroup" "test" {
  name = "%s"
}

resource "huaweicloud_compute_instance" "test" {
  name               = "%s"
  image_id           = data.huaweicloud_images_image.test.id
  flavor_id          = data.huaweicloud_compute_flavors.test.ids[0]
  security_group_ids = [huaweicloud_networking_secgroup.test.id]
  availability_zone  = data.huaweicloud_availability_zones.test.names[0]

  network {
    uuid = huaweicloud_vpc_subnet.test.id
  }
}

resource "huaweicloud_apig_instance" "test" {
  name                  = "%s"
  edition               = "BASIC"
  vpc_id                = huaweicloud_vpc.test.id
  subnet_id             = huaweicloud_vpc_subnet.test.id
  security_group_id     = huaweicloud_networking_secgroup.test.id
  enterprise_project_id = "%s"

  available_zones = [
    data.huaweicloud_availability_zones.test.names[0],
  ]
}
`, testAccApigInstance_base(rName), rName, rName, rName, HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccApigVpcChannel_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_apig_vpc_channel" "test" {
  name        = "%s"
  instance_id = huaweicloud_apig_instance.test.id
  name        = "Created by script"
  port        = 8080
  member_type = "ecs"
  algorithm   = "WRR"
  protocol    = "HTTPS"
  path        = "/"

  servers {
    instance_id = huaweicloud_compute_instance.test.id
  }
}
`, testAccApigVpcChannel_base(rName), rName)
}
