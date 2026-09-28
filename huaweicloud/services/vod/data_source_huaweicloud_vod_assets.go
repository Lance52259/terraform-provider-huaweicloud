package vod

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API VOD GET /v1.0/{project_id}/asset/list
func DataSourceAssets() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAssetsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the VOD assets. If omitted, the provider-level region will be used.`,
			},
			"asset_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				MaxItems:    10,
				Description: `Specifies the asset IDs used to query. A maximum of 10 IDs can be specified.`,
			},
			"status": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Specifies the asset statuses used to query.`,
			},
			"start_time": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the start time used to query, in the format yyyymmddhhmmss (UTC).`,
			},
			"end_time": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the end time used to query, in the format yyyymmddhhmmss (UTC).`,
			},
			"category_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: `Specifies the category ID used to query.`,
			},
			"tags": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the asset tags used to query. Multiple tags are separated by commas.`,
			},
			"query_string": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the fuzzy query string used to match asset title, description, or category name.`,
			},
			"media_types": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				MaxItems:    20,
				Description: `Specifies the media types used to query. A maximum of 20 types can be specified.`,
			},
			"assets": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `All VOD assets that match the filter parameters.`,
				Elem:        assetsSchema(),
			},
		},
	}
}

func assetsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"asset_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The asset ID.`,
			},
			"title": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The asset title.`,
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The asset description.`,
			},
			"duration": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The asset duration, in seconds.`,
			},
			"duration_ms": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The asset duration, in milliseconds.`,
			},
			"size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The asset size, in bytes.`,
			},
			"original_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The original playback URL of the asset.`,
			},
			"category": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The category name of the asset.`,
			},
			"covers": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The cover URLs of the asset.`,
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time of the asset, in the format yyyymmddhhmmss (UTC).`,
			},
			"asset_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The status of the asset.`,
			},
			"transcode_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The transcoding status of the asset.`,
			},
			"thumbnail_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The thumbnail status of the asset.`,
			},
			"review_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The review status of the asset.`,
			},
			"exec_desc": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The execution description of the asset tasks.`,
			},
			"media_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The media type of the asset.`,
			},
		},
	}
}

func buildAssetsQueryParams(d *schema.ResourceData, page, size int) string {
	params := url.Values{}
	params.Set("page", strconv.Itoa(page))
	params.Set("size", strconv.Itoa(size))

	for _, v := range d.Get("asset_ids").([]interface{}) {
		if id := fmt.Sprint(v); id != "" {
			params.Add("asset_id", id)
		}
	}
	for _, v := range d.Get("status").([]interface{}) {
		if status := fmt.Sprint(v); status != "" {
			params.Add("status", status)
		}
	}
	if v, ok := d.GetOk("start_time"); ok {
		params.Set("start_time", v.(string))
	}
	if v, ok := d.GetOk("end_time"); ok {
		params.Set("end_time", v.(string))
	}
	if v, ok := d.GetOk("category_id"); ok {
		params.Set("category_id", strconv.Itoa(v.(int)))
	}
	if v, ok := d.GetOk("tags"); ok {
		params.Set("tags", v.(string))
	}
	if v, ok := d.GetOk("query_string"); ok {
		params.Set("query_string", v.(string))
	}
	for _, v := range d.Get("media_types").([]interface{}) {
		if mediaType := fmt.Sprint(v); mediaType != "" {
			params.Add("media_type", mediaType)
		}
	}

	return "?" + params.Encode()
}

func flattenAssetCovers(covers []interface{}) []interface{} {
	if len(covers) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(covers))
	for _, cover := range covers {
		if coverURL := utils.PathSearch("cover_url", cover, "").(string); coverURL != "" {
			result = append(result, coverURL)
		}
	}
	return result
}

func flattenAssets(assets []interface{}) []map[string]interface{} {
	if len(assets) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(assets))
	for _, asset := range assets {
		result = append(result, map[string]interface{}{
			"asset_id":         utils.PathSearch("asset_id", asset, nil),
			"title":            utils.PathSearch("title", asset, nil),
			"description":      utils.PathSearch("description", asset, nil),
			"duration":         utils.PathSearch("duration", asset, nil),
			"duration_ms":      utils.PathSearch("duration_ms", asset, nil),
			"size":             utils.PathSearch("size", asset, nil),
			"original_url":     utils.PathSearch("original_url", asset, nil),
			"category":         utils.PathSearch("category", asset, nil),
			"covers":           flattenAssetCovers(utils.PathSearch("covers", asset, make([]interface{}, 0)).([]interface{})),
			"create_time":      utils.PathSearch("create_time", asset, nil),
			"asset_status":     utils.PathSearch("asset_status", asset, nil),
			"transcode_status": utils.PathSearch("transcode_status", asset, nil),
			"thumbnail_status": utils.PathSearch("thumbnail_status", asset, nil),
			"review_status":    utils.PathSearch("review_status", asset, nil),
			"exec_desc":        utils.PathSearch("exec_desc", asset, nil),
			"media_type":       utils.PathSearch("media_type", asset, nil),
		})
	}
	return result
}

func listAssets(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl  = "v1.0/{project_id}/asset/list"
		page     = 0
		pageSize = 100
		result   = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		requestPath := listPath + buildAssetsQueryParams(d, page, pageSize)
		resp, err := client.Request("GET", requestPath, &opt)
		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return nil, err
		}

		assets := utils.PathSearch("assets", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, assets...)

		total := int(utils.PathSearch("total", respBody, float64(0)).(float64))
		if len(assets) == 0 || len(result) >= total {
			break
		}
		page++
	}

	return result, nil
}

func dataSourceAssetsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("vod", region)
	if err != nil {
		return diag.Errorf("error creating VOD client: %s", err)
	}

	assets, err := listAssets(client, d)
	if err != nil {
		return diag.Errorf("error querying VOD assets: %s", err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("assets", flattenAssets(assets)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
