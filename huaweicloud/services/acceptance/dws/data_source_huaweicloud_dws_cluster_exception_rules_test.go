package dws

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataClusterExceptionRules_basic(t *testing.T) {
	var (
		all   = "data.huaweicloud_dws_cluster_exception_rules.all"
		dcAll = acceptance.InitDataSourceCheck(all)

		byName   = "data.huaweicloud_dws_cluster_exception_rules.filter_by_rule_name"
		dcByName = acceptance.InitDataSourceCheck(byName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDwsClusterId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataClusterExceptionRules_basic(),
				Check: resource.ComposeTestCheckFunc(
					// Without filter parameters.
					dcAll.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "rules.#", regexp.MustCompile(`^[0-9]+$`)),
					// Filter by rule name.
					dcByName.CheckResourceExists(),
					resource.TestCheckOutput("is_rule_name_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataClusterExceptionRules_basic() string {
	name := acceptance.RandomAccResourceName()

	return fmt.Sprintf(`
variable "exception_rule_configurations" {
  description = "The configurations of the exception rule."
  type        = list(object({
    key   = string
    value = string
  }))

  default = [
    {
      key   = "action"
      value = "penalty"
    },
	{
		key   = "blocktime"
		value = "300"
	},
	{
		key   = "elapsedtime"
		value = "400"
	},
	{
		key   = "allcputime"
		value = "500"
	},
  ]
}

resource "huaweicloud_dws_cluster_exception_rule" "test" {
  cluster_id = "%[1]s"
  name       = "%[2]s"

  dynamic "configurations" {
    for_each = var.exception_rule_configurations

    content {
      key   = configurations.value.key
      value = configurations.value.value
    }
  }
}

# Query all exception rules in the cluster
data "huaweicloud_dws_cluster_exception_rules" "all" {
  depends_on = [
    huaweicloud_dws_cluster_exception_rule.test,
  ]

  cluster_id = "%[1]s"
}

# Filter the exception rules with a specified name in the cluster (fuzzy matching)
locals {
  rule_name = huaweicloud_dws_cluster_exception_rule.test.name
}

data "huaweicloud_dws_cluster_exception_rules" "filter_by_rule_name" {
  # The behavior of parameter 'rule_name' of the exception rule resource is 'Required', means this parameter does not
  # have 'Know After Apply' behavior.
  depends_on = [
    huaweicloud_dws_cluster_exception_rule.test,
  ]

  cluster_id = "%[1]s"
  rule_name  = local.rule_name
}

locals {
  rule_name_filter_result = [for v in data.huaweicloud_dws_cluster_exception_rules.filter_by_rule_name.rules[*].name : v == local.rule_name]
}

output "is_rule_name_filter_useful" {
  value = length(local.rule_name_filter_result) > 0 && alltrue(local.rule_name_filter_result)
}
`, acceptance.HW_DWS_CLUSTER_ID, name)
}
