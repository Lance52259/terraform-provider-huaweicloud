package cceautopilot

import (
	"context"

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

func DataSourceCceAutopilotAddonInstances() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCceAutopilotAddonInstancesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
			},
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the cluster.`,
			},
			"api_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `API版本,固定值“v3”,该值不可修改。`,
			},
			"items": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `插件实例列表`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"api_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `API版本,固定值“v3”,该值不可修改。`,
						},
						"kind": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `API类型,固定值“Addon”,该值不可修改。`,
						},
						"metadata": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `基本信息,为集合类的元素类型,包含一组由不同名称定义的属性`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"alias": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `插件别名`,
									},
									"annotations": {
										Type:        schema.TypeMap,
										Computed:    true,
										Description: `插件注解,由key/value组成`,
									},
									"creation_timestamp": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `创建时间`,
									},
									"labels": {
										Type:        schema.TypeMap,
										Computed:    true,
										Description: `插件标签,key/value对格式,接口保留字段,填写不会生效`,
									},
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `插件名称`,
									},
									"uid": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `唯一id标识`,
									},
									"update_timestamp": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `更新时间`,
									},
								},
							},
						},
						"spec": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `spec是集合类的元素类型,内容为插件实例具体信息,实例的详细描述主体部分都在spec中给出`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"addon_template_labels": {
										Type:        schema.TypeList,
										Computed:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: `插件模板所属类型`,
									},
									"addon_template_logo": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `插件模板logo图片的地址`,
									},
									"addon_template_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `插件模板名称,如coredns`,
									},
									"addon_template_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `插件模板类型`,
									},
									"cluster_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `集群id`,
									},
									"description": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `插件模板描述`,
									},
									"values": {
										Type:        schema.TypeMap,
										Computed:    true,
										Description: `插件模板安装参数(各插件不同),请根据具体插件模板信息填写安装参数。`,
									},
									"version": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `插件模板版本号,如1.0.0`,
									},
								},
							},
						},
						"status": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `插件实例状态`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"reason": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `插件安装失败原因`,
									},
									"current_version": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: `当前插件实例使用的具体插件版本信息`,
										Elem:        itemsStatusCurrentVersionElem(),
									},
									"is_rollbackable": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: `是否支持回滚到插件升级前的插件版本`,
									},
									"message": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `安装错误详情`,
									},
									"previous_version": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `插件升级或回滚前的版本`,
									},
									"status": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `插件实例状态, 取值如下`,
									},
									"target_versions": {
										Type:        schema.TypeList,
										Computed:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: `此插件版本,支持升级的集群版本`,
									},
								},
							},
						},
					},
				},
			},
			"kind": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `API类型,固定值“Addon”,该值不可修改。`,
			},
		},
	}
}

// itemsStatusCurrentVersionElem
// The Elem of "items.status.currentVersion"
func itemsStatusCurrentVersionElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"creation_timestamp": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `创建时间`,
			},
			"input": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: `插件安装参数`,
			},
			"stable": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `是否为稳定版本`,
			},
			"support_versions": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `支持集群版本号`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"category": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `作用的集群类型`,
						},
						"cluster_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `支持的集群类型`,
						},
						"cluster_version": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `支持的集群版本(正则表达式)`,
						},
					},
				},
			},
			"translate": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: `供界面使用的翻译信息`,
			},
			"update_timestamp": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `更新时间`,
			},
			"version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `插件版本号`,
			},
		},
	}
}

type AddonInstancesDSWrapper struct {
	*schemas.ResourceDataWrapper
	Config *config.Config
}

func newAddonInstancesDSWrapper(d *schema.ResourceData, meta interface{}) *AddonInstancesDSWrapper {
	return &AddonInstancesDSWrapper{
		ResourceDataWrapper: schemas.NewSchemaWrapper(d),
		Config:              meta.(*config.Config),
	}
}

func dataSourceCceAutopilotAddonInstancesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	wrapper := newAddonInstancesDSWrapper(d, meta)
	lisAutAddInsRst, err := wrapper.ListAutopilotAddonInstances()
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	err = wrapper.listAutopilotAddonInstancesToSchema(lisAutAddInsRst)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// @API CCE GET /autopilot/v3/addons
func (w *AddonInstancesDSWrapper) ListAutopilotAddonInstances() (*gjson.Result, error) {
	client, err := w.NewClient(w.Config, "cce")
	if err != nil {
		return nil, err
	}

	uri := "/autopilot/v3/addons"
	params := map[string]any{
		"cluster_id": w.Get("cluster_id"),
	}
	params = utils.RemoveNil(params)
	return httphelper.New(client).
		Method("GET").
		URI(uri).
		Query(params).
		OffsetPager("instances", "offset", "limit", 100).
		Request().
		Result()
}

func (w *AddonInstancesDSWrapper) listAutopilotAddonInstancesToSchema(body *gjson.Result) error {
	d := w.ResourceData
	mErr := multierror.Append(nil,
		d.Set("region", w.Config.GetRegion(w.ResourceData)),
		d.Set("api_version", body.Get("apiVersion").Value()),
		d.Set("items", schemas.SliceToList(body.Get("items"),
			func(items gjson.Result) any {
				return map[string]any{
					"api_version": items.Get("apiVersion").Value(),
					"kind":       items.Get("kind").Value(),
					"metadata": schemas.SliceToList(items.Get("metadata"),
						func(metadata gjson.Result) any {
							return map[string]any{
								"alias":             metadata.Get("alias").Value(),
								"annotations":       metadata.Get("annotations").Value(),
								"creation_timestamp": metadata.Get("creationTimestamp").Value(),
								"labels":            metadata.Get("labels").Value(),
								"name":              metadata.Get("name").Value(),
								"uid":               metadata.Get("uid").Value(),
								"update_timestamp":   metadata.Get("updateTimestamp").Value(),
							}
						},
					),
					"spec": schemas.SliceToList(items.Get("spec"),
						func(spec gjson.Result) any {
							return map[string]any{
								"addon_template_labels": schemas.SliceToStrList(spec.Get("addonTemplateLabels")),
								"addon_template_logo":   spec.Get("addonTemplateLogo").Value(),
								"addon_template_name":   spec.Get("addonTemplateName").Value(),
								"addon_template_type":   spec.Get("addonTemplateType").Value(),
								"cluster_id":           spec.Get("clusterID").Value(),
								"description":         spec.Get("description").Value(),
								"values":              spec.Get("values").Value(),
								"version":             spec.Get("version").Value(),
							}
						},
					),
					"status": schemas.SliceToList(items.Get("status"),
						func(status gjson.Result) any {
							return map[string]any{
								"reason":          status.Get("Reason").Value(),
								"current_version":  w.setIteStaCurVersion(status),
								"is_rollbackable":  status.Get("isRollbackable").Value(),
								"message":         status.Get("message").Value(),
								"previous_version": status.Get("previousVersion").Value(),
								"status":          status.Get("status").Value(),
								"target_versions":  schemas.SliceToStrList(status.Get("targetVersions")),
							}
						},
					),
				}
			},
		)),
		d.Set("kind", body.Get("kind").Value()),
	)

	return mErr.ErrorOrNil()
}

func (*AddonInstancesDSWrapper) setIteStaCurVersion(status gjson.Result) any {
	return schemas.SliceToList(status.Get("currentVersion"), func(currentVersion gjson.Result) any {
		return map[string]any{
			"creation_timestamp": currentVersion.Get("creationTimestamp").Value(),
			"input":             currentVersion.Get("input").Value(),
			"stable":            currentVersion.Get("stable").Value(),
			"support_versions": schemas.SliceToList(currentVersion.Get("supportVersions"),
				func(supportVersions gjson.Result) any {
					return map[string]any{
						"category":       schemas.SliceToStrList(supportVersions.Get("category")),
						"cluster_type":    supportVersions.Get("clusterType").Value(),
						"cluster_version": schemas.SliceToStrList(supportVersions.Get("clusterVersion")),
					}
				},
			),
			"translate":       currentVersion.Get("translate").Value(),
			"update_timestamp": currentVersion.Get("updateTimestamp").Value(),
			"version":         currentVersion.Get("version").Value(),
		}
	})
}
