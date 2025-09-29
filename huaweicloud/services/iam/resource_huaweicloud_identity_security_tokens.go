package iam

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/chnsz/golangsdk"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API IAM POST /v3.0/OS-CREDENTIAL/securitytokens
func ResourceIdentitySecurityTokens() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIdentitySecurityTokensCreate,
		ReadContext:   resourceIdentitySecurityTokensRead,
		DeleteContext: resourceUserTokenDelete,

		Schema: map[string]*schema.Schema{
			"token": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"methods": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Authentication method, the content of this field is either [token] or [assume_role].",
			},
			"policy": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"agency_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"domain_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"domain_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"duration_seconds": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"session_user_name": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"expires_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"access": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"secret": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"securitytoken": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceIdentitySecurityTokensCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	iamClient := common.NewCustomClient(true, "https://iam.{region_id}.myhuaweicloud.com")
	err := createSecurityTokens(iamClient, d, cfg)
	if err != nil {
		return diag.FromErr(err)
	}
	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(id)
	return nil

}

func buildCreateSecurityTokensBodyParams(d *schema.ResourceData) map[string]interface{} {
	methods := d.Get("methods").(string)
	policy := ""
	if val, ok := d.GetOk("policy"); ok {
		policy = val.(string)
	}
	bodyParams := map[string]interface{}{}
	if methods == "token" {
		bodyParams = map[string]interface{}{
			"identity": map[string]interface{}{
				"methods": []string{
					methods,
				},
				"token": map[string]interface{}{
					"duration_seconds": d.Get("duration_seconds").(int),
				},
				"policy": policy,
			},
		}
	} else {
		bodyParams = map[string]interface{}{
			"identity": map[string]interface{}{
				"methods": []string{
					methods,
				},
				"session_user": map[string]interface{}{
					"domain_name":       d.Get("domain_name"),
					"agency_name":       d.Get("agency_name"),
					"duration_seconds":  d.Get("duration_seconds").(int),
					"session_user_name": d.Get("session_user_name"),
					"domain_id":         d.Get("domain_id"),
				},
				"policy": policy,
			},
		}
	}
	return bodyParams
}

func createSecurityTokens(iamClient *golangsdk.ServiceClient, d *schema.ResourceData, cfg *config.Config) error {
	basePath := iamClient.ResourceBase + "/v3.0/OS-CREDENTIAL/securitytokens"
	basePath = strings.ReplaceAll(basePath, "{region_id}", cfg.GetRegion(d))
	createTokenOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"X-Auth-Token": d.Get("token").(string),
		},
		JSONBody: map[string]interface{}{
			"auth": buildCreateSecurityTokensBodyParams(d),
		},
	}
	response, err := iamClient.Request("POST", basePath, &createTokenOpt)
	if err != nil {
		return fmt.Errorf("error getting identity security token: %s", err)
	}
	resp, err := utils.FlattenResponse(response)
	if err != nil {
		return fmt.Errorf("error flattening response: %s", err)
	}
	mErr := multierror.Append(
		d.Set("expires_at", utils.PathSearch("credential.expires_at", resp, nil)),
		d.Set("access", utils.PathSearch("credential.access", resp, nil)),
		d.Set("secret", utils.PathSearch("credential.secret", resp, nil)),
		d.Set("securitytoken", utils.PathSearch("credential.securitytoken", resp, nil)),
	)
	return mErr.ErrorOrNil()
}

func resourceIdentitySecurityTokensRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	product := "identity"
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("Error creating IAM client: %s", err)
	}

	expiresAt, err := time.ParseInLocation(`2006-01-02T15:04:05Z`, d.Get("expires_at").(string), time.UTC)
	if err != nil {
		diag.Errorf("error parsing expires at: %s", err)
	}
	if time.Now().After(expiresAt) {
		err = createToken(client, d)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	return nil
}
