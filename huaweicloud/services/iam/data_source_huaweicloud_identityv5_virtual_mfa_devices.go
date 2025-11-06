package iam

import (
	"context"
	"fmt"

	"github.com/chnsz/golangsdk"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API GET /v5/mfa-devices
const limit = 100

func DataSourceIdentityV5VirtualMfaDevices() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIdentityV5VirtualMfaDevicesRead,

		Schema: map[string]*schema.Schema{
			"user_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID of the IAM user.",
			},
			"marker": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"devices": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"serial_number": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"user_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceIdentityV5VirtualMfaDevicesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("iam", region)
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	var allDevices []interface{}
	for {
		listOpts := buildListMfaDevicesV5Params(d)
		resp, err := listMfaDevicesV5(client, listOpts)
		if err != nil {
			return diag.Errorf("error retrieving virtual MFA devices: %s", err)
		}

		devices := flattenListMfaDevicesV5Response(resp)
		allDevices = append(allDevices, devices...)

		nextMarker := utils.PathSearch("page_info.next_marker", resp, "").(string)
		if nextMarker == "" {
			break
		}
		if err = d.Set("marker", nextMarker); err != nil {
			return diag.Errorf("error setting nextmarker: %s", err)
		}
	}
	id, _ := uuid.GenerateUUID()
	d.SetId(id)
	mErr := multierror.Append(nil,
		d.Set("devices", allDevices),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting virtual MFA devices fields: %s", err)
	}
	return nil
}

func buildListMfaDevicesV5Params(d *schema.ResourceData) string {
	res := ""
	if v, ok := d.GetOk("user_id"); ok {
		res = fmt.Sprintf("%s&user_id=%v", res, v)
	}
	if v, ok := d.GetOk("marker"); ok {
		res = fmt.Sprintf("%s&marker=%v", res, v)
		res += fmt.Sprintf("limit=%v", limit)
	}
	if res != "" {
		res = "?" + res[1:]
	}
	return res
}

func listMfaDevicesV5(client *golangsdk.ServiceClient, param string) (interface{}, error) {
	url := fmt.Sprintf("%sv5/mfa-devices%s", client.Endpoint, param)
	reqOpt := &golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	r, err := client.Request("GET", url, reqOpt)
	if err != nil {
		return nil, err
	}
	respBody, err := utils.FlattenResponse(r)
	return respBody, err
}

func flattenListMfaDevicesV5Response(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	devices := utils.PathSearch("mfa_devices", resp, make([]interface{}, 0)).([]interface{})
	result := make([]interface{}, len(devices))
	for i, device := range devices {
		deviceMap := device.(map[string]interface{})
		result[i] = map[string]interface{}{
			"enabled":       utils.PathSearch("enabled", deviceMap, nil),
			"serial_number": utils.PathSearch("serial_number", deviceMap, nil),
			"user_id":       utils.PathSearch("user_id", deviceMap, nil),
		}
	}
	return result
}
