package iam

import (
	"context"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	iamModel "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/iam/v3/model"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// DataSourceIdentityDomainQuota
// @API IAM GET /v3.0/OS-QUOTA/domains/{domain_id}
func DataSourceIdentityDomainQuota() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceIdentityDomainQuotaRead,

		Schema: map[string]*schema.Schema{
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{"user", "group", "idp", "agency", "policy",
					"assigment_group_mp", "assigment_agency_mp", "assigment_group_ep", "assigment_user_ep", "mapping"},
					false),
				Description: "Query the type of quota. If this parameter is not filled in, all types of quotas will " +
					"be returned. The value range is as follows: user, group, idp, agency, policy, " +
					"assigment_group_mp, assigment_agency_mp, assigment_group_ep, assigment_user_ep, mapping",
			},

			"resources": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of resource",
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

func DataSourceIdentityDomainQuotaRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	iamV3Client, err := cfg.HcIamV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}
	requestType := convertStringToType(d.Get("type").(string))
	request := &iamModel.ShowDomainQuotaRequest{DomainId: cfg.DomainID, Type: &requestType}
	response, err := iamV3Client.ShowDomainQuota(request)
	if err != nil {
		return diag.Errorf("ShowDomainQuota error : %s", err)
	}
	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(id)
	result := flattenDomainQuotaResources(response.Quotas.Resources)
	return diag.FromErr(d.Set("resources", result))
}

func convertStringToType(typeStr string) iamModel.ShowDomainQuotaRequestType {
	switch typeStr {
	case "user":
		return iamModel.GetShowDomainQuotaRequestTypeEnum().USER
	case "group":
		return iamModel.GetShowDomainQuotaRequestTypeEnum().GROUP
	case "idp":
		return iamModel.GetShowDomainQuotaRequestTypeEnum().IDP
	case "agency":
		return iamModel.GetShowDomainQuotaRequestTypeEnum().AGENCY
	case "policy":
		return iamModel.GetShowDomainQuotaRequestTypeEnum().POLICY
	case "assigment_group_mp":
		return iamModel.GetShowDomainQuotaRequestTypeEnum().ASSIGMENT_GROUP_MP
	case "assigment_agency_mp":
		return iamModel.GetShowDomainQuotaRequestTypeEnum().ASSIGMENT_AGENCY_MP
	case "assigment_group_ep":
		return iamModel.GetShowDomainQuotaRequestTypeEnum().ASSIGMENT_GROUP_EP
	case "assigment_user_ep":
		return iamModel.GetShowDomainQuotaRequestTypeEnum().ASSIGMENT_USER_EP
	case "mapping":
		return iamModel.GetShowDomainQuotaRequestTypeEnum().MAPPING
	default:
		return iamModel.ShowDomainQuotaRequestType{}
	}
}

func flattenDomainQuotaResources(resources *[]iamModel.Resources) []map[string]interface{} {
	if len(*resources) < 1 {
		return nil
	}
	result := make([]map[string]interface{}, len(*resources))
	for i, resource := range *resources {
		result[i] = map[string]interface{}{
			"type":  resource.Type,
			"max":   resource.Max,
			"min":   resource.Min,
			"quota": resource.Quota,
			"used":  resource.Used,
		}
	}
	return result
}
