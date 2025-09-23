package iam

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccIdentityCheckAgencyRoleAssignment(t *testing.T) {
	var (
		agencyName1          = acceptance.RandomAccResourceName()
		dataSourceDomain     = "data.huaweicloud_identity_check_agency_role_assignment.test_domain"
		dataSourceProject    = "data.huaweicloud_identity_check_agency_role_assignment.test_project"
		dataSourceProjectAll = "data.huaweicloud_identity_check_agency_role_assignment.test_project_all"
		dataSourceDomainNot  = "data.huaweicloud_identity_check_agency_role_assignment.test_domain_not"
	)
	dcDomain := acceptance.InitDataSourceCheck(dataSourceDomain)
	dcProject := acceptance.InitDataSourceCheck(dataSourceProject)
	dcProjectAll := acceptance.InitDataSourceCheck(dataSourceProjectAll)
	dcDomainNot := acceptance.InitDataSourceCheck(dataSourceDomainNot)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAdminOnly(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityCheckAgencyRoleAssignmentDomain(agencyName1),
				Check: resource.ComposeTestCheckFunc(
					dcDomain.CheckResourceExists(),
					dcProject.CheckResourceExists(),
					dcProjectAll.CheckResourceExists(),
					dcDomainNot.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceDomain, "result", "true"),
					resource.TestCheckResourceAttr(dataSourceProject, "result", "true"),
					resource.TestCheckResourceAttr(dataSourceProjectAll, "result", "true"),
					resource.TestCheckResourceAttr(dataSourceDomainNot, "result", "false"),
				),
			},
		},
	})
}

func testAccIdentityCheckAgencyRoleAssignmentDomain(agencyName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identity_agency" "test" {
  name                  = "%s"
  description           = "This is a test agency"
  delegated_domain_name = "%s"

  project_role {
    project = "%s"
    roles   = ["CCE Administrator"]
  }
  domain_roles = [
    "Server Administrator",
    "Anti-DDoS Administrator",
  ]
  all_resources_roles = [
    "VPC Administrator"
  ]
}

data "huaweicloud_identity_check_agency_role_assignment" "test_domain" {
  agency_id  = huaweicloud_identity_agency.test.id
  role_id   = "61f48a1c78b64ed9b966cb1747337488"
  domain_id = "%s"
}

data "huaweicloud_identity_check_agency_role_assignment" "test_project" {
  agency_id   = huaweicloud_identity_agency.test.id
  role_id    = "45c78aba7ad541c4a4e3d64085fc0339"
  project_id = "%s"
}

data "huaweicloud_identity_check_agency_role_assignment" "test_project_all" {
  agency_id   = huaweicloud_identity_agency.test.id
  role_id    = "3e36e08c3bc049289790f496797b312b"
  project_id = "all"
}

data "huaweicloud_identity_check_agency_role_assignment" "test_domain_not" {
  agency_id  = huaweicloud_identity_agency.test.id
  role_id   = "3383644449894377b248351f66bf0f85"
  domain_id = "%s"
}
`, agencyName, acceptance.HW_DOMAIN_NAME, acceptance.HW_REGION_NAME,
		acceptance.HW_DOMAIN_ID, acceptance.HW_PROJECT_ID, acceptance.HW_DOMAIN_ID)
}
