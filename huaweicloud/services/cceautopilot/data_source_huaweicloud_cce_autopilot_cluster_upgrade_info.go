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
)

func DataSourceCceAutopilotClusterUpgradeInfo() *schema.Resource {
    return &schema.Resource{
        ReadContext: dataSourceCceAutopilotClusterUpgradeInfoRead,

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
            "api_version": {
                Type:        schema.TypeString,
                Computed:    true,
                Description: "API version.",
            },
            "annotations": {
                Type:        schema.TypeMap,
                Computed:    true,
                Description: "Resource annotations, which consist of key/value pairs.",
            },
            "creation_timestamp": {
                Type:        schema.TypeString,
                Computed:    true,
                Description: "Creation time.",
            },
            "labels": {
                Type:        schema.TypeMap,
                Computed:    true,
                Description: "Resource labels, in key/value pair format. These are reserved fields in the interface and will not take effect.",
            },
            "name": {
                Type:        schema.TypeString,
                Computed:    true,
                Description: "Resource name.",
            },
            "uid": {
                Type:        schema.TypeString,
                Computed:    true,
                Description: "Unique ID identifier.",
            },
            "update_timestamp": {
                Type:        schema.TypeString,
                Computed:    true,
                Description: "Update time.",
            },
            "kind": {
                Type:        schema.TypeString,
                Computed:    true,
                Description: "Type.",
            },
        },
    }
}

type CceAutopilotClusterUpgradeInfoDSWrapper struct {
    *schemas.ResourceDataWrapper
    Config *config.Config
}

func newCceAutopilotClusterUpgradeInfoDSWrapper(d *schema.ResourceData, meta interface{}) *CceAutopilotClusterUpgradeInfoDSWrapper {
    return &CceAutopilotClusterUpgradeInfoDSWrapper{
        ResourceDataWrapper: schemas.NewSchemaWrapper(d),
        Config:              meta.(*config.Config),
    }
}

func dataSourceCceAutopilotClusterUpgradeInfoRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
    wrapper := newCceAutopilotClusterUpgradeInfoDSWrapper(d, meta)
    showAutCluUpgInfRst, err := wrapper.ShowAutopilotClusterUpgradeInfo()
    if err != nil {
        return diag.FromErr(err)
    }

    id, err := uuid.GenerateUUID()
    if err != nil {
        return diag.FromErr(err)
    }
    d.SetId(id)

    err = wrapper.showAutopilotClusterUpgradeInfoToSchema(showAutCluUpgInfRst)
    if err != nil {
        return diag.FromErr(err)
    }

    return nil
}

// @API CCE GET /autopilot/v3/projects/{project_id}/clusters/{cluster_id}/upgradeinfo
func (w *CceAutopilotClusterUpgradeInfoDSWrapper) ShowAutopilotClusterUpgradeInfo() (*gjson.Result, error) {
    client, err := w.NewClient(w.Config, "cce")
    if err != nil {
        return nil, err
    }

    uri := "/autopilot/v3/projects/{project_id}/clusters/{cluster_id}/upgradeinfo"
    uri = strings.ReplaceAll(uri, "{project_id}", w.Get("project_id").(string))
    uri = strings.ReplaceAll(uri, "{cluster_id}", w.Get("cluster_id").(string))
    return httphelper.New(client).
        Method("GET").
        URI(uri).
        Request().
        Result()
}

func (w *CceAutopilotClusterUpgradeInfoDSWrapper) showAutopilotClusterUpgradeInfoToSchema(body *gjson.Result) error {
    d := w.ResourceData
    mErr := multierror.Append(nil,
        d.Set("region", w.Config.GetRegion(w.ResourceData)),
        d.Set("api_version", body.Get("apiVersion").Value()),
        d.Set("annotations", schemas.MapToStrMap(body.Get("metadata.annotations"))),
        d.Set("creation_timestamp", body.Get("metadata.creationTimestamp").Value()),
        d.Set("labels", schemas.MapToStrMap(body.Get("metadata.labels"))),
        d.Set("name", body.Get("metadata.name").Value()),
        d.Set("uid", body.Get("metadata.uid").Value()),
        d.Set("update_timestamp", body.Get("metadata.updateTimestamp").Value()),
        d.Set("kind", body.Get("kind").Value()),
	)
    return mErr.ErrorOrNil()
}
