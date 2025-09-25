package iam

import (
	"context"
	"fmt"

	"github.com/chnsz/golangsdk"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

const perPage = 50

// @API IAM GET /v3.0/OS-PERMISSION/role-assignments
func DataSourceIdentityRoleAssignments() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceIdentityRoleAssignmentsRead,

		Schema: map[string]*schema.Schema{
			"role_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"subject": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"subject_user_id": {
				Type:     schema.TypeString,
				Optional: true,
				ConflictsWith: []string{
					"subject_group_id",
					"subject_agency_id",
				},
			},
			"subject_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				ConflictsWith: []string{
					"subject_user_id",
					"subject_agency_id",
				},
			},
			"subject_agency_id": {
				Type:     schema.TypeString,
				Optional: true,
				ConflictsWith: []string{
					"subject_group_id",
					"subject_user_id",
				},
			},
			"scope": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "domain",
			},
			"scope_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				ConflictsWith: []string{
					"scope_domain_id",
					"scope_enterprise_projects_id",
				},
			},
			"scope_domain_id": {
				Type:     schema.TypeString,
				Optional: true,
				ConflictsWith: []string{
					"scope_project_id",
					"scope_enterprise_projects_id",
				},
			},
			"scope_enterprise_projects_id": {
				Type:     schema.TypeString,
				Optional: true,
				ConflictsWith: []string{
					"scope_project_id",
					"scope_domain_id",
				},
			},
			"is_inherited": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"include_group": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"role_assignments": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"role_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"agency_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"is_inherited": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func DataSourceIdentityRoleAssignmentsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	domainId := cfg.DomainID
	iamClient, err := cfg.IAMV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}
	listRecordPath := iamClient.Endpoint + "v3.0/OS-PERMISSION/role-assignments"
	listRecordPath += fmt.Sprintf("?per_page=%v", perPage)
	path := buildQueryRoleAssignmentsPath(d, listRecordPath, domainId)
	options := golangsdk.RequestOpts{KeepResponseBody: true}
	currentPage := 1
	result := make([]map[string]interface{}, 0)
	for {
		getPath := path + fmt.Sprintf("&page=%d", currentPage)
		response, err := iamClient.Request("GET", getPath, &options)
		if err != nil {
			return diag.Errorf("error listRoleAssignments: %s", err)
		}
		respBody, err := utils.FlattenResponse(response)
		result = flattenRoleAssignments(utils.PathSearch("role_assignments", respBody, nil).([]interface{}), result)
		totalNum := utils.PathSearch("total_num", respBody, nil).(float64)
		if float64(len(result)) >= totalNum {
			break
		}
		currentPage++
	}
	id, _ := uuid.GenerateUUID()
	d.SetId(id)
	if err = d.Set("role_assignments", result); err != nil {
		return diag.Errorf("error setting total_num fields: %s", err)
	}
	return nil
}

func flattenRoleAssignments(assignmentsBody []interface{}, results []map[string]interface{}) []map[string]interface{} {
	for _, assignment := range assignmentsBody {
		results = append(results, map[string]interface{}{
			"role_id":      utils.PathSearch("id", utils.PathSearch("role", assignment, nil), nil),
			"is_inherited": utils.PathSearch("is_inherited", assignment, nil),
			"user_id":      utils.PathSearch("id", utils.PathSearch("user", assignment, nil), nil),
			"group_id":     utils.PathSearch("id", utils.PathSearch("group", assignment, nil), nil),
			"agency_id":    utils.PathSearch("id", utils.PathSearch("agency", assignment, nil), nil),
		})
	}
	return results
}

func buildQueryRoleAssignmentsPath(d *schema.ResourceData, getPath string, domainId string) string {
	getPath += fmt.Sprintf("&domain_id=%s", domainId)
	if roleId := d.Get("role_id").(string); roleId != "" {
		getPath += fmt.Sprintf("&role_id=%s", roleId)
	}
	subject := d.Get("subject").(string)
	if subject != "" {
		getPath += fmt.Sprintf("&subject=%s", subject)
	}
	subjectUserId := d.Get("subject_user_id").(string)
	if subjectUserId != "" {
		getPath += fmt.Sprintf("&subject_user_id=%s", subjectUserId)
	}
	if subjectGroupId := d.Get("subject_group_id").(string); subjectGroupId != "" {
		getPath += fmt.Sprintf("&subject_group_id=%s", subjectGroupId)
	}
	if subjectAgencyId := d.Get("subject_agency_id").(string); subjectAgencyId != "" {
		getPath += fmt.Sprintf("&subject_agency_id=%s", subjectAgencyId)
	}
	scope := d.Get("scope").(string)
	if scope != "" {
		getPath += fmt.Sprintf("&scope=%s", scope)
	}
	if scopeProjectId := d.Get("scope_project_id").(string); scopeProjectId != "" {
		getPath += fmt.Sprintf("&scope_project_id=%s", scopeProjectId)
	}
	scopeDomainId := d.Get("scope_domain_id").(string)
	if scopeDomainId != "" {
		getPath += fmt.Sprintf("&scope_domain_id=%s", scopeDomainId)
	}
	if scopeEnterpriseProjectsId := d.Get("scope_enterprise_projects_id").(string); scopeEnterpriseProjectsId != "" {
		getPath += fmt.Sprintf("&scope_enterprise_projects_id=%s", scopeEnterpriseProjectsId)
	}
	if isInherited := d.Get("is_inherited").(bool); scope == "domain" || scopeDomainId != "" {
		getPath += fmt.Sprintf("&is_inherited=%t", isInherited)
	}
	if includeGroup := d.Get("include_group").(bool); subject == "user" || subjectUserId != "" {
		getPath += fmt.Sprintf("&include_group=%t", includeGroup)
	}
	return getPath
}
