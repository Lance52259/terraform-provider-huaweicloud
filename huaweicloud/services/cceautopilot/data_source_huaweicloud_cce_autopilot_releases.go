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

func DataSourceCCEAutopilotReleases() *schema.Resource {
    return &schema.Resource{
        ReadContext: dataSourceCCEAutopilotReleasesRead,

        Schema: map[string]*schema.Schema{
            "region": {
                Type:        schema.TypeString,
                Optional:    true,
                Computed:    true,
                Description: "Specifies the region in which to query the resource. If omitted, the provider-level region will be used.",
            },
            "chart_id": {
                Type:        schema.TypeString,
                Optional:    true,
                Description: "模板ID",
            },
            "cluster_id": {
                Type:        schema.TypeString,
                Required:    true,
                Description: "集群ID,获取方式请参见[如何获取接口URI中参数](cce_02_0271.xml)。",
            },
            "namespace": {
                Type:        schema.TypeString,
                Optional:    true,
                Description: "模板对应的命名空间",
            },
        },
    }
}

type CCEAutopilotReleasesDSWrapper struct {
    *schemas.ResourceDataWrapper
    Config *config.Config
}

func newCCEAutopilotReleasesDSWrapper(d *schema.ResourceData, meta interface{}) *CCEAutopilotReleasesDSWrapper {
    return &CCEAutopilotReleasesDSWrapper{
        ResourceDataWrapper: schemas.NewSchemaWrapper(d),
        Config:              meta.(*config.Config),
    }
}

func dataSourceCCEAutopilotReleasesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
    wrapper := newCCEAutopilotReleasesDSWrapper(d, meta)
    lisAutRelRst, err := wrapper.ListAutopilotReleases()
    if err != nil {
        return diag.FromErr(err)
    }

    id, err := uuid.GenerateUUID()
    if err != nil {
        return diag.FromErr(err)
    }
    d.SetId(id)

    err = wrapper.listAutopilotReleasesToSchema(lisAutRelRst)
    if err != nil {
        return diag.FromErr(err)
    }

    return nil
}

// @API CCE GET /autopilot/cam/v3/clusters/{cluster_id}/releases
func (w *CCEAutopilotReleasesDSWrapper) ListAutopilotReleases() (*gjson.Result, error) {
    client, err := w.NewClient(w.Config, "cce")
    if err != nil {
        return nil, err
    }

    uri := "/autopilot/cam/v3/clusters/{cluster_id}/releases"
    uri = strings.ReplaceAll(uri, "{cluster_id}", w.Get("cluster_id").(string))
    params := map[string]any{
        "chart_id":  w.Get("chart_id"),
        "namespace": w.Get("namespace"),
    }

    params = utils.RemoveNil(params)
    return httphelper.New(client).
        Method("GET").
        URI(uri).
        Query(params).
        Request().
        Result()
}

func (w *CCEAutopilotReleasesDSWrapper) listAutopilotReleasesToSchema(body *gjson.Result) error {
    d := w.ResourceData
    mErr := multierror.Append(nil,
        d.Set("region", w.Config.GetRegion(w.ResourceData)),
    )
    return mErr.ErrorOrNil()
}
