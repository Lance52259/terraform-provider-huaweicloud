// Generated by PMS #406
package dds

import (
	"context"
	"strings"

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

func DataSourceDdsRestoreCollections() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDdsRestoreCollectionsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the instance ID.`,
			},
			"db_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the database name.`,
			},
			"restore_time": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the restoration time point.`,
			},
			"collections": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Indicates the list of collections.`,
			},
		},
	}
}

type RestoreCollectionsDSWrapper struct {
	*schemas.ResourceDataWrapper
	Config *config.Config
}

func newRestoreCollectionsDSWrapper(d *schema.ResourceData, meta interface{}) *RestoreCollectionsDSWrapper {
	return &RestoreCollectionsDSWrapper{
		ResourceDataWrapper: schemas.NewSchemaWrapper(d),
		Config:              meta.(*config.Config),
	}
}

func dataSourceDdsRestoreCollectionsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	wrapper := newRestoreCollectionsDSWrapper(d, meta)
	lisResColRst, err := wrapper.ListRestoreCollections()
	if err != nil {
		return diag.FromErr(err)
	}

	err = wrapper.listRestoreCollectionsToSchema(lisResColRst)
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

// @API DDS GET /v3/{project_id}/instances/{instance_id}/restore-collection
func (w *RestoreCollectionsDSWrapper) ListRestoreCollections() (*gjson.Result, error) {
	client, err := w.NewClient(w.Config, "dds")
	if err != nil {
		return nil, err
	}

	uri := "/v3/{project_id}/instances/{instance_id}/restore-collection"
	uri = strings.ReplaceAll(uri, "{instance_id}", w.Get("instance_id").(string))
	params := map[string]any{
		"db_name":      w.Get("db_name"),
		"restore_time": w.Get("restore_time"),
	}
	params = utils.RemoveNil(params)
	return httphelper.New(client).
		Method("GET").
		URI(uri).
		Query(params).
		OffsetPager("collections", "offset", "limit", 0).
		Request().
		Result()
}

func (w *RestoreCollectionsDSWrapper) listRestoreCollectionsToSchema(body *gjson.Result) error {
	d := w.ResourceData
	mErr := multierror.Append(nil,
		d.Set("region", w.Config.GetRegion(w.ResourceData)),
		d.Set("collections", body.Get("collections").Value()),
	)
	return mErr.ErrorOrNil()
}