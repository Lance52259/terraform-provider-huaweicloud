// Generated by PMS #517
package gaussdb

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

func DataSourceGaussdbOpengaussProjectQuotas() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGaussdbOpengaussProjectQuotasRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the resource type used to filter quotas. Value options: **instance**.`,
			},
			"quotas": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `Indicates the instance quota of a tenant.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resources": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `Indicates the resource objects.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `Indicates the quota of a specified type.`,
									},
									"used": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: `Indicates the number of created resources.`,
									},
									"quota": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: `Indicates the maximum resource quota.`,
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

type OpengaussProjectQuotasDSWrapper struct {
	*schemas.ResourceDataWrapper
	Config *config.Config
}

func newOpengaussProjectQuotasDSWrapper(d *schema.ResourceData, meta interface{}) *OpengaussProjectQuotasDSWrapper {
	return &OpengaussProjectQuotasDSWrapper{
		ResourceDataWrapper: schemas.NewSchemaWrapper(d),
		Config:              meta.(*config.Config),
	}
}

func dataSourceGaussdbOpengaussProjectQuotasRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	wrapper := newOpengaussProjectQuotasDSWrapper(d, meta)
	showProjectQuotasRst, err := wrapper.ShowProjectQuotas()
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	err = wrapper.showProjectQuotasToSchema(showProjectQuotasRst)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// @API GaussDB GET /v3/{project_id}/project-quotas
func (w *OpengaussProjectQuotasDSWrapper) ShowProjectQuotas() (*gjson.Result, error) {
	client, err := w.NewClient(w.Config, "opengauss")
	if err != nil {
		return nil, err
	}

	uri := "/v3/{project_id}/project-quotas"
	params := map[string]any{
		"type": w.Get("type"),
	}
	params = utils.RemoveNil(params)
	return httphelper.New(client).
		Method("GET").
		URI(uri).
		Query(params).
		Request().
		Result()
}

func (w *OpengaussProjectQuotasDSWrapper) showProjectQuotasToSchema(body *gjson.Result) error {
	d := w.ResourceData
	mErr := multierror.Append(nil,
		d.Set("region", w.Config.GetRegion(w.ResourceData)),
		d.Set("quotas", schemas.ObjectToList(body.Get("quotas"),
			func(quotas gjson.Result) any {
				return map[string]any{
					"resources": schemas.SliceToList(quotas.Get("resources"),
						func(resources gjson.Result) any {
							return map[string]any{
								"type":  resources.Get("type").Value(),
								"used":  resources.Get("used").Value(),
								"quota": resources.Get("quota").Value(),
							}
						},
					),
				}
			},
		)),
	)
	return mErr.ErrorOrNil()
}
