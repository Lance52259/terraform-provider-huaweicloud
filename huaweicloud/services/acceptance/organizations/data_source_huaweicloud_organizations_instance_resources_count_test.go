package organizations

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceOrganizationsResourceInstancesCount_basic(t *testing.T) {
	dataSource := "data.huaweicloud_organizations_resource_instances_count.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckMultiAccount(t)
			acceptance.TestAccPreCheckOrganizationsOpen(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceOrganizationsResourceInstancesCount_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "total_count"),
					resource.TestCheckOutput("without_any_tag_filter_is_useful", "true"),
					resource.TestCheckOutput("tags_filter_is_useful", "true"),
					resource.TestCheckOutput("matches_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceOrganizationsResourceInstancesCount_basic() string {
	return `
data "huaweicloud_organizations_resource_instances_count" "test" {
  resource_type = "organizations:accounts"
}

data "huaweicloud_organizations_resource_instances_count" "without_any_tag_filter" {
  resource_type   = "organizations:accounts"
  without_any_tag = true
}
	
output "without_any_tag_filter_is_useful" {
  value = data.huaweicloud_organizations_resource_instances_count.without_any_tag_filter.total_count > 0 
}

data "huaweicloud_organizations_resource_instances_count" "tags_filter" {
  resource_type = "organizations:accounts"

  tags {
    key    = "key3"
    values = ["value3"]
  }
}

locals {
  tag_key   = "key3"
  tag_value = "value3"
}
	
output "tags_filter_is_useful" {
  value = data.huaweicloud_organizations_resource_instances_count.tags_filter.total_count > 0 
}

data "huaweicloud_organizations_resource_instances_count" "matches_filter" {
  resource_type = "organizations:accounts"

  matches {
    key   = "key3"
    value = "value3"
  }
}

locals {
  match_key   = "resource_name"
  match_value = "tf_test"
}
	
output "matches_filter_is_useful" {
  value = data.huaweicloud_organizations_resource_instances_count.matches_filter.total_count > 0 
}
`
}
