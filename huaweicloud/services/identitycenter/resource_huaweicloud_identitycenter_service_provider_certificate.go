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
			"state": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
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
		},
	}
}

func resourceIdentityCenterServiceProviderCertificateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		createSpCertificateHttpUrl = "v1/identity-stores/{identity_store_id}/saml-certificates"
		activeSpCertificateHttpUrl = "v1/identity-stores/{identity_store_id}/saml-certificates/{certificate_id}/active"
		product                    = "identitystore"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating Identity Center Client: %s", err)
	}

	createSpCertificatePath := client.Endpoint + createSpCertificateHttpUrl
	createSpCertificatePath = strings.ReplaceAll(createSpCertificatePath, "{identity_store_id}", fmt.Sprintf("%v", d.Get("identity_store_id")))

	createSpCertificateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	createSpCertificateResp, err := client.Request("POST", createSpCertificatePath, &createSpCertificateOpt)
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

	if status, ok := d.GetOk("state"); ok && status == "ACTIVE" {
		activeSpCertificatePath := client.Endpoint + activeSpCertificateHttpUrl
		activeSpCertificatePath = strings.ReplaceAll(activeSpCertificatePath, "{identity_store_id}", fmt.Sprintf("%v", d.Get("identity_store_id")))
		activeSpCertificatePath = strings.ReplaceAll(activeSpCertificatePath, "{certificate_id}", d.Id())

		activeSpCertificateOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
		}
		_, err = client.Request("POST", activeSpCertificatePath, &activeSpCertificateOpt)
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
		getSpCertificatesHttpUrl = "v1/identity-stores/{identity_store_id}/saml-certificates"
		getSpCertificatesProduct = "identitystore"
	)
	getSpCertificatesClient, err := cfg.NewServiceClient(getSpCertificatesProduct, region)
	if err != nil {
		return diag.Errorf("error creating Identity Center client: %s", err)
	}

	getSpCertificatesPath := getSpCertificatesClient.Endpoint + getSpCertificatesHttpUrl
	getSpCertificatesPath = strings.ReplaceAll(getSpCertificatesPath, "{identity_store_id}", fmt.Sprintf("%v", d.Get("identity_store_id")))

	getSpCertificatesOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getSpCertificatesResp, err := getSpCertificatesClient.Request("GET", getSpCertificatesPath, &getSpCertificatesOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving Identity Center service provider certificate.")
	}

	getSpCertificatesRespBody, err := utils.FlattenResponse(getSpCertificatesResp)
	if err != nil {
		return diag.FromErr(err)
	}

	certificate := utils.PathSearch(fmt.Sprintf("[?certificate_id =='%s']|[0]", d.Id()), getSpCertificatesRespBody, nil)

	if certificate == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "no sp certificate found.")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("state", utils.PathSearch("state", certificate, nil)),
		d.Set("x509certificate", utils.PathSearch("x509certificate", certificate, nil)),
		d.Set("expiry_date", utils.FormatTimeStampRFC3339(int64(utils.PathSearch("expiry_date", certificate, float64(0)).(float64))/1000, false)),
		d.Set("algorithm", utils.PathSearch("algorithm", certificate, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceIdentityCenterServiceProviderCertificateUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		activeSpCertificateHttpUrl = "v1/identity-stores/{identity_store_id}/saml-certificates/{certificate_id}/active"
		activeSpCertificateProduct = "identitystore"
	)

	if d.HasChange("state") {
		if d.Get("state").(string) == "ACTIVE" {
			activeSpCertificateClient, err := cfg.NewServiceClient(activeSpCertificateProduct, region)
			if err != nil {
				return diag.Errorf("error creating Identity Center Client: %s", err)
			}

			activeSpCertificatePath := activeSpCertificateClient.Endpoint + activeSpCertificateHttpUrl
			activeSpCertificatePath = strings.ReplaceAll(activeSpCertificatePath, "{identity_store_id}", fmt.Sprintf("%v", d.Get("identity_store_id")))
			activeSpCertificatePath = strings.ReplaceAll(activeSpCertificatePath, "{certificate_id}", d.Id())

			activeSpCertificateOpt := golangsdk.RequestOpts{
				KeepResponseBody: true,
			}
			_, err = activeSpCertificateClient.Request("POST", activeSpCertificatePath, &activeSpCertificateOpt)
			if err != nil {
				return diag.Errorf("error activing Identity Center identity provider certificate: %s", err)
			}
		} else {
			return diag.Errorf("state must be set to 'ACTIVE'")
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
	deleteCertificatePath = strings.ReplaceAll(deleteCertificatePath, "{identity_store_id}", fmt.Sprintf("%v", d.Get("identity_store_id")))
	deleteCertificatePath = strings.ReplaceAll(deleteCertificatePath, "{certificate_id}", d.Id())

	deleteCertificateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	_, err = deleteCertificateClient.Request("DELETE", deleteCertificatePath, &deleteCertificateOpt)
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
