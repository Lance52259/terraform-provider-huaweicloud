package iam

import (
	"context"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	iamv3 "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/iam/v3"
	iamModel "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/iam/v3/model"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// DataSourceIdentityRegions
// @API IAM GET /v3/regions
// @API IAM GET /v3/regions/{region_id}
func DataSourceIdentityRegions() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceIdentityRegionsRead,

		Schema: map[string]*schema.Schema{
			"region_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The id of the region to be queried.",
			},

			"regions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Region Id",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Region Type",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Region Description",
						},
						"link": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Resource link",
						},
						"locales": {
							Type:        schema.TypeMap,
							Computed:    true,
							Description: `Region Name`,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
		},
	}
}

func DataSourceIdentityRegionsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	iamV3Client, err := cfg.HcIamV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}
	regionId := d.Get("region_id").(string)
	if regionId == "" {
		return listRegions(iamV3Client, d)
	} else {
		return showRegion(iamV3Client, regionId, d)
	}
}

func listRegions(iamV3Client *iamv3.IamClient, d *schema.ResourceData) diag.Diagnostics {
	request := &iamModel.KeystoneListRegionsRequest{}
	response, err := iamV3Client.KeystoneListRegions(request)
	if err != nil {
		return diag.Errorf("KeystoneListRegions error : %s", err)
	}
	id, _ := uuid.GenerateUUID()
	d.SetId(id)
	regions := make([]map[string]interface{}, len(*response.Regions))
	for i, region := range *response.Regions {
		regions[i] = flattenRegion(region)
	}
	if err = d.Set("regions", regions); err != nil {
		return diag.Errorf("error setting regions fields: %s", err)
	}
	return nil
}

func showRegion(iamV3Client *iamv3.IamClient, regionId string, d *schema.ResourceData) diag.Diagnostics {
	request := &iamModel.KeystoneShowRegionRequest{RegionId: regionId}
	response, err := iamV3Client.KeystoneShowRegion(request)
	if err != nil {
		return diag.Errorf("showRegion error : %s", err)
	}
	id, _ := uuid.GenerateUUID()
	d.SetId(id)
	regions := make([]map[string]interface{}, 1)
	regions[0] = flattenRegion(*response.Region)
	if err = d.Set("regions", regions); err != nil {
		return diag.Errorf("error setting regions fields: %s", err)
	}
	return nil
}

func flattenRegion(regionModel iamModel.Region) map[string]interface{} {
	region := make(map[string]interface{})
	region["id"] = regionModel.Id
	region["type"] = regionModel.Type
	region["description"] = regionModel.Description
	region["link"] = regionModel.Links.Self

	locales := make(map[string]string)
	locales["zh-cn"] = regionModel.Locales.ZhCn
	locales["en-us"] = regionModel.Locales.EnUs
	if regionModel.Locales.PtBr != nil {
		locales["pt-br"] = *regionModel.Locales.PtBr
	}
	if regionModel.Locales.EsUs != nil {
		locales["es-us"] = *regionModel.Locales.EsUs
	}
	if regionModel.Locales.EsEs != nil {
		locales["es-es"] = *regionModel.Locales.EsEs
	}
	region["locales"] = locales
	return region
}
