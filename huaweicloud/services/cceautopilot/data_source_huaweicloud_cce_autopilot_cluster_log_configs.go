package cceautopilot

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

func DataSourceCceAutopilotClusterLogConfigs() *schema.Resource {
    return &schema.Resource{
        ReadContext: dataSourceCceAutopilotClusterLogConfigsRead,

        Schema: map[string]*schema.Schema{
            "region": {
                Type:        schema.TypeString,
                Optional:    true,
                Computed:    true,
                Description: "Specifies the region in which to query the resource. If omitted, the provider-level region will be used.",
            },
            "cluster_id": {
                Type:        schema.TypeString,
                Required:    true,
                Description: "Cluster ID, which can be obtained from the parameter in the interface URI.",
            },
			"project_id": {
                Type:        schema.TypeString,
                Required:    true,
                Description: "project ID, which can be obtained from the parameter in the interface URI.",
            },
            "type": {
                Type:        schema.TypeString,
                Optional:    true,
                Description: "Component Type.",
            },
            "ttl_in_days": {
                Type:        schema.TypeInt,
                Computed:    true,
                Description: "Storage duration.",
            },
            "log_configs": {
                Type:        schema.TypeList,
                Computed:    true,
                Description: `all log configs`,
                Elem: &schema.Resource{
                    Schema: map[string]*schema.Schema{
                        "name": {
                            Type:        schema.TypeString,
                            Computed:    true,
                            Description: `log service name`,
                        },
                        "enable": {
                            Type:        schema.TypeBool,
                            Computed:    true,
                            Description: `Collected or Not`,
                        },
                    },
                },
            },
        },
    }
}

type CceAutopilotClusterLogConfigsDSWrapper struct {
    *schemas.ResourceDataWrapper
    Config *config.Config
}

func newCceAutopilotClusterLogConfigsDSWrapper(d *schema.ResourceData, meta interface{}) *CceAutopilotClusterLogConfigsDSWrapper {
    return &CceAutopilotClusterLogConfigsDSWrapper{
        ResourceDataWrapper: schemas.NewSchemaWrapper(d),
        Config:              meta.(*config.Config),
    }
}

func dataSourceCceAutopilotClusterLogConfigsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
    wrapper := newCceAutopilotClusterLogConfigsDSWrapper(d, meta)
    showAutCluLogConfRst, err := wrapper.ShowAutopilotClusterLogConfigs()
    if err != nil {
        return diag.FromErr(err)
    }

    id, err := uuid.GenerateUUID()
    if err != nil {
        return diag.FromErr(err)
    }
    d.SetId(id)

    err = wrapper.showAutopilotClusterLogConfigsToSchema(showAutCluLogConfRst)
    if err != nil {
        return diag.FromErr(err)
    }

    return nil
}

// @API CCE GET /autopilot/v3/projects/{project_id}/cluster/{cluster_id}/log-configs
func (w *CceAutopilotClusterLogConfigsDSWrapper) ShowAutopilotClusterLogConfigs() (*gjson.Result, error) {
    client, err := w.NewClient(w.Config, "cce")
    if err != nil {
        return nil, err
    }

    uri := "/autopilot/v3/projects/{project_id}/cluster/{cluster_id}/log-configs"
    uri = strings.ReplaceAll(uri, "{project_id}", w.Get("project_id").(string))
    uri = strings.ReplaceAll(uri, "{cluster_id}", w.Get("cluster_id").(string))

    params := map[string]any{
        "type":  w.Get("type"),
    }
    
    params = utils.RemoveNil(params)
    return httphelper.New(client).
        Method("GET").
        URI(uri).
        Query(params).
        Request().
        Result()
}

func (w *CceAutopilotClusterLogConfigsDSWrapper) showAutopilotClusterLogConfigsToSchema(body *gjson.Result) error {
    d := w.ResourceData
    mErr := multierror.Append(nil,
        d.Set("region", w.Config.GetRegion(w.ResourceData)),
    )
    return mErr.ErrorOrNil()
}
