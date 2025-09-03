package iam

import (
	"context"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	iamModel "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/iam/v3/model"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// @IAM GET /v3.0/OS-CREDENTIAL/credentials
// @IAM GET /v3.0/OS-CREDENTIAL/credentials/{access_key}
func DataSourceIdentityAccessKey() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIdentityAccessKeyRead,
		Schema: map[string]*schema.Schema{
			"user_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"access_key"},
			},
			"access_key": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"user_id"},
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

	access_key := d.Get("access_key").(string)
	if access_key != "" {
		request := &iamModel.ShowPermanentAccessKeyRequest{
			AccessKey: access_key,
		}
		response, err := iamClient.ShowPermanentAccessKey(request)
		if err != nil {
			return diag.Errorf("iamV3Client response error : %s", err)
		}
		return diag.FromErr(d.Set("credentials", flattenShowAccessKeyResponse(response)))
	}
	userId := d.Get("user_id").(string)
	if userId == "" {
		userId = cfg.UserID
	}
	request := &iamModel.ListPermanentAccessKeysRequest{
		UserId: &userId,
	}
	response, err := iamClient.ListPermanentAccessKeys(request)
	if err != nil {
		return diag.Errorf("iamV3Client response error : %s", err)
	}
	return diag.FromErr(d.Set("credentials", flattenListAccessKeyResponse(response)))
}

func flattenListAccessKeyResponse(response *iamModel.ListPermanentAccessKeysResponse) []map[string]interface{} {
	res := make([]map[string]interface{}, len(*response.Credentials))
	for i, credential := range *response.Credentials {
		res[i] = map[string]interface{}{
			"user_id":     credential.UserId,
			"access":      credential.Access,
			"status":      credential.Status,
			"create_time": credential.CreateTime,
			"description": credential.Description,
		}
	}
	return res
}

func flattenShowAccessKeyResponse(response *iamModel.ShowPermanentAccessKeyResponse) []map[string]interface{} {
	res := make([]map[string]interface{}, 1)
	res[0] = map[string]interface{}{
		"user_id":       response.Credential.UserId,
		"access":        response.Credential.Access,
		"status":        response.Credential.Status,
		"create_time":   response.Credential.CreateTime,
		"description":   response.Credential.Description,
		"last_use_time": response.Credential.LastUseTime,
	}
	return res
}
