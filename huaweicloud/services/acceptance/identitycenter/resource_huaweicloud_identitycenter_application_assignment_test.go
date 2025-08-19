package identitycenter

import (
	"fmt"
	"github.com/hashicorp/go-uuid"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getApplicationAssignmentResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME

	var (
		listApplicationAssignmentsHttpUrl = "v1/instances/{instance_id}/applications/{application_instance_id}/assignments"
		listApplicationAssignmentsProduct = "identitycenter"
	)

	listApplicationAssignmentsClient, err := cfg.NewServiceClient(listApplicationAssignmentsProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating Identity Center Client: %s", err)
	}

	listApplicationAssignmentsPath := listApplicationAssignmentsClient.Endpoint + listApplicationAssignmentsHttpUrl
	listApplicationAssignmentsPath = strings.ReplaceAll(listApplicationAssignmentsPath, "{instance_id}",
		state.Primary.Attributes["instance_id"])
	listApplicationAssignmentsPath = strings.ReplaceAll(listApplicationAssignmentsPath, "{application_instance_id}",
		state.Primary.Attributes["application_instance_id"])

	listApplicationAssignmentsPath = listApplicationAssignmentsPath + buildGetApplicationAssignmentQueryParams("")

	listApplicationAssignmentsOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	principalId := state.Primary.Attributes["principal_id"]

	for {
		listApplicationAssignmentsResp, err := listApplicationAssignmentsClient.Request("GET", listApplicationAssignmentsPath, &listApplicationAssignmentsOpt)
		if err != nil {
			return nil, fmt.Errorf("error retrieving Identity Center application assignment: %s", err)
		}

		listApplicationAssignmentsRespBody, err := utils.FlattenResponse(listApplicationAssignmentsResp)
		if err != nil {
			return nil, err
		}

		assignments := utils.PathSearch("application_assignments", listApplicationAssignmentsRespBody, nil)
		if assignments == nil {
			return nil, fmt.Errorf("error get Identity Center application assignment")
		}

		for _, v := range assignments.([]interface{}) {
			if principalId == utils.PathSearch("principal_id", v, "") {
				return v, nil
			}
		}
		marker := utils.PathSearch("page_info.next_marker", listApplicationAssignmentsRespBody, nil)
		if marker == nil {
			break
		}

		listApplicationAssignmentsPath = listApplicationAssignmentsPath + buildGetApplicationAssignmentQueryParams(marker.(string))
	}

	return nil, fmt.Errorf("error get Identity Center application assignment")
}

func buildGetApplicationAssignmentQueryParams(marker string) string {
	res := "?limit=100"

	if marker != "" {
		res = fmt.Sprintf("%s&marker=%v", res, marker)
	}
	return res
}

func TestAccIdentityCenterApplicationAssignment_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	uuid, _ := uuid.GenerateUUID()
	rName := "huaweicloud_identitycenter_application_assignment.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getApplicationAssignmentResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckMultiAccount(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testApplicationAssignment_basic(name, uuid),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "instance_id",
						"data.huaweicloud_identitycenter_instance.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "application_instance_id",
						"huaweicloud_identitycenter_application_instance.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "principal_id",
						"huaweicloud_identitycenter_user.test", "id"),
					resource.TestCheckResourceAttr(rName, "principal_type", "USER"),
					resource.TestCheckResourceAttrSet(rName, "application_urn"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testIdentityCenterApplicationAssignmentImportState(rName),
			},
		},
	})
}

func testApplicationAssignment_basic(name string, uuid string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_identitycenter_catalog_applications" "test"{}

data "huaweicloud_identitycenter_application_templates" "test"{
  application_id = data.huaweicloud_identitycenter_catalog_applications.test.applications[0].application_id
}

resource "huaweicloud_identitycenter_application_instance" "test"{
  depends_on = [huaweicloud_identitycenter_user.test]
  name = "%[2]s"
  template_id = data.huaweicloud_identitycenter_application_templates.test.application_templates[0].template_id
  instance_id = data.huaweicloud_identitycenter_instance.test.id
  display_name = "create"
  description = "create"
  response_schema_config = "{\"properties\":{\"key1\":{\"attr_name_format\":\"urn:oasis:names:tc:SAML:2.0:attrname-format:basic\",\"include\":\"YES\"},\"key2\":{\"attr_name_format\":\"urn:oasis:names:tc:SAML:2.0:attrname-format:unspecified\",\"include\":\"YES\"},\"key5\":{\"attr_name_format\":\"urn:oasis:names:tc:SAML:2.0:attrname-format:basic\",\"include\":\"YES\"},\"key3\":{\"attr_name_format\":\"urn:oasis:names:tc:SAML:2.0:attrname-format:uri\",\"include\":\"YES\"},\"key4\":{\"attr_name_format\":\"urn:oasis:names:tc:SAML:2.0:attrname-format:basic\",\"include\":\"YES\"}},\"subject\":{\"name_id_format\":\"urn:oasis:names:tc:SAML:1.1:nameid-format:emailAddress\",\"include\":\"REQUIRED\"},\"supported_name_id_formats\":[]}"
  response_config = "{\"properties\":{\"key1\":{\"source\":[\"$${user:email}\"]},\"key2\":{\"source\":[\"$${user:familyName}\"]},\"key5\":{\"source\":[\"$${user:preferredUsername}\"]},\"key3\":{\"source\":[\"$${user:givenName}\"]},\"key4\":{\"source\":[\"$${user:familyName}\"]}},\"subject\":{\"source\":[\"$${user:name}\"]},\"relay_state\":null,\"ttl\":\"PT1H\"}"
  security_config {
    ttl = "P9M"
  }
  service_provider_config {
    audience = "https://create.com"
    require_request_signature = false
    consumers{
        binding = "urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST"
        default_value = true
        location = "https://create.com"
    }
  }
}

resource "huaweicloud_identitycenter_application_assignment" "test"{
  depends_on = [huaweicloud_identitycenter_application_instance.test]
  application_instance_id = huaweicloud_identitycenter_application_instance.test.id
  instance_id = data.huaweicloud_identitycenter_instance.test.id
  principal_id = huaweicloud_identitycenter_user.test.id
  principal_type = "USER"
}

`, testIdentityCenterUser_basic(name), uuid)
}

func testIdentityCenterApplicationAssignmentImportState(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", name, rs)
		}

		instanceID := rs.Primary.Attributes["instance_id"]
		if instanceID == "" {
			return "", fmt.Errorf("attribute (instance_id) of Resource (%s) not found: %s", name, rs)
		}
		applicationInstanceID := rs.Primary.Attributes["application_instance_id"]
		if applicationInstanceID == "" {
			return "", fmt.Errorf("attribute (application_instance_id) of Resource (%s) not found: %s", name, rs)
		}
		principalID := rs.Primary.Attributes["principal_id"]
		if principalID == "" {
			return "", fmt.Errorf("attribute (principal_id) of Resource (%s) not found: %s", name, rs)
		}
		return fmt.Sprintf("%s/%s/%s", instanceID, applicationInstanceID, principalID), nil
	}
}
