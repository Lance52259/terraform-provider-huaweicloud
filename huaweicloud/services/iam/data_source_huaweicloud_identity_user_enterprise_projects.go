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

// @API IAM GET /v3.0/OS-PERMISSION/users/{user_id}/enterprise-projects
func DataSourceIdentityUserEnterpriseProjects() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceIdentityUserEnterpriseProjectsRead,
		Schema: map[string]*schema.Schema{
			"user_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enterprise_projects": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"project_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func DataSourceIdentityUserEnterpriseProjectsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	iamClient, err := cfg.IdentityV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}
	userId := d.Get("user_id").(string)
	path := iamClient.Endpoint + "v3.0/OS-PERMISSION/users/{user_id}/enterprise-projects"
	path = strings.ReplaceAll(path, "{user_id}", userId)
	identityOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	response, err := iamClient.Request("GET", path, &identityOpt)
	if err != nil {
		return diag.Errorf("error getting identity enterprise-projects: %s", err)
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
	enterpriseProjectList := flattenUserEnterpriseProjectList(utils.PathSearch("\"enterprise-projects\"", respBody, make([]interface{}, 0)).([]interface{}))
	if e := d.Set("enterprise_projects", enterpriseProjectList); e != nil {
		return diag.FromErr(e)
	}
	return nil
}

func flattenUserEnterpriseProjectList(list []interface{}) []map[string]interface{} {
	res := make([]map[string]interface{}, len(list))
	for i, project := range list {
		res[i] = map[string]interface{}{
			"project_id": utils.PathSearch("projectId", project, ""),
		}
	}
	return res
}
