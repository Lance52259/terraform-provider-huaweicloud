package iam

import (
	"context"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	iamModel "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/iam/v3/model"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @IAM GET /v3.0/OS-CREDENTIAL/credentials

func DataSourceIdentityAccessKey() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIdentityAccessKeyRead,
		Schema: map[string]*schema.Schema{
			"user_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"access_key": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"credentials": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"access": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"last_use_time": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "field is only present when you are certified by access_key",
						},
					},
				},
			},
		},
	}
}

func dataSourceIdentityAccessKeyRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cfg := m.(*config.Config)
	iamClient, err := cfg.HcIamV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(id)

	userId := d.Get("user_id").(string)
	if userId == "" {
		userId = cfg.UserID
	}
	access_key := d.Get("access_key").(string)
	if access_key != "" {
		request := &iamModel.ShowPermanentAccessKeyRequest{
			AccessKey: access_key,
		}
		response, err := iamClient.ShowPermanentAccessKey(request)
		if err != nil {
			return diag.Errorf("iamV3Client response error : %s", err)
		}
		return diag.FromErr(d.Set("credentials", flattenAccessKey(response)))
	}
	request := &iamModel.ListPermanentAccessKeysRequest{
		UserId: &userId,
	}
	response, err := iamClient.ListPermanentAccessKeys(request)
	if err != nil {
		return diag.Errorf("iamV3Client response error : %s", err)
	}
	return diag.FromErr(d.Set("credentials", flattenKey(response)))
}

func flattenKey(response *iamModel.ListPermanentAccessKeysResponse) []map[string]interface{} {
	keyRaw := utils.PathSearch("credentials", response, nil)
	if keyRaw == nil {
		return make([]map[string]interface{}, 0)
	}
	key := keyRaw.(*[]iamModel.Credentials)

	res := make([]map[string]interface{}, len(*key))
	for i, v := range *key {
		res[i] = map[string]interface{}{
			"user_id":     v.UserId,
			"access":      v.Access,
			"status":      v.Status,
			"create_time": v.CreateTime,
			"description": v.Description,
		}
	}
	return res
}

func flattenAccessKey(response *iamModel.ShowPermanentAccessKeyResponse) []map[string]interface{} {
	keyRaw := utils.PathSearch("credential", response, nil)
	if keyRaw == nil {
		return make([]map[string]interface{}, 0)
	}
	key := keyRaw.(*iamModel.ShowCredential)

	res := make([]map[string]interface{}, 1)
	res[0] = map[string]interface{}{
		"user_id":       key.UserId,
		"access":        key.Access,
		"status":        key.Status,
		"create_time":   key.CreateTime,
		"description":   key.Description,
		"last_use_time": key.LastUseTime,
	}
	return res
}
