// Generated by PMS #523
package swr

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
)

func DataSourceSwrFeatureGates() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSwrFeatureGatesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
			},
			"enable_experience": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether the experience center is enabled.`,
			},
			"enable_hss_service": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether interconnection with HSS is enabled.`,
			},
			"enable_image_scan": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether image scanning is enabled.`,
			},
			"enable_sm3": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether SM algorithms are enabled.`,
			},
			"enable_image_sync": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether image synchronization is enabled.`,
			},
			"enable_cci_service": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether interconnection with CCI is enabled.`,
			},
			"enable_image_label": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether image tagging is enabled.`,
			},
			"enable_pipeline": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether pipeline is enabled.`,
			},
		},
	}
}

type FeatureGatesDSWrapper struct {
	*schemas.ResourceDataWrapper
	Config *config.Config
}

func newFeatureGatesDSWrapper(d *schema.ResourceData, meta interface{}) *FeatureGatesDSWrapper {
	return &FeatureGatesDSWrapper{
		ResourceDataWrapper: schemas.NewSchemaWrapper(d),
		Config:              meta.(*config.Config),
	}
}

func dataSourceSwrFeatureGatesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	wrapper := newFeatureGatesDSWrapper(d, meta)
	shoShaFeaGatRst, err := wrapper.ShowShareFeatureGates()
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	err = wrapper.showShareFeatureGatesToSchema(shoShaFeaGatRst)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// @API SWR GET /v2/manage/projects/{project_id}/feature-gates
func (w *FeatureGatesDSWrapper) ShowShareFeatureGates() (*gjson.Result, error) {
	client, err := w.NewClient(w.Config, "swr")
	if err != nil {
		return nil, err
	}

	uri := "/v2/manage/projects/{project_id}/feature-gates"
	return httphelper.New(client).
		Method("GET").
		URI(uri).
		Request().
		Result()
}

func (w *FeatureGatesDSWrapper) showShareFeatureGatesToSchema(body *gjson.Result) error {
	d := w.ResourceData
	mErr := multierror.Append(nil,
		d.Set("region", w.Config.GetRegion(w.ResourceData)),
		d.Set("enable_experience", body.Get("enable_experience").Value()),
		d.Set("enable_hss_service", body.Get("enable_hss_service").Value()),
		d.Set("enable_image_scan", body.Get("enable_image_scan").Value()),
		d.Set("enable_sm3", body.Get("enable_sm3").Value()),
		d.Set("enable_image_sync", body.Get("enable_image_sync").Value()),
		d.Set("enable_cci_service", body.Get("enable_cci_service").Value()),
		d.Set("enable_image_label", body.Get("enable_image_label").Value()),
		d.Set("enable_pipeline", body.Get("enable_pipeline").Value()),
	)
	return mErr.ErrorOrNil()
}