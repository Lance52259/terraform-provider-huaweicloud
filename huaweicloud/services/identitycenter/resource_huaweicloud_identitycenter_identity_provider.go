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

var identityCenterIdentityProviderNonUpdateParams = []string{"identity_store_id", "idp_saml_metadata", "idp_certificate"}

// @API IdentityCenter POST /v1/identity-stores/{identity_store_id}/external-idp
// @API IdentityCenter PUT /v1/identity-stores/{identity_store_id}/external-idp/{idp_id}
// @API IdentityCenter GET /v1/identity-stores/{identity_store_id}/external-idp
// @API IdentityCenter DELETE /v1/identity-stores/{identity_store_id}/external-idp/{idp_id}
// @API IdentityCenter POST /v1/identity-stores/{identity_store_id}/external-idp/{idp_id}/enable
// @API IdentityCenter POST /v1/identity-stores/{identity_store_id}/external-idp/{idp_id}/disable
func ResourceIdentityCenterIdentityProvider() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIdentityCenterIdentityProviderCreate,
		UpdateContext: resourceIdentityCenterIdentityProviderUpdate,
		ReadContext:   resourceIdentityCenterIdentityProviderRead,
		DeleteContext: resourceIdentityCenterIdentityProviderDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceIdentityCenterIdentityProviderImport,
		},

		CustomizeDiff: config.FlexibleForceNew(identityCenterIdentityProviderNonUpdateParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"identity_store_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"idp_saml_metadata": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"idp_certificate": {
				Type:         schema.TypeString,
				Optional:     true,
				RequiredWith: []string{"entity_id", "login_url"},
			},
			"entity_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"login_url": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"want_request_signed": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"is_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"idp_certificate_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"certificate_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceIdentityCenterIdentityProviderCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		createExternalIdPConfigurationForDirectoryHttpUrl = "v1/identity-stores/{identity_store_id}/external-idp"
		enableExternalIdPConfigurationForDirectoryHttpUrl = "v1/identity-stores/{identity_store_id}/external-idp/{idp_id}/enable"
		identityStoreProduct                              = "identitystore"
	)

	createExternalIdPConfigurationForDirectoryClient, err := cfg.NewServiceClient(identityStoreProduct, region)
	if err != nil {
		return diag.Errorf("error creating Identity Center client: %s", err)
	}

	createExternalIdPConfigurationForDirectoryPath := createExternalIdPConfigurationForDirectoryClient.Endpoint + createExternalIdPConfigurationForDirectoryHttpUrl
	createExternalIdPConfigurationForDirectoryPath = strings.ReplaceAll(createExternalIdPConfigurationForDirectoryPath, "{identity_store_id}",
		fmt.Sprintf("%v", d.Get("identity_store_id")))

	createExternalIdPConfigurationForDirectoryOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateExternalIdPConfigurationForDirectoryBodyParams(d)),
	}

	createExternalIdPConfigurationForDirectoryResp, err := createExternalIdPConfigurationForDirectoryClient.Request("POST", createExternalIdPConfigurationForDirectoryPath, &createExternalIdPConfigurationForDirectoryOpt)
	if err != nil {
		return diag.Errorf("error creating identity provider: %s", err)
	}

	createExternalIdPConfigurationForDirectoryRespBody, err := utils.FlattenResponse(createExternalIdPConfigurationForDirectoryResp)
	if err != nil {
		return diag.FromErr(err)
	}

	idpId := utils.PathSearch("idp_id", createExternalIdPConfigurationForDirectoryRespBody, "").(string)
	if idpId == "" {
		return diag.Errorf("unable to find the Identity Center identity provider ID from the API response")
	}

	d.SetId(idpId)

	enableExternalIdPConfigurationForDirectoryClient, err := cfg.NewServiceClient(identityStoreProduct, region)
	if err != nil {
		return diag.Errorf("error creating Identity Center client: %s", err)
	}

	enableExternalIdPConfigurationForDirectoryPath := enableExternalIdPConfigurationForDirectoryClient.Endpoint + enableExternalIdPConfigurationForDirectoryHttpUrl
	enableExternalIdPConfigurationForDirectoryPath = strings.ReplaceAll(enableExternalIdPConfigurationForDirectoryPath, "{identity_store_id}",
		fmt.Sprintf("%v", d.Get("identity_store_id")))
	enableExternalIdPConfigurationForDirectoryPath = strings.ReplaceAll(enableExternalIdPConfigurationForDirectoryPath, "{idp_id}", d.Id())

	enableExternalIdPConfigurationForDirectoryOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = enableExternalIdPConfigurationForDirectoryClient.Request("POST", enableExternalIdPConfigurationForDirectoryPath, &enableExternalIdPConfigurationForDirectoryOpt)
	if err != nil {
		return diag.Errorf("error enabling identity provider: %s", err)
	}

	return resourceIdentityCenterIdentityProviderRead(ctx, d, meta)
}

func buildCreateExternalIdPConfigurationForDirectoryBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"idp_saml_metadata": utils.ValueIgnoreEmpty(d.Get("idp_saml_metadata")),
		"idp_certificate":   utils.ValueIgnoreEmpty(d.Get("idp_certificate")),
		"idp_saml_config": map[string]interface{}{
			"entity_id": utils.ValueIgnoreEmpty(d.Get("entity_id")),
			"login_url": utils.ValueIgnoreEmpty(d.Get("login_url")),
		},
	}
	return bodyParams
}

func resourceIdentityCenterIdentityProviderRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		listExternalIdPConfigurationsForDirectoryHttpUrl = "v1/identity-stores/{identity_store_id}/external-idp"
		listExternalIdPConfigurationsForDirectoryProduct = "identitystore"
	)
	listExternalIdPConfigurationsForDirectoryClient, err := cfg.NewServiceClient(listExternalIdPConfigurationsForDirectoryProduct, region)
	if err != nil {
		return diag.Errorf("error creating Identity Center client: %s", err)
	}

	listExternalIdPConfigurationsForDirectoryPath := listExternalIdPConfigurationsForDirectoryClient.Endpoint + listExternalIdPConfigurationsForDirectoryHttpUrl
	listExternalIdPConfigurationsForDirectoryPath = strings.ReplaceAll(listExternalIdPConfigurationsForDirectoryPath, "{identity_store_id}",
		fmt.Sprintf("%v", d.Get("identity_store_id")))

	listExternalIdPConfigurationsForDirectoryOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	listExternalIdPConfigurationsForDirectoryResp, err := listExternalIdPConfigurationsForDirectoryClient.Request("GET", listExternalIdPConfigurationsForDirectoryPath, &listExternalIdPConfigurationsForDirectoryOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving identity provider.")
	}

	listExternalIdPConfigurationsForDirectoryRespBody, err := utils.FlattenResponse(listExternalIdPConfigurationsForDirectoryResp)
	if err != nil {
		return diag.FromErr(err)
	}
	associations, _ := utils.PathSearch("associations", listExternalIdPConfigurationsForDirectoryRespBody, make([]interface{}, 0)).([]interface{})
	if len(associations) == 0 {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "no identity provider found.")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("entity_id", utils.PathSearch("associations|[0].idp_saml_config.entity_id", listExternalIdPConfigurationsForDirectoryRespBody, nil)),
		d.Set("login_url", utils.PathSearch("associations|[0].idp_saml_config.login_url", listExternalIdPConfigurationsForDirectoryRespBody, nil)),
		d.Set("is_enabled", utils.PathSearch("associations|[0].is_enabled", listExternalIdPConfigurationsForDirectoryRespBody, nil)),
		d.Set("want_request_signed", utils.PathSearch("associations|[0].idp_saml_config.want_request_signed", listExternalIdPConfigurationsForDirectoryRespBody, nil)),
		d.Set("idp_certificate_ids", flattenIdentityCenterIdentityProviderCertificate(listExternalIdPConfigurationsForDirectoryRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenIdentityCenterIdentityProviderCertificate(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("associations|[0].idp_certificate_ids", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	certs := make([]interface{}, len(curArray))
	for i, v := range curArray {
		certs[i] = map[string]interface{}{
			"certificate_id": utils.PathSearch("certificate_id", v, ""),
			"status":         utils.PathSearch("status", v, ""),
		}
	}
	return certs
}

func resourceIdentityCenterIdentityProviderUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		updateExternalIdPConfigurationsForDirectoryHttpUrl = "v1/identity-stores/{identity_store_id}/external-idp/{idp_id}"
		updateExternalIdPConfigurationsForDirectoryProduct = "identitystore"
	)
	updateExternalIdPConfigurationsForDirectoryClient, err := cfg.NewServiceClient(updateExternalIdPConfigurationsForDirectoryProduct, region)
	if err != nil {
		return diag.Errorf("error creating Identity Center client: %s", err)
	}

	updateExternalIdPConfigurationsForDirectoryChanges := []string{
		"entity_id",
		"login_url",
	}

	if d.HasChanges(updateExternalIdPConfigurationsForDirectoryChanges...) {
		updateExternalIdPConfigurationsForDirectoryPath := updateExternalIdPConfigurationsForDirectoryClient.Endpoint + updateExternalIdPConfigurationsForDirectoryHttpUrl
		updateExternalIdPConfigurationsForDirectoryPath = strings.ReplaceAll(updateExternalIdPConfigurationsForDirectoryPath, "{identity_store_id}",
			fmt.Sprintf("%v", d.Get("identity_store_id")))
		updateExternalIdPConfigurationsForDirectoryPath = strings.ReplaceAll(updateExternalIdPConfigurationsForDirectoryPath, "{idp_id}", d.Id())

		updateExternalIdPConfigurationsForDirectoryOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			JSONBody:         utils.RemoveNil(buildUpdateExternalIdPConfigurationsForDirectoryBodyParams(d)),
		}

		_, err = updateExternalIdPConfigurationsForDirectoryClient.Request("PUT", updateExternalIdPConfigurationsForDirectoryPath, &updateExternalIdPConfigurationsForDirectoryOpt)
		if err != nil {
			return diag.Errorf("error updating identity provider: %s", err)
		}
	}

	return resourceIdentityCenterIdentityProviderRead(ctx, d, meta)
}

func buildUpdateExternalIdPConfigurationsForDirectoryBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"entity_id": utils.ValueIgnoreEmpty(d.Get("entity_id")),
		"login_url": utils.ValueIgnoreEmpty(d.Get("login_url")),
	}
	return bodyParams
}

func resourceIdentityCenterIdentityProviderDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		deleteExternalIdPConfigurationForDirectoryHttpUrl  = "v1/identity-stores/{identity_store_id}/external-idp/{idp_id}"
		disableExternalIdPConfigurationForDirectoryHttpUrl = "v1/identity-stores/{identity_store_id}/external-idp/{idp_id}/disable"
		identityStoreProduct                               = "identitystore"
	)

	disableExternalIdPConfigurationForDirectoryClient, err := cfg.NewServiceClient(identityStoreProduct, region)
	if err != nil {
		return diag.Errorf("error creating Identity Center client: %s", err)
	}

	disableExternalIdPConfigurationForDirectoryPath := disableExternalIdPConfigurationForDirectoryClient.Endpoint + disableExternalIdPConfigurationForDirectoryHttpUrl
	disableExternalIdPConfigurationForDirectoryPath = strings.ReplaceAll(disableExternalIdPConfigurationForDirectoryPath, "{identity_store_id}",
		fmt.Sprintf("%v", d.Get("identity_store_id")))
	disableExternalIdPConfigurationForDirectoryPath = strings.ReplaceAll(disableExternalIdPConfigurationForDirectoryPath, "{idp_id}", d.Id())

	disableExternalIdPConfigurationForDirectoryOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	_, err = disableExternalIdPConfigurationForDirectoryClient.Request("POST", disableExternalIdPConfigurationForDirectoryPath, &disableExternalIdPConfigurationForDirectoryOpt)
	if err != nil {
		return diag.Errorf("error deleting identity provider: %s", err)
	}

	deleteExternalIdPConfigurationForDirectoryClient, err := cfg.NewServiceClient(identityStoreProduct, region)
	if err != nil {
		return diag.Errorf("error creating Identity Center client: %s", err)
	}

	deleteExternalIdPConfigurationForDirectoryPath := deleteExternalIdPConfigurationForDirectoryClient.Endpoint + deleteExternalIdPConfigurationForDirectoryHttpUrl
	deleteExternalIdPConfigurationForDirectoryPath = strings.ReplaceAll(deleteExternalIdPConfigurationForDirectoryPath, "{identity_store_id}",
		fmt.Sprintf("%v", d.Get("identity_store_id")))
	deleteExternalIdPConfigurationForDirectoryPath = strings.ReplaceAll(deleteExternalIdPConfigurationForDirectoryPath, "{idp_id}", d.Id())

	deleteExternalIdPConfigurationForDirectoryOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	_, err = deleteExternalIdPConfigurationForDirectoryClient.Request("DELETE", deleteExternalIdPConfigurationForDirectoryPath, &deleteExternalIdPConfigurationForDirectoryOpt)
	if err != nil {
		return diag.Errorf("error deleting identity provider: %s", err)
	}

	return nil

}

func resourceIdentityCenterIdentityProviderImport(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid id format, must be <identity_store_id>/<idp_id>")
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
