package apig

import (
	"regexp"
	"time"

	"github.com/huaweicloud/golangsdk/openstack/apigw/v2/apis"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

var (
	matching = map[string]string{
		"Prefix": "SWA",
		"Exact":  "NORMAL",
	}
)

func ResourceApigApiV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceApigApiV2Create,
		Read:   resourceApigApiV2Read,
		Update: resourceApigApiV2Update,
		Delete: resourceApigApiV2Delete,
		Importer: &schema.ResourceImporter{
			State: resourceApigInstanceSubResourceImportState,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(40 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile("^([\u4e00-\u9fa5A-Za-z][\u4e00-\u9fa5A-Za-z_0-9]{2,63})$"),
					"The name contains of 3 to 64 characters, starting with a letter. Only letters, digits, "+
						"hyphens (-) and underscore (_) are allowed."),
			},
			"version": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(1, 16),
			},
			"request_protocol": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"HTTP", "HTTPS", "BOTH",
				}, false),
			},
			"request_method": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"GET", "POST", "PUT", "DELETE", "HEAD", "PATCH", "OPTIONS", "ANY",
				}, false),
			},
			"request_path": {
				Type:     schema.TypeString,
				Required: true,
			},
			"security_authentication": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "NONE",
				ValidateFunc: validation.StringInSlice([]string{
					"NONE", "APP", "IAM", "AUTHORIZER",
				}, false),
			},
			"cors": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"matching": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"Exact", "Prefix",
				}, false),
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile("^[^<>]{1,255}$"),
					"The description contain a maximum of 255 characters, "+
						"and the angle brackets (< and >) are not allowed."),
			},
			"body_description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(1, 20480),
			},
			"success_response": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(1, 20480),
			},
			"failure_response": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(1, 20480),
			},
			"response_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"request_params": {
				Type:     schema.TypeSet,
				Required: true,
				MaxItems: 50,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"location": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "PATH",
							ValidateFunc: validation.StringInSlice([]string{
								"PATH", "HEADER", "QUERY",
							}, false),
						},
						"type": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "STRING",
							ValidateFunc: validation.StringInSlice([]string{
								"STRING", "NUMBER",
							}, false),
						},
						"is_required": {
							Type:     schema.TypeBool,
							Required: true,
						},
						"maximum": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"minimum": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"example": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringMatch(regexp.MustCompile("^[^<>]{1,255}$"),
								"The description contain a maximum of 255 characters, "+
									"and the angle brackets (< and >) are not allowed."),
						},
						"default": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringMatch(regexp.MustCompile("^[^<>]{1,255}$"),
								"The description contain a maximum of 255 characters, "+
									"and the angle brackets (< and >) are not allowed."),
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringMatch(regexp.MustCompile("^[^<>]{1,255}$"),
								"The description contain a maximum of 255 characters, "+
									"and the angle brackets (< and >) are not allowed."),
						},
					},
				},
			},
			"backend_params": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"location": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"PATH", "QUERY", "HEADER",
							}, false),
						},
						"req_param": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"constant_params": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"location": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"PATH", "QUERY", "HEADER",
							}, false),
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"mock": {
				Type:         schema.TypeList,
				Optional:     true,
				MaxItems:     1,
				ExactlyOneOf: []string{"func_graph", "web"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"description": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"content": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"version": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"func_graph": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"urn": {
							Type:     schema.TypeString,
							Required: true,
						},
						"timeout": {
							Type:     schema.TypeString,
							Required: true,
						},
						"type": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "sync",
							ValidateFunc: validation.StringInSlice([]string{
								"async", "sync",
							}, false),
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"version": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"web": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"path": {
							Type:     schema.TypeString,
							Required: true,
						},
						"host_header": {
							Type:          schema.TypeString,
							Optional:      true,
							ConflictsWith: []string{"web.0.backend_address"},
						},
						"vpc_channel_id": {
							Type:         schema.TypeString,
							Optional:     true,
							AtLeastOneOf: []string{"web.0.backend_address"},
						},
						"backend_address": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"request_method": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								"GET", "POST", "PUT", "DELETE", "HEAD", "PATCH", "OPTIONS", "ANY",
							}, false),
						},
						"protocol": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "HTTPS",
							ValidateFunc: validation.StringInSlice([]string{
								"HTTP", "HTTPS",
							}, false),
						},
						"timeout": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      5000,
							ValidateFunc: validation.IntBetween(1, 600000),
						},
						"ssl_enable": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
					},
				},
			},
			"mock_policy": {
				Type:          schema.TypeSet,
				MaxItems:      5,
				Optional:      true,
				ConflictsWith: []string{"func_graph", "web", "func_graph_policy", "web_policy"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"response": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringLenBetween(8, 2048),
						},
						"effective_mode": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"ALL", "ANY",
							}, false),
						},
						"backend_params": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Required: true,
									},
									"location": {
										Type:     schema.TypeString,
										Required: true,
										ValidateFunc: validation.StringInSlice([]string{
											"PATH", "QUERY", "HEADER",
										}, false),
									},
									"req_param": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
						"constant_params": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Required: true,
									},
									"location": {
										Type:     schema.TypeString,
										Required: true,
										ValidateFunc: validation.StringInSlice([]string{
											"PATH", "QUERY", "HEADER",
										}, false),
									},
									"value": {
										Type:     schema.TypeString,
										Required: true,
									},
									"description": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"conditions": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"source": {
										Type:     schema.TypeString,
										Optional: true,
										Default:  "param",
										ValidateFunc: validation.StringInSlice([]string{
											"param", "source",
										}, false),
									},
									"param_name": {
										Type:         schema.TypeString,
										Optional:     true,
										RequiredWith: []string{},
									},
									"type": {
										Type:     schema.TypeString,
										Optional: true,
										ValidateFunc: validation.StringInSlice([]string{
											"exact", "enum", "pattern",
										}, false),
									},
									"value": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
					},
				},
			},
			"func_graph_policy": {
				Type:          schema.TypeSet,
				MaxItems:      5,
				Optional:      true,
				ConflictsWith: []string{"mock", "web", "mock_policy", "web_policy"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"function_urn": {
							Type:     schema.TypeString,
							Required: true,
						},
						"invocation_mode": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "sync",
							ValidateFunc: validation.StringInSlice([]string{
								"async", "sync",
							}, false),
						},
						"version": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								"HTTP", "HTTPS",
							}, false),
						},
						"effective_mode": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"ALL", "ANY",
							}, false),
						},
						"timeout": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      5000,
							ValidateFunc: validation.IntBetween(1, 600000),
						},
						"backend_params": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Required: true,
									},
									"location": {
										Type:     schema.TypeString,
										Required: true,
										ValidateFunc: validation.StringInSlice([]string{
											"PATH", "QUERY", "HEADER",
										}, false),
									},
									"req_param": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
						"constant_params": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Required: true,
									},
									"location": {
										Type:     schema.TypeString,
										Required: true,
										ValidateFunc: validation.StringInSlice([]string{
											"PATH", "QUERY", "HEADER",
										}, false),
									},
									"value": {
										Type:     schema.TypeString,
										Required: true,
									},
									"description": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"conditions": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"source": {
										Type:     schema.TypeString,
										Optional: true,
										Default:  "param",
										ValidateFunc: validation.StringInSlice([]string{
											"param", "source",
										}, false),
									},
									"param_name": {
										Type:         schema.TypeString,
										Optional:     true,
										RequiredWith: []string{},
									},
									"type": {
										Type:     schema.TypeString,
										Optional: true,
										ValidateFunc: validation.StringInSlice([]string{
											"exact", "enum", "pattern",
										}, false),
									},
									"value": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
					},
				},
			},
			"web_policy": {
				Type:          schema.TypeSet,
				MaxItems:      5,
				Optional:      true,
				ConflictsWith: []string{"mock", "func_graph", "mock_policy", "func_graph_policy"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"host_header": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"vpc_channel_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"backend_address": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"request_protocol": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								"HTTP", "HTTPS",
							}, false),
						},
						"method": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"GET", "POST", "PUT", "DELETE", "HEAD", "PATCH", "OPTIONS", "ANY",
							}, false),
						},
						"effective_mode": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"ALL", "ANY",
							}, false),
						},
						"path": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"timeout": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      5000,
							ValidateFunc: validation.IntBetween(1, 600000),
						},
						"backend_params": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Required: true,
									},
									"location": {
										Type:     schema.TypeString,
										Required: true,
										ValidateFunc: validation.StringInSlice([]string{
											"PATH", "QUERY", "HEADER",
										}, false),
									},
									"req_param": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
						"constant_params": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Required: true,
									},
									"location": {
										Type:     schema.TypeString,
										Required: true,
										ValidateFunc: validation.StringInSlice([]string{
											"PATH", "QUERY", "HEADER",
										}, false),
									},
									"value": {
										Type:     schema.TypeString,
										Required: true,
									},
									"description": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"conditions": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"source": {
										Type:     schema.TypeString,
										Optional: true,
										Default:  "param",
										ValidateFunc: validation.StringInSlice([]string{
											"param", "source",
										}, false),
									},
									"param_name": {
										Type:         schema.TypeString,
										Optional:     true,
										RequiredWith: []string{},
									},
									"type": {
										Type:     schema.TypeString,
										Optional: true,
										ValidateFunc: validation.StringInSlice([]string{
											"exact", "enum", "pattern",
										}, false),
									},
									"value": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
					},
				},
			},
			"register_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"update_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"publish_env_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"publish_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildApigApiMock(mock []interface{}) apis.Mock {
	mockMap := mock[0].(map[string]interface{})
	return apis.Mock{
		Description:   mockMap["description"].(string),
		ResultContent: mockMap["content"].(string),
		Version:       mockMap["version"].(string),
	}
}

func buildApigApiFuncGraph(funcGraph []interface{}) apis.FuncGraph {
	funcGraphMap := funcGraph[0].(map[string]interface{})
	return apis.FuncGraph{
		Timeout:        funcGraphMap["timeout"].(int),
		InvocationType: funcGraphMap["type"].(string),
		FunctionUrn:    funcGraphMap["urn"].(string),
		Description:    funcGraphMap["description"].(string),
		Version:        funcGraphMap["version"].(string),
	}
}

func buildApigApiWeb(web []interface{}) apis.Web {
	webMap := web[0].(map[string]interface{})
	result := apis.Web{
		ReqUri:          webMap["path"].(string),
		ReqMethod:       webMap["request_method"].(string),
		ReqProtocol:     webMap["protocol"].(string),
		Timeout:         webMap["timeout"].(int),
		EnableClientSsl: webMap["ssl_enable"].(bool),
	}
	if chanId, ok := webMap["vpc_channel_id"]; ok && chanId != "" {
		result.VpcChannelEnable = 1
		result.VpcChannelInfo = apis.VpcChannel{
			VpcChannelId: chanId.(string),
		}
	} else {
		result.VpcChannelEnable = 2
		result.UrlDomain = webMap["backend_address"].(string)
	}

	return result
}

func buildApiRequestParameters(requests *schema.Set) []apis.ReqParamBase {
	result := make([]apis.ReqParamBase, requests.Len())
	for i, v := range requests.List() {
		paramMap := v.(map[string]interface{})
		paramType := paramMap["type"].(string)
		param := apis.ReqParamBase{
			Name:     paramMap["name"].(string),
			Location: paramMap["location"].(string),
		}
		if paramType == "NUMBER" {
			param.MaxNum = paramMap["maximum"].(int)
			param.MinNum = paramMap["minimum"].(int)
		} else if paramType == "STRING" {
			param.MaxSize = paramMap["maximum"].(int)
			param.MinSize = paramMap["minimum"].(int)
		}
		param.Type = paramType

		if paramMap["is_required"].(bool) {
			param.Required = 1
		} else {
			param.Required = 2
		}
		result[i] = param
	}
	return result
}

func buildApiBackendParameters(backends, constants *schema.Set) []apis.BackendParamBase {
	result := make([]apis.BackendParamBase, 0, backends.Len()+constants.Len())
	for _, v := range backends.List() {
		paramMap := v.(map[string]interface{})
		result = append(result, apis.BackendParamBase{
			Origin:   "REQUEST",
			Name:     paramMap["name"].(string),
			Location: paramMap["location"].(string),
			Value:    paramMap["req_param"].(string),
		})
	}
	for _, v := range constants.List() {
		paramMap := v.(map[string]interface{})
		result = append(result, apis.BackendParamBase{
			Origin:      "CONSTANT",
			Name:        paramMap["name"].(string),
			Location:    paramMap["location"].(string),
			Value:       paramMap["value"].(string),
			Description: paramMap["description"].(string),
		})
	}

	return result
}

func buildApiPolicyConditions(conditions *schema.Set) []apis.ApiConditionBase {
	result := make([]apis.ApiConditionBase, conditions.Len())
	for i, v := range conditions.List() {
		conditionMap := v.(map[string]interface{})
		result[i] = apis.ApiConditionBase{
			ReqParamName:    conditionMap["param_name"].(string),
			ConditionOrigin: conditionMap["source"].(string),
			ConditionType:   conditionMap["type"].(string),
			ConditionValue:  conditionMap["value"].(string),
		}
	}
	return result
}

func buildApigApiMockPolicy(mocks *schema.Set) []apis.PolicyMock {
	result := make([]apis.PolicyMock, mocks.Len())
	for i, policy := range mocks.List() {
		pm := policy.(map[string]interface{})
		result[i] = apis.PolicyMock{
			Name:          pm["name"].(string),
			ResultContent: pm["response"].(string),
			EffectMode:    pm["effective_mode"].(string),
			BackendParams: buildApiBackendParameters(pm["backend_params"].(*schema.Set),
				pm["constant_params"].(*schema.Set)),
			Conditions: buildApiPolicyConditions(pm["conditions"].(*schema.Set)),
		}
	}
	return result
}

func buildApigApiFuncGraphPolicy(policies *schema.Set) []apis.PolicyFuncGraph {
	result := make([]apis.PolicyFuncGraph, policies.Len())
	for i, policy := range policies.List() {
		pm := policy.(map[string]interface{})
		result[i] = apis.PolicyFuncGraph{
			Name:           pm["name"].(string),
			FunctionUrn:    pm["function_urn"].(string),
			InvocationType: pm["invocation_mode"].(string),
			EffectMode:     pm["effective_mode"].(string),
			Timeout:        pm["timeout"].(int),
			BackendParams: buildApiBackendParameters(pm["backend_params"].(*schema.Set),
				pm["constant_params"].(*schema.Set)),
			Conditions: buildApiPolicyConditions(pm["conditions"].(*schema.Set)),
		}
	}
	return result
}

func buildApigApiWebPolicy(policies *schema.Set) []apis.PolicyWeb {
	result := make([]apis.PolicyWeb, policies.Len())
	for i, policy := range policies.List() {
		pm := policy.(map[string]interface{})
		wp := apis.PolicyWeb{
			Name:        pm["name"].(string),
			ReqProtocol: pm["request_protocol"].(string),
			ReqMethod:   pm["method"].(string),
			ReqUri:      pm["path"].(string),
			EffectMode:  pm["effective_mode"].(string),
			Timeout:     pm["timeout"].(int),
			UrlDomain:   pm["host_header"].(string),
			BackendParams: buildApiBackendParameters(pm["backend_params"].(*schema.Set),
				pm["constant_params"].(*schema.Set)),
			Conditions: buildApiPolicyConditions(pm["conditions"].(*schema.Set)),
		}
		if chanId, ok := pm["vpc_channel_id"]; ok {
			if chanId != "" {
				wp.VpcChannelInfo = apis.VpcChannel{
					VpcChannelId: pm["vpc_channel_id"].(string),
				}
				wp.VpcChannelEnable = 1
			} else {
				wp.VpcChannelEnable = 2
			}
		}
		result[i] = wp
	}
	return result
}

func buildApigApiParameters(d *schema.ResourceData) (apis.ApiOpts, error) {
	opt := apis.ApiOpts{
		Type:                2,
		GroupId:             d.Get("group_id").(string),
		Name:                d.Get("name").(string),
		Version:             d.Get("version").(string),
		ReqProtocol:         d.Get("request_protocol").(string),
		ReqMethod:           d.Get("request_method").(string),
		ReqUri:              d.Get("request_path").(string),
		Cors:                d.Get("cors").(bool),
		AuthType:            d.Get("security_authentication").(string),
		MatchMode:           d.Get("matching").(string),
		Description:         d.Get("description").(string),
		BodyDescription:     d.Get("body_description").(string),
		ResultNormalSample:  d.Get("success_response").(string),
		ResultFailureSample: d.Get("failure_response").(string),
		ResponseId:          d.Get("response_id").(string),
	}
	// build match mode
	v, ok := matching[d.Get("matching").(string)]
	if !ok {
		return opt, fmtp.Errorf("Unable to extract match mode")
	}
	opt.MatchMode = v

	opt.ReqParams = buildApiRequestParameters(d.Get("request_params").(*schema.Set))
	opt.BackendParams = buildApiBackendParameters(d.Get("backend_params").(*schema.Set),
		d.Get("constant_params").(*schema.Set))

	// build backend (one of the mock, function graph and web) server
	if m, ok := d.GetOk("mock"); ok {
		opt.BackendType = "MOCK"
		opt.MockInfo = buildApigApiMock(m.([]interface{}))
		mp := d.Get("mock_policy").(*schema.Set)
		opt.PolicyMocks = buildApigApiMockPolicy(mp)
	} else if fg, ok := d.GetOk("func_graph"); ok {
		opt.BackendType = "FUNCTION"
		opt.FuncInfo = buildApigApiFuncGraph(fg.([]interface{}))
		fgp := d.Get("func_graph_policy").(*schema.Set)
		opt.PolicyFunctions = buildApigApiFuncGraphPolicy(fgp)
	} else {
		opt.BackendType = "HTTP"
		web := d.Get("web").([]interface{})
		opt.WebInfo = buildApigApiWeb(web)
		wp := d.Get("web_policy").(*schema.Set)
		opt.PolicyWebs = buildApigApiWebPolicy(wp)
	}

	return opt, nil
}

func resourceApigApiV2Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	opt, err := buildApigApiParameters(d)
	if err != nil {
		return fmtp.Errorf("Unable to build the Api parameter: %s", err)
	}
	client, err := config.ApigV2Client(config.GetRegion(d)) // client
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud APIG v2 client: %s", err)
	}
	instanceId := d.Get("instance_id").(string)
	resp, err := apis.Create(client, instanceId, opt).Extract()
	if err != nil {
		return fmtp.Errorf("Error creating APIG v2 API: %s", err)
	}
	logp.Printf("[Lance] The ID 1 is %s", resp.Id)
	d.SetId(resp.Id)
	logp.Printf("[Lance] The ID 2 is %s", d.Id())
	return resourceApigApiV2Read(d, meta)
}

func getApigAPiReqParams(reqParams []apis.ReqParamResp) []map[string]interface{} {
	result := make([]map[string]interface{}, len(reqParams))
	for i, v := range reqParams {
		param := map[string]interface{}{
			"name":        v.Name,
			"location":    v.Location,
			"type":        v.Type,
			"example":     v.SampleValue,
			"default":     v.DefaultValue,
			"description": v.Description,
		}
		if v.Type == "NUMBER" {
			param["maximum"] = v.MaxNum
			param["minimum"] = v.MinNum
		} else if v.Type == "STRING" {
			param["maximum"] = v.MaxSize
			param["minimum"] = v.MinSize
		}
		if v.Required == 1 {
			param["is_required"] = true
		} else if v.Required == 2 {
			param["is_required"] = false
		}
		result[i] = param
	}
	return result
}

func getApigApiBackendParams(backendParams []apis.BackendParamResp) ([]map[string]interface{},
	[]map[string]interface{}) {
	backendResult := make([]map[string]interface{}, 0)
	constantResult := make([]map[string]interface{}, 0)
	for _, v := range backendParams {
		origin := v.Origin
		if origin == "REQUEST" {
			param := map[string]interface{}{
				"name":      v.Name,
				"location":  v.Location,
				"req_param": v.Value,
			}
			backendResult = append(backendResult, param)
		}
		if origin == "CONSTANT" {
			param := map[string]interface{}{
				"name":        v.Name,
				"location":    v.Location,
				"value":       v.Value,
				"description": v.Description,
			}
			constantResult = append(constantResult, param)
		}
	}
	return backendResult, constantResult
}

func getApigApiPolicyConditions(conditions []apis.ApiConditionBase) []map[string]interface{} {
	result := make([]map[string]interface{}, len(conditions))
	for i, v := range conditions {
		result[i] = map[string]interface{}{
			"source":     v.ConditionOrigin,
			"param_name": v.ReqParamName,
			"type":       v.ConditionType,
			"value":      v.ConditionValue,
		}
	}
	return result
}

func setApigApiMatchMode(d *schema.ResourceData, mode string) error {
	result, ok := matching[mode]
	if !ok {
		return fmtp.Errorf("The matching mode is invalid")
	}
	return d.Set("matching", result)
}

func setApigApiMock(d *schema.ResourceData, mockResp apis.Mock) error {
	result := []map[string]interface{}{
		{
			"description": mockResp.Description,
			"content":     mockResp.ResultContent,
			"version":     mockResp.Version,
		},
	}
	return d.Set("mock", result)
}

func setApigApiFuncGraph(d *schema.ResourceData, funcResp apis.FuncGraph) error {
	result := []map[string]interface{}{
		{
			"urn":         funcResp.FunctionUrn,
			"timeout":     funcResp.Timeout,
			"type":        funcResp.InvocationType,
			"description": funcResp.Description,
			"version":     funcResp.Version,
		},
	}
	return d.Set("func_graph", result)
}

func setApigApiWeb(d *schema.ResourceData, webResp apis.Web) error {
	result := make([]map[string]interface{}, 1)
	web := map[string]interface{}{
		"path":           webResp.ReqUri,
		"request_method": webResp.ReqMethod,
		"protocol":       webResp.ReqProtocol,
		"timeout":        webResp.Timeout,
		"ssl_enable":     d.Get("web.0.ssl_enable"),
	}
	if webResp.VpcChannelInfo.VpcChannelId != "" {
		web["vpc_channel_id"] = webResp.VpcChannelInfo.VpcChannelId
		web["host_header"] = webResp.VpcChannelInfo.VpcChannelProxyHost
	} else {
		web["backend_address"] = webResp.UrlDomain
	}
	return d.Set("web", result)
}

func setApigApiMockPolicy(d *schema.ResourceData, policies []apis.PolicyMockResp) error {
	result := make([]map[string]interface{}, len(policies))
	for i, policy := range policies {
		mp := map[string]interface{}{
			"name":           policy.Name,
			"response":       policy.ResultContent,
			"effective_mode": policy.EffectMode,
		}
		backendParams, constantParams := getApigApiBackendParams(policy.BackendParams)
		mp["backend_params"] = backendParams
		mp["constant_params"] = constantParams
		mp["conditions"] = getApigApiPolicyConditions(policy.Conditions)

		result[i] = mp
	}
	return d.Set("mock_policy", result)
}

func setApigApiFuncGraphPolicy(d *schema.ResourceData, policies []apis.PolicyFuncGraphResp) error {
	result := make([]map[string]interface{}, len(policies))
	for i, policy := range policies {
		fgp := map[string]interface{}{
			"name":            policy.Name,
			"function_urn":    policy.FunctionUrn,
			"version":         policy.Version,
			"invocation_mode": policy.InvocationType,
			"effective_mode":  policy.EffectMode,
			"timeout":         policy.Timeout,
		}
		backendParams, constantParams := getApigApiBackendParams(policy.BackendParams)
		fgp["backend_params"] = backendParams
		fgp["constant_params"] = constantParams
		fgp["conditions"] = getApigApiPolicyConditions(policy.Conditions)

		result[i] = fgp
	}
	return d.Set("func_graph_policy", result)
}

func setApigApiWebPolicy(d *schema.ResourceData, policies []apis.PolicyWebResp) error {
	result := make([]map[string]interface{}, len(policies))
	for i, policy := range policies {
		wp := map[string]interface{}{
			"name":             policy.Name,
			"request_protocol": policy.ReqProtocol,
			"method":           policy.ReqMethod,
			"effective_mode":   policy.EffectMode,
			"path":             policy.ReqUri,
			"timeout":          policy.Timeout,
		}
		if policy.VpcChannelInfo.VpcChannelId != "" {
			wp["vpc_channel_id"] = policy.VpcChannelInfo.VpcChannelId
			wp["host_header"] = policy.VpcChannelInfo.VpcChannelProxyHost
		} else {
			wp["backend_address"] = policy.UrlDomain
		}
		backendParams, constantParams := getApigApiBackendParams(policy.BackendParams)
		wp["backend_params"] = backendParams
		wp["constant_params"] = constantParams
		wp["conditions"] = getApigApiPolicyConditions(policy.Conditions)

		result[i] = wp
	}
	return d.Set("web_policy", result)
}

func setApigApiParameters(d *schema.ResourceData, config *config.Config, resp *apis.ApiResp) error {
	backendParams, constantParams := getApigApiBackendParams(resp.BackendParams)
	mErr := multierror.Append(nil,
		d.Set("region", config.GetRegion(d)),
		d.Set("name", resp.Name),
		d.Set("version", resp.Version),
		d.Set("request_protocol", resp.ReqProtocol),
		d.Set("request_method", resp.ReqMethod),
		d.Set("request_path", resp.ReqUri),
		d.Set("security_authentication", resp.AuthType),
		d.Set("cors", resp.Cors),
		d.Set("description", resp.Description),
		d.Set("body_description", resp.BodyDescription),
		d.Set("success_response", resp.ResultNormalSample),
		d.Set("failure_response", resp.ResultFailureSample),
		d.Set("response_id", resp.ResponseId),
		d.Set("request_params", getApigAPiReqParams(resp.ReqParams)),
		d.Set("backend_params", backendParams),
		d.Set("constant_params", constantParams),
		setApigApiMatchMode(d, resp.MatchMode),
		setApigApiMock(d, resp.MockInfo),
		setApigApiMockPolicy(d, resp.PolicyMocks),
		setApigApiFuncGraph(d, resp.FuncInfo),
		setApigApiFuncGraphPolicy(d, resp.PolicyFunctions),
		setApigApiWeb(d, resp.WebInfo),
		setApigApiWebPolicy(d, resp.PolicyWebs),
		d.Set("publish_env_id", resp.RunEnvId),
		d.Set("publish_id", resp.PublishId),
		d.Set("register_time", resp.RegisterTime),
		d.Set("update_time", resp.UpdateTime),
	)
	if mErr.ErrorOrNil() != nil {
		return mErr
	}
	return nil
}

func resourceApigApiV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	client, err := config.ApigV2Client(config.GetRegion(d)) // client
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud APIG client: %s", err)
	}

	instanceId := d.Get("instance_id").(string)
	resp, err := apis.Get(client, instanceId, d.Id()).Extract()
	if err != nil {
		return fmtp.Errorf("Error getting api information from server: %s", err)
	}
	return setApigApiParameters(d, config, resp)
}

func resourceApigApiV2Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	client, err := config.ApigV2Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud APIG v2 client: %s", err)
	}
	instanceId := d.Get("instance_id").(string)
	opt, err := buildApigApiParameters(d)
	if err != nil {
		return fmtp.Errorf("Unable to build the Api parameter: %s", err)
	}
	_, err = apis.Update(client, instanceId, d.Id(), opt).Extract()
	if err != nil {
		return fmtp.Errorf("Error updating APIG v2 API: %s", err)
	}

	return resourceApigApiV2Read(d, meta)
}

func resourceApigApiV2Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	client, err := config.ApigV2Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud APIG v2 client: %s", err)
	}
	instanceId := d.Get("instance_id").(string)
	if err = apis.Delete(client, instanceId, d.Id()).ExtractErr(); err != nil {
		return fmtp.Errorf("Unable to delete the APIG v2 dedicated instance (%s): %s", d.Id(), err)
	}
	d.SetId("")
	return nil
}
