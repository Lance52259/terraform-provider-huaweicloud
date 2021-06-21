package huaweicloud

import (
	"fmt"
	"log"
	"regexp"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/huaweicloud/golangsdk/openstack/apigw/v2/apigroups"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func ResourceApigApiGroupV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceApigApiGroupV2Create,
		Read:   resourceApigApiGroupV2Read,
		Update: resourceApigApiGroupV2Update,
		Delete: resourceApigApiGroupV2Delete,
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

func buildApigApiGroupParameters(d *schema.ResourceData) apigroups.GroupOpts {
	return apigroups.GroupOpts{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
	}
}

func resourceApigApiGroupV2Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	client, err := config.ApigV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud APIG v2 client: %s", err)
	}
	opts := buildApigApiGroupParameters(d)
	log.Printf("[DEBUG] Create Options: %#v", opts)
	instanceId := d.Get("instance_id").(string)
	resp, err := apigroups.Create(client, instanceId, opts).Extract()
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud APIG group: %s", err)
	}
	d.SetId(resp.Id)
	return resourceApigApiGroupV2Read(d, meta)
}

func setApigApiGroupParamters(d *schema.ResourceData, config *config.Config, resp *apigroups.Group) error {
	mErr := multierror.Append(nil,
		d.Set("region", GetRegion(d, config)),
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

func resourceApigApiGroupV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	client, err := config.ApigV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud APIG v2 client: %s", err)
	}
	instanceId := d.Get("instance_id").(string)
	resp, err := apigroups.Get(client, instanceId, d.Id()).Extract()

	return setApigApiGroupParamters(d, config, resp)
}

func resourceApigApiGroupV2Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	client, err := config.ApigV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud APIG v2 client: %s", err)
	}
	opt := apigroups.GroupOpts{}
	if d.HasChange("name") {
		opt.Name = d.Get("name").(string)
	}
	if d.HasChange("description") {
		opt.Description = d.Get("description").(string)
	}
	log.Printf("[DEBUG] Update Options: %#v", opt)
	instanceId := d.Get("instance_id").(string)
	_, err = apigroups.Update(client, instanceId, d.Id(), opt).Extract()
	if err != nil {
		return fmt.Errorf("Error updating HuaweiCloud APIG group (%s): %s", d.Id(), err)
	}

	return resourceApigApiGroupV2Read(d, meta)
}

func resourceApigApiGroupV2Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	client, err := config.ApigV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud APIG v2 client: %s", err)
	}
	instanceId := d.Get("instance_id").(string)
	err = apigroups.Delete(client, instanceId, d.Id()).ExtractErr()
	if err != nil {
		return fmt.Errorf("Error deleting HuaweiCloud APIG group from the instance (%s): %s", instanceId, err)
	}
	d.SetId("")
	return nil
}
