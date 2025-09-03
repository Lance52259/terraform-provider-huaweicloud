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

var identityCenterTenantNonUpdateParams = []string{"identity_store_id"}

// @API IdentityCenter POST /v1/identity-stores/{identity_store_id}/provision-tenant
// @API IdentityCenter GET /v1/identity-stores/{identity_store_id}/provision-tenant
// @API IdentityCenter DELETE /v1/identity-stores/{identity_store_id}/tenant/{tenant_id}
func ResourceIdentityCenterTenant() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIdentityCenterTenantCreate,
		UpdateContext: resourceIdentityCenterTenantUpdate,
		ReadContext:   resourceIdentityCenterTenantRead,
		DeleteContext: resourceIdentityCenterTenantDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceIdentityCenterTenantImport,
		},

		CustomizeDiff: config.FlexibleForceNew(identityCenterTenantNonUpdateParams),

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
			"creation_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"scim_endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceIdentityCenterTenantCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		createProvisioningTenantHttpUrl = "v1/identity-stores/{identity_store_id}/provision-tenant"
		createProvisioningTenantProduct = "identitystore"
	)
	createProvisioningTenantClient, err := cfg.NewServiceClient(createProvisioningTenantProduct, region)
	if err != nil {
		return diag.Errorf("error creating Identity Center Client: %s", err)
	}

	createProvisioningTenantPath := createProvisioningTenantClient.Endpoint + createProvisioningTenantHttpUrl
	createProvisioningTenantPath = strings.ReplaceAll(createProvisioningTenantPath, "{identity_store_id}",
		fmt.Sprintf("%v", d.Get("identity_store_id")))

	createProvisioningTenantOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	createProvisioningTenantResp, err := createProvisioningTenantClient.Request("POST",
		createProvisioningTenantPath, &createProvisioningTenantOpt)
	if err != nil {
		return diag.Errorf("error creating Identity Center tenant: %s", err)
	}

	createProvisioningTenantRespBody, err := utils.FlattenResponse(createProvisioningTenantResp)
	if err != nil {
		return diag.FromErr(err)
	}

	tenantId := utils.PathSearch("tenant_id", createProvisioningTenantRespBody, "").(string)
	if tenantId == "" {
		return diag.Errorf("unable to find the Identity Center tenant ID from the API response")
	}
	d.SetId(tenantId)

	return resourceIdentityCenterTenantRead(ctx, d, meta)
}

func resourceIdentityCenterTenantRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		listProvisioningTenantsHttpUrl = "v1/identity-stores/{identity_store_id}/provision-tenant"
		listProvisioningTenantsProduct = "identitystore"
	)
	listProvisioningTenantsClient, err := cfg.NewServiceClient(listProvisioningTenantsProduct, region)
	if err != nil {
		return diag.Errorf("error creating Identity Center client: %s", err)
	}

	listProvisioningTenantsPath := listProvisioningTenantsClient.Endpoint + listProvisioningTenantsHttpUrl
	listProvisioningTenantsPath = strings.ReplaceAll(listProvisioningTenantsPath, "{identity_store_id}",
		fmt.Sprintf("%v", d.Get("identity_store_id")))

	listProvisioningTenantsOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	listProvisioningTenantsResp, err := listProvisioningTenantsClient.Request("GET", listProvisioningTenantsPath, &listProvisioningTenantsOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving Identity Center provisioning tenant.")
	}

	listProvisioningTenantsRespBody, err := utils.FlattenResponse(listProvisioningTenantsResp)
	if err != nil {
		return diag.FromErr(err)
	}
	tenants, _ := utils.PathSearch("provisioning_tenants", listProvisioningTenantsRespBody, make([]interface{}, 0)).([]interface{})
	if len(tenants) == 0 {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "no tenants found.")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("creation_time", utils.FormatTimeStampRFC3339(
			int64(utils.PathSearch("provisioning_tenants|[0].creation_time", listProvisioningTenantsRespBody, float64(0)).(float64))/1000, false)),
		d.Set("scim_endpoint", utils.PathSearch("provisioning_tenants|[0].scim_endpoint", listProvisioningTenantsRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceIdentityCenterTenantUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceIdentityCenterTenantDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		deleteProvisioningTenantHttpUrl = "v1/identity-stores/{identity_store_id}/tenant/{tenant_id}"
		deleteProvisioningTenantProduct = "identitystore"
	)
	deleteProvisioningTenantClient, err := cfg.NewServiceClient(deleteProvisioningTenantProduct, region)
	if err != nil {
		return diag.Errorf("error creating Identity Center Client: %s", err)
	}

	deleteProvisioningTenantPath := deleteProvisioningTenantClient.Endpoint + deleteProvisioningTenantHttpUrl
	deleteProvisioningTenantPath = strings.ReplaceAll(deleteProvisioningTenantPath, "{identity_store_id}",
		fmt.Sprintf("%v", d.Get("identity_store_id")))
	deleteProvisioningTenantPath = strings.ReplaceAll(deleteProvisioningTenantPath, "{tenant_id}", d.Id())

	deleteIdentityCenterUserOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	_, err = deleteProvisioningTenantClient.Request("DELETE", deleteProvisioningTenantPath,
		&deleteIdentityCenterUserOpt)
	if err != nil {
		return diag.Errorf("error deleting Identity Center Tenant: %s", err)
	}

	return nil
}

func resourceIdentityCenterTenantImport(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid id format, must be <identity_store_id>/<tenant_id>")
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
