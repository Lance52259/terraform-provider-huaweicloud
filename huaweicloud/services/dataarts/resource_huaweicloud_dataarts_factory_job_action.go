package dataarts

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var (
	actionResourceNotFoundCodes = []string{
		"DLF.0819", // The workspace ID does not exist.
	}
	jobActionNonUpdatableParams = []string{
		"job_name",
		"process_type",
		"workspace_id",
	}
)

// @API DataArtsStudio POST /v1/{project_id}/jobs/{job_name}/start
// @API DataArtsStudio POST /v1/{project_id}/jobs/{job_name}/stop
// @API DataArtsStudio GET /v1/{project_id}/jobs
func ResourceFactoryJobAction() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFactoryJobActionCreate,
		ReadContext:   resourceFactoryJobActionRead,
		UpdateContext: resourceFactoryJobActionUpdate,
		DeleteContext: resourceFactoryJobActionDelete,

		CustomizeDiff: config.FlexibleForceNew(jobActionNonUpdatableParams),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the job is located.`,
			},

			// Required parameters.
			"job_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of the job to be performed.`,
			},
			"process_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The type of the job to be performed.`,
			},
			"action": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The action type of the job to be performed.`,
			},

			// Optional parameters.
			"workspace_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The ID of the workspace to which the job belongs.`,
			},

			// Attribute.
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The current status of the job.`,
			},
		},
	}
}

func buildFactoryRequestMoreHeaders(workspaceId string) map[string]string {
	results := map[string]string{
		"Content-Type": "application/json",
	}

	if workspaceId != "" {
		results["workspace"] = workspaceId
	}

	return results
}

func startJob(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	var (
		httpUrl     = "v1/{project_id}/jobs/{job_name}/start"
		workspaceId = d.Get("workspace_id").(string)
		jobName     = d.Get("job_name").(string)
	)

	actionPath := client.Endpoint + httpUrl
	actionPath = strings.ReplaceAll(actionPath, "{project_id}", client.ProjectID)
	actionPath = strings.ReplaceAll(actionPath, "{job_name}", jobName)

	actionOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildFactoryRequestMoreHeaders(workspaceId),
	}

	_, err := client.Request("POST", actionPath, &actionOpts)
	return err
}

func stopJob(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	var (
		httpUrl     = "v1/{project_id}/jobs/{job_name}/stop"
		workspaceId = d.Get("workspace_id").(string)
		jobName     = d.Get("job_name").(string)
	)

	actionPath := client.Endpoint + httpUrl
	actionPath = strings.ReplaceAll(actionPath, "{project_id}", client.ProjectID)
	actionPath = strings.ReplaceAll(actionPath, "{job_name}", jobName)

	actionOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildFactoryRequestMoreHeaders(workspaceId),
	}

	_, err := client.Request("POST", actionPath, &actionOpts)
	return err
}

func getJobByName(client *golangsdk.ServiceClient, workspaceId, jobName, jobType string) (interface{}, error) {
	// The maximum value of limit is 100.
	httpUrl := "v1/{project_id}/jobs?limit=100"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	// The job_name field supports fuzzy matching.
	getPath = fmt.Sprintf("%s&jobName=%v&jobType=%v", getPath, jobName, jobType)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"workspace": workspaceId,
		},
	}

	offset := 0
	for {
		listPathWithOffset := fmt.Sprintf("%s&offset=%d", getPath, offset)
		requestResp, err := client.Request("GET", listPathWithOffset, &getOpt)
		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}

		jobs := utils.PathSearch("jobs", respBody, make([]interface{}, 0)).([]interface{})
		if len(jobs) < 1 {
			break
		}

		job := utils.PathSearch(fmt.Sprintf("[?name=='%s']|[0]", jobName), jobs, nil)
		if job != nil {
			return job, nil
		}
		offset += len(jobs)
	}

	return nil, golangsdk.ErrDefault404{
		ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
			Method:    "GET",
			URL:       "/v1/{project_id}/jobs",
			RequestId: "NONE",
			Body:      []byte(fmt.Sprintf("the job (%s) is not found", jobName)),
		},
	}
}

func jobStateRefreshFunc(client *golangsdk.ServiceClient, workspaceId, jobName, jobType string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		respBody, err := getJobByName(client, workspaceId, jobName, jobType)
		if err != nil {
			return respBody, "ERROR", err
		}

		var (
			jobStatus           = utils.PathSearch("status", respBody, "").(string)
			unexpectedJobStatus = []string{"EXCEPTION"}
		)
		if utils.StrSliceContains(unexpectedJobStatus, jobStatus) {
			return respBody, "ERROR", fmt.Errorf("unexpected job status (%s)", jobStatus)
		}

		if utils.StrSliceContains([]string{"NORMAL", "STOPPED", "SCHEDULING"}, jobStatus) {
			return respBody, "COMPLETED", nil
		}

		return respBody, "PENDING", nil
	}
}

func doActionJob(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData, timeout time.Duration) error {
	var (
		workspaceId = d.Get("workspace_id").(string)
		jobName     = d.Get("job_name").(string)
		processType = d.Get("process_type").(string)
		actionType  = d.Get("action").(string)
		err         error
	)

	switch actionType {
	case "start":
		err = startJob(client, d)
	case "stop":
		err = stopJob(client, d)
	default:
		return fmt.Errorf("invalid action type (%s)", actionType)
	}

	if err != nil {
		return err
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      jobStateRefreshFunc(client, workspaceId, jobName, processType),
		Timeout:      timeout,
		Delay:        10 * time.Second,
		PollInterval: 20 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for the job (%s) action to become completed: %s", jobName, err)
	}
	return nil
}

func resourceFactoryJobActionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		jobName = d.Get("job_name").(string)
	)
	client, err := cfg.NewServiceClient("dataarts-dlf", region)
	if err != nil {
		return diag.Errorf("error creating DataArts client: %s", err)
	}

	err = doActionJob(ctx, client, d, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("unable to operate status of the job (%s): %s", jobName, err)
	}

	d.SetId(jobName)

	return resourceFactoryJobActionRead(ctx, d, meta)
}

func resourceFactoryJobActionRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		workspaceId = d.Get("workspace_id").(string)
		jobName     = d.Get("job_name").(string)
		jobType     = d.Get("process_type").(string)
	)

	client, err := cfg.NewServiceClient("dataarts-dlf", region)
	if err != nil {
		return diag.Errorf("error creating DataArts client: %s", err)
	}

	job, err := getJobByName(client, workspaceId, jobName, jobType)
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "error_code", actionResourceNotFoundCodes...),
			fmt.Sprintf("job (%s) not found: %s", jobName, err))
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("job_name", utils.PathSearch("name", job, nil)),
		d.Set("status", utils.PathSearch("status", job, nil)),
		d.Set("action", jobType),
	)
	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error saving the fields of the job action: %s", err)
	}
	return nil
}

func resourceFactoryJobActionUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("dataarts-dlf", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating DataArts client: %s", err)
	}

	if d.HasChange("action") {
		jobName := d.Get("name").(string)
		err = doActionJob(ctx, client, d, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return diag.Errorf("error updating DataArts job (%s) status: %s", jobName, err)
		}
	}

	return resourceFactoryJobActionRead(ctx, d, meta)
}

func resourceFactoryJobActionDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is only a one-time action resource for changing job status. Deleting this resource will
not change the current status, but will only remove the resource information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
