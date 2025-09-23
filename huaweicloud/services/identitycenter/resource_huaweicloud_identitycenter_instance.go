package identitycenter

import (
	"context"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
	"strings"
	"time"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// @API IdentityCenter POST /v1/service/start
// @API IdentityCenter POST /v1/instances/{instance_id}/alias
// @API IdentityCenter GET  /v1/instances
// @API IdentityCenter GET  /v1/identity-center-service/status
// @API IdentityCenter POST /v1/service/delete
func ResourceIdentityCenterInstance() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIdentityCenterInstanceCreate,
		UpdateContext: resourceIdentityCenterInstanceUpdate,
		ReadContext:   resourceIdentityCenterInstanceRead,
		DeleteContext: resourceIdentityCenterInstanceDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"identity_store_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"alias": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"instance_urn": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceIdentityCenterInstanceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		startIdentityCenterHttpUrl            = "v1/service/start"
		getIdentityCenterServiceStatusHttpUrl = "v1/identity-center-service/status"
		createAliasHttpUrl                    = "v1/instances/{instance_id}/alias"
		listInstancesHttpUrl                  = "v1/instances"
		identityCenterProduct                 = "identitycenter"
	)
	client, err := cfg.NewServiceClient(identityCenterProduct, region)
	if err != nil {
		return diag.Errorf("error creating IdentityCenter client: %s", err)
	}

	startIdentityCenterPath := client.Endpoint + startIdentityCenterHttpUrl

	startIdentityCenterOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	_, err = client.Request("POST", startIdentityCenterPath, &startIdentityCenterOpt)
	if err != nil {
		return diag.Errorf("error creating IdentityCenter instance: %s", err)
	}

	uuId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuId)

	stateConf := &resource.StateChangeConf{
		Pending: []string{"disable"},
		Target:  []string{"enable"},
		Refresh: identityCenterServiceStatusRefreshFunc(getIdentityCenterServiceStatusHttpUrl,
			"serviceStatus", client),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        1 * time.Second,
		PollInterval: 5 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for creating IAM Identity Center instance to complete: %s", err)
	}

	listInstancesClient, err := cfg.NewServiceClient(identityCenterProduct, region)
	if err != nil {
		return diag.Errorf("error creating Identity Center Client: %s", err)
	}

	listInstancesPath := listInstancesClient.Endpoint + listInstancesHttpUrl

	listInstancesOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	listInstancesResp, err := listInstancesClient.Request("GET",
		listInstancesPath, &listInstancesOpt)

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving Identity Center Instance")
	}

	listInstancesRespBody, err := utils.FlattenResponse(listInstancesResp)
	if err != nil {
		return diag.FromErr(err)
	}

	instances, _ := utils.PathSearch("instances", listInstancesRespBody, make([]interface{}, 0)).([]interface{})
	if len(instances) == 0 {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "no instance found.")
	}

	instanceId := utils.PathSearch("instances|[0].instance_id", listInstancesRespBody, nil).(string)

	d.SetId(instanceId)

	if _, ok := d.GetOk("alias"); ok {
		createAliasClient, err := cfg.NewServiceClient(identityCenterProduct, region)
		if err != nil {
			return diag.Errorf("error creating IdentityCenter client: %s", err)
		}

		createAliasPath := createAliasClient.Endpoint + createAliasHttpUrl
		createAliasPath = strings.ReplaceAll(createAliasPath, "{instance_id}", instanceId)
		createAliasOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			JSONBody:         utils.RemoveNil(buildCreateAliasBodyParams(d)),
		}
		_, err = client.Request("POST", createAliasPath, &createAliasOpt)
		if err != nil {
			return diag.Errorf("error creating IdentityCenter instance alias: %s", err)
		}
	}

	return resourceIdentityCenterInstanceRead(ctx, d, meta)
}

func buildCreateAliasBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"alias": utils.ValueIgnoreEmpty(d.Get("alias")),
	}
	return bodyParams
}

func identityCenterServiceStatusRefreshFunc(getRequestStatusHttpUrl, searchExpression string,
	client *golangsdk.ServiceClient) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		getIdentityCenterServiceStatusPath := client.Endpoint + getRequestStatusHttpUrl

		getIdentityCenterServiceStatusOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
		}
		getIdentityCenterServiceStatusResp, err := client.Request("GET", getIdentityCenterServiceStatusPath, &getIdentityCenterServiceStatusOpt)
		if err != nil {
			return nil, "", err
		}
		getIdentityCenterServiceStatusRespBody, err := utils.FlattenResponse(getIdentityCenterServiceStatusResp)
		if err != nil {
			return nil, "", err
		}

		status := utils.PathSearch(searchExpression, getIdentityCenterServiceStatusRespBody, "").(string)
		if status == "enable" {
			return getIdentityCenterServiceStatusRespBody, status, nil
		}
		return getIdentityCenterServiceStatusRespBody, "disable", nil
	}
}

func resourceIdentityCenterInstanceRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		listInstancesHttpUrl = "v1/instances"
		listInstancesProduct = "identitycenter"
	)

	var mErr *multierror.Error

	listInstancesClient, err := cfg.NewServiceClient(listInstancesProduct, region)
	if err != nil {
		return diag.Errorf("error creating Identity Center Client: %s", err)
	}

	listInstancesPath := listInstancesClient.Endpoint + listInstancesHttpUrl

	listInstancesOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	listInstancesResp, err := listInstancesClient.Request("GET",
		listInstancesPath, &listInstancesOpt)

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving Identity Center Instance")
	}

	listInstancesRespBody, err := utils.FlattenResponse(listInstancesResp)
	if err != nil {
		return diag.FromErr(err)
	}

	instances, _ := utils.PathSearch("instances", listInstancesRespBody, make([]interface{}, 0)).([]interface{})
	if len(instances) == 0 {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "no instance found.")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("identity_store_id", utils.PathSearch("instances|[0].identity_store_id",
			listInstancesRespBody, nil)),
		d.Set("instance_urn", utils.PathSearch("instances|[0].instance_urn",
			listInstancesRespBody, nil)),
		d.Set("alias", utils.PathSearch("instances|[0].alias",
			listInstancesRespBody, nil)),
		d.Set("instance_id", utils.PathSearch("instances|[0].instance_id",
			listInstancesRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceIdentityCenterInstanceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		createAliasHttpUrl    = "v1/instances/{instance_id}/alias"
		listInstancesHttpUrl  = "v1/instances"
		identityCenterProduct = "identitycenter"
	)

	if d.HasChange("alias") {
		listInstancesClient, err := cfg.NewServiceClient(identityCenterProduct, region)
		if err != nil {
			return diag.Errorf("error creating Identity Center Client: %s", err)
		}

		listInstancesPath := listInstancesClient.Endpoint + listInstancesHttpUrl

		listInstancesOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
		}
		listInstancesResp, err := listInstancesClient.Request("GET",
			listInstancesPath, &listInstancesOpt)

		if err != nil {
			return common.CheckDeletedDiag(d, err, "error retrieving Identity Center Instance")
		}

		listInstancesRespBody, err := utils.FlattenResponse(listInstancesResp)
		if err != nil {
			return diag.FromErr(err)
		}

		instances, _ := utils.PathSearch("instances", listInstancesRespBody, make([]interface{}, 0)).([]interface{})
		if len(instances) == 0 {
			return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "no instance found.")
		}

		instanceId := utils.PathSearch("instances|[0].instance_id", listInstancesRespBody, nil).(string)

		createAliasClient, err := cfg.NewServiceClient(identityCenterProduct, region)
		if err != nil {
			return diag.Errorf("error creating IdentityCenter client: %s", err)
		}

		createAliasPath := createAliasClient.Endpoint + createAliasHttpUrl
		createAliasPath = strings.ReplaceAll(createAliasPath, "{instance_id}", instanceId)
		createAliasOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			JSONBody:         utils.RemoveNil(buildCreateAliasBodyParams(d)),
		}

		_, err = createAliasClient.Request("POST", createAliasPath, &createAliasOpt)

		if err != nil {
			return diag.Errorf("error creating IdentityCenter instance alias: %s", err)
		}
	}

	return resourceIdentityCenterInstanceRead(ctx, d, meta)
}

func resourceIdentityCenterInstanceDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		deleteIdentityCenterHttpUrl = "v1/service/delete"
		deleteIdentityCenterProduct = "identitycenter"
	)
	deleteIdentityCenterClient, err := cfg.NewServiceClient(deleteIdentityCenterProduct, region)
	if err != nil {
		return diag.Errorf("error creating IdentityCenter client: %s", err)
	}

	deleteIdentityCenterPath := deleteIdentityCenterClient.Endpoint + deleteIdentityCenterHttpUrl

	deleteIdentityCenterOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	_, err = deleteIdentityCenterClient.Request("POST", deleteIdentityCenterPath, &deleteIdentityCenterOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting IdentityCenter instance")
	}

	return nil
}
