package iam

import (
	"context"
	"github.com/chnsz/golangsdk"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
	"net/url"
)

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
	interfaceParam := d.Get("interface").(string)
	serviceId := d.Get("service_id").(string)
	u, err := url.Parse(getIdentityEndpointBasePath)
	if err != nil {
		return diag.Errorf("error update query path: %s", err)
	}
	query := u.Query()
	if serviceId != "" {
		query.Add("service_id", serviceId)
	}
	if interfaceParam != "" {
		query.Add("interface", interfaceParam)
	}
	u.RawQuery = query.Encode()
	response, err := client.Request("GET", u.String(), &getIdentityEndpointOpt)
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
		endpoints := flattEndpoints([]interface{}{utils.PathSearch("endpoint", respBody, nil)})
		if e := d.Set("endpoints", endpoints); e != nil {
			return diag.FromErr(e)
		}
	} else {
		endpoint := flattEndpoints(utils.PathSearch("endpoints", respBody, nil).([]interface{}))
		if e := d.Set("endpoints", endpoint); e != nil {
			return diag.FromErr(e)
		}
	}
	return nil
}

func flattEndpoints(endPointBody []interface{}) []map[string]interface{} {
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
			"link":       utils.PathSearch("self", utils.PathSearch("links", endpoint, nil), nil),
		}
	}
	return endpoints
}
