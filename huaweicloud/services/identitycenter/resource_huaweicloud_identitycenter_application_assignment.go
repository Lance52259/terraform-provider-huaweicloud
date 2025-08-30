package identitycenter

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var identityCenterApplicationAssignmentNonUpdateParams = []string{"instance_id", "application_instance_id", "principal_id", "principal_type"}

// @API IdentityCenter POST /v1/instances/{instance_id}/applications/{application_instance_id}/assignments/create
// @API IdentityCenter GET /v1/instances/{instance_id}/applications/{application_instance_id}/assignments
// @API IdentityCenter POST /v1/instances/{instance_id}/applications/{application_instance_id}/assignments/delete
func ResourceIdentityCenterApplicationAssignment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIdentityCenterApplicationAssignmentCreate,
		ReadContext:   resourceIdentityCenterApplicationAssignmentRead,
		DeleteContext: resourceIdentityCenterApplicationAssignmentDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceIdentityCenterApplicationAssignmentImport,
		},

		CustomizeDiff: config.FlexibleForceNew(identityCenterApplicationAssignmentNonUpdateParams),

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"application_instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"principal_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"principal_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"application_urn": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceIdentityCenterApplicationAssignmentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		createApplicationAssignmentHttpUrl = "v1/instances/{instance_id}/applications/{application_instance_id}/assignments/create"
		createApplicationAssignmentProduct = "identitycenter"
	)

	createApplicationAssignmentClient, err := cfg.NewServiceClient(createApplicationAssignmentProduct, region)
	if err != nil {
		return diag.Errorf("error creating Identity Center Client: %s", err)
	}

	createApplicationAssignmentPath := createApplicationAssignmentClient.Endpoint + createApplicationAssignmentHttpUrl
	createApplicationAssignmentPath = strings.ReplaceAll(createApplicationAssignmentPath, "{instance_id}",
		fmt.Sprintf("%v", d.Get("instance_id")))
	createApplicationAssignmentPath = strings.ReplaceAll(createApplicationAssignmentPath, "{application_instance_id}",
		fmt.Sprintf("%v", d.Get("application_instance_id")))

	createApplicationAssignmentOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateApplicationAssignmentBodyParams(d)),
	}

	_, err = createApplicationAssignmentClient.Request("POST",
		createApplicationAssignmentPath, &createApplicationAssignmentOpt)
	if err != nil {
		return diag.Errorf("error creating Identity Center application asssignment: %s", err)
	}

	instanceId := d.Get("instance_id")
	applicationInstanceId := d.Get("application_instance_id")
	principalId := d.Get("principal_id")
	d.SetId(fmt.Sprintf("%v/%v/%v", instanceId, applicationInstanceId, principalId))

	return resourceIdentityCenterApplicationAssignmentRead(ctx, d, meta)
}

func buildCreateApplicationAssignmentBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"principal_id":   utils.ValueIgnoreEmpty(d.Get("principal_id")),
		"principal_type": utils.ValueIgnoreEmpty(d.Get("principal_type")),
	}
	return bodyParams
}

func resourceIdentityCenterApplicationAssignmentRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		listApplicationAssignmentsHttpUrl = "v1/instances/{instance_id}/applications/{application_instance_id}/assignments"
		listApplicationAssignmentsProduct = "identitycenter"
	)

	listApplicationAssignmentsClient, err := cfg.NewServiceClient(listApplicationAssignmentsProduct, region)
	if err != nil {
		return diag.Errorf("error creating Identity Center client: %s", err)
	}

	listApplicationAssignmentsPath := listApplicationAssignmentsClient.Endpoint + listApplicationAssignmentsHttpUrl
	listApplicationAssignmentsPath = strings.ReplaceAll(listApplicationAssignmentsPath, "{instance_id}",
		fmt.Sprintf("%v", d.Get("instance_id")))
	listApplicationAssignmentsPath = strings.ReplaceAll(listApplicationAssignmentsPath, "{application_instance_id}",
		fmt.Sprintf("%v", d.Get("application_instance_id")))

	listApplicationAssignmentsPath = listApplicationAssignmentsPath + buildListApplicationAssignmentsQueryParams("")

	listApplicationAssignmentsOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	var assignment interface{}
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 3 {
		return diag.Errorf("invalid format of Id: %s", d.Id())
	}
	principalId := parts[2]

listApplicationAssignmentsLoop:
	for {
		listApplicationAssignmentsResp, err := listApplicationAssignmentsClient.Request("GET", listApplicationAssignmentsPath, &listApplicationAssignmentsOpt)
		if err != nil {
			return common.CheckDeletedDiag(d, err, "error retrieving Identity Center application assignments.")
		}

		listApplicationAssignmentsRespBody, err := utils.FlattenResponse(listApplicationAssignmentsResp)
		if err != nil {
			return diag.FromErr(err)
		}

		assignments := utils.PathSearch("application_assignments", listApplicationAssignmentsRespBody, nil)

		if assignments == nil {
			return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "")
		}

		for _, v := range assignments.([]interface{}) {
			if principalId == utils.PathSearch("principal_id", v, "") {
				assignment = v
				break listApplicationAssignmentsLoop
			}
		}
		marker := utils.PathSearch("page_info.next_marker", listApplicationAssignmentsRespBody, nil)
		if marker == nil {
			break
		}
		listApplicationAssignmentsPath = listApplicationAssignmentsPath + buildListApplicationAssignmentsQueryParams(marker.(string))
	}
	if assignment == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("principal_id", utils.PathSearch("principal_id", assignment, nil)),
		d.Set("principal_type", utils.PathSearch("principal_type", assignment, nil)),
		d.Set("application_urn", utils.PathSearch("application_urn", assignment, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildListApplicationAssignmentsQueryParams(marker string) string {
	res := "?limit=100"

	if marker != "" {
		res = fmt.Sprintf("%s&marker=%v", res, marker)
	}
	return res
}

func resourceIdentityCenterApplicationAssignmentDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		deleteApplicationAssignmentHttpUrl = "v1/instances/{instance_id}/applications/{application_instance_id}/assignments/delete"
		deleteApplicationAssignmentProduct = "identitycenter"
	)
	deleteApplicationAssignmentClient, err := cfg.NewServiceClient(deleteApplicationAssignmentProduct, region)
	if err != nil {
		return diag.Errorf("error creating Identity Center Client: %s", err)
	}

	deleteApplicationAssignmentPath := deleteApplicationAssignmentClient.Endpoint + deleteApplicationAssignmentHttpUrl
	deleteApplicationAssignmentPath = strings.ReplaceAll(deleteApplicationAssignmentPath, "{instance_id}",
		fmt.Sprintf("%v", d.Get("instance_id")))
	deleteApplicationAssignmentPath = strings.ReplaceAll(deleteApplicationAssignmentPath, "{application_instance_id}",
		fmt.Sprintf("%v", d.Get("application_instance_id")))

	deleteApplicationAssignmentOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateApplicationAssignmentBodyParams(d)),
	}

	_, err = deleteApplicationAssignmentClient.Request("POST", deleteApplicationAssignmentPath,
		&deleteApplicationAssignmentOpt)
	if err != nil {
		return diag.Errorf("error deleting Identity Center application assignment: %s", err)
	}

	return nil
}

func resourceIdentityCenterApplicationAssignmentImport(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid id format, must be " +
			"<instance_id>/<application_instance_id>/<principal_id>")
	}

	d.SetId(fmt.Sprintf("%s/%s/%s", parts[0], parts[1], parts[2]))
	mErr := multierror.Append(nil,
		d.Set("instance_id", parts[0]),
		d.Set("application_instance_id", parts[1]),
		d.Set("principal_id", parts[2]),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return nil, fmt.Errorf("failed to set value to state when import, %s", err)
	}
	return []*schema.ResourceData{d}, nil
}
