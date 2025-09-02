package iam

import (
	"context"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	iamv3 "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/iam/v3"
	iamModel "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/iam/v3/model"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// DataSourceIamIdentityProviderProtocols
// @API IAM GET /v3/OS-FEDERATION/identity_providers/{idp_id}/protocols
// @API IAM GET /v3/OS-FEDERATION/identity_providers/{idp_id}/protocols/{protocol_id}
func DataSourceIamIdentityProviderProtocols() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIdentityProviderProtocolsRead,

		Schema: map[string]*schema.Schema{
			"provider_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of an identity provider.",
			},
			"protocol_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"saml", "oidc"}, false),
				Description:  "The protocol id. The valid values are `saml` and `oidc`",
			},

			"protocols": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Protocol Information List.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Protocol Id",
						},
						"mapping_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Mapping Id",
						},
						"links": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `The links of protocol.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"self": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `resource link.`,
									},
									"identity_provider": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `identity provider resource link.`,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceIdentityProviderProtocolsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	iamV3Client, err := cfg.HcIamV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}
	protocolId := d.Get("protocol_id").(string)
	if protocolId == "" {
		return listProviderProtocols(iamV3Client, d)
	} else {
		return showProviderProtocol(iamV3Client, d)
	}
}

func listProviderProtocols(iamV3Client *iamv3.IamClient, d *schema.ResourceData) diag.Diagnostics {
	request := &iamModel.KeystoneListProtocolsRequest{IdpId: d.Get("provider_id").(string)}
	response, err := iamV3Client.KeystoneListProtocols(request)
	if err != nil {
		return diag.Errorf("KeystoneListProtocols error : %s", err)
	}
	result := make([]map[string]interface{}, len(*response.Protocols))
	for i, protocol := range *response.Protocols {
		result[i] = map[string]interface{}{
			"id":         protocol.Id,
			"mapping_id": protocol.MappingId,
		}
		links := make([]map[string]interface{}, 1)
		links[0] = map[string]interface{}{
			"self":              protocol.Links.Self,
			"identity_provider": protocol.Links.IdentityProvider,
		}
		result[i]["links"] = links
	}
	id, _ := uuid.GenerateUUID()
	d.SetId(id)
	return diag.FromErr(d.Set("protocols", result))
}

func showProviderProtocol(iamV3Client *iamv3.IamClient, d *schema.ResourceData) diag.Diagnostics {
	request := &iamModel.KeystoneShowProtocolRequest{
		IdpId:      d.Get("provider_id").(string),
		ProtocolId: d.Get("protocol_id").(string),
	}
	response, err := iamV3Client.KeystoneShowProtocol(request)
	if err != nil {
		return diag.Errorf("KeystoneShowProtocol error : %s", err)
	}
	result := make([]map[string]interface{}, 1)
	result[0] = map[string]interface{}{
		"id":         response.Protocol.Id,
		"mapping_id": response.Protocol.MappingId,
	}
	links := make([]map[string]interface{}, 1)
	links[0] = map[string]interface{}{
		"self":              response.Protocol.Links.Self,
		"identity_provider": response.Protocol.Links.IdentityProvider,
	}
	result[0]["links"] = links
	id, _ := uuid.GenerateUUID()
	d.SetId(id)
	return diag.FromErr(d.Set("protocols", result))
}
