package huaweicloud

import (
	"github.com/huaweicloud/golangsdk/openstack/apigw/v2/apis"
	"regexp"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
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
			State: schema.ImportStatePassthrough,
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
			"simple_authentication": { // Not support yet.
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			}, // Not support yet.
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
				Computed:     true,
				ValidateFunc: validation.StringLenBetween(1, 20480),
			},
			"normal_sample": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(1, 20480),
			}, // Cannot find on console.
			"failure_sample": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(1, 20480),
			}, // Cannot find on console.
			"response_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"request_params": {
				Type:     schema.TypeList,
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
							Optional: true,
							Default:  true,
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
			"mock": {
				Type:     schema.TypeList,
				MaxItems: 1,
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
			"function_graph": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"timeout": {
							Type:     schema.TypeString,
							Required: true,
						},
						"type": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"async", "sync",
							}, false),
						},
						"urn": {
							Type:     schema.TypeString,
							Required: true,
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
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"path": {
							Type:     schema.TypeString,
							Required: true,
						},
						"vpc_channel_id": {
							Type:     schema.TypeString,
							Required: true,
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
						"description": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringMatch(regexp.MustCompile("^[^<>]{1,255}$"),
								"The description contain a maximum of 255 characters, "+
									"and the angle brackets (< and >) are not allowed."),
						},
						"timeout": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      5000,
							ValidateFunc: validation.IntBetween(1, 600000),
						},
						"enable_ssl": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},
			"web_policy": {
				Type:     schema.TypeList,
				MaxItems: 5,
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
						"request_protocol": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								"HTTP", "HTTPS",
							}, false),
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringMatch(regexp.MustCompile("^[^<>]{1,255}$"),
								"The description contain a maximum of 255 characters, "+
									"and the angle brackets (< and >) are not allowed."),
						},
						"method": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"GET", "POST", "PUT", "DELETE", "HEAD", "PATCH", "OPTIONS", "ANY",
							}, false),
						},
						"uri": {
							Type:     schema.TypeString,
							Optional: true,
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
							Type:     schema.TypeString,
							Optional: true,
						},
						"vpc_channel_id": {
							Type:     schema.TypeString,
							Optional: true,
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
							Type: schema.TypeSet,
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
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
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
	return apis.Web{
		ReqUri:          webMap["path"].(string),
		ReqMethod:       webMap["request_method"].(string),
		ReqProtocol:     webMap["protocol"].(string),
		Timeout:         webMap["timeout"].(int),
		EnableClientSsl: webMap["ssl_enable"].(bool),
	}
}

func buildApiBackendParameters(backends, constants []interface{}) []apis.BackendParamBase {
	result := make([]apis.BackendParamBase, 0, len(backends)+len(constants))
	for _, v := range backends {
		paramMap := v.(map[string]interface{})
		result = append(result, apis.BackendParamBase{
			Origin:   "REQUEST",
			Name:     paramMap["name"].(string),
			Location: paramMap["location"].(string),
			Value:    paramMap["req_param"].(string),
		})
	}
	for _, v := range backends {
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

func buildApiPolicyConditions(conditions []interface{}) []apis.ApiConditionBase {
	result := make([]apis.ApiConditionBase, len(conditions))
	for i, v := range conditions {
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

func buildApigApiMockPolicy(mocks []interface{}) []apis.PolicyMock {
	result := make([]apis.PolicyMock, len(mocks))
	//for i, mock := range mocks {
	//	mockMap := mock.(map[string]interface{})
	//	result[i] = apis.PolicyMock{
	//		Name:         mockMap["name"].(string),
	//		ReqProtocol:  mockMap["request_protocol"].(string),
	//		ReqMethod:    mockMap["method"].(string),
	//		ReqUri:       mockMap["path"].(string),
	//		EffectMode:   mockMap["effective_mode"].(string),
	//		Timeout:      mockMap["timeout"].(int),
	//		BackendParams: buildApiBackendParameters(policyMap["backend_params"].([]interface{}),
	//			policyMap["constant_params"].([]interface{})),
	//	}
	//
	//}
	return result
}

func buildApigApiFuncGraphPolicy(policies []interface{}) []apis.PolicyWeb {
	result := make([]apis.PolicyWeb, len(policies))
	//for i, policy := range policies {
	//	policyMap := policy.(map[string]interface{})
	//	result[i] = apis.PolicyWeb{
	//		Name:         policyMap["name"].(string),
	//		ReqProtocol:  policyMap["request_protocol"].(string),
	//		ReqMethod:    policyMap["method"].(string),
	//		ReqUri:       policyMap["path"].(string),
	//		EffectMode:   policyMap["effective_mode"].(string),
	//		Timeout:      policyMap["timeout"].(int),
	//		BackendParams: buildApiBackendParameters(policyMap["backend_params"].([]interface{}),
	//			policyMap["constant_params"].([]interface{})),
	//	}
	//
	//}
	return result
}

func buildApigApiWebPolicy(policies []interface{}) []apis.PolicyWeb {
	result := make([]apis.PolicyWeb, len(policies))
	for i, policy := range policies {
		policyMap := policy.(map[string]interface{})
		result[i] = apis.PolicyWeb{
			Name:        policyMap["name"].(string),
			ReqProtocol: policyMap["request_protocol"].(string),
			ReqMethod:   policyMap["method"].(string),
			ReqUri:      policyMap["path"].(string),
			EffectMode:  policyMap["effective_mode"].(string),
			Timeout:     policyMap["timeout"].(int),
			BackendParams: buildApiBackendParameters(policyMap["backend_params"].([]interface{}),
				policyMap["constant_params"].([]interface{})),
			Conditions: buildApiPolicyConditions(policyMap["conditions"].([]interface{})),
		}

	}
	return result
}

func buildApigApiParameters(d *schema.ResourceData) (apis.ApiOpts, error) {
	opt := apis.ApiOpts{
		Name:                d.Get("name").(string),
		Type:                d.Get("type").(string),
		ReqProtocol:         d.Get("req_protocol").(string),
		ReqMethod:           d.Get("request_method").(string),
		ReqUri:              d.Get("request_path").(string),
		Cors:                d.Get("cors").(bool),
		AuthType:            d.Get("authentication").(string),
		MatchMode:           d.Get("matching").(string),
		Description:         d.Get("description").(string),
		BodyDescription:     d.Get("request_body").(string),
		ResultNormalSample:  d.Get("normal_sample").(string),
		ResultFailureSample: d.Get("failure_sample").(string),
		ResponseId:          d.Get("response_id").(string),
	}
	// build match mode
	v, ok := matching[d.Get("matching").(string)]
	if !ok {
		return opt, fmtp.Errorf("Unable to extract match mode")
	}
	opt.MatchMode = v
	// build backend (one of the mock, function graph and web) server
	if m, ok := d.GetOk("mock"); ok {
		opt.MockInfo = buildApigApiMock(m.([]interface{}))
		mp := d.Get("mock_policy")
		opt.PolicyMocks = buildApigApiMockPolicy(mp.([]interface{}))
	} else if fg, ok := d.GetOk("function_graph"); ok {
		opt.FuncInfo = buildApigApiFuncGraph(fg.([]interface{}))

	} else {
		web := d.Get("web").([]interface{})
		opt.WebInfo = buildApigApiWeb(web)
		wp := d.Get("web_policy")
		opt.PolicyWebs = buildApigApiWebPolicy(wp.([]interface{}))
	}

	return opt, nil
}

func resourceApigApiV2Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	opt, err := buildApigApiParameters(d)
	if err != nil {
		return fmtp.Errorf("Unable to build the Api parameter: %s", err)
	}
	client, err := config.ApigV2Client(GetRegion(d, config)) // client
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud APIG v2 client: %s", err)
	}
	instanceId := d.Get("instance_id").(string)
	_, err = apis.Create(client, instanceId, opt).Extract()
	if err != nil {
		return fmtp.Errorf("Error creating APIG v2 API: %s", err)
	}
	return resourceApigInstanceV2Read(d, meta)
}

func setApigApiWeb(d *schema.ResourceData, webResp apis.Web) error {
	result := make([]map[string]interface{}, 1)
	result[0] = map[string]interface{}{
		"path":           webResp.ReqUri,
		"vpc_channel_id": webResp.VpcChannelInfo.VpcChannelId,
		"request_method": webResp.ReqMethod,
		"protocol":       webResp.ReqProtocol,
		"description":    webResp.Description,
		"timeout":        webResp.Timeout,
		"enable_ssl":     webResp.EnableClientSsl,
	}
	return d.Set("web", result)
}

func setApigApiWebPolicy(d *schema.ResourceData, webResp apis.PolicyWeb) error {
	result := make([]map[string]interface{}, 1)
	result[0] = map[string]interface{}{
		"path":           webResp.ReqUri,
		"vpc_channel_id": webResp.VpcChannelInfo.VpcChannelId,
		"request_method": webResp.ReqMethod,
		"protocol":       webResp.ReqProtocol,
		"description":    webResp.Description,
		"timeout":        webResp.Timeout,
		"enable_ssl":     webResp.EnableClientSsl,
	}
	return d.Set("web", result)
}

func setApigApiParameters(d *schema.ResourceData, config *config.Config, resp *apis.ApiResp) error {
	mErr := multierror.Append(nil,
		d.Set("region", GetRegion(d, config)),
		d.Set("name", resp.Name),
		d.Set("type", resp.Type),
		d.Set("req_protocol", resp.ReqProtocol),
		d.Set("request_method", resp.ReqMethod),
		d.Set("request_path", resp.ReqUri),
		d.Set("cors", resp.Cors),
		d.Set("matching", resp.MatchMode),
		d.Set("description", resp.Description),
		d.Set("request_body", resp.BodyDescription),
		d.Set("normal_sample", resp.ResultNormalSample),
		d.Set("failure_sample", resp.ResultFailureSample),
		d.Set("response_id", resp.ResponseId),
	)
	if mErr.ErrorOrNil() != nil {
		return mErr
	}
	return nil
}

func resourceApigApiV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	client, err := config.ApigV2Client(GetRegion(d, config)) // client
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
	client, err := config.ApigV2Client(GetRegion(d, config))
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

	return resourceApigInstanceV2Read(d, meta)
}

func resourceApigApiV2Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	client, err := config.ApigV2Client(GetRegion(d, config))
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
