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

func getApplicationCertificateResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME

	var (
		listApplicationInstanceCertificatesHttpUrl = "v1/instances/{instance_id}/application-instances/{application_instance_id}/certificates"
		listApplicationInstanceCertificatesProduct = "identitycenter"
	)
	listApplicationInstanceCertificatesClient, err := cfg.NewServiceClient(listApplicationInstanceCertificatesProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating Identity Center Client: %s", err)
	}

	listApplicationInstanceCertificatesPath := listApplicationInstanceCertificatesClient.Endpoint + listApplicationInstanceCertificatesHttpUrl
	listApplicationInstanceCertificatesPath = strings.ReplaceAll(listApplicationInstanceCertificatesPath, "{instance_id}",
		state.Primary.Attributes["instance_id"])
	listApplicationInstanceCertificatesPath = strings.ReplaceAll(listApplicationInstanceCertificatesPath, "{application_instance_id}",
		state.Primary.Attributes["application_instance_id"])

	listApplicationInstanceCertificatesPath = listApplicationInstanceCertificatesPath + buildGetApplicationCertificateQueryParams("")

	listApplicationAssignmentsOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	certificateId := state.Primary.ID

	for {
		listApplicationInstanceCertificatesResp, err := listApplicationInstanceCertificatesClient.Request("GET", listApplicationInstanceCertificatesPath, &listApplicationAssignmentsOpt)

		if err != nil {
			return nil, fmt.Errorf("error retrieving Identity Center application assignment: %s", err)
		}

		listApplicationInstanceCertificatesRespBody, err := utils.FlattenResponse(listApplicationInstanceCertificatesResp)
		if err != nil {
			return nil, err
		}

		certificates := utils.PathSearch("application_instance_certificates", listApplicationInstanceCertificatesRespBody, nil)

		if certificates == nil {
			return nil, fmt.Errorf("error get Identity Center application assignment")
		}

		for _, v := range certificates.([]interface{}) {
			if certificateId == utils.PathSearch("certificate_id", v, "") {
				return v, nil
			}
		}
		marker := utils.PathSearch("page_info.next_marker", listApplicationInstanceCertificatesRespBody, nil)
		if marker == nil {
			break
		}

		listApplicationInstanceCertificatesPath = listApplicationInstanceCertificatesPath + buildGetApplicationCertificateQueryParams(marker.(string))
	}

	return nil, fmt.Errorf("error get Identity Center application certificate")
}

func buildGetApplicationCertificateQueryParams(marker string) string {
	res := "?limit=100"

	if marker != "" {
		res = fmt.Sprintf("%s&marker=%v", res, marker)
	}
	return res
}

func TestAccIdentityCenterApplicationCertificate_basic(t *testing.T) {
	var obj interface{}

	uuid, _ := uuid.GenerateUUID()
	rName := "huaweicloud_identitycenter_application_certificate.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getApplicationCertificateResourceFunc,
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
				Config: testApplicationCertificate_basic(uuid),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "instance_id",
						"data.huaweicloud_identitycenter_instance.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "application_instance_id",
						"huaweicloud_identitycenter_application_instance.test", "id"),
					resource.TestCheckResourceAttrSet(rName, "application_instance_certificate.#"),
					resource.TestCheckResourceAttrSet(rName, "application_instance_certificate.0.algorithm"),
					resource.TestCheckResourceAttrSet(rName, "application_instance_certificate.0.certificate"),
					resource.TestCheckResourceAttrSet(rName, "application_instance_certificate.0.certificate_id"),
					resource.TestCheckResourceAttrSet(rName, "application_instance_certificate.0.expiry_date"),
					resource.TestCheckResourceAttrSet(rName, "application_instance_certificate.0.status"),
					resource.TestCheckResourceAttrSet(rName, "application_instance_certificate.0.key_size"),
					resource.TestCheckResourceAttrSet(rName, "application_instance_certificate.0.issue_date"),
				),
			},
			{
				ResourceName:            rName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateIdFunc:       testIdentityCenterApplicationCertificateImportState(rName),
				ImportStateVerifyIgnore: []string{"application_certificate_status"},
			},
		},
	})
}

func testApplicationCertificate_basic(uuid string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_identitycenter_application_certificate" "test"{
  depends_on = [huaweicloud_identitycenter_application_instance.test]
  application_instance_id = huaweicloud_identitycenter_application_instance.test.id
  instance_id = data.huaweicloud_identitycenter_instance.test.id
}

`, testApplicationInstance_basic(uuid))
}

func testIdentityCenterApplicationCertificateImportState(name string) resource.ImportStateIdFunc {
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

		return fmt.Sprintf("%s/%s/%s", instanceID, applicationInstanceID, rs.Primary.ID), nil
	}
}
