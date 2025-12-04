package cce

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

// @API CCE GET /api/v3/access-policies
func DataSourceCCEAccessPolicies() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCCEAccessPoliciesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"kind": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"api_version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"access_policy_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

type CCEAccessPoliciesDSWrapper struct {
	*schemas.ResourceDataWrapper
	Config *config.Config
}

func newCCEAccessPoliciesDSWrapper(d *schema.ResourceData, meta interface{}) *CCEAccessPoliciesDSWrapper {
	return &CCEAccessPoliciesDSWrapper{
		ResourceDataWrapper: schemas.NewSchemaWrapper(d),
		Config:              meta.(*config.Config),
	}
}

func dataSourceCCEAccessPoliciesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	wrapper := newCCEAccessPoliciesDSWrapper(d, meta)
	accessPoliciesRst, err := wrapper.getAccessPolicies()
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	err = wrapper.accessPoliciesToSchema(accessPoliciesRst)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// @API CCE GET /api/v3/access-policies
func (w *CCEAccessPoliciesDSWrapper) getAccessPolicies() (*gjson.Result, error) {
	client, err := w.NewClient(w.Config, "cce")
	if err != nil {
		return nil, err
	}

	uri := "/api/v3/access-policies"
	return httphelper.New(client).
		Method("GET").
		URI(uri).
		Request().
		Result()
}

func (w *CCEAccessPoliciesDSWrapper) accessPoliciesToSchema(body *gjson.Result) error {
	d := w.ResourceData
	mErr := multierror.Append(nil,
		d.Set("region", w.Config.GetRegion(w.ResourceData)),
	)
	return mErr.ErrorOrNil()
}
