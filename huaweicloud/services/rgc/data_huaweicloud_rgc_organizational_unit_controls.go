package rgc

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API RGC GET /v1/governance/managed-organizational-units/{managed_organizational_unit_id}/controls
func DataSourceOrganizationalUnitControls() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceOrganizationalUnitControlsRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"managed_organizational_unit_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"control_summaries": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     organizationalUnitControlSchema(),
			},
		},
	}
}

func organizationalUnitControlSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"manage_account_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"control_identifier": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"control_objective": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"behavior": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"owner": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"regional_preference": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"guidance": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"service": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"implementation": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}

	return &sc
}

func dataSourceOrganizationalUnitControlsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	managedOuId := d.Get("managed_organizational_unit_id").(string)
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	listOrganizationalUnitControlHttpUrl := "v1/governance/managed-organizational-units/{managed_organizational_unit_id}/controls"
	listOrganizationalUnitControlProduct := "rgc"
	listOrganizationalUnitControlClient, err := cfg.NewServiceClient(listOrganizationalUnitControlProduct, region)
	if err != nil {
		return diag.Errorf("error creating RGC client: %s", err)
	}

	listOrganizationalUnitControlPath := listOrganizationalUnitControlClient.Endpoint + listOrganizationalUnitControlHttpUrl
	listOrganizationalUnitControlPath = strings.ReplaceAll(listOrganizationalUnitControlPath, "{managed_organizational_unit_id}", managedOuId)

	listOrganizationalUnitControlOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	var OrganizationalUnitControl []interface{}
	var marker string
	var queryPath string

	for {
		queryPath = listOrganizationalUnitControlPath + buildListOrganizationalUnitControlQueryParams(marker)
		listOrganizationalUnitControlResp, err := listOrganizationalUnitControlClient.Request("GET", queryPath, &listOrganizationalUnitControlOpt)
		if err != nil {
			return diag.Errorf("error retrieving RGC organizational unit controls: %s", err)
		}

		listOrganizationalUnitControlRespBody, err := utils.FlattenResponse(listOrganizationalUnitControlResp)
		if err != nil {
			return diag.FromErr(err)
		}

		onePageOrganizationalUnitOrganizationalUnitControl := FlattenOrganizationalUnitControlResp(listOrganizationalUnitControlRespBody)
		OrganizationalUnitControl = append(OrganizationalUnitControl, onePageOrganizationalUnitOrganizationalUnitControl...)
		marker = utils.PathSearch("page_info.next_marker", listOrganizationalUnitControlRespBody, "").(string)
		if marker == "" {
			break
		}
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(generateUUID)

	mErr := multierror.Append(nil,
		d.Set("control_summaries", OrganizationalUnitControl),
		d.Set("region", region),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildListOrganizationalUnitControlQueryParams(marker string) string {
	// the default value of limit is 200
	res := "?limit=200"

	if marker != "" {
		res = fmt.Sprintf("%s&marker=%v", res, marker)
	}

	return res
}

func FlattenOrganizationalUnitControlResp(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	curJson := utils.PathSearch("control_summaries", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"manage_account_id":   utils.PathSearch("manage_account_id", v, nil),
			"control_identifier":  utils.PathSearch("control_identifier", v, nil),
			"state":               utils.PathSearch("state", v, nil),
			"version":             utils.PathSearch("version", v, nil),
			"name":                utils.PathSearch("name", v, nil),
			"description":         utils.PathSearch("description", v, nil),
			"control_objective":   utils.PathSearch("control_objective", v, nil),
			"behavior":            utils.PathSearch("behavior", v, nil),
			"owner":               utils.PathSearch("owner", v, nil),
			"regional_preference": utils.PathSearch("regional_preference", v, nil),
			"guidance":            utils.PathSearch("guidance", v, nil),
			"service":             utils.PathSearch("service", v, nil),
			"implementation":      utils.PathSearch("implementation", v, nil),
		})
	}
	return rst
}
