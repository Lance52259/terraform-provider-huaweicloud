package iam

import (
	"context"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	iamModel "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/iam/v3/model"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @IAM GET /v3.0/OS-CREDENTIAL/credentials
func DataSourceIdentityCredentials() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIdentityCredentialsRead,
		Schema: map[string]*schema.Schema{
			"user_id": {
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
					},
				},
			},
		},
	}
}

func dataSourceIdentityCredentialsRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cfg := m.(*config.Config)
	iamV3Client, err := cfg.HcIamV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	userId := d.Get("user_id").(string)
	request := &iamModel.ListPermanentAccessKeysRequest{
		UserId: &userId,
	}
	if request == nil {
		return diag.Errorf("request error is nil")
	}
	response, err := iamV3Client.ListPermanentAccessKeys(request)
	if err != nil {
		return diag.Errorf("iamV3Client response error : %s", err)
	}
	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(id)
	mErr := multierror.Append(
		nil,
		d.Set("credentials", flattenKey(response)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenKey(response *iamModel.ListPermanentAccessKeysResponse) []map[string]interface{} {
	keyRaw := utils.PathSearch("credentials", response, nil)
	if keyRaw == nil {
		return nil
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
