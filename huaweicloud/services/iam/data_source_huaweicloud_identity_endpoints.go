package iam

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

const limit = 10

// @API IAM GET /v3/endpoints
func DataSourceIdentityEndpoints() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIdentityEndpointsRead,

		Schema: map[string]*schema.Schema{
			"endpoint_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"interface", "service_id"},
				Description:   "The terminal node ID to be queried",
			},
			"interface": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Terminal node plane",
			},
			"service_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID of service",
			},
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
						"enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"service_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"link": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceIdentityEndpointsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.IAMV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating Identity Endpoint client: %s", err)
	}
	getIdentityEndpointBasePath := client.Endpoint + "v3/endpoints"
	getIdentityEndpointOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	endpointId := d.Get("endpoint_id").(string)
	if endpointId != "" {
		getIdentityEndpointBasePath = getIdentityEndpointBasePath + "/" + endpointId
	}
	getIdentityEndpointBasePath += fmt.Sprintf("?limit=%v", limit)
	path := buildQueryParamPath(d, getIdentityEndpointBasePath)
	response, err := client.Request("GET", path, &getIdentityEndpointOpt)
	if err != nil {
		return diag.Errorf("error getting identity endpoints: %s", err)
	}
	respBody, err := utils.FlattenResponse(response)
	if err != nil {
		return diag.FromErr(err)
	}
	id, _ := uuid.GenerateUUID()
	d.SetId(id)
	if endpointId != "" {
		endpoints := flattenEndpointsList([]interface{}{utils.PathSearch("endpoint", respBody, make([]interface{}, 0))})
		if e := d.Set("endpoints", endpoints); e != nil {
			return diag.FromErr(e)
		}
	} else {
		endpoint := flattenEndpointsList(utils.PathSearch("endpoints", respBody, make([]interface{}, 0)).([]interface{}))
		if e := d.Set("endpoints", endpoint); e != nil {
			return diag.FromErr(e)
		}
	}
	return nil
}

func flattenEndpointsList(endPointBody []interface{}) []map[string]interface{} {
	endpoints := make([]map[string]interface{}, len(endPointBody))
	for i, endpoint := range endPointBody {
		endpoints[i] = map[string]interface{}{
			"service_id": utils.PathSearch("service_id", endpoint, ""),
			"region_id":  utils.PathSearch("region_id", endpoint, ""),
			"id":         utils.PathSearch("id", endpoint, ""),
			"interface":  utils.PathSearch("interface", endpoint, ""),
			"region":     utils.PathSearch("region", endpoint, ""),
			"url":        utils.PathSearch("url", endpoint, ""),
			"enabled":    utils.PathSearch("enabled", endpoint, ""),
			"link":       utils.PathSearch("self", utils.PathSearch("links", endpoint, ""), ""),
		}
	}
	return endpoints
}

func buildQueryParamPath(d *schema.ResourceData, path string) string {
	if interfaceParam, ok := d.GetOk("interface"); ok {
		path += fmt.Sprintf("&interface=%s", interfaceParam)
	}
	if serviceId, ok := d.GetOk("service_id"); ok {
		path += fmt.Sprintf("&service_id=%s", serviceId)
	}
	return path
}
