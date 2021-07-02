package huaweicloud

import (
	"github.com/huaweicloud/golangsdk/openstack/apigw/v2/channels"
	"regexp"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

var (
	balanceStrategy = map[string]int{
		"WRR":         1,
		"WLC":         2,
		"SH":          3,
		"URI hashing": 4,
	}
	channelStatus = map[int]string{
		1: "Normal",
		2: "Abnormal",
	}
)

func ResourceApigVpcChannelV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceApigVpcChannelV2Create,
		Read:   resourceApigVpcChannelV2Read,
		Update: resourceApigVpcChannelV2Update,
		Delete: resourceApigVpcChannelV2Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
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
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile("^([\u4e00-\u9fa5A-Za-z][\u4e00-\u9fa5A-Za-z-_0-9]{2,63})$"),
					"The name contains of 3 to 64 characters, starting with a letter. Only letters, digits, "+
						"hyphens (-) and underscore (_) are allowed."),
			},
			"port": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(1, 65535),
			},
			"member_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "ecs",
				ValidateFunc: validation.StringInSlice([]string{
					"ecs", "ip",
				}, false),
			},
			"algorithm": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "WRR",
				ValidateFunc: validation.StringInSlice([]string{
					"WRR", "WLC", "SH", "URI hashing",
				}, false),
			},
			"protocol": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "TCP",
				ValidateFunc: validation.StringInSlice([]string{
					"TCP", "HTTP", "HTTPS",
				}, false),
			},
			"path": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"healthy_threshold": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      2,
				ValidateFunc: validation.IntBetween(2, 10),
			},
			"unhealthy_threshold": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      5,
				ValidateFunc: validation.IntBetween(2, 10),
			},
			"timeout": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      5,
				ValidateFunc: validation.IntBetween(2, 30),
			},
			"interval": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      10,
				ValidateFunc: validation.IntBetween(5, 300),
			},
			"http_codes": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"members": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 10,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip_address": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"weight": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      1,
							ValidateFunc: validation.IntBetween(1, 100),
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

func buildApigVpcChannelHealthConfig(d *schema.ResourceData) (channels.VpcHealthConfig, error) {
	conf := channels.VpcHealthConfig{
		Protocol:          d.Get("protocol").(string),
		Path:              d.Get("path").(string),
		Port:              d.Get("port").(int),
		ThresholdNormal:   d.Get("healthy_threshold").(int),
		ThresholdAbnormal: d.Get("unhealthy_threshold").(int),
		Timeout:           d.Get("timeout").(int),
		TimeInterval:      d.Get("interval").(int),
	}
	// The parameter of http codes is required if protocol is set to http or https.
	if val := d.Get("protocol"); val.(string) == "HTTP" || val.(string) == "HTTPS" {
		if codes, ok := d.GetOk("http_codes"); ok {
			conf.HttpCodes = codes.(string)
		} else {
			return conf, fmtp.Errorf("The http code cannot be empty if protocol is http or https")
		}
	}
	return conf, nil
}

func buildApigVpcChannelMembers(d *schema.ResourceData) ([]channels.MemberInfo, error) {
	members := d.Get("members").([]interface{})
	memberType := d.Get("member_type").(string)
	result := make([]channels.MemberInfo, len(members))
	for i, v := range members {
		member := v.(map[string]interface{})
		info := channels.MemberInfo{
			Weight: member["weight"].(int),
		}
		if id, ok := member["id"]; memberType == "ecs" && ok {
			info.EcsId = id.(string)
		} else if addr, ok := member["ip_address"]; memberType == "ip" && ok {
			info.Host = addr.(string)
		} else {
			return result, fmtp.Errorf("The members are wrong, please check whether your member" +
				" type corresponds to the parameters in members.")
		}
		result[i] = info
	}
	return result, nil
}

func buildApigVpcChannelParameters(d *schema.ResourceData) (channels.ChannelOpts, error) {
	opt := channels.ChannelOpts{
		Name:       d.Get("name").(string),
		Port:       d.Get("port").(int),
		MemberType: d.Get("member_type").(string),
		Type:       2, // The type 1 (private network ELB channel) is to be deprecated.
	}
	// Backend servers
	members, err := buildApigVpcChannelMembers(d)
	if err != nil {
		return opt, err
	}
	opt.Members = members
	// Healthy check config
	conf, err := buildApigVpcChannelHealthConfig(d)
	if err != nil {
		return opt, err
	}
	opt.VpcHealthConfig = conf
	// algorithm
	v, ok := balanceStrategy[d.Get("algorithm").(string)]
	if ok {
		opt.BalanceStrategy = v
	} else {
		return opt, fmtp.Errorf("The value of algorithm is invalid")
	}
	return opt, nil
}

func resourceApigVpcChannelV2Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	opts, err := buildApigVpcChannelParameters(d)
	if err != nil {
		return fmtp.Errorf("Error craeting APIG v2 dedicated instance options: %s", err)
	}
	logp.Printf("[DEBUG] Create APIG v2 dedicated instance options: %#v", opts)

	client, err := config.ApigV2Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud APIG v2 client: %s", err)
	}
	instanceId := d.Get("instance_id").(string)
	v, err := channels.Create(client, instanceId, opts).Extract()
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud APIG v2 dedicated instance: %s", err)
	}
	d.SetId(v.Id)
	return resourceApigInstanceV2Read(d, meta)
}

func getApigVpcChannelAlgorithm(resp channels.VpcChannel) (string, error) {
	algorithm := resp.BalanceStrategy
	for k, v := range balanceStrategy {
		if v == algorithm {
			return k, nil
		}
	}
	return "", fmtp.Errorf("Unable to extract the algorithm")
}

func getApigVpcChannelMembers(resp channels.VpcChannel) []map[string]interface{} {
	members := make([]map[string]interface{}, len(resp.Members))
	for i, v := range resp.Members {
		members[i] = map[string]interface{}{
			"ip_address": v.Host,
			"id":         v.EcsId,
			"weight":     v.Weight,
		}
	}
	return members
}

func setApigVpcChannelParameters(d *schema.ResourceData, config *config.Config, resp channels.VpcChannel) error {
	algorithm, err := getApigVpcChannelAlgorithm(resp)
	if err != nil {
		return err
	}
	status, ok := channelStatus[resp.Status]
	if !ok {
		return fmtp.Errorf("The response status is invalid")
	}
	mErr := multierror.Append(nil,
		d.Set("region", GetRegion(d, config)),
		d.Set("name", resp.Name),
		d.Set("port", resp.Port),
		d.Set("member_type", resp.MemberType),
		d.Set("algorithm", algorithm),
		d.Set("protocol", resp.VpcHealthConfig.Protocol),
		d.Set("path", resp.VpcHealthConfig.Protocol),
		d.Set("healthy_threshold", resp.VpcHealthConfig.ThresholdNormal),
		d.Set("unhealthy_threshold", resp.VpcHealthConfig.ThresholdAbnormal),
		d.Set("timeout", resp.VpcHealthConfig.Timeout),
		d.Set("interval", resp.VpcHealthConfig.TimeInterval),
		d.Set("http_codes", resp.VpcHealthConfig.HttpCodes),
		d.Set("members", getApigVpcChannelMembers(resp)),
		d.Set("create_time", resp.CreateTime),
		d.Set("status", status),
	)
	if mErr.ErrorOrNil() != nil {
		return mErr
	}
	return nil
}

func resourceApigVpcChannelV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	client, err := config.ApigV2Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud APIG client: %s", err)
	}
	instanceId := d.Get("instance_id").(string)
	resp, err := channels.Get(client, instanceId, d.Id()).Extract()
	if err != nil {
		return fmtp.Errorf("Error getting APIG v2 dedicated instance (%s) form server: %s", d.Id(), err)
	}
	return setApigVpcChannelParameters(d, config, *resp)
}

func resourceApigVpcChannelV2Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	client, err := config.ApigV2Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud APIG v2 client: %s", err)
	}

	opt, err := buildApigVpcChannelParameters(d)
	if err != nil {
		return fmtp.Errorf("Unable to get the update options of APIG v2 vpc channel: %s", err)
	}
	instanceId := d.Get("instance_id").(string)
	_, err = channels.Update(client, instanceId, d.Id(), opt).Extract()
	if err != nil {
		return fmtp.Errorf("Error updating APIG v2 vpc channel: %s", err)
	}
	return resourceApigInstanceV2Read(d, meta)
}

func resourceApigVpcChannelV2Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	client, err := config.ApigV2Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud APIG v2 client: %s", err)
	}
	instanceId := d.Get("instance_id").(string)
	if err = channels.Delete(client, instanceId, d.Id()).ExtractErr(); err != nil {
		return fmtp.Errorf("Unable to delete the APIG v2 vpc channel (%s): %s", d.Id(), err)
	}
	d.SetId("")
	return nil
}
