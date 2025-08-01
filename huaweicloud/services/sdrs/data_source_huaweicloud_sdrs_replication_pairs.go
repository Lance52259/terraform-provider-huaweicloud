package sdrs

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

// @API SDRS GET /v1/{project_id}/replications
func DataSourceReplicationPairs() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceReplicationPairsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region where the SDRS replication pairs are located. If omitted, the provider-level region will be used.`,
			},
			"availability_zone": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the AZ of the current production site of the protection group.`,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the name of the replication pair. Supports fuzzy query.`,
			},
			"protected_instance_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the ID of the protected instance bound to the replication pair.`,
			},
			"protected_instance_ids": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the list of protected instance IDs (URL-encoded).`,
			},
			"query_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the query type of the current production site. Valid values: "status_abnormal", "general"`,
			},
			"server_group_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the protection group ID.`,
			},
			"server_group_ids": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the list of protection group IDs (URL-encoded).`,
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the status of the replication pair.`,
			},
			// Computed list
			"replication_pairs": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The replication pair list.`,
				Elem:        dataReplicationPairsSchema(),
			},
		},
	}
}

func dataReplicationPairsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the replication pair.`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the replication pair.`,
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The description of the replication pair.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The status of the replication pair.`,
			},
			"volume_ids": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The IDs of the volumes used by the replication pair.`,
			},
			"attachment": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The attachment information.`,
				Elem:        dataAttachmentSchema(),
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time of the replication pair.`,
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The update time of the replication pair.`,
			},
			"replication_model": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The replication model of the replication pair. Default: "hypermetro"`,
			},
			"progress": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The synchronization progress, in percentage.`,
			},
			"failure_detail": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The error code when the replication pair status is error.`,
			},
			"record_metadata": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The record metadata information.`,
				Elem:        dataRecordMetadataSchema(),
			},
			"server_group_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the protection group.`,
			},
			"fault_level": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The fault level of the replication pair.`,
			},
			"priority_station": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The current production site of the replication pair.`,
			},
			"replication_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The replication status of the replication pair.`,
			},
		},
	}
}

func dataAttachmentSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"protected_instance": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the protected instance.`,
			},
			"device": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The device name on the server.`,
			},
		},
	}
}

func dataRecordMetadataSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"multiattach": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether the volume supports multi-attach.`,
			},
			"bootable": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether the volume is a system disk.`,
			},
			"volume_size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The size of the volume in GB.`,
			},
			"volume_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The type of the volume.`,
			},
		},
	}
}

func buildReplicationPairsQueryParams(d *schema.ResourceData, offset int) string {
	res := ""

	if v, ok := d.GetOk("availability_zone"); ok {
		res = fmt.Sprintf("%s&availability_zone=%v", res, v)
	}

	if v, ok := d.GetOk("name"); ok {
		res = fmt.Sprintf("%s&name=%v", res, v)
	}

	if v, ok := d.GetOk("protected_instance_id"); ok {
		res = fmt.Sprintf("%s&protected_instance_id=%v", res, v)
	}

	if v, ok := d.GetOk("protected_instance_ids"); ok {
		res = fmt.Sprintf("%s&protected_instance_ids=%v", res, v)
	}

	if v, ok := d.GetOk("query_type"); ok {
		res = fmt.Sprintf("%s&query_type=%v", res, v)
	}

	if v, ok := d.GetOk("server_group_id"); ok {
		res = fmt.Sprintf("%s&server_group_id=%v", res, v)
	}

	if v, ok := d.GetOk("server_group_ids"); ok {
		res = fmt.Sprintf("%s&server_group_ids=%v", res, v)
	}

	if v, ok := d.GetOk("status"); ok {
		res = fmt.Sprintf("%s&status=%v", res, v)
	}

	if offset > 0 {
		res = fmt.Sprintf("%s&offset=%v", res, offset)
	}

	if res != "" {
		res = "?" + res[1:]
	}

	return res
}

func dataSourceReplicationPairsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg      = meta.(*config.Config)
		region   = cfg.GetRegion(d)
		product  = "sdrs"
		apiPath  = "v1/{project_id}/replications"
		offset   = 0
		allPairs []interface{}
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating SDRS client: %s", err)
	}

	basePath := strings.ReplaceAll(client.Endpoint+apiPath, "{project_id}", client.ProjectID)
	reqOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	for {
		eachPath := basePath + buildReplicationPairsQueryParams(d, offset)
		resp, err := client.Request("GET", eachPath, &reqOpts)
		if err != nil {
			return diag.Errorf("error retrieving SDRS replication pairs: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		pairs := utils.PathSearch("replications", respBody, make([]interface{}, 0)).([]interface{})
		if len(pairs) == 0 {
			break
		}
		allPairs = append(allPairs, pairs...)
		offset += len(pairs)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(id)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("replication_pairs", flattenDataReplicationPairs(allPairs)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenDataReplicationPairs(rawPairs []interface{}) []interface{} {
	result := make([]interface{}, 0, len(rawPairs))
	for _, v := range rawPairs {
		result = append(result, map[string]interface{}{
			"id":                 utils.PathSearch("id", v, nil),
			"name":               utils.PathSearch("name", v, nil),
			"description":        utils.PathSearch("description", v, nil),
			"status":             utils.PathSearch("status", v, nil),
			"volume_ids":         utils.PathSearch("volume_ids", v, nil),
			"attachment":         flattenAttachmentAttribute(v),
			"created_at":         utils.PathSearch("created_at", v, nil),
			"updated_at":         utils.PathSearch("updated_at", v, nil),
			"replication_model":  utils.PathSearch("replication_model", v, nil),
			"progress":           utils.PathSearch("progress", v, nil),
			"failure_detail":     utils.PathSearch("failure_detail", v, nil),
			"record_metadata":    flattenRecordMetadataAttribute(v),
			"server_group_id":    utils.PathSearch("server_group_id", v, nil),
			"fault_level":        utils.PathSearch("fault_level", v, nil),
			"priority_station":   utils.PathSearch("priority_station", v, nil),
			"replication_status": utils.PathSearch("replication_status", v, nil),
		})
	}
	return result
}

func flattenAttachmentAttribute(resp interface{}) []interface{} {
	rawArray := utils.PathSearch("attachment", resp, make([]interface{}, 0)).([]interface{})
	if len(rawArray) == 0 {
		return nil
	}

	res := make([]interface{}, 0, len(rawArray))
	for _, v := range rawArray {
		res = append(res, map[string]interface{}{
			"protected_instance": utils.PathSearch("protected_instance", v, nil),
			"device":             utils.PathSearch("device", v, nil),
		})
	}
	return res
}

func flattenRecordMetadataAttribute(resp interface{}) []interface{} {
	rawMap := utils.PathSearch("record_metadata", resp, nil)
	if rawMap == nil {
		return nil
	}
	return []interface{}{map[string]interface{}{
		"multiattach": utils.PathSearch("multiattach", rawMap, nil),
		"bootable":    utils.PathSearch("bootable", rawMap, nil),
		"volume_size": utils.PathSearch("volume_size", rawMap, nil),
		"volume_type": utils.PathSearch("volume_type", rawMap, nil),
	}}
}
