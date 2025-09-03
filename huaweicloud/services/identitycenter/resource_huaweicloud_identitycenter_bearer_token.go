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

var identityCenterBearerTokenNonUpdateParams = []string{"identity_store_id", "tenant_id"}

// @API IdentityCenter POST /v1/identity-stores/{identity_store_id}/tenant/{tenant_id}/bearer-token
// @API IdentityCenter GET /v1/identity-stores/{identity_store_id}/tenant/{tenant_id}/bearer-token
// @API IdentityCenter DELETE /v1/identity-stores/{identity_store_id}/tenant/{tenant_id}/bearer-token/{token_id}
func ResourceIdentityCenterBearerToken() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIdentityCenterBearerTokenCreate,
		UpdateContext: resourceIdentityCenterBearerTokenUpdate,
		ReadContext:   resourceIdentityCenterBearerTokenRead,
		DeleteContext: resourceIdentityCenterBearerTokenDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceIdentityCenterBearerTokenImport,
		},

		CustomizeDiff: config.FlexibleForceNew(identityCenterBearerTokenNonUpdateParams),

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
			"tenant_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"token": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"creation_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"expiration_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceIdentityCenterBearerTokenCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		createBearerTokenHttpUrl = "v1/identity-stores/{identity_store_id}/tenant/{tenant_id}/bearer-token"
		createBearerTokenProduct = "identitystore"
	)
	createBearerTokenClient, err := cfg.NewServiceClient(createBearerTokenProduct, region)
	if err != nil {
		return diag.Errorf("error creating Identity Center Client: %s", err)
	}

	createBearerTokenPath := createBearerTokenClient.Endpoint + createBearerTokenHttpUrl
	createBearerTokenPath = strings.ReplaceAll(createBearerTokenPath, "{identity_store_id}",
		fmt.Sprintf("%v", d.Get("identity_store_id")))
	createBearerTokenPath = strings.ReplaceAll(createBearerTokenPath, "{tenant_id}",
		fmt.Sprintf("%v", d.Get("tenant_id")))

	createBearerTokenOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	createBearerTokenResp, err := createBearerTokenClient.Request("POST",
		createBearerTokenPath, &createBearerTokenOpt)
	if err != nil {
		return diag.Errorf("error creating Identity Center bearer token: %s", err)
	}

	createBearerTokenRespBody, err := utils.FlattenResponse(createBearerTokenResp)
	if err != nil {
		return diag.FromErr(err)
	}

	tokenId := utils.PathSearch("token_id", createBearerTokenRespBody, "").(string)
	if tokenId == "" {
		return diag.Errorf("unable to find the Identity Center bearer token ID from the API response")
	}
	d.SetId(tokenId)

	return resourceIdentityCenterBearerTokenRead(ctx, d, meta)
}

func resourceIdentityCenterBearerTokenRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		listBearerTokensHttpUrl = "v1/identity-stores/{identity_store_id}/tenant/{tenant_id}/bearer-token"
		listBearerTokensProduct = "identitystore"
	)
	listBearerTokensClient, err := cfg.NewServiceClient(listBearerTokensProduct, region)
	if err != nil {
		return diag.Errorf("error creating Identity Center client: %s", err)
	}

	listBearerTokensPath := listBearerTokensClient.Endpoint + listBearerTokensHttpUrl
	listBearerTokensPath = strings.ReplaceAll(listBearerTokensPath, "{identity_store_id}",
		fmt.Sprintf("%v", d.Get("identity_store_id")))
	listBearerTokensPath = strings.ReplaceAll(listBearerTokensPath, "{tenant_id}",
		fmt.Sprintf("%v", d.Get("tenant_id")))

	listBearerTokensOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	listBearerTokensResp, err := listBearerTokensClient.Request("GET", listBearerTokensPath, &listBearerTokensOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving Identity Center bearer token.")
	}

	listBearerTokensRespBody, err := utils.FlattenResponse(listBearerTokensResp)
	if err != nil {
		return diag.FromErr(err)
	}
	tokens := flattenBearerTokens(listBearerTokensRespBody, d)
	if len(tokens) == 0 {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "no bearer token found.")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("token", tokens),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenBearerTokens(resp interface{}, d *schema.ResourceData) []interface{} {
	if resp == nil {
		return nil
	}
	jsonPath := fmt.Sprintf("bearer_tokens[?token_id =='%s']", d.Id())
	curJson := utils.PathSearch(jsonPath, resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	tokens := make([]interface{}, len(curArray))
	for i, v := range curArray {
		tokens[i] = map[string]interface{}{
			"creation_time": utils.FormatTimeStampRFC3339(
				int64(utils.PathSearch("creation_time", v, float64(0)).(float64))/1000, false),
			"expiration_time": utils.FormatTimeStampRFC3339(
				int64(utils.PathSearch("expiration_time", v, float64(0)).(float64))/1000, false),
		}
	}
	return tokens
}

func resourceIdentityCenterBearerTokenUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceIdentityCenterBearerTokenDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		deleteBearerTokenHttpUrl = "v1/identity-stores/{identity_store_id}/tenant/{tenant_id}/bearer-token/{token_id}"
		deleteBearerTokenProduct = "identitystore"
	)
	deleteBearerTokenClient, err := cfg.NewServiceClient(deleteBearerTokenProduct, region)
	if err != nil {
		return diag.Errorf("error creating Identity Center Client: %s", err)
	}

	deleteBearerTokenPath := deleteBearerTokenClient.Endpoint + deleteBearerTokenHttpUrl
	deleteBearerTokenPath = strings.ReplaceAll(deleteBearerTokenPath, "{identity_store_id}",
		fmt.Sprintf("%v", d.Get("identity_store_id")))
	deleteBearerTokenPath = strings.ReplaceAll(deleteBearerTokenPath, "{tenant_id}",
		fmt.Sprintf("%v", d.Get("tenant_id")))
	deleteBearerTokenPath = strings.ReplaceAll(deleteBearerTokenPath, "{token_id}", d.Id())

	deleteBearerTokenOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	_, err = deleteBearerTokenClient.Request("DELETE", deleteBearerTokenPath,
		&deleteBearerTokenOpt)
	if err != nil {
		return diag.Errorf("error deleting Identity Center bearer token: %s", err)
	}

	return nil
}

func resourceIdentityCenterBearerTokenImport(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid id format, must be <identity_store_id>/<tenant_id>/<token_id>")
	}
	d.SetId(parts[2])
	mErr := multierror.Append(nil,
		d.Set("identity_store_id", parts[0]),
		d.Set("tenant_id", parts[1]),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return nil, fmt.Errorf("failed to set value to state when import, %s", err)
	}
	return []*schema.ResourceData{d}, nil
}
