package iam

import (
	"context"
	"github.com/chnsz/golangsdk"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
	"strings"
)

// @API IAM GET /v3/auth/catalog
func DataSourceIdentityCatalog() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIdentityCatalogRead,

		Schema: map[string]*schema.Schema{
			"project_token": {
				Type:     schema.TypeString,
				Required: true,
			},
			"catalog": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"endpoints": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"interface": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"region": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"region_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"url": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"link": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceIdentityCatalogRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	iamClient := common.NewCustomClient(true, "https://iam.{region_id}.myhuaweicloud.com")
	getIdentityCatalogBasePath := iamClient.ResourceBase + "v3/auth/catalog"
	getIdentityCatalogBasePath = strings.ReplaceAll(getIdentityCatalogBasePath, "{region_id}", cfg.GetRegion(d))
	getIdentityCatalogOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"X-Auth-Token": d.Get("project_token").(string),
		},
	}
	response, err := iamClient.Request("GET", getIdentityCatalogBasePath, &getIdentityCatalogOpt)
	if err != nil {
		return diag.Errorf("error getting identity catalog: %s", err)
	}
	respBody, err := utils.FlattenResponse(response)
	if err != nil {
		return diag.FromErr(err)
	}
	id, _ := uuid.GenerateUUID()
	d.SetId(id)
	catalog := flattenCatalogList(utils.PathSearch("catalog", respBody, nil).([]interface{}))
	if e := d.Set("catalog", catalog); e != nil {
		return diag.FromErr(e)
	}
	if e := d.Set("link", utils.PathSearch("self", utils.PathSearch("links", respBody, nil), nil)); e != nil {
		return diag.FromErr(e)
	}
	return nil
}

func flattenCatalogList(catalogBody []interface{}) []map[string]interface{} {
	catalog := make([]map[string]interface{}, len(catalogBody))
	for i, c := range catalogBody {
		catalog[i] = map[string]interface{}{
			"endpoints": flattenEndpointsList(utils.PathSearch("endpoints", c, nil).([]interface{})),
			"name":      utils.PathSearch("name", c, nil),
			"id":        utils.PathSearch("id", c, nil),
			"type":      utils.PathSearch("type", c, nil),
		}
	}
	return catalog
}

func flattenEndpointsList(endpointsBody []interface{}) []map[string]interface{} {
	endpoints := make([]map[string]interface{}, len(endpointsBody))
	for i, endpoint := range endpointsBody {
		endpoints[i] = map[string]interface{}{
			"region":    utils.PathSearch("region", endpoint, nil),
			"interface": utils.PathSearch("interface", endpoint, nil),
			"id":        utils.PathSearch("id", endpoint, nil),
			"region_id": utils.PathSearch("region_id", endpoint, nil),
			"url":       utils.PathSearch("url", endpoint, nil),
		}
	}
	return endpoints
}
