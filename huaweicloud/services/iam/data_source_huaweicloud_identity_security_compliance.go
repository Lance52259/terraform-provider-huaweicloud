package iam

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	iamv3 "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/iam/v3"
	iamModel "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/iam/v3/model"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// DataSourceIdentitySecurityCompliance
// @API IAM GET /v3/domains/{domain_id}/config/security_compliance
// @API IAM GET /v3/domains/{domain_id}/config/security_compliance/{option}
func DataSourceIdentitySecurityCompliance() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIdentitySecurityComplianceRead,

		Schema: map[string]*schema.Schema{
			"option": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"password_regex", "password_regex_description"}, false),
				Description:  "Specifies the query type, which supports password_regex or password_regex_description.",
			},

			"password_regex": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"password_regex_description": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceIdentitySecurityComplianceRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	iamV3Client, err := cfg.HcIamV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}
	optionStr := d.Get("option").(string)
	if optionStr == "" {
		return showSecurityCompliance(iamV3Client, d, cfg)
	} else {
		return showSecurityComplianceByOption(iamV3Client, optionStr, d, cfg)
	}
}

func showSecurityCompliance(iamV3Client *iamv3.IamClient, d *schema.ResourceData, cfg *config.Config) diag.Diagnostics {
	request := &iamModel.KeystoneShowSecurityComplianceRequest{DomainId: cfg.DomainID}
	response, err := iamV3Client.KeystoneShowSecurityCompliance(request)
	if err != nil {
		return diag.Errorf("showSecurityCompliance error : %s", err)
	}
	setIdWithUUID(d)
	mErr := multierror.Append(nil,
		d.Set("password_regex", response.Config.SecurityCompliance.PasswordRegex),
		d.Set("password_regex_description", response.Config.SecurityCompliance.PasswordRegexDescription),
	)
	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting SecurityCompliance fields: %s", err)
	}
	return nil
}

func showSecurityComplianceByOption(iamV3Client *iamv3.IamClient, optionStr string, d *schema.ResourceData, cfg *config.Config) diag.Diagnostics {
	var option iamModel.KeystoneShowSecurityComplianceByOptionRequestOption
	if optionStr == "password_regex" {
		option = iamModel.GetKeystoneShowSecurityComplianceByOptionRequestOptionEnum().PASSWORD_REGEX
	} else if optionStr == "password_regex_description" {
		option = iamModel.GetKeystoneShowSecurityComplianceByOptionRequestOptionEnum().PASSWORD_REGEX_DESCRIPTION
	} else {
		return diag.Errorf("%s option is not supportted", optionStr)
	}
	request := &iamModel.KeystoneShowSecurityComplianceByOptionRequest{DomainId: cfg.DomainID, Option: option}
	response, err := iamV3Client.KeystoneShowSecurityComplianceByOption(request)
	if err != nil {
		return diag.Errorf("showSecurityComplianceByOption error : %s", err)
	}
	setIdWithUUID(d)
	mErr := multierror.Append(nil,
		d.Set("password_regex", response.Config.PasswordRegex),
		d.Set("password_regex_description", response.Config.PasswordRegexDescription),
	)
	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting SecurityComplianceByOption fields: %s", err)
	}
	return nil
}

func setIdWithUUID(d *schema.ResourceData) diag.Diagnostics {
	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(id)
	return nil
}
