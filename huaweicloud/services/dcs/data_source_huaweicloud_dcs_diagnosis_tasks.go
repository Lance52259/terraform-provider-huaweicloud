// Generated by PMS #153
package dcs

import (
	"context"
	"strings"

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

func DataSourceDcsDiagnosisTasks() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDcsDiagnosisTasksRead,

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
				Description: `Specifies the ID of the DCS instance.`,
			},
			"task_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the ID of the diagnosis task.`,
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the status of the diagnosis task.`,
			},
			"begin_time": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the start time of the diagnosis task, in RFC3339 format.`,
			},
			"end_time": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the end time of the diagnosis task, in RFC3339 format.`,
			},
			"node_num": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the number of diagnosed nodes.`,
			},
			"diagnosis_tasks": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `Indicates the list of diagnosis reports.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the diagnosis task ID.`,
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the diagnosis task status.`,
						},
						"begin_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the start time of the diagnosis task, in RFC3339 format.`,
						},
						"end_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the end time of the diagnosis task, in RFC3339 format.`,
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the time when the diagnosis report is created.`,
						},
						"node_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `Indicates the number of diagnosed nodes.`,
						},
						"abnormal_item_sum": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `Indicates the total number of abnormal diagnosis items.`,
						},
						"failed_item_sum": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `Indicates the total number of failed diagnosis items.`,
						},
					},
				},
			},
		},
	}
}

type DiagnosisTasksDSWrapper struct {
	*schemas.ResourceDataWrapper
	Config *config.Config
}

func newDiagnosisTasksDSWrapper(d *schema.ResourceData, meta interface{}) *DiagnosisTasksDSWrapper {
	return &DiagnosisTasksDSWrapper{
		ResourceDataWrapper: schemas.NewSchemaWrapper(d),
		Config:              meta.(*config.Config),
	}
}

func dataSourceDcsDiagnosisTasksRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	wrapper := newDiagnosisTasksDSWrapper(d, meta)
	lisDiaTasRst, err := wrapper.ListDiagnosisTasks()
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	err = wrapper.listDiagnosisTasksToSchema(lisDiaTasRst)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// @API DCS GET /v2/{project_id}/instances/{instance_id}/diagnosis
func (w *DiagnosisTasksDSWrapper) ListDiagnosisTasks() (*gjson.Result, error) {
	client, err := w.NewClient(w.Config, "dcs")
	if err != nil {
		return nil, err
	}

	uri := "/v2/{project_id}/instances/{instance_id}/diagnosis"
	uri = strings.ReplaceAll(uri, "{instance_id}", w.Get("instance_id").(string))
	return httphelper.New(client).
		Method("GET").
		URI(uri).
		OffsetPager("diagnosis_report_list", "offset", "limit", 0).
		Filter(
			filters.New().From("diagnosis_report_list").
				Where("report_id", "=", w.Get("task_id")).
				Where("status", "=", w.Get("status")).
				Where("begin_time", "=", w.Get("begin_time")).
				Where("end_time", "=", w.Get("end_time")).
				Where("node_num", "=", w.GetToInt("node_num")),
		).
		Request().
		Result()
}

func (w *DiagnosisTasksDSWrapper) listDiagnosisTasksToSchema(body *gjson.Result) error {
	d := w.ResourceData
	mErr := multierror.Append(nil,
		d.Set("region", w.Config.GetRegion(w.ResourceData)),
		d.Set("diagnosis_tasks", schemas.SliceToList(body.Get("diagnosis_report_list"),
			func(diagnosisTask gjson.Result) any {
				return map[string]any{
					"id":                diagnosisTask.Get("report_id").Value(),
					"status":            diagnosisTask.Get("status").Value(),
					"begin_time":        diagnosisTask.Get("begin_time").Value(),
					"end_time":          diagnosisTask.Get("end_time").Value(),
					"created_at":        diagnosisTask.Get("created_at").Value(),
					"node_num":          diagnosisTask.Get("node_num").Value(),
					"abnormal_item_sum": diagnosisTask.Get("abnormal_item_sum").Value(),
					"failed_item_sum":   diagnosisTask.Get("failed_item_sum").Value(),
				}
			},
		)),
	)
	return mErr.ErrorOrNil()
}