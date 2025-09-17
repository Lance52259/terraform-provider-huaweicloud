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
			"domain_id": {
				Type:     schema.TypeString,
				Required: true,
			},
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
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntAtLeast(1),
			},
			"per_page": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(1, 50),
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
						"scope": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"is_inherited": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
			"total_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func DataSourceIdentityRoleAssignmentsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
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
	domainId := d.Get("domain_id").(string)
	query.Add("domain_id", domainId)
	roleId := d.Get("role_id").(string)
	if roleId != "" {
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
	subjectGroupId := d.Get("subject_group_id").(string)
	if subjectGroupId != "" {
		query.Add("subject.group_id", subjectGroupId)
	}
	subjectAgencyId := d.Get("subject_agency_id").(string)
	if subjectAgencyId != "" {
		query.Add("subject.agency_id", subjectAgencyId)
	}
	scope := d.Get("scope").(string)
	if scope != "" {
		query.Add("scope", scope)
	}
	scopeProjectId := d.Get("scope_project_id").(string)
	if scopeProjectId != "" {
		query.Add("scope.project_id", scopeProjectId)
	}
	scopeDomainId := d.Get("scope_domain_id").(string)
	if scopeDomainId != "" {
		query.Add("scope.domain_id", scopeDomainId)
	}
	scopeEnterpriseProjectsId := d.Get("scope_enterprise_projects_id").(string)
	if scopeEnterpriseProjectsId != "" {
		query.Add("scope.enterprise_projects_id", scopeEnterpriseProjectsId)
	}
	isInherited := d.Get("is_inherited").(bool)
	if scope == "domain" || scopeDomainId != "" {
		query.Add("is_inherited", strconv.FormatBool(isInherited))
	}
	includeGroup := d.Get("include_group").(bool)
	if subject == "user" || subjectUserId != "" {
		query.Add("include_group", strconv.FormatBool(includeGroup))
	}
	page := d.Get("page").(int)
	if page > 0 {
		query.Add("page", strconv.Itoa(page))
	}
	perPage := d.Get("per_page").(int)
	if perPage > 0 {
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
	if err = d.Set("total_num", utils.PathSearch("total_num", respBody, nil)); err != nil {
		return diag.Errorf("error setting total_num fields: %s", err)
	}
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
			"scope":        showAssignmentsScope(utils.PathSearch("scope", assignment, nil)),
			"user_id":      utils.PathSearch("id", utils.PathSearch("user", assignment, nil), nil),
			"group_id":     utils.PathSearch("id", utils.PathSearch("group", assignment, nil), nil),
			"agency_id":    utils.PathSearch("id", utils.PathSearch("agency", assignment, nil), nil),
		}
	}
	return assignments
}

func showAssignmentsScope(scopeBody interface{}) map[string]interface{} {
	if scopeBody == nil {
		return nil
	}
	scope := map[string]interface{}{
		"project_id":            utils.PathSearch("id", utils.PathSearch("project", scopeBody, nil), nil),
		"domain_id":             utils.PathSearch("id", utils.PathSearch("domain", scopeBody, nil), nil),
		"enterprise_project_id": utils.PathSearch("id", utils.PathSearch("enterprise_project", scopeBody, nil), nil),
	}
	return scope

}
