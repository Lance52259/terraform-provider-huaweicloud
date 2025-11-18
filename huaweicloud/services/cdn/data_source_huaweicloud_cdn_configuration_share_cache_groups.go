package cdn

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CDN GET /v1.0/cdn/configuration/share-cache-groups
func DataSourceConfigurationShareCacheGroups() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceConfigurationShareCacheGroupsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the share cache groups are located.`,
			},

			// Attributes.
			"share_cache_groups": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        shareCacheGroupSchema(),
				Description: `The list of the share cache groups.`,
			},
		},
	}
}

func shareCacheGroupSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the share cache group.`,
			},
			"group_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the share cache group.`,
			},
			"primary_domain": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The primary domain name.`,
			},
			"share_cache_records": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of associated domain names.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"domain_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The associated domain name.`,
						},
					},
				},
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time of the share cache group, in RFC3339 format.`,
			},
		},
	}
}

func flattenShareCacheGroups(groups []interface{}) []map[string]interface{} {
	if len(groups) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(groups))
	for _, group := range groups {
		result = append(result, map[string]interface{}{
			"id":             utils.PathSearch("id", group, nil),
			"group_name":     utils.PathSearch("group_name", group, nil),
			"primary_domain": utils.PathSearch("primary_domain", group, nil),
			"share_cache_records": flattenShareCacheRecords(utils.PathSearch("share_cache_records", group,
				make([]interface{}, 0)).([]interface{})),
			"create_time": utils.FormatTimeStampRFC3339(int64(utils.PathSearch("create_time", group, float64(0)).(float64))/1000, false),
		})
	}

	return result
}

func dataSourceConfigurationShareCacheGroupsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)

	client, err := cfg.NewServiceClient("cdn", "")
	if err != nil {
		return diag.Errorf("error creating CDN client: %s", err)
	}

	groups, err := listShareCacheGroups(client)
	if err != nil {
		return diag.Errorf("error querying CDN share cache groups: %s", err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	mErr := multierror.Append(nil,
		d.Set("region", cfg.GetRegion(d)),
		d.Set("share_cache_groups", flattenShareCacheGroups(groups)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
