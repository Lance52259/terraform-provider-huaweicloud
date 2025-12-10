package cceautopilot

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

func DataSourceCceAutopilotChart() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCceAutopilotChartRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
			},
		},
	}
}

type ChartDSWrapper struct {
	*schemas.ResourceDataWrapper
	Config *config.Config
}

func newChartDSWrapper(d *schema.ResourceData, meta interface{}) *ChartDSWrapper {
	return &ChartDSWrapper{
		ResourceDataWrapper: schemas.NewSchemaWrapper(d),
		Config:              meta.(*config.Config),
	}
}

func dataSourceCceAutopilotChartRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	wrapper := newChartDSWrapper(d, meta)
	lisAutChaRst, err := wrapper.ListAutopilotCharts()
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	err = wrapper.listAutopilotChartsToSchema(lisAutChaRst)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// @API CCE GET /autopilot/v2/charts
func (w *ChartDSWrapper) ListAutopilotCharts() (*gjson.Result, error) {
	client, err := w.NewClient(w.Config, "cce")
	if err != nil {
		return nil, err
	}

	uri := "/autopilot/v2/charts"
	return httphelper.New(client).
		Method("GET").
		URI(uri).
		Request().
		Result()
}

func (w *ChartDSWrapper) listAutopilotChartsToSchema(body *gjson.Result) error {
	d := w.ResourceData
	mErr := multierror.Append(nil,
		d.Set("region", w.Config.GetRegion(w.ResourceData)),
	)
	return mErr.ErrorOrNil()
}
