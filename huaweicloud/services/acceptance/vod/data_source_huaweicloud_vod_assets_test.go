package vod

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceAssets_basic(t *testing.T) {
	var (
		name           = acceptance.RandomAccResourceNameWithDash()
		dataSourceName = "data.huaweicloud_vod_assets.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckVODMediaAsset(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAssets_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "assets.0.asset_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "assets.0.title"),
					resource.TestCheckResourceAttrSet(dataSourceName, "assets.0.media_type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "assets.0.asset_status"),
					resource.TestCheckOutput("is_asset_ids_filter_useful", "true"),
					resource.TestCheckOutput("is_query_string_filter_useful", "true"),
					resource.TestCheckOutput("is_media_types_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataSourceAssets_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_vod_assets" "test" {
  depends_on = [huaweicloud_vod_media_asset.test]
}

locals {
  asset_id = huaweicloud_vod_media_asset.test.id
  title    = huaweicloud_vod_media_asset.test.name
}

data "huaweicloud_vod_assets" "asset_ids_filter" {
  asset_ids = [local.asset_id]
}

output "is_asset_ids_filter_useful" {
  value = length(data.huaweicloud_vod_assets.asset_ids_filter.assets) >= 1 && alltrue([
    for v in data.huaweicloud_vod_assets.asset_ids_filter.assets[*].asset_id : v == local.asset_id
  ])
}

data "huaweicloud_vod_assets" "query_string_filter" {
  query_string = local.title
}

output "is_query_string_filter_useful" {
  value = length(data.huaweicloud_vod_assets.query_string_filter.assets) >= 1 && anytrue([
    for v in data.huaweicloud_vod_assets.query_string_filter.assets[*].title : v == local.title
  ])
}

data "huaweicloud_vod_assets" "media_types_filter" {
  media_types = ["MP4"]
}

output "is_media_types_filter_useful" {
  value = length(data.huaweicloud_vod_assets.media_types_filter.assets) >= 1 && alltrue([
    for v in data.huaweicloud_vod_assets.media_types_filter.assets[*].media_type : v == "MP4"
  ])
}
`, testAccMediaAsset_obs(name))
}
