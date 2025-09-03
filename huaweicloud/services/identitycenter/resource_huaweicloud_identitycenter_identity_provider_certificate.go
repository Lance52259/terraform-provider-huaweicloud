package identitycenter

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var identityCenterIdentityProviderCertificateNonUpdateParams = []string{"identity_store_id", "idp_id", "x509_certificate_in_pem", "certificate_use"}

// @API IdentityCenter POST /v1/identity-stores/{identity_store_id}/external-idp/{idp_id}/certificate
// @API IdentityCenter GET /v1/identity-stores/{identity_store_id}/external-idp/{idp_id}/certificate
// @API IdentityCenter DELETE /v1/identity-stores/{identity_store_id}/external-idp/{idp_id}/certificate/{certificate_id}
func ResourceIdentityCenterIdentityProviderCertificate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIdentityCenterIdentityProviderCertificateCreate,
		UpdateContext: resourceIdentityCenterIdentityProviderCertificateUpdate,
		ReadContext:   resourceIdentityCenterIdentityProviderCertificateRead,
		DeleteContext: resourceIdentityCenterIdentityProviderCertificateDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceIdentityCenterIdentityProviderCertificateImport,
		},

		CustomizeDiff: config.FlexibleForceNew(identityCenterIdentityProviderCertificateNonUpdateParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"identity_store_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"idp_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"x509_certificate_in_pem": {
				Type:     schema.TypeString,
				Required: true,
			},
			"certificate_use": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"idp_certificate": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"issuer_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"not_after": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"not_before": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"public_key": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"serial_number_string": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"signature_algorithm_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"subject_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"version": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"x509_certificate": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceIdentityCenterIdentityProviderCertificateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		importExternalIdPCertificateHttpUrl = "v1/identity-stores/{identity_store_id}/external-idp/{idp_id}/certificate"
		importExternalIdPCertificateProduct = "identitystore"
	)
	importExternalIdPCertificateClient, err := cfg.NewServiceClient(importExternalIdPCertificateProduct, region)
	if err != nil {
		return diag.Errorf("error creating Identity Center Client: %s", err)
	}

	importExternalIdPCertificatePath := importExternalIdPCertificateClient.Endpoint + importExternalIdPCertificateHttpUrl
	importExternalIdPCertificatePath = strings.ReplaceAll(importExternalIdPCertificatePath, "{identity_store_id}",
		fmt.Sprintf("%v", d.Get("identity_store_id")))
	importExternalIdPCertificatePath = strings.ReplaceAll(importExternalIdPCertificatePath, "{idp_id}",
		fmt.Sprintf("%v", d.Get("idp_id")))

	importExternalIdPCertificateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildImportExternalIdPCertificateBodyParams(d)),
	}

	importExternalIdPCertificateResp, err := importExternalIdPCertificateClient.Request("POST",
		importExternalIdPCertificatePath, &importExternalIdPCertificateOpt)
	if err != nil {
		return diag.Errorf("error creating Identity Center identity provider certificate: %s", err)
	}

	importExternalIdPCertificateRespBody, err := utils.FlattenResponse(importExternalIdPCertificateResp)
	if err != nil {
		return diag.FromErr(err)
	}

	certificateId := utils.PathSearch("certificate_id", importExternalIdPCertificateRespBody, "").(string)
	if certificateId == "" {
		return diag.Errorf("unable to find the Identity Center identity provider certificate ID from the API response")
	}
	d.SetId(certificateId)

	return resourceIdentityCenterIdentityProviderCertificateRead(ctx, d, meta)
}

func buildImportExternalIdPCertificateBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"x509_certificate_in_pem": utils.ValueIgnoreEmpty(d.Get("x509_certificate_in_pem")),
		"certificate_use":         utils.ValueIgnoreEmpty(d.Get("certificate_use")),
	}
	return bodyParams
}

func resourceIdentityCenterIdentityProviderCertificateRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		listExternalIdPCertificatesHttpUrl = "v1/identity-stores/{identity_store_id}/external-idp/{idp_id}/certificate"
		listExternalIdPCertificatesProduct = "identitystore"
	)
	listExternalIdPCertificatesClient, err := cfg.NewServiceClient(listExternalIdPCertificatesProduct, region)
	if err != nil {
		return diag.Errorf("error creating Identity Center client: %s", err)
	}

	listExternalIdPCertificatesPath := listExternalIdPCertificatesClient.Endpoint + listExternalIdPCertificatesHttpUrl
	listExternalIdPCertificatesPath = strings.ReplaceAll(listExternalIdPCertificatesPath, "{identity_store_id}",
		fmt.Sprintf("%v", d.Get("identity_store_id")))
	listExternalIdPCertificatesPath = strings.ReplaceAll(listExternalIdPCertificatesPath, "{idp_id}",
		fmt.Sprintf("%v", d.Get("idp_id")))

	listExternalIdPCertificatesOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	listExternalIdPCertificatesResp, err := listExternalIdPCertificatesClient.Request("GET", listExternalIdPCertificatesPath, &listExternalIdPCertificatesOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving Identity Center identity provider certificate.")
	}

	listExternalIdPCertificatesRespBody, err := utils.FlattenResponse(listExternalIdPCertificatesResp)
	if err != nil {
		return diag.FromErr(err)
	}

	certificates := flattenIdpCertificates(listExternalIdPCertificatesRespBody, d)

	if len(certificates) == 0 {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "no idp certificates found.")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("idp_certificate", certificates),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenIdpCertificates(resp interface{}, d *schema.ResourceData) []interface{} {
	if resp == nil {
		return nil
	}
	jsonPath := fmt.Sprintf("idp_certificates[?certificate_id =='%s']", d.Id())
	curJson := utils.PathSearch(jsonPath, resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	certs := make([]interface{}, len(curArray))
	for i, v := range curArray {
		certs[i] = map[string]interface{}{
			"issuer_name": utils.PathSearch("issuer_name", v, nil),
			"not_after": utils.FormatTimeStampRFC3339(
				int64(utils.PathSearch("not_after", v, float64(0)).(float64))/1000, false),
			"not_before": utils.FormatTimeStampRFC3339(
				int64(utils.PathSearch("not_before", v, float64(0)).(float64))/1000, false),
			"public_key":               utils.PathSearch("public_key", v, nil),
			"serial_number_string":     utils.PathSearch("serial_number_string", v, nil),
			"subject_name":             utils.PathSearch("subject_name", v, nil),
			"version":                  utils.PathSearch("version", v, nil),
			"x509_certificate":         utils.PathSearch("x509_Certificate_in_pem", v, nil),
			"signature_algorithm_name": utils.PathSearch("signature_algorithm_name", v, nil),
		}
	}
	return certs
}

func resourceIdentityCenterIdentityProviderCertificateUpdate(context.Context, *schema.ResourceData, interface{}) diag.Diagnostics {
	return nil
}

func resourceIdentityCenterIdentityProviderCertificateDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		deleteExternalIdPCertificateHttpUrl = "v1/identity-stores/{identity_store_id}/external-idp/{idp_id}/certificate/{certificate_id}"
		deleteExternalIdPCertificateProduct = "identitystore"
	)
	deleteExternalIdPCertificateClient, err := cfg.NewServiceClient(deleteExternalIdPCertificateProduct, region)
	if err != nil {
		return diag.Errorf("error creating Identity Center Client: %s", err)
	}

	deleteExternalIdPCertificatePath := deleteExternalIdPCertificateClient.Endpoint + deleteExternalIdPCertificateHttpUrl
	deleteExternalIdPCertificatePath = strings.ReplaceAll(deleteExternalIdPCertificatePath, "{identity_store_id}",
		fmt.Sprintf("%v", d.Get("identity_store_id")))
	deleteExternalIdPCertificatePath = strings.ReplaceAll(deleteExternalIdPCertificatePath, "{idp_id}",
		fmt.Sprintf("%v", d.Get("idp_id")))
	deleteExternalIdPCertificatePath = strings.ReplaceAll(deleteExternalIdPCertificatePath, "{certificate_id}", d.Id())

	deleteBearerTokenOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	_, err = deleteExternalIdPCertificateClient.Request("DELETE", deleteExternalIdPCertificatePath,
		&deleteBearerTokenOpt)
	if err != nil {
		return diag.Errorf("error deleting Identity Center identity provider certificate: %s", err)
	}

	return nil
}

func resourceIdentityCenterIdentityProviderCertificateImport(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid id format, must be <identity_store_id>/<idp_id>/<certificate_id>")
	}
	d.SetId(parts[2])
	mErr := multierror.Append(nil,
		d.Set("identity_store_id", parts[0]),
		d.Set("idp_id", parts[1]),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return nil, fmt.Errorf("failed to set value to state when import, %s", err)
	}
	return []*schema.ResourceData{d}, nil
}
