package iam

import (
	"context"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	iamModel "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/iam/v3/model"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// @API GET /v3/auth/tokens
func DataSourceIdentityToken() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIdentityTokenRead,

		Schema: map[string]*schema.Schema{
			"x_subject_token": {
				Type:     schema.TypeString,
				Required: true,
			},
			"no_catalog": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"token": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"expires_at": {

							Type:     schema.TypeString,
							Computed: true,
						},
						"issued_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"methods": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"catalog": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"endpoints": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"id": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"interface": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"region": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"region_id": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"url": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"type": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"domain": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"project": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"domain": {
										Type:     schema.TypeMap,
										Computed: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"roles": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"user": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"domain": {
										Type:     schema.TypeMap,
										Computed: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"password_expires_at": {
										Type:     schema.TypeString,
										Computed: true,
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

func dataSourceIdentityTokenRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	iamV3Client, err := cfg.HcIamV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}
	xSubjectToken := d.Get("x_subject_token").(string)
	if xSubjectToken == "" {
		return diag.Errorf("missing required parameters x_subject_token")
	}
	request := &iamModel.KeystoneValidateTokenRequest{
		XSubjectToken: xSubjectToken,
	}

	noCatalog := d.Get("no_catalog")
	noCatalogStr := ""
	if noCatalog != nil {
		noCatalogStr = noCatalog.(string)
		if noCatalogStr != "" {
			request.Nocatalog = &noCatalogStr
		}
	}
	response, err := iamV3Client.KeystoneValidateToken(request)
	if err != nil {
		return diag.Errorf("KeystoneValidateToken error : %s", err)
	}
	result := make([]map[string]interface{}, 1)
	result[0] = map[string]interface{}{
		"expires_at": response.Token.ExpiresAt,
		"issued_at":  response.Token.IssuedAt,
		"methods":    response.Token.Methods,
		//"catalog":    flattenCatalog(response.Token.Catalog),
		"domain":  flattenDomain(response.Token.Domain),
		"project": flattenProject(response.Token.Project),
		"roles":   flattenRoles(response.Token.Roles),
		"user":    flattenUser(response.Token.User),
	}
	if noCatalogStr == "" {
		result[0]["catalog"] = flattenCatalog(response.Token.Catalog)
	}
	id, _ := uuid.GenerateUUID()
	d.SetId(id)
	return diag.FromErr(d.Set("token", result))
}

func flattenCatalog(catalogs []iamModel.TokenCatalog) []map[string]interface{} {
	res := make([]map[string]interface{}, len(catalogs))
	for i, catalog := range catalogs {
		res[i] = map[string]interface{}{
			"endpoints": flattenEndpoints(catalog.Endpoints),
			"id":        catalog.Id,
			"name":      catalog.Name,
			"type":      catalog.Type,
		}
	}
	return res
}

func flattenEndpoints(endpoints []iamModel.TokenCatalogEndpoint) []map[string]interface{} {
	res := make([]map[string]interface{}, len(endpoints))
	for i, endpoint := range endpoints {
		res[i] = map[string]interface{}{
			"id":        endpoint.Id,
			"interface": endpoint.Interface,
			"region":    endpoint.Region,
			"region_id": endpoint.RegionId,
			"url":       endpoint.Url,
		}
	}
	return res
}

func flattenDomain(domain *iamModel.TokenDomainResult) map[string]interface{} {
	if domain == nil {
		return nil
	}
	res := map[string]interface{}{
		"id":   domain.Id,
		"name": domain.Name,
	}
	return res
}

func flattenProject(projects *iamModel.TokenProjectResult) []map[string]interface{} {
	if projects == nil {
		return nil
	}
	res := []map[string]interface{}{
		{
			"domain": flattenProjectDomain(projects.Domain),
			"id":     projects.Id,
			"name":   projects.Name,
		},
	}
	return res
}

func flattenProjectDomain(domain *iamModel.TokenProjectDomainResult) map[string]interface{} {
	if domain == nil {
		return nil
	}
	res := map[string]interface{}{
		"id":   domain.Id,
		"name": domain.Name,
	}
	return res
}

func flattenRoles(roles []iamModel.TokenRole) []map[string]interface{} {
	res := make([]map[string]interface{}, len(roles))
	for i, role := range roles {
		res[i] = map[string]interface{}{
			"id":   role.Id,
			"name": role.Name,
		}
	}
	return res
}

func flattenUser(user *iamModel.TokenUserResult) []map[string]interface{} {
	if user == nil {
		return nil
	}
	res := []map[string]interface{}{
		{
			"domain":              flattenUserDomain(user.Domain),
			"id":                  user.Id,
			"name":                user.Name,
			"password_expires_at": user.PasswordExpiresAt,
		},
	}
	return res
}

func flattenUserDomain(domain *iamModel.TokenUserDomainResult) map[string]interface{} {
	if domain == nil {
		return nil
	}
	res := map[string]interface{}{
		"id":   domain.Id,
		"name": domain.Name,
	}
	return res
}
