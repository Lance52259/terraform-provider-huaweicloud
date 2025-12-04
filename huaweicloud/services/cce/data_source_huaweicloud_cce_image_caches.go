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

// @API CCE GET /v5/imagecaches
func DataSourceCCEImageCaches() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCCEImageCachesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

type CCEImageCachesDSWrapper struct {
	*schemas.ResourceDataWrapper
	Config *config.Config
}

func newCCEImageCachesDSWrapper(d *schema.ResourceData, meta interface{}) *CCEImageCachesDSWrapper {
	return &CCEImageCachesDSWrapper{
		ResourceDataWrapper: schemas.NewSchemaWrapper(d),
		Config:              meta.(*config.Config),
	}
}

func dataSourceCCEImageCachesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	wrapper := newCCEImageCachesDSWrapper(d, meta)
	imageCachesRst, err := wrapper.getImageCaches()
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	err = wrapper.imageCachesToSchema(imageCachesRst)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// @API CCE GET /v5/imagecaches
func (w *CCEImageCachesDSWrapper) getImageCaches() (*gjson.Result, error) {
	client, err := w.NewClient(w.Config, "cce")
	if err != nil {
		return nil, err
	}

	uri := "/v5/imagecaches"
	return httphelper.New(client).
		Method("GET").
		URI(uri).
		Request().
		Result()
}

func (w *CCEImageCachesDSWrapper) imageCachesToSchema(body *gjson.Result) error {
	d := w.ResourceData
	mErr := multierror.Append(nil,
		d.Set("region", w.Config.GetRegion(w.ResourceData)),
	)
	return mErr.ErrorOrNil()
}
