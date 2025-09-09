package iam

import (
	"context"
	"github.com/chnsz/golangsdk"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
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

func dataSourceIdentityAccessKeyRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	iamClient, err := cfg.IAMV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(id)
	getAccessKeyPath := iamClient.Endpoint + "v3.0/OS-CREDENTIAL/credentials"
	accessKey := d.Get("access_key").(string)
	if accessKey != "" {
		getAccessKeyPath = getAccessKeyPath + "/" + accessKey
		return ListAccessKey(getAccessKeyPath, iamClient, d)
	}
	userId := d.Get("user_id").(string)
	return ShowAccessKey(userId, getAccessKeyPath, iamClient, d)
}

func ListAccessKey(path string, iamClient *golangsdk.ServiceClient, d *schema.ResourceData) diag.Diagnostics {
	options := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	response, err := iamClient.Request("GET", path, &options)
	if err != nil {
		return diag.Errorf("iamV3Client response error : %s", err)
	}
	respBody, err := utils.FlattenResponse(response)
	if err != nil {
		return diag.FromErr(err)
	}
	id, _ := uuid.GenerateUUID()
	d.SetId(id)
	credentialsModel := utils.PathSearch("credential", respBody, nil)
	res := make(map[string]interface{})
	res["user_id"] = utils.PathSearch("user_id", credentialsModel, nil)
	res["access"] = utils.PathSearch("access", credentialsModel, nil)
	res["status"] = utils.PathSearch("status", credentialsModel, nil)
	res["create_time"] = utils.PathSearch("create_time", credentialsModel, nil)
	res["description"] = utils.PathSearch("description", credentialsModel, nil)
	res["last_use_time"] = utils.PathSearch("last_use_time", credentialsModel, nil)
	credential := append(make([]interface{}, 0, 1), res)
	if err = d.Set("credentials", credential); err != nil {
		return diag.Errorf("error setting credentials fields: %s", err)
	}
	return nil
}

func ShowAccessKey(userId string, path string, iamClient *golangsdk.ServiceClient, d *schema.ResourceData) diag.Diagnostics {
	options := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	if userId != "" {
		path = path + "?user_id=" + userId
	}
	response, err := iamClient.Request("GET", path, &options)
	if err != nil {
		return diag.Errorf("iamV3Client response error : %s", err)
	}
	respBody, err := utils.FlattenResponse(response)
	if err != nil {
		return diag.FromErr(err)
	}
	id, _ := uuid.GenerateUUID()
	d.SetId(id)
	credentialsModel := utils.PathSearch("credentials", respBody, nil).([]interface{})
	credentials := make([]map[string]interface{}, len(credentialsModel))
	for i, credential := range credentialsModel {
		credentials[i] = map[string]interface{}{
			"user_id":     utils.PathSearch("user_id", credential, nil),
			"access":      utils.PathSearch("access", credential, nil),
			"status":      utils.PathSearch("status", credential, nil),
			"create_time": utils.PathSearch("create_time", credential, nil),
			"description": utils.PathSearch("description", credential, nil),
		}
	}
	if err = d.Set("credentials", credentials); err != nil {
		return diag.Errorf("error setting credentials fields: %s", err)
	}
	return nil
}
