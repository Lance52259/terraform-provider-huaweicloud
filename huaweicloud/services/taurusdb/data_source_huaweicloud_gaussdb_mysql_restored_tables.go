// Generated by PMS #373
package taurusdb

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

func DataSourceGaussdbMysqlRestoredTables() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGaussdbMysqlRestoredTablesRead,

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
				Description: `Specifies the ID of the GaussDB MySQL instance,`,
			},
			"restore_time": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the backup time, in timestamp format.`,
			},
			"last_table_info": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies  whether data is restored to the most recent table.`,
			},
			"database_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the database name, which is used for fuzzy match.`,
			},
			"table_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the table name, which is used for fuzzy match.`,
			},
			"databases": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `Indicates the database information.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the database name.`,
						},
						"total_tables": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `Indicates the total number of tables.`,
						},
						"tables": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `Indicates the table information.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `Indicates the table name.`,
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

type MysqlRestoredTablesDSWrapper struct {
	*schemas.ResourceDataWrapper
	Config *config.Config
}

func newMysqlRestoredTablesDSWrapper(d *schema.ResourceData, meta interface{}) *MysqlRestoredTablesDSWrapper {
	return &MysqlRestoredTablesDSWrapper{
		ResourceDataWrapper: schemas.NewSchemaWrapper(d),
		Config:              meta.(*config.Config),
	}
}

func dataSourceGaussdbMysqlRestoredTablesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	wrapper := newMysqlRestoredTablesDSWrapper(d, meta)
	showRestoreTablesRst, err := wrapper.ShowRestoreTables()
	if err != nil {
		return diag.FromErr(err)
	}

	err = wrapper.showRestoreTablesToSchema(showRestoreTablesRst)
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

// @API GaussDBforMySQL GET /v3/{project_id}/instances/{instance_id}/backups/restore/tables
func (w *MysqlRestoredTablesDSWrapper) ShowRestoreTables() (*gjson.Result, error) {
	client, err := w.NewClient(w.Config, "gaussdb")
	if err != nil {
		return nil, err
	}

	uri := "/v3/{project_id}/instances/{instance_id}/backups/restore/tables"
	uri = strings.ReplaceAll(uri, "{instance_id}", w.Get("instance_id").(string))
	params := map[string]any{
		"restore_time":    w.Get("restore_time"),
		"last_table_info": w.Get("last_table_info"),
		"database_name":   w.Get("database_name"),
		"table_name":      w.Get("table_name"),
	}
	params = utils.RemoveNil(params)
	return httphelper.New(client).
		Method("GET").
		URI(uri).
		Query(params).
		Request().
		Result()
}

func (w *MysqlRestoredTablesDSWrapper) showRestoreTablesToSchema(body *gjson.Result) error {
	d := w.ResourceData
	mErr := multierror.Append(nil,
		d.Set("region", w.Config.GetRegion(w.ResourceData)),
		d.Set("databases", schemas.SliceToList(body.Get("databases"),
			func(databases gjson.Result) any {
				return map[string]any{
					"name":         databases.Get("name").Value(),
					"total_tables": databases.Get("total_tables").Value(),
					"tables": schemas.SliceToList(databases.Get("tables"),
						func(tables gjson.Result) any {
							return map[string]any{
								"name": tables.Get("name").Value(),
							}
						},
					),
				}
			},
		)),
	)
	return mErr.ErrorOrNil()
}
