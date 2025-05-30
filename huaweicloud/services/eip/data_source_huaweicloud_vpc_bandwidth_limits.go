// Generated by PMS #572
package eip

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tidwall/gjson"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/httphelper"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/schemas"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func DataSourceBandwidthLimits() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceBandwidthLimitsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
			},
			"charge_mode": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the bandwidth charge mode.`,
			},
			"eip_bandwidth_limits": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `Indicates the bandwidth limit list.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the bandwidth type ID.`,
						},
						"charge_mode": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the bandwidth charging mode.`,
						},
						"min_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `Indicates the minimum size that can be purchased for this type of bandwidth.`,
						},
						"max_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `Indicates the maximum size that can be purchased for this type of bandwidth.`,
						},
						"ext_limit": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `Indicates the additional restriction information.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"min_ingress_size": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: `Indicates the minimum cloud access rate limit.`,
									},
									"max_ingress_size": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: `Indicates the maximum cloud access rate limit.`,
									},
									"ratio_95peak": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: `Indicates the 95 Minimum charging rate.`,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

type VpcBandwidthLimitsDSWrapper struct {
	*schemas.ResourceDataWrapper
	Config *config.Config
}

func newVpcBandwidthLimitsDSWrapper(d *schema.ResourceData, meta interface{}) *VpcBandwidthLimitsDSWrapper {
	return &VpcBandwidthLimitsDSWrapper{
		ResourceDataWrapper: schemas.NewSchemaWrapper(d),
		Config:              meta.(*config.Config),
	}
}

func dataSourceBandwidthLimitsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	wrapper := newVpcBandwidthLimitsDSWrapper(d, meta)
	lisBanLimRst, err := wrapper.ListBandwidthsLimit()
	if err != nil {
		return diag.FromErr(err)
	}

	err = wrapper.listBandwidthsLimitToSchema(lisBanLimRst)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)
	return nil
}

// @API EIP GET /v3/{project_id}/eip/eip-bandwidth-limits
func (w *VpcBandwidthLimitsDSWrapper) ListBandwidthsLimit() (*gjson.Result, error) {
	client, err := w.NewClient(w.Config, "vpcv3")
	if err != nil {
		return nil, err
	}

	uri := "/v3/{project_id}/eip/eip-bandwidth-limits"
	params := map[string]any{
		"charge_mode": w.Get("charge_mode"),
	}
	params = utils.RemoveNil(params)
	return httphelper.New(client).
		Method("GET").
		URI(uri).
		Query(params).
		OffsetPager("eip_bandwidth_limits", "offset", "offset", 0).
		Request().
		Result()
}

func (w *VpcBandwidthLimitsDSWrapper) listBandwidthsLimitToSchema(body *gjson.Result) error {
	d := w.ResourceData
	mErr := multierror.Append(nil,
		d.Set("region", w.Config.GetRegion(w.ResourceData)),
		d.Set("eip_bandwidth_limits", schemas.SliceToList(body.Get("eip_bandwidth_limits"),
			func(eipBanLimits gjson.Result) any {
				return map[string]any{
					"id":          eipBanLimits.Get("id").Value(),
					"charge_mode": eipBanLimits.Get("charge_mode").Value(),
					"min_size":    eipBanLimits.Get("min_size").Value(),
					"max_size":    eipBanLimits.Get("max_size").Value(),
					"ext_limit": schemas.SliceToList(eipBanLimits.Get("ext_limit"),
						func(extLimit gjson.Result) any {
							return map[string]any{
								"min_ingress_size": extLimit.Get("min_ingress_size").Value(),
								"max_ingress_size": extLimit.Get("max_ingress_size").Value(),
								"ratio_95peak":     extLimit.Get("ratio_95peak").Value(),
							}
						},
					),
				}
			},
		)),
	)
	return mErr.ErrorOrNil()
}
