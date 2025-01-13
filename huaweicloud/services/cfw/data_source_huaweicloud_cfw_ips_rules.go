// Generated by PMS #508
package cfw

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

func DataSourceCfwIpsRules() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCfwIpsRulesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
			},
			"object_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the protected object ID.`,
			},
			"ips_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the IPS rule ID.`,
			},
			"ips_name_like": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the IPS rule name.`,
			},
			"ips_status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the IPS rule status.`,
			},
			"is_updated_ips_rule_queried": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Specifies whether to check for new update rules.`,
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the enterprise project ID.`,
			},
			"records": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The IPS rule list.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ips_group": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The IPS rule group.`,
						},
						"ips_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The IPS rule ID.`,
						},
						"ips_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The IPS rule name.`,
						},
						"ips_rules_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The IPS rule type.`,
						},
						"affected_application": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The application affected by the rule.`,
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The creation time.`,
						},
						"ips_cve": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The CVE.`,
						},
						"default_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The default status of the IPS rule.`,
						},
						"ips_level": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The risk level.`,
						},
						"ips_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The current status of the IPS rule.`,
						},
					},
				},
			},
		},
	}
}

type IpsRulesDSWrapper struct {
	*schemas.ResourceDataWrapper
	Config *config.Config
}

func newIpsRulesDSWrapper(d *schema.ResourceData, meta interface{}) *IpsRulesDSWrapper {
	return &IpsRulesDSWrapper{
		ResourceDataWrapper: schemas.NewSchemaWrapper(d),
		Config:              meta.(*config.Config),
	}
}

func dataSourceCfwIpsRulesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	wrapper := newIpsRulesDSWrapper(d, meta)
	listIpsRules1Rst, err := wrapper.ListIpsRules1()
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	err = wrapper.listIpsRules1ToSchema(listIpsRules1Rst)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// @API CFW GET /v1/{project_id}/ips-rule
func (w *IpsRulesDSWrapper) ListIpsRules1() (*gjson.Result, error) {
	client, err := w.NewClient(w.Config, "cfw")
	if err != nil {
		return nil, err
	}

	uri := "/v1/{project_id}/ips-rule"
	params := map[string]any{
		"ips_id":                      w.Get("ips_id"),
		"ips_name_like":               w.Get("ips_name_like"),
		"ips_status":                  w.Get("ips_status"),
		"is_updated_ips_rule_queried": w.Get("is_updated_ips_rule_queried"),
		"object_id":                   w.Get("object_id"),
		"enterprise_project_id":       w.Get("enterprise_project_id"),
	}
	params = utils.RemoveNil(params)
	return httphelper.New(client).
		Method("GET").
		URI(uri).
		Query(params).
		OffsetPager("data.records", "offset", "limit", 1000).
		Request().
		Result()
}

func (w *IpsRulesDSWrapper) listIpsRules1ToSchema(body *gjson.Result) error {
	d := w.ResourceData
	mErr := multierror.Append(nil,
		d.Set("region", w.Config.GetRegion(w.ResourceData)),
		d.Set("records", schemas.SliceToList(body.Get("data.records"),
			func(records gjson.Result) any {
				return map[string]any{
					"ips_group":            records.Get("ips_group").Value(),
					"ips_id":               records.Get("ips_id").Value(),
					"ips_name":             records.Get("ips_name").Value(),
					"ips_rules_type":       records.Get("ips_rules_type").Value(),
					"affected_application": records.Get("affected_application").Value(),
					"create_time":          records.Get("create_time").Value(),
					"ips_cve":              records.Get("ips_cve").Value(),
					"default_status":       records.Get("default_status").Value(),
					"ips_level":            records.Get("ips_level").Value(),
					"ips_status":           records.Get("ips_status").Value(),
				}
			},
		)),
	)
	return mErr.ErrorOrNil()
}