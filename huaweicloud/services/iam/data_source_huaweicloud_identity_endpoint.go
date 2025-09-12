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

// @API IAM GET /v3/endpoints
func DataSourceIdentityEndpoint() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIdentityEndpointRead,

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
			"token": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "IAM user authentication token ",
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
						"links": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"links": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}
func dataSourceIdentityEndpointRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	var (
		getIdentityEndpointHttpUrl = "v3/endpoints"
		product                    = "iam"
	)
	client, err := cfg.NewServiceClient(product, cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating Identity Endpoint client: %s", err)
	}
	getIdentityEndpointBasePath := client.Endpoint + getIdentityEndpointHttpUrl
	getIdentityEndpointOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"X-Auth-Token": d.Get("token").(string),
			"Content-Type": "application/json;charset=utf8",
		},
	}
	endpointId := d.Get("endpoint_id").(string)
	if endpointId != "" {
		getIdentityEndpointBasePath = getIdentityEndpointBasePath + "/" + endpointId
	}
	interfaceParam := d.Get("interface").(string)
	serviceId := d.Get("service_id").(string)
	if interfaceParam != "" && serviceId != "" {
		getIdentityEndpointBasePath = getIdentityEndpointBasePath + "?service_id=" + serviceId + "&interface=" + interfaceParam
	} else if interfaceParam != "" {
		getIdentityEndpointBasePath = getIdentityEndpointBasePath + "?interface=" + interfaceParam
	} else {
		getIdentityEndpointBasePath = getIdentityEndpointBasePath + "?service_id=" + serviceId
	}

	response, err := client.Request("GET", getIdentityEndpointBasePath, &getIdentityEndpointOpt)
	if err != nil {
		return diag.Errorf("error get identity token: %s", err)
	}
	respBody, err := utils.FlattenResponse(response)
	if err != nil {
		return diag.FromErr(err)
	}
	id, _ := uuid.GenerateUUID()
	d.SetId(id)
	if endpointId != "" {
		endpoint := flattenEndpoint(append(make([]interface{}, 0, 1), utils.PathSearch("endpoint", respBody, nil)))
		if e := d.Set("endpoints", endpoint); e != nil {
			return diag.FromErr(e)
		}
	} else {
		endpoint := flattenEndpoint(utils.PathSearch("endpoints", respBody, nil).([]interface{}))
		if e := d.Set("endpoints", endpoint); e != nil {
			return diag.FromErr(e)
		}
		links := flattenLinks(utils.PathSearch("links", respBody, nil))
		if e := d.Set("links", links); e != nil {
			return diag.FromErr(e)
		}
	}
	return nil
}

func flattenEndpoint(endPointBody []interface{}) []map[string]interface{} {

	endpoints := make([]map[string]interface{}, len(endPointBody))
	for i, endpoint := range endPointBody {
		endpoints[i] = map[string]interface{}{
			"service_id": utils.PathSearch("service_id", endpoint, nil),
			"region_id":  utils.PathSearch("region_id", endpoint, nil),
			"id":         utils.PathSearch("id", endpoint, nil),
			"interface":  utils.PathSearch("interface", endpoint, nil),
			"region":     utils.PathSearch("region", endpoint, nil),
			"url":        utils.PathSearch("url", endpoint, nil),
			"enabled":    utils.PathSearch("enabled", endpoint, nil),
			"links":      flattenLinks(utils.PathSearch("links", endpoint, nil)),
		}
	}
	return endpoints
}

func flattenLinks(linkBody interface{}) map[string]interface{} {
	if linkBody == nil {
		return nil
	}
	res := map[string]interface{}{
		"next":     utils.PathSearch("next", linkBody, nil),
		"previous": utils.PathSearch("previous", linkBody, nil),
		"self":     utils.PathSearch("self", linkBody, nil).(string),
	}
	return res
}
