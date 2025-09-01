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

var identityCenterApplicationCertificateNonUpdateParams = []string{"instance_id", "application_instance_id"}

// @API IdentityCenter POST /v1/instances/{instance_id}/application-instances/{application_instance_id}/certificates
// @API IdentityCenter POST /v1/instances/{instance_id}/application-instances/{application_instance_id}/certificates/{certificate_id}
// @API IdentityCenter GET /v1/instances/{instance_id}/application-instances/{application_instance_id}/certificates
// @API IdentityCenter DELETE /v1/instances/{instance_id}/application-instances/{application_instance_id}/certificates/{certificate_id}
func ResourceIdentityCenterApplicationCertificate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIdentityCenterApplicationCertificateCreate,
		UpdateContext: resourceIdentityCenterApplicationCertificateUpdate,
		ReadContext:   resourceIdentityCenterApplicationCertificateRead,
		DeleteContext: resourceIdentityCenterApplicationCertificateDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceIdentityCenterApplicationCertificateImport,
		},

		CustomizeDiff: config.FlexibleForceNew(identityCenterApplicationCertificateNonUpdateParams),

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"application_instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"application_certificate_status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"application_instance_certificate": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"algorithm": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"certificate": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"certificate_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"expiry_date": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"key_size": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"issue_date": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceIdentityCenterApplicationCertificateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		createApplicationInstanceCertificateHttpUrl       = "v1/instances/{instance_id}/application-instances/{application_instance_id}/certificates"
		updateApplicationInstanceActiveCertificateHttpUrl = "v1/instances/{instance_id}/application-instances/{application_instance_id}/certificates/{certificate_id}"
		createApplicationInstanceCertificateProduct       = "identitycenter"
	)

	client, err := cfg.NewServiceClient(createApplicationInstanceCertificateProduct, region)
	if err != nil {
		return diag.Errorf("error creating Identity Center Client: %s", err)
	}

	createApplicationInstanceCertificatePath := client.Endpoint + createApplicationInstanceCertificateHttpUrl
	createApplicationInstanceCertificatePath = strings.ReplaceAll(createApplicationInstanceCertificatePath, "{instance_id}",
		fmt.Sprintf("%v", d.Get("instance_id")))
	createApplicationInstanceCertificatePath = strings.ReplaceAll(createApplicationInstanceCertificatePath, "{application_instance_id}",
		fmt.Sprintf("%v", d.Get("application_instance_id")))

	createApplicationInstanceCertificateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	createApplicationInstanceCertificateResp, err := client.Request("POST",
		createApplicationInstanceCertificatePath, &createApplicationInstanceCertificateOpt)
	if err != nil {
		return diag.Errorf("error creating Identity Center application certificate: %s", err)
	}

	createApplicationInstanceCertificateRespBody, err := utils.FlattenResponse(createApplicationInstanceCertificateResp)
	if err != nil {
		return diag.FromErr(err)
	}

	certificateId := utils.PathSearch("application_instance_certificate.certificate_id", createApplicationInstanceCertificateRespBody, "").(string)
	if certificateId == "" {
		return diag.Errorf("unable to find the Identity Center application certificate ID from the API response")
	}
	d.SetId(certificateId)

	if d.HasChange("application_certificate_status") {
		if d.Get("application_certificate_status").(string) == "active" {
			updateApplicationInstanceActiveCertificatePath := client.Endpoint + updateApplicationInstanceActiveCertificateHttpUrl
			updateApplicationInstanceActiveCertificatePath = strings.ReplaceAll(updateApplicationInstanceActiveCertificatePath, "{instance_id}",
				fmt.Sprintf("%v", d.Get("instance_id")))
			updateApplicationInstanceActiveCertificatePath = strings.ReplaceAll(updateApplicationInstanceActiveCertificatePath, "{application_instance_id}",
				fmt.Sprintf("%v", d.Get("application_instance_id")))
			updateApplicationInstanceActiveCertificatePath = strings.ReplaceAll(updateApplicationInstanceActiveCertificatePath, "{certificate_id}", d.Id())

			updateApplicationInstanceActiveCertificateOpt := golangsdk.RequestOpts{
				KeepResponseBody: true,
			}
			_, err = client.Request("PUT", updateApplicationInstanceActiveCertificatePath,
				&updateApplicationInstanceActiveCertificateOpt)
			if err != nil {
				return diag.Errorf("error updating Identity Center application certificate status: %s", err)
			}
		}
	}

	return resourceIdentityCenterApplicationCertificateRead(ctx, d, meta)
}

func resourceIdentityCenterApplicationCertificateRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		listApplicationInstanceCertificatesHttpUrl = "v1/instances/{instance_id}/application-instances/{application_instance_id}/certificates"
		listApplicationInstanceCertificatesProduct = "identitycenter"
	)

	listApplicationInstanceCertificatesClient, err := cfg.NewServiceClient(listApplicationInstanceCertificatesProduct, region)
	if err != nil {
		return diag.Errorf("error creating Identity Center client: %s", err)
	}

	listApplicationInstanceCertificatesPath := listApplicationInstanceCertificatesClient.Endpoint + listApplicationInstanceCertificatesHttpUrl
	listApplicationInstanceCertificatesPath = strings.ReplaceAll(listApplicationInstanceCertificatesPath, "{instance_id}",
		fmt.Sprintf("%v", d.Get("instance_id")))
	listApplicationInstanceCertificatesPath = strings.ReplaceAll(listApplicationInstanceCertificatesPath, "{application_instance_id}",
		fmt.Sprintf("%v", d.Get("application_instance_id")))

	listApplicationInstanceCertificatesOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	listApplicationInstanceCertificatesResp, err := listApplicationInstanceCertificatesClient.Request("GET", listApplicationInstanceCertificatesPath, &listApplicationInstanceCertificatesOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving Identity Center service provider certificate.")
	}

	listApplicationInstanceCertificatesRespBody, err := utils.FlattenResponse(listApplicationInstanceCertificatesResp)
	if err != nil {
		return diag.FromErr(err)
	}

	certificates := flattenApplicationCertificates(listApplicationInstanceCertificatesRespBody, d)
	if len(certificates) == 0 {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "no application certificates found.")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("application_instance_certificate", certificates),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenApplicationCertificates(resp interface{}, d *schema.ResourceData) []interface{} {
	if resp == nil {
		return nil
	}
	jsonPath := fmt.Sprintf("application_instance_certificates[?certificate_id =='%s']", d.Id())
	curJson := utils.PathSearch(jsonPath, resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	certs := make([]interface{}, len(curArray))
	for i, v := range curArray {
		certs[i] = map[string]interface{}{
			"algorithm":      utils.PathSearch("algorithm", v, ""),
			"certificate":    utils.PathSearch("certificate", v, ""),
			"certificate_id": utils.PathSearch("certificate_id", v, ""),
			"status":         utils.PathSearch("status", v, ""),
			"key_size":       utils.PathSearch("key_size", v, ""),
			"expiry_date": utils.FormatTimeStampRFC3339(
				int64(utils.PathSearch("expiry_date", v, float64(0)).(float64))/1000, false),
			"issue_date": utils.FormatTimeStampRFC3339(
				int64(utils.PathSearch("issue_date", v, float64(0)).(float64))/1000, false),
		}
	}
	return certs
}

func resourceIdentityCenterApplicationCertificateUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		updateApplicationInstanceActiveCertificateHttpUrl = "v1/instances/{instance_id}/application-instances/{application_instance_id}/certificates/{certificate_id}"
		updateApplicationInstanceActiveCertificateProduct = "identitycenter"
	)

	if d.HasChange("application_certificate_status") {
		if d.Get("application_certificate_status").(string) == "active" {
			updateApplicationInstanceActiveCertificateClient, err := cfg.NewServiceClient(updateApplicationInstanceActiveCertificateProduct, region)
			if err != nil {
				return diag.Errorf("error creating Identity Center Client: %s", err)
			}

			updateApplicationInstanceActiveCertificatePath := updateApplicationInstanceActiveCertificateClient.Endpoint + updateApplicationInstanceActiveCertificateHttpUrl
			updateApplicationInstanceActiveCertificatePath = strings.ReplaceAll(updateApplicationInstanceActiveCertificatePath, "{instance_id}",
				fmt.Sprintf("%v", d.Get("instance_id")))
			updateApplicationInstanceActiveCertificatePath = strings.ReplaceAll(updateApplicationInstanceActiveCertificatePath, "{application_instance_id}",
				fmt.Sprintf("%v", d.Get("application_instance_id")))
			updateApplicationInstanceActiveCertificatePath = strings.ReplaceAll(updateApplicationInstanceActiveCertificatePath, "{certificate_id}", d.Id())

			updateApplicationInstanceActiveCertificateOpt := golangsdk.RequestOpts{
				KeepResponseBody: true,
			}
			_, err = updateApplicationInstanceActiveCertificateClient.Request("PUT", updateApplicationInstanceActiveCertificatePath,
				&updateApplicationInstanceActiveCertificateOpt)
			if err != nil {
				return diag.Errorf("error updating Identity Center application certificate status: %s", err)
			}
		}
	}

	return resourceIdentityCenterApplicationCertificateRead(ctx, d, meta)
}

func resourceIdentityCenterApplicationCertificateDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		deleteApplicationInstanceCertificateHttpUrl = "v1/instances/{instance_id}/application-instances/{application_instance_id}/certificates/{certificate_id}"
		deleteApplicationInstanceCertificateProduct = "identitycenter"
	)

	deleteApplicationInstanceCertificateClient, err := cfg.NewServiceClient(deleteApplicationInstanceCertificateProduct, region)
	if err != nil {
		return diag.Errorf("error creating Identity Center Client: %s", err)
	}

	deleteApplicationInstanceCertificatePath := deleteApplicationInstanceCertificateClient.Endpoint + deleteApplicationInstanceCertificateHttpUrl
	deleteApplicationInstanceCertificatePath = strings.ReplaceAll(deleteApplicationInstanceCertificatePath, "{instance_id}",
		fmt.Sprintf("%v", d.Get("instance_id")))
	deleteApplicationInstanceCertificatePath = strings.ReplaceAll(deleteApplicationInstanceCertificatePath, "{application_instance_id}",
		fmt.Sprintf("%v", d.Get("application_instance_id")))
	deleteApplicationInstanceCertificatePath = strings.ReplaceAll(deleteApplicationInstanceCertificatePath, "{certificate_id}", d.Id())

	deleteApplicationInstanceCertificateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = deleteApplicationInstanceCertificateClient.Request("DELETE", deleteApplicationInstanceCertificatePath,
		&deleteApplicationInstanceCertificateOpt)
	if err != nil {
		return diag.Errorf("error deleting Identity Center application certificate: %s", err)
	}

	return nil
}

func resourceIdentityCenterApplicationCertificateImport(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid id format, must be <instance_id>/<application_instance_id>/<certificate_id>")
	}
	d.SetId(parts[2])
	mErr := multierror.Append(nil,
		d.Set("instance_id", parts[0]),
		d.Set("application_instance_id", parts[1]),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return nil, fmt.Errorf("failed to set value to state when import, %s", err)
	}
	return []*schema.ResourceData{d}, nil
}
