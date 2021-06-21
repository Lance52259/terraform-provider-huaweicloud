package apig

import (
	"log"
	"regexp"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/apigw/v2/apigroups"
	"github.com/huaweicloud/golangsdk/openstack/apigw/v2/environments"
	"github.com/huaweicloud/golangsdk/openstack/apigw/v2/responses"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
)

func ResourceApigGroupV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceApigGroupV2Create,
		Read:   resourceApigGroupV2Read,
		Update: resourceApigGroupV2Update,
		Delete: resourceApigGroupV2Delete,
		Importer: &schema.ResourceImporter{
			State: resourceApigInstanceSubResourceImportState,
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
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile("^[\u4e00-\u9fa5A-Za-z][\u4e00-\u9fa5A-Za-z_0-9]{2,63}$"),
					"The name consists of 3 to 64 characters, starting with a letter. "+
						"Only letters, digits and underscores (_) are allowed. "+
						"Chinese characters must be in UTF-8 or Unicode format."),
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile("^[^<>]{1,255}$"),
					"The description contain a maximum of 255 characters, "+
						"and the angle brackets (< and >) are not allowed."),
			},
			"environments": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"variables": {
							Type:     schema.TypeSet,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Required: true,
										ValidateFunc: validation.StringMatch(
											regexp.MustCompile("^[A-Za-z][\\w_-]{2,31}$"),
											"The name consists of 3 to 32 characters, starting with a letter. "+
												"Only letters, digits, hyphens (-) and underscores (_) are allowed."),
									},
									"value": {
										Type:     schema.TypeString,
										Required: true,
										ValidateFunc: validation.StringMatch(
											regexp.MustCompile("^[\\w:/.-]{1,255}$"),
											"The value consists of 1 to 255 characters, only letters, digit and "+
												"following special characters are allowed: _-/.:"),
									},
									"variable_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"environment_id": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"custom_responses": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"responses": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"error_type": {
										Type:     schema.TypeString,
										Required: true,
										ValidateFunc: validation.StringInSlice([]string{
											"AUTH_FAILURE", "AUTH_HEADER_MISSING", "AUTHORIZER_FAILURE", "AUTHORIZER_CONF_FAILURE",
											"AUTHORIZER_IDENTITIES_FAILURE", "BACKEND_UNAVAILABLE", "BACKEND_TIMEOUT", "THROTTLED",
											"UNAUTHORIZED", "ACCESS_DENIED", "NOT_FOUND", "REQUEST_PARAMETERS_FAILURE", "DEFAULT_4XX",
											"DEFAULT_5XX",
										}, false),
									},
									"body": {
										Type:     schema.TypeString,
										Required: true,
									},
									"status_code": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntBetween(200, 599),
									},
								},
							},
						},
						"response_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"registraion_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"update_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func createApigGroupEnvironmentVariables(client *golangsdk.ServiceClient, instanceId, groupId string,
	environmentSet *schema.Set) error {
	for _, env := range environmentSet.List() {
		envMap := env.(map[string]interface{})
		envId := envMap["environment_id"].(string)
		for _, v := range envMap["variables"].(*schema.Set).List() {
			variable := v.(map[string]interface{})
			opt := environments.CreateVariableOpts{
				Name:    variable["name"].(string),
				Value:   variable["value"].(string),
				GroupId: groupId,
				EnvId:   envId,
			}
			if _, err := environments.CreateVariable(client, instanceId, opt).Extract(); err != nil {
				return err
			}
		}
	}
	return nil
}

func removeApigGroupEnvironmentVariables(client *golangsdk.ServiceClient, instanceId string,
	environmentSet *schema.Set) error {
	for _, env := range environmentSet.List() {
		envMap := env.(map[string]interface{})
		for _, v := range envMap["variables"].(*schema.Set).List() {
			variable := v.(map[string]interface{})
			err := environments.DeleteVariable(client, instanceId, variable["variable_id"].(string)).ExtractErr()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func createApigGroupCustomResponses(client *golangsdk.ServiceClient, instanceId, groupId string,
	resps *schema.Set) error {
	for _, response := range resps.List() {
		respMap := response.(map[string]interface{})
		respName := respMap["name"].(string)
		customRespRules := make(map[string]responses.ResponseInfo)
		for _, v := range respMap["responses"].(*schema.Set).List() {
			respRule := v.(map[string]interface{})
			errorType := respRule["error_type"].(string)
			customRespRules[errorType] = responses.ResponseInfo{
				Body:   respRule["body"].(string),
				Status: respRule["status_code"].(int),
			}
		}
		opt := responses.ResponseOpts{
			Name:      respName,
			Responses: customRespRules,
		}
		_, err := responses.Create(client, instanceId, groupId, opt).Extract()
		if err != nil {
			return err
		}
	}
	return nil
}

func deleteApigGroupCustomResponses(client *golangsdk.ServiceClient, instanceId, groupId string,
	resps *schema.Set) error {
	for _, response := range resps.List() {
		respMap := response.(map[string]interface{})
		respId := respMap["response_id"].(string)
		err := responses.Delete(client, instanceId, groupId, respId).ExtractErr()
		if err != nil {
			return err
		}
	}
	return nil
}

func resourceApigGroupV2Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	client, err := config.ApigV2Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud APIG v2 client: %s", err)
	}
	opt := apigroups.GroupOpts{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
	}
	log.Printf("[DEBUG] Create Options: %#v", opt)
	instanceId := d.Get("instance_id").(string)
	resp, err := apigroups.Create(client, instanceId, opt).Extract()
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud APIG group: %s", err)
	}
	d.SetId(resp.Id)

	if environments, ok := d.GetOk("environments"); ok {
		err = createApigGroupEnvironmentVariables(client, instanceId, d.Id(), environments.(*schema.Set))
		if err != nil {
			return fmtp.Errorf("Binding environment variables failed: %s", err)
		}
	}
	if customResps, ok := d.GetOk("custom_responses"); ok {
		err := createApigGroupCustomResponses(client, instanceId, d.Id(), customResps.(*schema.Set))
		if err != nil {
			return fmtp.Errorf("Creating custom responses failed: %s", err)
		}
	}
	return resourceApigGroupV2Read(d, meta)
}

func setApigGroupEnvironmentVariables(d *schema.ResourceData, variables []environments.Variable) error {
	environmentMap := make(map[string]interface{})
	for _, variable := range variables {
		if val, ok := environmentMap[variable.EnvId]; !ok {
			environment := make([]map[string]interface{}, 0)
			varMap := map[string]interface{}{
				"name":        variable.Name,
				"value":       variable.Value,
				"variable_id": variable.Id,
			}
			environmentMap[variable.EnvId] = append(environment, varMap)
		} else {
			varMap := map[string]interface{}{
				"name":        variable.Name,
				"value":       variable.Value,
				"variable_id": variable.Id,
			}
			environmentMap[variable.EnvId] = append(val.([]map[string]interface{}), varMap)
		}
	}
	result := make([]map[string]interface{}, 0, len(environmentMap))
	for k, v := range environmentMap {
		envMap := map[string]interface{}{
			"variables":      v,
			"environment_id": k,
		}
		result = append(result, envMap)
	}

	if len(result) == 0 {
		return d.Set("environments", nil)
	}
	return d.Set("environments", result)
}

func setApigGroupCustomResponses(d *schema.ResourceData, customResps []responses.Response) error {
	result := make([]map[string]interface{}, 0, len(customResps))
	for _, rules := range customResps {
		respRules := make([]map[string]interface{}, len(rules.Responses))
		for errorType, rule := range rules.Responses {
			if rule.IsDefault {
				// The IsDefault of the modified response will be marked as false,
				// record these responses and skip other unmodified responses (IsDefault is true).
				continue
			}
			ruleMap := map[string]interface{}{
				"error_type":  errorType,
				"body":        rule.Body,
				"status_code": rule.Status,
			}
			respRules = append(respRules, ruleMap)
		}
		if len(respRules) == 0 {
			continue
		}
		result = append(result, map[string]interface{}{
			"name":        rules.Name,
			"responses":   respRules,
			"response_id": rules.Id,
		})
	}
	if len(result) == 0 {
		return d.Set("custom_responses", nil)
	}
	return d.Set("custom_responses", result)
}

func setApigGroupParamters(d *schema.ResourceData, config *config.Config, resp *apigroups.Group) error {
	mErr := multierror.Append(nil,
		d.Set("region", config.GetRegion(d)),
		d.Set("name", resp.Name),
		d.Set("description", resp.Description),
		d.Set("registraion_time", resp.RegistraionTime),
		d.Set("update_time", resp.UpdateTime),
	)
	if mErr.ErrorOrNil() != nil {
		return mErr
	}
	return nil
}

func getApigGroupEnvironmentVariables(d *schema.ResourceData,
	client *golangsdk.ServiceClient) ([]environments.Variable, error) {
	instanceId := d.Get("instance_id").(string)
	listOpt := environments.ListVariablesOpts{
		GroupId: d.Id(),
	}
	pages, err := environments.ListVariables(client, instanceId, listOpt).AllPages()
	if err != nil {
		return []environments.Variable{}, fmtp.Errorf("Error getting environment variable list from server: %s", err)
	}
	result, err := environments.ExtractVariables(pages)
	if err != nil {
		return []environments.Variable{}, fmtp.Errorf("Error extract variables: %s", err)
	}
	return result, nil
}

func getApigGroupCustomResponses(d *schema.ResourceData,
	client *golangsdk.ServiceClient) ([]responses.Response, error) {
	instanceId := d.Get("instance_id").(string)
	pages, err := responses.List(client, instanceId, d.Id(), responses.ListOpts{}).AllPages()
	if err != nil {
		return []responses.Response{}, fmtp.Errorf("Error getting custom responses list from server: %s", err)
	}
	result, err := responses.ExtractResponses(pages)
	if err != nil {
		return []responses.Response{}, fmtp.Errorf("Error extract custom responses: %s", err)
	}
	return result, nil
}

func resourceApigGroupV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	client, err := config.ApigV2Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud APIG v2 client: %s", err)
	}
	instanceId := d.Get("instance_id").(string)
	resp, err := apigroups.Get(client, instanceId, d.Id()).Extract()
	if err != nil {
		return fmtp.Errorf("Error getting APIG v2 group: %s", err)
	}
	if err = setApigGroupParamters(d, config, resp); err != nil {
		return fmtp.Errorf("Error saving group to state: %s", err)
	}
	// Saving environment variables to state file.
	variables, err := getApigGroupEnvironmentVariables(d, client)
	if err != nil {
		return err
	}
	if err = setApigGroupEnvironmentVariables(d, variables); err != nil {
		return fmtp.Errorf("Error saving variables to state: %s", err)
	}
	// Saving custom responses to state file.
	responses, err := getApigGroupCustomResponses(d, client)
	if err != nil {
		return err
	}
	if err = setApigGroupCustomResponses(d, responses); err != nil {
		return fmtp.Errorf("Error saving environments to state: %s", err)
	}
	return nil
}

// Since set is implemented based on map, the update of each element will produce two operation processes:
// Delete and Add.
func updateApigResponsesCustomResps(d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	oldRaws, newRaws := d.GetChange("custom_responses")
	addRaws := newRaws.(*schema.Set).Difference(oldRaws.(*schema.Set))
	removeRaws := oldRaws.(*schema.Set).Difference(newRaws.(*schema.Set))
	instanceId := d.Get("instance_id").(string)
	for _, v := range removeRaws.List() {
		resp := v.(map[string]interface{})
		err := responses.Delete(client, instanceId, d.Id(), resp["response_id"].(string)).ExtractErr()
		if err != nil {
			return fmtp.Errorf("Failed to delete response (%s): %s", resp["response_id"].(string), err)
		}
	}
	err := createApigGroupCustomResponses(client, instanceId, d.Id(), addRaws)
	if err != nil {
		return fmtp.Errorf("Creating custom responses failed: %s", err)
	}

	return nil
}

func updateApigGroupEnvironmentVariables(d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	oldRaws, newRaws := d.GetChange("environments")
	addRaws := newRaws.(*schema.Set).Difference(oldRaws.(*schema.Set))
	removeRaws := oldRaws.(*schema.Set).Difference(newRaws.(*schema.Set))
	instanceId := d.Get("instance_id").(string)
	if err := removeApigGroupEnvironmentVariables(client, instanceId, removeRaws); err != nil {
		return err
	}
	if err := createApigGroupEnvironmentVariables(client, instanceId, d.Id(), addRaws); err != nil {
		return err
	}
	return nil
}

func resourceApigGroupV2Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	client, err := config.ApigV2Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud APIG v2 client: %s", err)
	}
	opt := apigroups.GroupOpts{}
	if d.HasChange("name") {
		opt.Name = d.Get("name").(string)
	}
	if d.HasChange("description") {
		opt.Description = d.Get("description").(string)
	}
	if opt != (apigroups.GroupOpts{}) {
		log.Printf("[DEBUG] Update Options: %#v", opt)
		instanceId := d.Get("instance_id").(string)
		_, err = apigroups.Update(client, instanceId, d.Id(), opt).Extract()
		if err != nil {
			return fmtp.Errorf("Error updating HuaweiCloud APIG group (%s): %s", d.Id(), err)
		}
	}
	if d.HasChange("environments") {
		if err := updateApigGroupEnvironmentVariables(d, client); err != nil {
			return fmtp.Errorf("Error updating HuaweiCloud APIG environment variables for the group (%s): %s", d.Id(), err)
		}
	}
	if d.HasChange("custom_responses") {
		err := updateApigResponsesCustomResps(d, client)
		if err != nil {
			return fmtp.Errorf("Creating custom responses failed: %s", err)
		}
	}
	return resourceApigGroupV2Read(d, meta)
}

func resourceApigGroupV2Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	client, err := config.ApigV2Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud APIG v2 client: %s", err)
	}
	instanceId := d.Get("instance_id").(string)
	err = apigroups.Delete(client, instanceId, d.Id()).ExtractErr()
	if err != nil {
		return fmtp.Errorf("Error deleting HuaweiCloud APIG group from the instance (%s): %s", instanceId, err)
	}
	d.SetId("")
	return nil
}
