package identitycenter

import (
	"context"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

var identityCenterRegisteredRegionParams = []string{"region_id"}

// @API IdentityCenter POST /v1/register-regions
// @API IdentityCenter GET  /v1/registered-regions
func ResourceIdentityCenterRegisteredRegion() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIdentityCenterRegisteredRegionCreate,
		UpdateContext: resourceIdentityCenterRegisteredRegionUpdate,
		ReadContext:   resourceIdentityCenterRegisteredRegionRead,
		DeleteContext: resourceIdentityCenterRegisteredRegionDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: config.FlexibleForceNew(identityCenterRegisteredRegionParams),

		Schema: map[string]*schema.Schema{
			"region_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func resourceIdentityCenterRegisteredRegionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	regionId := d.Get("region_id").(string)

	var (
		registerRegionHttpUrl = "v1/register-regions"
		registerRegionProduct = "identitycenter"
	)
	client, err := cfg.NewServiceClient(registerRegionProduct, region)
	if err != nil {
		return diag.Errorf("error creating IdentityCenter client: %s", err)
	}

	registerRegionPath := client.Endpoint + registerRegionHttpUrl

	registerRegionOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildRegisterRegionBodyParams(d)),
	}

	_, err = client.Request("POST", registerRegionPath, &registerRegionOpt)
	if err != nil {
		return diag.Errorf("error registering IdentityCenter region: %s", err)
	}

	d.SetId(regionId)

	return resourceIdentityCenterRegisteredRegionRead(ctx, d, meta)
}

func buildRegisterRegionBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"region_id": utils.ValueIgnoreEmpty(d.Get("region_id")),
	}
	return bodyParams
}

func resourceIdentityCenterRegisteredRegionRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		describeRegisteredRegionsHttpUrl = "v1/registered-regions"
		describeRegisteredRegionsProduct = "identitycenter"
	)
	describeRegisteredRegionsClient, err := cfg.NewServiceClient(describeRegisteredRegionsProduct, region)
	if err != nil {
		return diag.Errorf("error creating Identity Center Client: %s", err)
	}

	describeRegisteredRegionsPath := describeRegisteredRegionsClient.Endpoint + describeRegisteredRegionsHttpUrl

	describeRegisteredRegionsOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	describeRegisteredRegionsResp, err := describeRegisteredRegionsClient.Request("GET", describeRegisteredRegionsPath, &describeRegisteredRegionsOpt)

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving Identity Center registered region")
	}

	describeRegisteredRegionsRespBody, err := utils.FlattenResponse(describeRegisteredRegionsResp)
	if err != nil {
		return diag.FromErr(err)
	}

	regions, _ := utils.PathSearch("regions", describeRegisteredRegionsRespBody, make([]interface{}, 0)).([]interface{})
	if len(regions) == 0 {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "no region found.")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region_id", utils.PathSearch("regions|[0].region_id", describeRegisteredRegionsRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceIdentityCenterRegisteredRegionUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceIdentityCenterRegisteredRegionDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	errorMsg := "Deleting IdentityCenter region resource is not supported. The resource is only removed from the state."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
