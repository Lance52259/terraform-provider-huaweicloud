package iam

import (
	"context"
	"strings"

	"github.com/chnsz/golangsdk"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API IAM POST /v5/users/{user_id}/access-keys
// @API IAM DELETE /v5/users/{user_id}/access-keys/{access_key_id}
// @API IAM GET /v5/users/{user_id}/access-keys
// @API IAM GET /v5/users/{user_id}/access-keys/{access_key_id}/last-used
// @API IAM PUT /v5/users/{user_id}/access-keys/{access_key_id}
var accessKeyV5NonUpdatableParams = []string{"user_id"}

func ResourceIdentityAccessKey() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIdentityV5AccessKeyCreate,
		ReadContext:   resourceIdentityV5AccessKeyRead,
		UpdateContext: resourceIdentityV5AccessKeyUpdate,
		DeleteContext: resourceIdentityV5AccessKeyDelete,

		CustomizeDiff: config.FlexibleForceNew(accessKeyV5NonUpdatableParams),

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"user_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"last_used_at": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"secret_access_key": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
		},
	}
}

func resourceIdentityV5AccessKeyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	iamClient, err := cfg.IAMNoVersionClient(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	userId := d.Get("user_id").(string)
	createAccessKeyHttpUrl := "v5/users/{user_id}/access-keys"
	createAccessKeyPath := iamClient.Endpoint + createAccessKeyHttpUrl
	createAccessKeyPath = strings.ReplaceAll(createAccessKeyPath, "{user_id}", userId)
	createAccessKeyOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	createAccessKeyResp, err := iamClient.Request("POST", createAccessKeyPath, &createAccessKeyOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error create IAM user access key")
	}
	createAccessKeyBody, err := utils.FlattenResponse(createAccessKeyResp)
	if err != nil {
		return diag.FromErr(err)
	}
	access := utils.PathSearch("access_key.access_key_id", createAccessKeyBody, "").(string)
	if access == "" {
		return common.CheckDeletedDiag(d, err, "error get IAM access key : access key is not found in API response")
	}
	d.SetId(access)
	secret := utils.PathSearch("access_key.secret_access_key", createAccessKeyBody, "").(string)
	if secret == "" {
		return common.CheckDeletedDiag(d, err, "error get IAM secret key : secret key is not found in API response")
	}
	if err := d.Set("secret_access_key", secret); err != nil {
		return diag.Errorf("error setting IAM field : %s", err)
	}
	return resourceIdentityV5AccessKeyRead(ctx, d, meta)
}

func resourceIdentityV5AccessKeyRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	iamClient, err := cfg.IAMNoVersionClient(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	getAccessKeyHttpUrl := "v5/users/{user_id}/access-keys"
	getAccessKeyPath := iamClient.Endpoint + getAccessKeyHttpUrl
	getAccessKeyPath = strings.ReplaceAll(getAccessKeyPath, "{user_id}", d.Get("user_id").(string))
	getAccessKeyOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getAccessKeyResp, err := iamClient.Request("GET", getAccessKeyPath, &getAccessKeyOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error get IAM user access key")
	}
	getAccessKeyRespBody, err := utils.FlattenResponse(getAccessKeyResp)
	if err != nil {
		return diag.FromErr(err)
	}
	accessKey := utils.PathSearch("access_keys", getAccessKeyRespBody, nil).([]interface{})
	if accessKey == nil {
		return common.CheckDeletedDiag(d, err, "error getting IAM access keys: keys is not found in API response")
	}
	body := checkAccessExist(d, accessKey)
	if body == nil {
		return common.CheckDeletedDiag(d, err, "error getting IAM access key: key is not found in API response")
	}
	mErr := multierror.Append(nil,
		d.Set("created_at", utils.PathSearch("created_at", body, nil)),
		d.Set("status", utils.PathSearch("status", body, nil)),
	)
	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting IAM access key fields: %s", err)
	}
	accessKeyId := utils.PathSearch("access_key_id", body, nil).(string)
	if accessKeyId != "" {
		getLastUsedHttpUrl := "v5/users/{user_id}/access-keys/{access_key_id}/last-used"
		getLastUsedPath := iamClient.Endpoint + getLastUsedHttpUrl
		getLastUsedPath = strings.ReplaceAll(getLastUsedPath, "{user_id}", d.Get("user_id").(string))
		getLastUsedPath = strings.ReplaceAll(getLastUsedPath, "{access_key_id}", accessKeyId)
		getLastUsedOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
		}
		getLastUsedResp, err := iamClient.Request("GET", getLastUsedPath, &getLastUsedOpt)
		if err != nil {
			return diag.Errorf("error get IAM User last login time: %s", err)
		}
		lastUsed, err := utils.FlattenResponse(getLastUsedResp)
		if err != nil {
			return diag.FromErr(err)
		}
		if e := d.Set("last_used_at", utils.PathSearch("access_key_last_used.last_used_at", lastUsed, nil)); e != nil {
			return diag.FromErr(e)
		}
	}
	d.SetId(accessKeyId)
	return nil
}

func checkAccessExist(d *schema.ResourceData, arrays []interface{}) interface{} {
	accessKey := d.Id()
	var res interface{}
	for _, array := range arrays {
		if accessKey == utils.PathSearch("access_key_id", array, nil).(string) {
			res = array
			break
		}
	}
	return res
}

func resourceIdentityV5AccessKeyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	iamClient, err := cfg.IAMNoVersionClient(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	accessKeyId := d.Id()
	if accessKeyId == "" {
		return diag.Errorf("error getting IAM access key id")
	}

	if d.HasChanges("status") {
		updateAccessKeyHttpUrl := "v5/users/{user_id}/access-keys/{access_key_id}"
		updateAccessKeyPath := iamClient.Endpoint + updateAccessKeyHttpUrl
		updateAccessKeyPath = strings.ReplaceAll(updateAccessKeyPath, "{user_id}", d.Get("user_id").(string))
		updateAccessKeyPath = strings.ReplaceAll(updateAccessKeyPath, "{access_key_id}", accessKeyId)
		updateAccessKeyOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			JSONBody: map[string]interface{}{
				"status": d.Get("status").(string),
			},
		}
		_, err := iamClient.Request("PUT", updateAccessKeyPath, &updateAccessKeyOpt)
		if err != nil {
			return diag.Errorf("error updating IAM access key: %s", err)
		}
	}
	return resourceIdentityV5AccessKeyRead(ctx, d, meta)
}

func resourceIdentityV5AccessKeyDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	iamClient, err := cfg.IAMNoVersionClient(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}
	deleteAccessKeyHttpUrl := "v5/users/{user_id}/access-keys/{access_key_id}"
	deleteAccessKeyPath := iamClient.Endpoint + deleteAccessKeyHttpUrl
	deleteAccessKeyPath = strings.ReplaceAll(deleteAccessKeyPath, "{user_id}", d.Get("user_id").(string))
	deleteAccessKeyPath = strings.ReplaceAll(deleteAccessKeyPath, "{access_key_id}", d.Id())
	deleteAccessKeyOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	_, err = iamClient.Request("DELETE", deleteAccessKeyPath, &deleteAccessKeyOpt)
	if err != nil {
		return diag.Errorf("error deleting IAM access key: %s", err)
	}
	return nil
}
