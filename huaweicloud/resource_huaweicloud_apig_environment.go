package huaweicloud

import (
	"fmt"
	"log"
	"regexp"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/apigw/v2/environments"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func ResourceApigEnvironmentV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceApigEnvironmentV2Create,
		Read:   resourceApigEnvironmentV2Read,
		Update: resourceApigEnvironmentV2Update,
		Delete: resourceApigEnvironmentV2Delete,
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
					regexp.MustCompile("^[A-Za-z][\\w0-9]{2,63}$"),
					"The name consists of 3 to 64 characters, starting with a letter. "+
						"Only letters, digits and underscores (_) are allowed."),
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile("^[^<>]{1,255}$"),
					"The description contain a maximum of 255 characters, "+
						"and the angle brackets (< and >) are not allowed."),
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildApigEnvironmentParameters(d *schema.ResourceData) environments.EnvironmentOpts {
	return environments.EnvironmentOpts{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
	}
}

func resourceApigEnvironmentV2Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	client, err := config.ApigV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud APIG v2 client: %s", err)
	}
	opts := buildApigEnvironmentParameters(d)
	log.Printf("[DEBUG] Create Options: %#v", opts)
	instanceId := d.Get("instance_id").(string)
	resp, err := environments.Create(client, instanceId, opts).Extract()
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud APIG environment: %s", err)
	}
	d.SetId(resp.Id)
	return resourceApigEnvironmentV2Read(d, meta)
}

func setApigEnvironmentParamters(d *schema.ResourceData, config *config.Config, resp *environments.Environment) error {
	mErr := multierror.Append(nil,
		d.Set("region", GetRegion(d, config)),
		d.Set("name", resp.Name),
		d.Set("description", resp.Description),
		d.Set("create_time", resp.CreateTime),
	)
	if mErr.ErrorOrNil() != nil {
		return mErr
	}
	return nil
}

func getEnvironmentFormServer(client *golangsdk.ServiceClient, instanceId, envId string) (*environments.Environment, error) {
	allPages, err := environments.List(client, instanceId, environments.ListOpts{}).AllPages()
	if err != nil {
		return nil, err
	}
	envs, err := environments.ExtractEnvironments(allPages)
	if err != nil {
		return nil, err
	}
	for _, v := range envs {
		if v.Id == envId {
			return &v, nil
		}
	}
	return nil, fmt.Errorf("The environment does not exist")
}

func resourceApigEnvironmentV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	client, err := config.ApigV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud APIG v2 client: %s", err)
	}
	instanceId := d.Get("instance_id").(string)
	env, err := getEnvironmentFormServer(client, instanceId, d.Id())
	if err != nil {
		return fmt.Errorf("Unable to get the environment (%s) form server: %s", d.Id(), err)
	}
	return setApigEnvironmentParamters(d, config, env)
}

func resourceApigEnvironmentV2Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	client, err := config.ApigV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud APIG v2 client: %s", err)
	}
	opt := environments.EnvironmentOpts{}
	if d.HasChange("name") {
		opt.Name = d.Get("name").(string)
	}
	if d.HasChange("description") {
		opt.Description = d.Get("description").(string)
	}
	log.Printf("[DEBUG] Update Options: %#v", opt)
	instanceId := d.Get("instance_id").(string)
	_, err = environments.Update(client, instanceId, d.Id(), opt).Extract()
	if err != nil {
		return fmt.Errorf("Error updating HuaweiCloud APIG environment (%s): %s", d.Id(), err)
	}

	return resourceApigEnvironmentV2Read(d, meta)
}

func resourceApigEnvironmentV2Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	client, err := config.ApigV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud APIG v2 client: %s", err)
	}
	instanceId := d.Get("instance_id").(string)
	err = environments.Delete(client, instanceId, d.Id()).ExtractErr()
	if err != nil {
		return fmt.Errorf("Error deleting HuaweiCloud APIG environment from the instance (%s): %s", instanceId, err)
	}
	d.SetId("")
	return nil
}
