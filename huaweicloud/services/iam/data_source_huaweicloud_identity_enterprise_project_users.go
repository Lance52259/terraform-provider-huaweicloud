package iam

import (
	"context"
	"strings"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API IAM GET /v3.0/OS-PERMISSION/enterprise-projects/{enterprise_project_id}/users
func DataSourceIdentityEnterpriseProjectUsers() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceIdentityEnterpriseProjectUsersRead,
		Schema: map[string]*schema.Schema{
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"users": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"policy_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"domain_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"lastest_policy_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func DataSourceIdentityEnterpriseProjectUsersRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	iamClient, err := cfg.IdentityV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}
	enterpriseProjectId := d.Get("enterprise_project_id").(string)
	path := iamClient.Endpoint + "v3.0/OS-PERMISSION/enterprise-projects/{enterprise_project_id}/users"
	path = strings.ReplaceAll(path, "{enterprise_project_id}", enterpriseProjectId)
	identityOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	response, err := iamClient.Request("GET", path, &identityOpt)
	if err != nil {
		return diag.Errorf("error getting identity enterprise-projects users: %s", err)
	}
	respBody, err := utils.FlattenResponse(response)
	if err != nil {
		return diag.FromErr(err)
	}
	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(id)
	usersList := flattenEnterpriseProjectUsersList(utils.PathSearch("users", respBody, make([]interface{}, 0)).([]interface{}))
	if e := d.Set("users", usersList); e != nil {
		return diag.FromErr(e)
	}
	return nil
}

func flattenEnterpriseProjectUsersList(list []interface{}) []map[string]interface{} {
	res := make([]map[string]interface{}, len(list))
	for i, user := range list {
		res[i] = map[string]interface{}{
			"policy_num":          utils.PathSearch("policy_num", user, nil),
			"description":         utils.PathSearch("description", user, ""),
			"domain_id":           utils.PathSearch("domain_id", user, ""),
			"id":                  utils.PathSearch("id", user, ""),
			"name":                utils.PathSearch("name", user, ""),
			"enabled":             utils.PathSearch("enabled", user, nil),
			"lastest_policy_time": utils.PathSearch("lastest_policy_time", user, nil),
		}
	}
	return res
}
