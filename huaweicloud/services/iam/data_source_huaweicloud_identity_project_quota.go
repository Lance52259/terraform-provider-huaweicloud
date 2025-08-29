package iam

import (
	"context"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	iamModel "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/iam/v3/model"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// DataSourceIdentityProjectQuota
// @API IAM GET /v3.0/OS-QUOTA/projects/{project_id}
func DataSourceIdentityProjectQuota() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceIdentityProjectQuotaRead,

		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
			},

			"resources": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of resource, must be `project`",
						},
						"max": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Maximum Quota",
						},
						"min": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Minimum Quota",
						},
						"quota": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Current quota",
						},
						"used": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Used quota",
						},
					},
				},
			},
		},
	}
}

func DataSourceIdentityProjectQuotaRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	iamV3Client, err := cfg.HcIamV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}
	request := &iamModel.ShowProjectQuotaRequest{ProjectId: d.Get("project_id").(string)}
	response, err := iamV3Client.ShowProjectQuota(request)
	if err != nil {
		return diag.Errorf("ShowProjectQuota error : %s", err)
	}
	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(id)
	result := flattenProjectQuotaResources(response.Quotas.Resources)
	return diag.FromErr(d.Set("resources", result))
}

func flattenProjectQuotaResources(resources *[]iamModel.Resources) []map[string]interface{} {
	if len(*resources) < 1 {
		return nil
	}
	result := make([]map[string]interface{}, len(*resources))
	for i, resource := range *resources {
		result[i] = map[string]interface{}{
			"type":  "project",
			"max":   resource.Max,
			"min":   resource.Min,
			"quota": resource.Quota,
			"used":  resource.Used,
		}
	}
	return result
}
