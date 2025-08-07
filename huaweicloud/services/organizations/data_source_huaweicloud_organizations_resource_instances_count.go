package organizations

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API Organizations POST /v1/organizations/{resource_type}/resource-instances/count
func DataSourceOrganizationsResourceInstancesCount() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceOrganizationsResourceInstancesCountRead,
		Schema: map[string]*schema.Schema{
			"resource_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"without_any_tag": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"tags": {
				Type:     schema.TypeList,
				Elem:     organizationsResourceInstancesCountTagSchema(),
				Optional: true,
			},
			"matches": {
				Type:     schema.TypeList,
				Elem:     organizationsResourceInstancesCountMatchSchema(),
				Optional: true,
			},
			"total_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func organizationsResourceInstancesCountTagSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"key": {
				Type:     schema.TypeString,
				Required: true,
			},
			"values": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
			},
		},
	}
	return &sc
}

func organizationsResourceInstancesCountMatchSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"key": {
				Type:     schema.TypeString,
				Required: true,
			},
			"value": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
	return &sc
}

func dataSourceOrganizationsResourceInstancesCountRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v1/organizations/{resource_type}/resource-instances/count"
		product = "organizations"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating Organizations client: %s", err)
	}

	listBasePath := client.Endpoint + httpUrl
	listBasePath = strings.ReplaceAll(listBasePath, "{resource_type}", d.Get("resource_type").(string))
	listsOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	listsOpt.JSONBody = utils.RemoveNil(buildGetResourceInstancesCountBodyParams(d))

	listPath := listBasePath
	listResp, err := client.Request("POST", listPath, &listsOpt)
	if err != nil {
		return diag.Errorf("error retrieving Organizations resource instances: %s", err)
	}

	listRespBody, err := utils.FlattenResponse(listResp)
	if err != nil {
		return diag.FromErr(err)
	}

	totalCount := utils.PathSearch("total_count", listRespBody, float64(0)).(float64)

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr := multierror.Append(nil,
		d.Set("total_count", totalCount),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildGetResourceInstancesCountBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"without_any_tag": utils.ValueIgnoreEmpty(d.Get("without_any_tag")),
		"tags":            buildGetResourceInstanceCountTagsBodyParams(d),
		"matches":         buildGetResourceInstanceCountMatchesBodyParams(d),
	}
	return bodyParams
}

func buildGetResourceInstanceCountTagsBodyParams(d *schema.ResourceData) []map[string]interface{} {
	tagsRawParams := d.Get("tags").([]interface{})
	if len(tagsRawParams) == 0 {
		return nil
	}
	res := make([]map[string]interface{}, 0, len(tagsRawParams))
	for _, v := range tagsRawParams {
		raw := v.(map[string]interface{})
		res = append(res, map[string]interface{}{
			"key":    raw["key"],
			"values": raw["values"],
		})
	}
	return res
}

func buildGetResourceInstanceCountMatchesBodyParams(d *schema.ResourceData) []map[string]interface{} {
	matchesRawParams := d.Get("matches").([]interface{})
	if len(matchesRawParams) == 0 {
		return nil
	}
	res := make([]map[string]interface{}, 0, len(matchesRawParams))
	for _, v := range matchesRawParams {
		raw := v.(map[string]interface{})
		res = append(res, map[string]interface{}{
			"key":   raw["key"],
			"value": raw["value"],
		})
	}
	return res
}
