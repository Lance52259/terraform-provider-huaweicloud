package rgc

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API RGC GET /v1/governance/managed-accounts/{managed_account_id}/external-config-rule-compliances
func DataSourceExternalConfigRuleCompliances() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceExternalConfigRuleCompliancesRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"managed_account_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"config_rule_compliances": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"rule_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"control_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceExternalConfigRuleCompliancesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var listConfigRuleCompliancesProduct = "rgc"
	listConfigRuleCompliancesClient, err := cfg.NewServiceClient(listConfigRuleCompliancesProduct, region)
	if err != nil {
		return diag.Errorf("error creating RGC client: %s", err)
	}

	listConfigRuleCompliancesRespBody, err := listExternalConfigRuleCompliances(listConfigRuleCompliancesClient, d)

	if err != nil {
		return diag.Errorf("error retrieving RGC external config rule compliance: %s", err)
	}

	randUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randUUID)

	mErr = multierror.Append(mErr,
		d.Set("config_rule_compliances", utils.PathSearch("config_rule_compliances", listConfigRuleCompliancesRespBody, nil)),
		d.Set("region", region),
	)

	return diag.FromErr(mErr.ErrorOrNil())

}

func listExternalConfigRuleCompliances(client *golangsdk.ServiceClient, d *schema.ResourceData) (interface{}, error) {
	managedAccountId := d.Get("managed_account_id").(string)
	var (
		listExternalConfigRuleCompliancesHttpUrl = "v1/governance/managed-accounts/{managed_account_id}/external-config-rule-compliances"
	)
	listExternalConfigRuleCompliancesHttpPath := client.Endpoint + listExternalConfigRuleCompliancesHttpUrl
	listExternalConfigRuleCompliancesHttpPath = strings.ReplaceAll(listExternalConfigRuleCompliancesHttpPath, "{managed_account_id}", managedAccountId)

	listConfigRuleCompliancesHttpOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	listExternalConfigRuleCompliancesHttpResp, err := client.Request("GET", listExternalConfigRuleCompliancesHttpPath, &listConfigRuleCompliancesHttpOpt)
	if err != nil {
		return nil, err
	}
	listExternalConfigRuleCompliancesRespBody, err := utils.FlattenResponse(listExternalConfigRuleCompliancesHttpResp)
	if err != nil {
		return nil, err
	}
	return listExternalConfigRuleCompliancesRespBody, nil
}
