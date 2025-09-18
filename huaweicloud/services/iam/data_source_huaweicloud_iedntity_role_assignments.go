package iam

import (
	"context"
	"github.com/chnsz/golangsdk"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
	"net/url"
	"strconv"
)

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
				ValidateFunc: validation.StringInSlice([]string{
					"user",
					"group",
					"agency",
				}, false),
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
				ValidateFunc: validation.StringInSlice([]string{
					"project",
					"domain",
					"enterprise_project",
				}, false),
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
			"page": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"per_page": {
				Type:     schema.TypeInt,
				Optional: true,
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
	u, err := url.Parse(listRecordPath)
	if err != nil {
		return diag.Errorf("error update query path: %s", err)
	}
	query := u.Query()
	query.Add("domain_id", domainId)
	if roleId := d.Get("role_id").(string); roleId != "" {
		query.Add("role_id", roleId)
	}
	subject := d.Get("subject").(string)
	if subject != "" {
		query.Add("subject", subject)
	}
	subjectUserId := d.Get("subject_user_id").(string)
	if subjectUserId != "" {
		query.Add("subject.user_id", subjectUserId)
	}
	if subjectGroupId := d.Get("subject_group_id").(string); subjectGroupId != "" {
		query.Add("subject.group_id", subjectGroupId)
	}
	if subjectAgencyId := d.Get("subject_agency_id").(string); subjectAgencyId != "" {
		query.Add("subject.agency_id", subjectAgencyId)
	}
	scope := d.Get("scope").(string)
	if scope != "" {
		query.Add("scope", scope)
	}
	if scopeProjectId := d.Get("scope_project_id").(string); scopeProjectId != "" {
		query.Add("scope.project_id", scopeProjectId)
	}
	scopeDomainId := d.Get("scope_domain_id").(string)
	if scopeDomainId != "" {
		query.Add("scope.domain_id", scopeDomainId)
	}
	if scopeEnterpriseProjectsId := d.Get("scope_enterprise_projects_id").(string); scopeEnterpriseProjectsId != "" {
		query.Add("scope.enterprise_projects_id", scopeEnterpriseProjectsId)
	}
	if isInherited := d.Get("is_inherited").(bool); scope == "domain" || scopeDomainId != "" {
		query.Add("is_inherited", strconv.FormatBool(isInherited))
	}
	if includeGroup := d.Get("include_group").(bool); subject == "user" || subjectUserId != "" {
		query.Add("include_group", strconv.FormatBool(includeGroup))
	}
	if page := d.Get("page").(int); page > 0 {
		query.Add("page", strconv.Itoa(page))
	}
	if perPage := d.Get("per_page").(int); perPage > 0 {
		query.Add("per_page", strconv.Itoa(perPage))
	}
	u.RawQuery = query.Encode()
	options := golangsdk.RequestOpts{KeepResponseBody: true}
	response, err := iamClient.Request("GET", u.String(), &options)
	if err != nil {
		return diag.Errorf("error listRoleAssignments: %s", err)
	}
	respBody, err := utils.FlattenResponse(response)
	id, _ := uuid.GenerateUUID()
	d.SetId(id)
	if err = d.Set("role_assignments", showRoleAssignments(utils.PathSearch("role_assignments", respBody, nil).([]interface{}))); err != nil {
		return diag.Errorf("error setting total_num fields: %s", err)
	}
	return nil
}

func showRoleAssignments(assignmentsBody []interface{}) []map[string]interface{} {
	assignments := make([]map[string]interface{}, len(assignmentsBody))
	for i, assignment := range assignmentsBody {
		assignments[i] = map[string]interface{}{
			"role_id":      utils.PathSearch("id", utils.PathSearch("role", assignment, nil), nil),
			"is_inherited": utils.PathSearch("is_inherited", assignment, nil),
			"user_id":      utils.PathSearch("id", utils.PathSearch("user", assignment, nil), nil),
			"group_id":     utils.PathSearch("id", utils.PathSearch("group", assignment, nil), nil),
			"agency_id":    utils.PathSearch("id", utils.PathSearch("agency", assignment, nil), nil),
		}
	}
	return assignments
}
