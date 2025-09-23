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

var identityCenterServiceProviderCertificateNonUpdateParams = []string{"identity_store_id"}

// @API IdentityCenter POST /v1/identity-stores/{identity_store_id}/saml-certificates
// @API IdentityCenter POST /v1/identity-stores/{identity_store_id}/saml-certificates/{certificate_id}/active
// @API IdentityCenter GET /v1/identity-stores/{identity_store_id}/saml-certificates
// @API IdentityCenter DELETE /v1/identity-stores/{identity_store_id}/saml-certificates/{certificate_id}
func ResourceIdentityCenterServiceProviderCertificate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIdentityCenterServiceProviderCertificateCreate,
		UpdateContext: resourceIdentityCenterServiceProviderCertificateUpdate,
		ReadContext:   resourceIdentityCenterServiceProviderCertificateRead,
		DeleteContext: resourceIdentityCenterServiceProviderCertificateDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceIdentityCenterServiceProviderCertificateImport,
		},

		CustomizeDiff: config.FlexibleForceNew(identityCenterServiceProviderCertificateNonUpdateParams),

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
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"certificate": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"x509certificate": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"algorithm": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"expiry_date": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"state": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"certificate_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"sp_certificate_status": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceIdentityCenterServiceProviderCertificateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		createSpCertificateHttpUrl = "v1/identity-stores/{identity_store_id}/saml-certificates"
		activeSpCertificateHttpUrl = "v1/identity-stores/{identity_store_id}/saml-certificates/{certificate_id}/active"
		createSpCertificateProduct = "identitystore"
	)
	createSpCertificateClient, err := cfg.NewServiceClient(createSpCertificateProduct, region)
	if err != nil {
		return diag.Errorf("error creating Identity Center Client: %s", err)
	}

	createSpCertificatePath := createSpCertificateClient.Endpoint + createSpCertificateHttpUrl
	createSpCertificatePath = strings.ReplaceAll(createSpCertificatePath, "{identity_store_id}",
		fmt.Sprintf("%v", d.Get("identity_store_id")))

	createSpCertificateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	createSpCertificateResp, err := createSpCertificateClient.Request("POST",
		createSpCertificatePath, &createSpCertificateOpt)
	if err != nil {
		return diag.Errorf("error creating Identity Center service provider certificate: %s", err)
	}

	createSpCertificateRespBody, err := utils.FlattenResponse(createSpCertificateResp)
	if err != nil {
		return diag.FromErr(err)
	}

	certificateId := utils.PathSearch("certificate_id", createSpCertificateRespBody, "").(string)
	if certificateId == "" {
		return diag.Errorf("unable to find the Identity Center service provider certificate ID from the API response")
	}
	d.SetId(certificateId)

	if status, ok := d.GetOk("sp_certificate_status"); ok && status == "active" {
		activeSpCertificateClient, err := cfg.NewServiceClient(createSpCertificateProduct, region)
		if err != nil {
			return diag.Errorf("error creating Identity Center Client: %s", err)
		}

		activeSpCertificatePath := activeSpCertificateClient.Endpoint + activeSpCertificateHttpUrl
		activeSpCertificatePath = strings.ReplaceAll(activeSpCertificatePath, "{identity_store_id}",
			fmt.Sprintf("%v", d.Get("identity_store_id")))
		activeSpCertificatePath = strings.ReplaceAll(activeSpCertificatePath, "{certificate_id}", d.Id())

		activeSpCertificateOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
		}
		_, err = activeSpCertificateClient.Request("POST", activeSpCertificatePath,
			&activeSpCertificateOpt)
		if err != nil {
			return diag.Errorf("error activing Identity Center identity provider certificate: %s", err)
		}
	}

	return resourceIdentityCenterServiceProviderCertificateRead(ctx, d, meta)
}

func resourceIdentityCenterServiceProviderCertificateRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		listSpCertificatesHttpUrl = "v1/identity-stores/{identity_store_id}/saml-certificates"
		listSpCertificatesProduct = "identitystore"
	)
	listSpCertificatesClient, err := cfg.NewServiceClient(listSpCertificatesProduct, region)
	if err != nil {
		return diag.Errorf("error creating Identity Center client: %s", err)
	}

	listSpCertificatesPath := listSpCertificatesClient.Endpoint + listSpCertificatesHttpUrl
	listSpCertificatesPath = strings.ReplaceAll(listSpCertificatesPath, "{identity_store_id}",
		fmt.Sprintf("%v", d.Get("identity_store_id")))

	listSpCertificatesOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	listSpCertificatesResp, err := listSpCertificatesClient.Request("GET", listSpCertificatesPath, &listSpCertificatesOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving Identity Center service provider certificate.")
	}

	listSpCertificatesRespBody, err := utils.FlattenResponse(listSpCertificatesResp)
	if err != nil {
		return diag.FromErr(err)
	}

	certificates := flattenSpCertificates(listSpCertificatesRespBody, d)

	if len(certificates) == 0 {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "no sp certificate found.")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("certificate", certificates),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenSpCertificates(resp interface{}, d *schema.ResourceData) []interface{} {
	if resp == nil {
		return nil
	}
	jsonPath := fmt.Sprintf("[?certificate_id =='%s']", d.Id())
	curJson := utils.PathSearch(jsonPath, resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	certs := make([]interface{}, len(curArray))
	for i, v := range curArray {
		certs[i] = map[string]interface{}{
			"x509certificate": utils.PathSearch("x509certificate", v, nil),
			"expiry_date": utils.FormatTimeStampRFC3339(
				int64(utils.PathSearch("expiry_date", v, float64(0)).(float64))/1000, false),
			"algorithm":      utils.PathSearch("algorithm", v, nil),
			"state":          utils.PathSearch("state", v, nil),
			"certificate_id": utils.PathSearch("certificate_id", v, nil),
		}
	}
	return certs
}

func resourceIdentityCenterServiceProviderCertificateUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		activeSpCertificateHttpUrl = "v1/identity-stores/{identity_store_id}/saml-certificates/{certificate_id}/active"
		activeSpCertificateProduct = "identitystore"
	)

	if d.HasChange("sp_certificate_status") {
		if d.Get("sp_certificate_status").(string) == "active" {
			activeSpCertificateClient, err := cfg.NewServiceClient(activeSpCertificateProduct, region)
			if err != nil {
				return diag.Errorf("error creating Identity Center Client: %s", err)
			}

			activeSpCertificatePath := activeSpCertificateClient.Endpoint + activeSpCertificateHttpUrl
			activeSpCertificatePath = strings.ReplaceAll(activeSpCertificatePath, "{identity_store_id}",
				fmt.Sprintf("%v", d.Get("identity_store_id")))
			activeSpCertificatePath = strings.ReplaceAll(activeSpCertificatePath, "{certificate_id}", d.Id())

			activeSpCertificateOpt := golangsdk.RequestOpts{
				KeepResponseBody: true,
			}
			_, err = activeSpCertificateClient.Request("POST", activeSpCertificatePath,
				&activeSpCertificateOpt)
			if err != nil {
				return diag.Errorf("error activing Identity Center identity provider certificate: %s", err)
			}
		}
	}
	return resourceIdentityCenterServiceProviderCertificateRead(ctx, d, meta)
}

func resourceIdentityCenterServiceProviderCertificateDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		deleteCertificateHttpUrl = "v1/identity-stores/{identity_store_id}/saml-certificates/{certificate_id}"
		deleteCertificateProduct = "identitystore"
	)
	deleteCertificateClient, err := cfg.NewServiceClient(deleteCertificateProduct, region)
	if err != nil {
		return diag.Errorf("error creating Identity Center Client: %s", err)
	}

	deleteCertificatePath := deleteCertificateClient.Endpoint + deleteCertificateHttpUrl
	deleteCertificatePath = strings.ReplaceAll(deleteCertificatePath, "{identity_store_id}",
		fmt.Sprintf("%v", d.Get("identity_store_id")))
	deleteCertificatePath = strings.ReplaceAll(deleteCertificatePath, "{certificate_id}", d.Id())

	deleteCertificateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	_, err = deleteCertificateClient.Request("DELETE", deleteCertificatePath,
		&deleteCertificateOpt)
	if err != nil {
		return diag.Errorf("error deleting Identity Center service provider certificate: %s", err)
	}

	return nil
}

func resourceIdentityCenterServiceProviderCertificateImport(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid id format, must be <identity_store_id>/<certificate_id>")
	}
	d.SetId(parts[1])
	mErr := multierror.Append(nil,
		d.Set("identity_store_id", parts[0]),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return nil, fmt.Errorf("failed to set value to state when import, %s", err)
	}
	return []*schema.ResourceData{d}, nil
}
