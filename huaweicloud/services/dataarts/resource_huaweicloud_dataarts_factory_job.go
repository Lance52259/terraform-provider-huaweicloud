// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product DataArts
// ---------------------------------------------------------------

package dataarts

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jmespath/go-jmespath"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DataArtsStudio POST /v1/{project_id}/jobs
// @API DataArtsStudio GET /v1/{project_id}/jobs/{job_name}
// @API DataArtsStudio PUT /v1/{project_id}/jobs/{job_name}
// @API DataArtsStudio DELETE /v1/{project_id}/jobs/{job_name}
func ResourceFactoryJob() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFactoryJobCreate,
		UpdateContext: resourceFactoryJobUpdate,
		ReadContext:   resourceFactoryJobRead,
		DeleteContext: resourceFactoryJobDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceFactoryJobImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Job name.`,
			},
			"nodes": {
				Type:        schema.TypeList,
				Elem:        jobNodeSchema(),
				Required:    true,
				Description: `Node definition.`,
			},
			"schedule": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Elem:        jobScheduleSchema(),
				Required:    true,
				Description: `Scheduling configuration.`,
			},
			"process_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Job type.`,
			},
			"workspace_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: `The workspace ID.`,
			},
			"params": {
				Type:        schema.TypeList,
				Elem:        jobParamSchema(),
				Optional:    true,
				Computed:    true,
				Description: `Job parameter definition.`,
			},
			"directory": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Path of a job in the directory tree.`,
			},
			"log_path": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The OBS path where job execution logs are stored.`,
			},
			"basic_config": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Elem:        jobBasicConfigSchema(),
				Optional:    true,
				Computed:    true,
				Description: `Baisc job information.`,
			},
		},
	}
}

func jobNodeSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Node name.`,
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Node type.`,
			},
			"location": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Elem:     jobNodeLocationSchema(),
				Required: true,
			},
			"pre_node_name": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Computed:    true,
				Description: `Name of the previous node on which the current node depends.`,
			},
			"conditions": {
				Type:        schema.TypeList,
				Elem:        jobNodeConditionSchema(),
				Optional:    true,
				Computed:    true,
				Description: `Node execution condition.`,
			},
			"properties": {
				Type:        schema.TypeList,
				Elem:        jobNodePropertySchema(),
				Required:    true,
				Description: `Node property. Each type of node has its own property definition.`,
			},
			"polling_interval": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `Interval at which node running results are checked.`,
			},
			"max_execution_time": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `Maximum execution time of a node.`,
			},
			"retry_times": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `Number of the node retries.`,
			},
			"retry_interval": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `Interval at which a retry is performed upon a failure.`,
			},
			"fail_policy": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Node failure policy.`,
			},
			"event_trigger": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Elem:        jobNodeEventTriggerSchema(),
				Optional:    true,
				Computed:    true,
				Description: `Event trigger for the real-time job node.`,
			},
			"cron_trigger": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Elem:        jobNodeCronTriggerSchema(),
				Optional:    true,
				Computed:    true,
				Description: `Cron trigger for the real-time job node`,
			},
		},
	}
	return &sc
}

func jobNodeLocationSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"x": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `Position of the node on the horizontal axis of the job canvas.`,
			},
			"y": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `Position of the node on the vertical axis of the job canvas.`,
			},
		},
	}
	return &sc
}

func jobNodeConditionSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"pre_node_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Name of the previous node on which the current node depends.`,
			},
			"expression": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `EL expression.`,
			},
		},
	}
	return &sc
}

func jobNodePropertySchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Property name.`,
			},
			"value": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Property value.`,
			},
		},
	}
	return &sc
}

func jobNodeEventTriggerSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"event_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Event type.`,
			},
			"channel": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `DIS stream name.`,
			},
			"fail_policy": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Job failure policy.`,
			},
			"concurrent": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `Number of the concurrently scheduled jobs.`,
			},
			"read_policy": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Access policy.`,
			},
		},
	}
	return &sc
}

func jobNodeCronTriggerSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"start_time": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Scheduling start time.`,
			},
			"end_time": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Scheduling end time.`,
			},
			"expression": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Cron expression in the format of **<second><minute><hour><day><month><week>**.`,
			},
			"expression_time_zone": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Time zone corresponding to the Cron expression.`,
			},
			"period": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Job execution interval consisting of a time and time unit.`,
			},
			"depend_pre_period": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
				Description: `Indicates whether to depend on the execution result of the current
                 job's dependent job in the previous scheduling period.`,
			},
			"depend_jobs": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Elem:     jobCronTriggerDependJobsSchema(),
				Optional: true,
				Computed: true,
			},
			"concurrent": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `Number of concurrent executions allowed.`,
			},
		},
	}
	return &sc
}

func jobCronTriggerDependJobsSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"jobs": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Required:    true,
				Description: `A list of dependent jobs. Only the existing jobs can be depended on.`,
			},
			"depend_period": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Dependency period.`,
			},
			"depend_fail_policy": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Dependency job failure policy.`,
			},
		},
	}
	return &sc
}

func jobScheduleSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Scheduling type.`,
			},
			"cron": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Elem:     jobScheduleCronSchema(),
				Optional: true,
				Computed: true,
			},
			"event": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Elem:     jobScheduleEventSchema(),
				Optional: true,
				Computed: true,
			},
		},
	}
	return &sc
}

func jobScheduleCronSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"start_time": {
				Type:     schema.TypeString,
				Required: true,
				Description: `Scheduling start time in the format of **yyyy-MM-dd'T'HH:mm:ssZ**, 
                which is an ISO 8601 time format.`,
			},
			"end_time": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: `Scheduling end time in the format of **yyyy-MM-dd'T'HH:mm:ssZ**, 
                which is an ISO 8601 time format.`,
			},
			"expression": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Cron expression in the format of **<second><minute><hour><day><month><week>**.`,
			},
			"expression_time_zone": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Time zone corresponding to the Cron expression.`,
			},
			"depend_pre_period": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
				Description: `Indicates whether to depend on the execution result of 
                the current job's dependent job in the previous scheduling period.`,
			},
			"depend_jobs": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Elem:     jobCronDependJobsSchema(),
				Optional: true,
				Computed: true,
			},
		},
	}
	return &sc
}

func jobCronDependJobsSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"jobs": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Required:    true,
				Description: `A list of dependent jobs. Only the existing jobs can be depended on.`,
			},
			"depend_period": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Dependency period.`,
			},
			"depend_fail_policy": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Dependency job failure policy.`,
			},
		},
	}
	return &sc
}

func jobScheduleEventSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"event_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Event type.`,
			},
			"channel": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `DIS stream name.`,
			},
			"fail_policy": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Job failure policy.`,
			},
			"concurrent": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `Number of the concurrently scheduled jobs.`,
			},
			"read_policy": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Access policy.`,
			},
		},
	}
	return &sc
}

func jobParamSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Name of a parameter.`,
			},
			"value": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Value of the parameter.`,
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Parameter type.`,
			},
		},
	}
	return &sc
}

func jobBasicConfigSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"owner": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Job owner.`,
			},
			"priority": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `Job priority.`,
			},
			"execute_user": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Job execution user. The value must be an existing user.`,
			},
			"instance_timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `Maximum execution time of a job instance.`,
			},
			"custom_fields": {
				Type:        schema.TypeMap,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Computed:    true,
				Description: `Custom fields.`,
			},
		},
	}
	return &sc
}

func resourceFactoryJobCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		createJobHttpUrl = "v1/{project_id}/jobs"
		createJobProduct = "dataarts-dlf"
	)
	createJobClient, err := cfg.NewServiceClient(createJobProduct, region)
	if err != nil {
		return diag.Errorf("error creating DataArts client: %s", err)
	}

	createJobPath := createJobClient.Endpoint + createJobHttpUrl
	createJobPath = strings.ReplaceAll(createJobPath, "{project_id}", createJobClient.ProjectID)

	createJobOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			204,
		},
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	if v, ok := d.GetOk("workspace_id"); ok {
		createJobOpt.MoreHeaders["workspace"] = v.(string)
	}

	createJobOpt.JSONBody = utils.RemoveNil(buildCreateJobBodyParams(d))
	_, err = createJobClient.Request("POST", createJobPath, &createJobOpt)
	if err != nil {
		return diag.Errorf("error creating Job: %s", err)
	}

	d.SetId(d.Get("name").(string))

	return resourceFactoryJobRead(ctx, d, meta)
}

func buildCreateJobBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":        d.Get("name"),
		"nodes":       buildCreateJobRequestBodyNode(d.Get("nodes")),
		"schedule":    buildCreateJobRequestBodySchedule(d.Get("schedule")),
		"params":      buildCreateJobRequestBodyParam(d.Get("params")),
		"directory":   utils.ValueIgnoreEmpty(d.Get("directory")),
		"processType": d.Get("process_type"),
		"logPath":     utils.ValueIgnoreEmpty(d.Get("log_path")),
		"basicConfig": buildCreateJobRequestBodyBasicConfig(d.Get("basic_config")),
	}
	return bodyParams
}

func buildCreateJobRequestBodyNode(rawParams interface{}) []map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}

		rst := make([]map[string]interface{}, len(rawArray))
		for i, v := range rawArray {
			raw := v.(map[string]interface{})
			rst[i] = map[string]interface{}{
				"name":             utils.ValueIgnoreEmpty(raw["name"]),
				"type":             utils.ValueIgnoreEmpty(raw["type"]),
				"location":         buildNodeLocation(raw["location"]),
				"preNodeName":      utils.ValueIgnoreEmpty(raw["pre_node_name"]),
				"conditions":       buildNodeCondition(raw["conditions"]),
				"properties":       buildNodeProperty(raw["properties"]),
				"pollingInterval":  utils.ValueIgnoreEmpty(raw["polling_interval"]),
				"maxExecutionTime": utils.ValueIgnoreEmpty(raw["max_execution_time"]),
				"retryTimes":       utils.ValueIgnoreEmpty(raw["retry_times"]),
				"retryInterval":    utils.ValueIgnoreEmpty(raw["retry_interval"]),
				"failPolicy":       utils.ValueIgnoreEmpty(raw["fail_policy"]),
				"eventTrigger":     buildNodeEventTrigger(raw["event_trigger"]),
				"cronTrigger":      buildNodeCronTrigger(raw["cron_trigger"]),
			}
		}
		return rst
	}
	return nil
}

func buildNodeLocation(rawParams interface{}) map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}
		raw, ok := rawArray[0].(map[string]interface{})
		if !ok {
			return nil
		}

		params := map[string]interface{}{
			"x": utils.ValueIgnoreEmpty(raw["x"]),
			"y": utils.ValueIgnoreEmpty(raw["y"]),
		}
		return params
	}
	return nil
}

func buildNodeCondition(rawParams interface{}) []map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}

		rst := make([]map[string]interface{}, len(rawArray))
		for i, v := range rawArray {
			raw := v.(map[string]interface{})
			rst[i] = map[string]interface{}{
				"preNodeName": utils.ValueIgnoreEmpty(raw["pre_node_name"]),
				"expression":  utils.ValueIgnoreEmpty(raw["expression"]),
			}
		}
		return rst
	}
	return nil
}

func buildNodeProperty(rawParams interface{}) []map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}

		rst := make([]map[string]interface{}, len(rawArray))
		for i, v := range rawArray {
			raw := v.(map[string]interface{})
			rst[i] = map[string]interface{}{
				"name":  utils.ValueIgnoreEmpty(raw["name"]),
				"value": utils.ValueIgnoreEmpty(raw["value"]),
			}
		}
		return rst
	}
	return nil
}

func buildNodeEventTrigger(rawParams interface{}) map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}
		raw, ok := rawArray[0].(map[string]interface{})
		if !ok {
			return nil
		}

		params := map[string]interface{}{
			"eventType":  utils.ValueIgnoreEmpty(raw["event_type"]),
			"channel":    utils.ValueIgnoreEmpty(raw["channel"]),
			"failPolicy": utils.ValueIgnoreEmpty(raw["fail_policy"]),
			"concurrent": utils.ValueIgnoreEmpty(raw["concurrent"]),
			"readPolicy": utils.ValueIgnoreEmpty(raw["read_policy"]),
		}
		return params
	}
	return nil
}

func buildNodeCronTrigger(rawParams interface{}) map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}
		raw, ok := rawArray[0].(map[string]interface{})
		if !ok {
			return nil
		}

		params := map[string]interface{}{
			"startTime":          utils.ValueIgnoreEmpty(raw["start_time"]),
			"endTime":            utils.ValueIgnoreEmpty(raw["end_time"]),
			"expression":         utils.ValueIgnoreEmpty(raw["expression"]),
			"expressionTimeZone": utils.ValueIgnoreEmpty(raw["expression_time_zone"]),
			"period":             utils.ValueIgnoreEmpty(raw["period"]),
			"dependPrePeriod":    utils.ValueIgnoreEmpty(raw["depend_pre_period"]),
			"dependJobs":         buildCronTriggerDependJobs(raw["depend_jobs"]),
			"concurrent":         utils.ValueIgnoreEmpty(raw["concurrent"]),
		}
		return params
	}
	return nil
}

func buildCronTriggerDependJobs(rawParams interface{}) map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}
		raw, ok := rawArray[0].(map[string]interface{})
		if !ok {
			return nil
		}

		params := map[string]interface{}{
			"jobs":             utils.ValueIgnoreEmpty(raw["jobs"]),
			"dependPeriod":     utils.ValueIgnoreEmpty(raw["depend_period"]),
			"dependFailPolicy": utils.ValueIgnoreEmpty(raw["depend_fail_policy"]),
		}
		return params
	}
	return nil
}

func buildCreateJobRequestBodySchedule(rawParams interface{}) map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}
		raw, ok := rawArray[0].(map[string]interface{})
		if !ok {
			return nil
		}

		params := map[string]interface{}{
			"type":  utils.ValueIgnoreEmpty(raw["type"]),
			"cron":  buildScheduleCron(raw["cron"]),
			"event": buildScheduleEvent(raw["event"]),
		}
		return params
	}
	return nil
}

func buildScheduleCron(rawParams interface{}) map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}
		raw, ok := rawArray[0].(map[string]interface{})
		if !ok {
			return nil
		}

		params := map[string]interface{}{
			"startTime":          utils.ValueIgnoreEmpty(raw["start_time"]),
			"endTime":            utils.ValueIgnoreEmpty(raw["end_time"]),
			"expression":         utils.ValueIgnoreEmpty(raw["expression"]),
			"expressionTimeZone": utils.ValueIgnoreEmpty(raw["expression_time_zone"]),
			"dependPrePeriod":    utils.ValueIgnoreEmpty(raw["depend_pre_period"]),
			"dependJobs":         buildCronDependJobs(raw["depend_jobs"]),
		}
		return params
	}
	return nil
}

func buildCronDependJobs(rawParams interface{}) map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}
		raw, ok := rawArray[0].(map[string]interface{})
		if !ok {
			return nil
		}

		params := map[string]interface{}{
			"jobs":             raw["jobs"],
			"dependPeriod":     utils.ValueIgnoreEmpty(raw["depend_period"]),
			"dependFailPolicy": utils.ValueIgnoreEmpty(raw["depend_fail_policy"]),
		}
		return params
	}
	return nil
}

func buildScheduleEvent(rawParams interface{}) map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}
		raw, ok := rawArray[0].(map[string]interface{})
		if !ok {
			return nil
		}

		params := map[string]interface{}{
			"eventType":  utils.ValueIgnoreEmpty(raw["event_type"]),
			"channel":    utils.ValueIgnoreEmpty(raw["channel"]),
			"failPolicy": utils.ValueIgnoreEmpty(raw["fail_policy"]),
			"concurrent": utils.ValueIgnoreEmpty(raw["concurrent"]),
			"readPolicy": utils.ValueIgnoreEmpty(raw["read_policy"]),
		}
		return params
	}
	return nil
}

func buildCreateJobRequestBodyParam(rawParams interface{}) []map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}

		rst := make([]map[string]interface{}, len(rawArray))
		for i, v := range rawArray {
			raw := v.(map[string]interface{})
			rst[i] = map[string]interface{}{
				"name":  utils.ValueIgnoreEmpty(raw["name"]),
				"value": utils.ValueIgnoreEmpty(raw["value"]),
				"type":  utils.ValueIgnoreEmpty(raw["type"]),
			}
		}
		return rst
	}
	return nil
}

func buildCreateJobRequestBodyBasicConfig(rawParams interface{}) map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}
		raw, ok := rawArray[0].(map[string]interface{})
		if !ok {
			return nil
		}

		params := map[string]interface{}{
			"owner":           utils.ValueIgnoreEmpty(raw["owner"]),
			"priority":        utils.ValueIgnoreEmpty(raw["priority"]),
			"executeUser":     utils.ValueIgnoreEmpty(raw["execute_user"]),
			"instanceTimeout": utils.ValueIgnoreEmpty(raw["instance_timeout"]),
			"customFields":    utils.ValueIgnoreEmpty(raw["custom_fields"]),
		}
		return params
	}
	return nil
}

func resourceFactoryJobRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		getJobHttpUrl = "v1/{project_id}/jobs/{job_name}"
		getJobProduct = "dataarts-dlf"
	)
	getJobClient, err := cfg.NewServiceClient(getJobProduct, region)
	if err != nil {
		return diag.Errorf("error creating DataArts client: %s", err)
	}

	getJobPath := getJobClient.Endpoint + getJobHttpUrl
	getJobPath = strings.ReplaceAll(getJobPath, "{project_id}", getJobClient.ProjectID)
	getJobPath = strings.ReplaceAll(getJobPath, "{job_name}", d.Id())

	getJobOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
		MoreHeaders: map[string]string{"Content-Type": "application/json"},
	}

	if v, ok := d.GetOk("workspace_id"); ok {
		getJobOpt.MoreHeaders["workspace"] = v.(string)
	}

	getJobResp, err := getJobClient.Request("GET", getJobPath, &getJobOpt)

	if err != nil {
		return common.CheckDeletedDiag(d, parseFactoryJobNotFoundError(err), "error retrieving Job")
	}

	getJobRespBody, err := utils.FlattenResponse(getJobResp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("name", getJobRespBody, nil)),
		d.Set("nodes", flattenGetJobResponseBodyNode(getJobRespBody)),
		d.Set("schedule", flattenGetJobResponseBodySchedule(utils.PathSearch("schedule",
			getJobRespBody, make(map[string]interface{})).(map[string]interface{}))),
		d.Set("params", flattenGetJobResponseBodyParam(utils.PathSearch("params",
			getJobRespBody, make([]interface{}, 0)).([]interface{}))),
		d.Set("directory", utils.PathSearch("directory", getJobRespBody, nil)),
		d.Set("process_type", utils.PathSearch("processType", getJobRespBody, nil)),
		d.Set("log_path", utils.PathSearch("logPath", getJobRespBody, nil)),
		d.Set("basic_config", flattenGetJobResponseBodyBasicConfig(utils.PathSearch("basicConfig",
			getJobRespBody, make(map[string]interface{})).(map[string]interface{}))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenGetJobResponseBodyNode(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("nodes", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"name":               utils.PathSearch("name", v, nil),
			"type":               utils.PathSearch("type", v, nil),
			"location":           flattenNodeLocation(utils.PathSearch("location", v, make(map[string]interface{})).(map[string]interface{})),
			"pre_node_name":      utils.PathSearch("preNodeName", v, nil),
			"conditions":         flattenNodeConditions(utils.PathSearch("conditions", v, make([]interface{}, 0)).([]interface{})),
			"properties":         flattenNodeProperties(utils.PathSearch("properties", v, make([]interface{}, 0)).([]interface{})),
			"polling_interval":   utils.PathSearch("pollingInterval", v, nil),
			"max_execution_time": utils.PathSearch("maxExecutionTime", v, nil),
			"retry_times":        utils.PathSearch("retryTimes", v, nil),
			"retry_interval":     utils.PathSearch("retryInterval", v, nil),
			"fail_policy":        utils.PathSearch("failPolicy", v, nil),
			"event_trigger":      flattenNodeEventTrigger(utils.PathSearch("eventTrigger", v, make(map[string]interface{})).(map[string]interface{})),
			"cron_trigger":       flattenNodeCronTrigger(utils.PathSearch("cron_trigger", v, make(map[string]interface{})).(map[string]interface{})),
		})
	}
	return rst
}

func flattenNodeLocation(location map[string]interface{}) []interface{} {
	if len(location) < 1 {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"x": int(utils.PathSearch("to_number(x)", location, 0.0).(float64)),
			"y": int(utils.PathSearch("to_number(y)", location, 0.0).(float64)),
		},
	}
}

func flattenNodeConditions(conditions []interface{}) []interface{} {
	if len(conditions) < 1 {
		return nil
	}

	rst := make([]interface{}, 0, len(conditions))
	for _, v := range conditions {
		rst = append(rst, map[string]interface{}{
			"pre_node_name": utils.PathSearch("preNodeName", v, nil),
			"expression":    utils.PathSearch("expression", v, nil),
		})
	}
	return rst
}

func flattenNodeProperties(properties []interface{}) []interface{} {
	if len(properties) < 1 {
		return nil
	}

	rst := make([]interface{}, 0, len(properties))
	for _, v := range properties {
		rst = append(rst, map[string]interface{}{
			"name":  utils.PathSearch("name", v, nil),
			"value": utils.PathSearch("value", v, nil),
		})
	}
	return rst
}

func flattenNodeEventTrigger(trigger map[string]interface{}) []interface{} {
	if len(trigger) < 1 {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"event_type":  utils.PathSearch("eventType", trigger, nil),
			"channel":     utils.PathSearch("channel", trigger, nil),
			"fail_policy": utils.PathSearch("failPolicy", trigger, nil),
			"concurrent":  utils.PathSearch("concurrent", trigger, nil),
			"read_policy": utils.PathSearch("readPolicy", trigger, nil),
		},
	}
}

func flattenNodeCronTrigger(trigger map[string]interface{}) []interface{} {
	if len(trigger) < 1 {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"start_time":           utils.PathSearch("startTime", trigger, nil),
			"end_time":             utils.PathSearch("endTime", trigger, nil),
			"expression":           utils.PathSearch("expression", trigger, nil),
			"expression_time_zone": utils.PathSearch("expressionTimeZone", trigger, nil),
			"period":               utils.PathSearch("period", trigger, nil),
			"depend_pre_period":    utils.PathSearch("dependPrePeriod", trigger, nil),
			"depend_jobs": flattenCronTriggerDependJobs(utils.PathSearch("dependJobs",
				trigger, make(map[string]interface{})).(map[string]interface{})),
			"concurrent": utils.PathSearch("concurrent", trigger, nil),
		},
	}
}

func flattenCronTriggerDependJobs(jobs map[string]interface{}) []interface{} {
	if len(jobs) < 1 {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"jobs":               utils.PathSearch("jobs", jobs, nil),
			"depend_period":      utils.PathSearch("dependPeriod", jobs, nil),
			"depend_fail_policy": utils.PathSearch("dependFailPolicy", jobs, nil),
		},
	}
}

func flattenGetJobResponseBodySchedule(schedule map[string]interface{}) []interface{} {
	if len(schedule) < 1 {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"type": utils.PathSearch("type", schedule, nil),
			"cron": flattenScheduleCron(utils.PathSearch("cron",
				schedule, make(map[string]interface{})).(map[string]interface{})),
			"event": flattenScheduleEvent(utils.PathSearch("event",
				schedule, make(map[string]interface{})).(map[string]interface{})),
		},
	}
}

func flattenScheduleCron(schedule map[string]interface{}) []interface{} {
	if len(schedule) < 1 {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"start_time":           utils.PathSearch("startTime", schedule, nil),
			"end_time":             utils.PathSearch("endTime", schedule, nil),
			"expression":           utils.PathSearch("expression", schedule, nil),
			"expression_time_zone": utils.PathSearch("expressionTimeZone", schedule, nil),
			"depend_pre_period":    utils.PathSearch("dependPrePeriod", schedule, nil),
			"depend_jobs": flattenCronDependJobs(utils.PathSearch("dependJobs",
				schedule, make(map[string]interface{})).(map[string]interface{})),
		},
	}
}

func flattenCronDependJobs(jobs map[string]interface{}) []interface{} {
	if len(jobs) < 1 {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"jobs":               utils.PathSearch("jobs", jobs, nil),
			"depend_period":      utils.PathSearch("dependPeriod", jobs, nil),
			"depend_fail_policy": utils.PathSearch("dependFailPolicy", jobs, nil),
		},
	}
}

func flattenScheduleEvent(event map[string]interface{}) []interface{} {
	if len(event) < 1 {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"event_type":  utils.PathSearch("eventType", event, nil),
			"channel":     utils.PathSearch("channel", event, nil),
			"fail_policy": utils.PathSearch("failPolicy", event, nil),
			"concurrent":  utils.PathSearch("concurrent", event, nil),
			"read_policy": utils.PathSearch("readPolicy", event, nil),
		},
	}
}

func flattenGetJobResponseBodyParam(params []interface{}) []interface{} {
	if len(params) < 1 {
		return nil
	}

	rst := make([]interface{}, 0, len(params))
	for _, v := range params {
		rst = append(rst, map[string]interface{}{
			"name":  utils.PathSearch("name", v, nil),
			"value": utils.PathSearch("value", v, nil),
			"type":  utils.PathSearch("type", v, nil),
		})
	}
	return rst
}

func flattenGetJobResponseBodyBasicConfig(basicConfig map[string]interface{}) []interface{} {
	if len(basicConfig) < 1 {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"owner":            utils.PathSearch("owner", basicConfig, nil),
			"priority":         utils.PathSearch("priority", basicConfig, nil),
			"execute_user":     utils.PathSearch("executeUser", basicConfig, nil),
			"instance_timeout": utils.PathSearch("instanceTimeout", basicConfig, nil),
			"custom_fields":    utils.PathSearch("customFields", basicConfig, nil),
		},
	}
}

func resourceFactoryJobUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	updateJobChanges := []string{
		"name",
		"nodes",
		"schedule",
		"params",
		"directory",
		"process_type",
		"log_path",
		"basic_config",
	}

	if d.HasChanges(updateJobChanges...) {
		var (
			updateJobHttpUrl = "v1/{project_id}/jobs/{job_name}"
			updateJobProduct = "dataarts-dlf"
		)
		updateJobClient, err := cfg.NewServiceClient(updateJobProduct, region)
		if err != nil {
			return diag.Errorf("error creating DataArts client: %s", err)
		}

		updateJobPath := updateJobClient.Endpoint + updateJobHttpUrl
		updateJobPath = strings.ReplaceAll(updateJobPath, "{project_id}", updateJobClient.ProjectID)
		updateJobPath = strings.ReplaceAll(updateJobPath, "{job_name}", d.Id())

		updateJobOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			OkCodes: []int{
				204,
			},
			MoreHeaders: map[string]string{
				"Content-Type": "application/json",
			},
		}

		if v, ok := d.GetOk("workspace_id"); ok {
			updateJobOpt.MoreHeaders["workspace"] = v.(string)
		}

		updateJobOpt.JSONBody = utils.RemoveNil(buildUpdateJobBodyParams(d))
		_, err = updateJobClient.Request("PUT", updateJobPath, &updateJobOpt)
		if err != nil {
			return diag.Errorf("error updating Job: %s", err)
		}
	}
	return resourceFactoryJobRead(ctx, d, meta)
}

func buildUpdateJobBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":        d.Get("name"),
		"nodes":       buildCreateJobRequestBodyNode(d.Get("nodes")),
		"schedule":    buildCreateJobRequestBodySchedule(d.Get("schedule")),
		"params":      buildCreateJobRequestBodyParam(d.Get("params")),
		"directory":   utils.ValueIgnoreEmpty(d.Get("directory")),
		"processType": d.Get("process_type"),
		"logPath":     utils.ValueIgnoreEmpty(d.Get("log_path")),
		"basicConfig": buildCreateJobRequestBodyBasicConfig(d.Get("basic_config")),
	}
	return bodyParams
}

func resourceFactoryJobDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		deleteJobHttpUrl = "v1/{project_id}/jobs/{job_name}"
		deleteJobProduct = "dataarts-dlf"
	)
	deleteJobClient, err := cfg.NewServiceClient(deleteJobProduct, region)
	if err != nil {
		return diag.Errorf("error creating DataArts client: %s", err)
	}

	deleteJobPath := deleteJobClient.Endpoint + deleteJobHttpUrl
	deleteJobPath = strings.ReplaceAll(deleteJobPath, "{project_id}", deleteJobClient.ProjectID)
	deleteJobPath = strings.ReplaceAll(deleteJobPath, "{job_name}", d.Id())

	deleteJobOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			204,
		},
		MoreHeaders: map[string]string{"Content-Type": "application/json"},
	}

	if v, ok := d.GetOk("workspace_id"); ok {
		deleteJobOpt.MoreHeaders["workspace"] = v.(string)
	}

	_, err = deleteJobClient.Request("DELETE", deleteJobPath, &deleteJobOpt)
	if err != nil {
		return diag.Errorf("error deleting Job: %s", err)
	}

	return nil
}

func resourceFactoryJobImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import id, must be <workspace_id>/<name>")
	}

	d.Set("workspace_id", parts[0])
	d.SetId(parts[1])

	return []*schema.ResourceData{d}, nil
}

func parseFactoryJobNotFoundError(respErr error) error {
	var apiErr interface{}
	if errCode, ok := respErr.(golangsdk.ErrDefault400); ok {
		pErr := json.Unmarshal(errCode.Body, &apiErr)
		if pErr != nil {
			return pErr
		}
		errCode, err := jmespath.Search(`error_code`, apiErr)
		if err != nil {
			return fmt.Errorf("error parse error_code from response body: %s", err.Error())
		}

		if errCode == `DLF.0100` {
			return golangsdk.ErrDefault404{}
		}
	}
	return respErr
}
