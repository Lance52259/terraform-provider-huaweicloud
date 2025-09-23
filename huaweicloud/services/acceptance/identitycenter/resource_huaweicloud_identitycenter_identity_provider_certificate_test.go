package identitycenter

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getIdentityCenterIdentityProviderCertificateResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME

	var (
		listExternalIdPCertificatesHttpUrl = "v1/identity-stores/{identity_store_id}/external-idp/{idp_id}/certificate"
		listExternalIdPCertificatesProduct = "identitystore"
	)
	listExternalIdPCertificatesClient, err := cfg.NewServiceClient(listExternalIdPCertificatesProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating Identity Center Client: %s", err)
	}

	listExternalIdPCertificatesPath := listExternalIdPCertificatesClient.Endpoint + listExternalIdPCertificatesHttpUrl
	listExternalIdPCertificatesPath = strings.ReplaceAll(listExternalIdPCertificatesPath, "{identity_store_id}",
		state.Primary.Attributes["identity_store_id"])
	listExternalIdPCertificatesPath = strings.ReplaceAll(listExternalIdPCertificatesPath, "{idp_id}",
		state.Primary.Attributes["idp_id"])

	listExternalIdPCertificatesOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	certificateId := state.Primary.ID

	listExternalIdPCertificatesResp, err := listExternalIdPCertificatesClient.Request("GET",
		listExternalIdPCertificatesPath, &listExternalIdPCertificatesOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving Identity Center identity provider certificate: %s", err)
	}

	listExternalIdPCertificatesRespBody, err := utils.FlattenResponse(listExternalIdPCertificatesResp)
	if err != nil {
		return nil, err
	}

	certificates := utils.PathSearch("idp_certificates", listExternalIdPCertificatesRespBody, nil)

	if certificates == nil {
		return nil, fmt.Errorf("error get Identity Center identity provider certificate")
	}

	for _, v := range certificates.([]interface{}) {
		if certificateId == utils.PathSearch("certificate_id", v, "") {
			return v, nil
		}
	}

	return nil, fmt.Errorf("error get Identity Center identity provider certificate")
}

func TestAccIdentityProviderCertificate_basic(t *testing.T) {
	var obj interface{}

	rName := "huaweicloud_identitycenter_identity_provider_certificate.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getIdentityCenterIdentityProviderCertificateResourceFunc,
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
				Config: testIdentityCenterIdentityProviderCertificate_basic(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "identity_store_id",
						"data.huaweicloud_identitycenter_instance.test", "identity_store_id"),
					resource.TestCheckResourceAttrSet(rName, "idp_certificate.0.issuer_name"),
					resource.TestCheckResourceAttrSet(rName, "idp_certificate.0.not_after"),
					resource.TestCheckResourceAttrSet(rName, "idp_certificate.0.not_before"),
					resource.TestCheckResourceAttrSet(rName, "idp_certificate.0.public_key"),
					resource.TestCheckResourceAttrSet(rName, "idp_certificate.0.serial_number_string"),
					resource.TestCheckResourceAttrSet(rName, "idp_certificate.0.subject_name"),
					resource.TestCheckResourceAttrSet(rName, "idp_certificate.0.version"),
					resource.TestCheckResourceAttrSet(rName, "idp_certificate.0.x509_certificate"),
				),
			},
			{
				ResourceName:            rName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateIdFunc:       testIdentityCenterIdentityProviderCertificateImportState(rName),
				ImportStateVerifyIgnore: []string{"x509_certificate_in_pem", "certificate_use"},
			},
		},
	})
}

func testIdentityCenterIdentityProviderCertificate_basic() string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_identitycenter_identity_provider_certificate" "test"{
  identity_store_id       = data.huaweicloud_identitycenter_instance.test.identity_store_id
  idp_id                  = huaweicloud_identitycenter_identity_provider.test.id
  x509_certificate_in_pem = "-----BEGIN CERTIFICATE-----MIIDUTCCAjmgAwIBAgIQAP9wc90YPLxcirh7/qTyBTANBgkqhkiG9w0BAQsFADAUMRIwEAYDVQQDDAliYWlkdS5jb20wHhcNMjUwOTAzMDI0NzAzWhcNMzUwOTAxMDI0NzAzWjAUMRIwEAYDVQQDDAliYWlkdS5jb20wggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQDZGZlyBZTE6I3RCsCdOTNvdT0rn0dzhO3squ0owQjpO5APLEVtu3rOeImkOpGIKrsx1s0Y17F+/DP2eB+YTkIOtkn01+XK9SUerL2kcAJQFMbNBtyihOxGPvq5S55NneZGY7UaPab+C9ugurmqi9xzH1NYdiDgb5Aa58t7XL2G0mnd4OkRYd43fY4lpo9jhmGhOUKZBFequ1TgV6EQg3FPTNpDzUN2skucmoU0Rz0Btz1bKdKtQ7tCHffFIPEI63dR6Rmi7CQT5WEm8f9j5UPhxFWC7f65PuaLbcTomccBbUXVm0JxaDI3L06yZyQ4U9ZKnLOSGSaZGs5hgtmeZlCNAgMBAAGjgZ4wgZswHQYDVR0OBBYEFIpds0vTmKvtGLhem/BxatDeyC+7MA4GA1UdDwEB/wQEAwIEsDAMBgNVHRMBAf8EAjAAMDsGA1UdJQQ0MDIGCCsGAQUFBwMCBggrBgEFBQcDAQYIKwYBBQUHAwMGCCsGAQUFBwMEBggrBgEFBQcDCDAfBgNVHSMEGDAWgBSKXbNL05ir7Ri4XpvwcWrQ3sgvuzANBgkqhkiG9w0BAQsFAAOCAQEARbgGdia5b0QNDxQyitntXq+Gn5JAK5Lx/5JYL0/RsgJ7uT6kVveMN6ySyDclZW2/Pvf3weCJdiz0h8NIAY0TIKd9NNJ53YCAyGTzi/BNPXAwAFJztyEwGtTIBCSIHwNmHifYzfrEFBdy33LY6xBO+W98d9NyyOppstFbRtgz4WCEdGJxRDNQ2h4oZJcIloDj54WXFyEulibbieC4oIyVP58j2MXUZwXYrhfnlir/qtaQTudjcA43+YorkTP2CBDCONm9vjINy7mDF7dTdFDjuUMyWPqokuvqLVB7zHZpKu/QhfsOBNMKgxTiNfHgqQe+EFxwvhXnxXZnkmd7F1pXqw==-----END CERTIFICATE-----"
  certificate_use         = "SIGNING"
}
`, testIdentityCenterIdentityProvider_create_without_metadata_basic())
}

func testIdentityCenterIdentityProviderCertificateImportState(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", name, rs)
		}

		identityStoreId := rs.Primary.Attributes["identity_store_id"]
		if identityStoreId == "" {
			return "", fmt.Errorf("attribute (identity_store_id) of Resource (%s) not found: %s", name, rs)
		}
		idpID := rs.Primary.Attributes["idp_id"]
		if idpID == "" {
			return "", fmt.Errorf("attribute (idp_id) of Resource (%s) not found: %s", name, rs)
		}
		certificateID := rs.Primary.ID
		return fmt.Sprintf("%s/%s/%s", identityStoreId, idpID, certificateID), nil
	}
}
