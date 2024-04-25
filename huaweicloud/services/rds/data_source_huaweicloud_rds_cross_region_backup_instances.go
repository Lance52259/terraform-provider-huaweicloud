// Generated by PMS #126
package rds

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tidwall/gjson"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/filters"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/httphelper"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/schemas"
)

func DataSourceRdsCrossRegionBackupInstances() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRdsCrossRegionBackupInstancesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the ID of the instance.`,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the name of the instance.`,
			},
			"source_region": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the source backup region.`,
			},
			"source_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the project ID of the source backup region.`,
			},
			"destination_region": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the region where the cross-region backup is located.`,
			},
			"destination_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the project ID of the target backup region.`,
			},
			"keep_days": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: `Specifies the number of days to retain cross-region backups.`,
			},
			"backup_instances": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `Indicates the list of instances for which cross-region backups are created.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the ID of the instance.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the name of the instance.`,
						},
						"source_region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the source backup region.`,
						},
						"source_project_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the project ID of the source backup region.`,
						},
						"destination_region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the region where the cross-region backup is located.`,
						},
						"destination_project_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the project ID of the target backup region.`,
						},
						"datastore": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `Indicates the database information.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"version": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `Indicates the database engine version.`,
									},
									"type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `Indicates the database engine.`,
									},
								},
							},
						},
						"keep_days": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `Indicates the number of days to retain cross-region backups.`,
						},
					},
				},
			},
		},
	}
}

type CrossRegionBackupInstancesDSWrapper struct {
	*schemas.ResourceDataWrapper
	Config *config.Config
}

func newCrossRegionBackupInstancesDSWrapper(d *schema.ResourceData, meta interface{}) *CrossRegionBackupInstancesDSWrapper {
	return &CrossRegionBackupInstancesDSWrapper{
		ResourceDataWrapper: schemas.NewSchemaWrapper(d),
		Config:              meta.(*config.Config),
	}
}

func dataSourceRdsCrossRegionBackupInstancesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	wrapper := newCrossRegionBackupInstancesDSWrapper(d, meta)
	lisOffSitInsRst, err := wrapper.ListOffSiteInstances()
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	err = wrapper.listOffSiteInstancesToSchema(lisOffSitInsRst)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// @API RDS GET /v3/backups/offsite-backup-instance
func (w *CrossRegionBackupInstancesDSWrapper) ListOffSiteInstances() (*gjson.Result, error) {
	client, err := w.NewClient(w.Config, "rds")
	if err != nil {
		return nil, err
	}

	uri := "/v3/backups/offsite-backup-instance"
	return httphelper.New(client).
		Method("GET").
		URI(uri).
		OffsetPager("offsite_backup_instances", "offset", "limit", 0).
		Filter(
			filters.New().From("offsite_backup_instances").
				Where("id", "=", w.Get("instance_id")).
				Where("name", "=", w.Get("name")).
				Where("source_region", "=", w.Get("source_region")).
				Where("source_project_id", "=", w.Get("source_project_id")).
				Where("destination_region", "=", w.Get("destination_region")).
				Where("destination_project_id", "=", w.Get("destination_project_id")).
				Where("keep_days", "=", w.Get("keep_days")),
		).
		Request().
		Result()
}

func (w *CrossRegionBackupInstancesDSWrapper) listOffSiteInstancesToSchema(body *gjson.Result) error {
	d := w.ResourceData
	mErr := multierror.Append(nil,
		d.Set("region", w.Config.GetRegion(w.ResourceData)),
		d.Set("backup_instances", schemas.SliceToList(body.Get("offsite_backup_instances"),
			func(backupInstance gjson.Result) any {
				return map[string]any{
					"id":                     backupInstance.Get("id").Value(),
					"name":                   backupInstance.Get("name").Value(),
					"source_region":          backupInstance.Get("source_region").Value(),
					"source_project_id":      backupInstance.Get("source_project_id").Value(),
					"destination_region":     backupInstance.Get("destination_region").Value(),
					"destination_project_id": backupInstance.Get("destination_project_id").Value(),
					"datastore": schemas.SliceToList(backupInstance.Get("datastore"),
						func(datastore gjson.Result) any {
							return map[string]any{
								"version": datastore.Get("version").Value(),
								"type":    datastore.Get("type").Value(),
							}
						},
					),
					"keep_days": backupInstance.Get("keep_days").Value(),
				}
			},
		)),
	)
	return mErr.ErrorOrNil()
}