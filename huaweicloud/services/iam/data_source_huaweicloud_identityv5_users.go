package iam

import (
	"context"
	"fmt"
	"strings"

	"github.com/chnsz/golangsdk"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API GET /v5/users
// @API GET /v5/users/{user_id}
// @API GET /v5/users/{user_id}/last-login
func DataSourceIdentityV5Users() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIdentityV5UserRead,

		Schema: map[string]*schema.Schema{
			"group_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID of the IAM user.",
			},
			"user_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"users": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"is_root_user": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"created_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"urn": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"user_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tags": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"tag_key": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"tag_value": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"last_login_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceIdentityV5UserRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("iam", region)
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	var allUsers []interface{}
	var marker string
	var path string
	if userId, ok := d.GetOk("user_id"); ok {
		getUserHttpUrl := "v5/users/{user_id}"
		getUserPath := client.Endpoint + getUserHttpUrl
		getUserPath = strings.ReplaceAll(getUserPath, "{user_id}", userId.(string))
		getUserOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
		}
		getUserResp, err := client.Request("GET", getUserPath, &getUserOpt)
		if err != nil {
			return common.CheckDeletedDiag(d, err, "error get IAM user")
		}
		getUserRespBody, err := utils.FlattenResponse(getUserResp)
		if err != nil {
			return diag.FromErr(err)
		}
		user := utils.PathSearch("user", getUserRespBody, nil)
		if user == nil {
			return common.CheckDeletedDiag(d, err, "error getting IAM user: user is not found in API response")
		}

		tags := listUserTagInfo(user)
		getLastLoginHttpUrl := "v5/users/{user_id}/last-login"
		getLastLoginPath := client.Endpoint + getLastLoginHttpUrl
		getLastLoginPath = strings.ReplaceAll(getLastLoginPath, "{user_id}", userId.(string))
		getLastLoginOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
		}
		getLastLoginResp, err := client.Request("GET", getLastLoginPath, &getLastLoginOpt)
		if err != nil {
			return diag.Errorf("error get IAM User last login time: %s", err)
		}
		userLastLogin, err := utils.FlattenResponse(getLastLoginResp)
		if err != nil {
			return diag.FromErr(err)
		}
		users := flattenShowUserV5Response(user, tags, userLastLogin)
		allUsers = append(allUsers, users)
	} else {
		for {
			path = fmt.Sprintf("%sv5/users", client.Endpoint) + buildListUsersV5Params(d, marker)
			reqOpt := &golangsdk.RequestOpts{
				KeepResponseBody: true,
			}
			r, err := client.Request("GET", path, reqOpt)
			if err != nil {
				return common.CheckDeletedDiag(d, err, "error retrieving users")
			}
			resp, err := utils.FlattenResponse(r)
			if err != nil {
				return diag.FromErr(err)
			}
			users := flattenListUsersV5Response(resp)
			allUsers = append(allUsers, users...)

			marker = utils.PathSearch("page_info.next_marker", resp, "").(string)
			if marker == "" {
				break
			}
		}
	}

	id, _ := uuid.GenerateUUID()
	d.SetId(id)
	mErr := multierror.Append(nil,
		d.Set("users", allUsers),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting users fields: %s", err)
	}
	return nil
}

func buildListUsersV5Params(d *schema.ResourceData, marker string) string {
	res := "?limit=100"
	if v, ok := d.GetOk("group_id"); ok {
		res = fmt.Sprintf("%s&group_id=%v", res, v)
	}
	if marker != "" {
		res = fmt.Sprintf("%s&marker=%v", res, marker)
	}
	return res
}

func flattenListUsersV5Response(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	users := utils.PathSearch("users", resp, make([]interface{}, 0)).([]interface{})
	result := make([]interface{}, len(users))
	for i, user := range users {
		result[i] = map[string]interface{}{
			"enabled":      utils.PathSearch("enabled", user, nil),
			"user_name":    utils.PathSearch("user_name", user, nil),
			"description":  utils.PathSearch("description", user, nil),
			"is_root_user": utils.PathSearch("is_root_user", user, nil),
			"created_at":   utils.PathSearch("created_at", user, nil),
			"urn":          utils.PathSearch("urn", user, nil),
			"user_id":      utils.PathSearch("user_id", user, nil),
		}
	}
	return result
}

func flattenShowUserV5Response(user interface{}, tags []interface{}, userLastLogin interface{}) interface{} {
	result := map[string]interface{}{
		"enabled":       utils.PathSearch("enabled", user, nil),
		"user_name":     utils.PathSearch("user_name", user, nil),
		"description":   utils.PathSearch("description", user, nil),
		"is_root_user":  utils.PathSearch("is_root_user", user, nil),
		"created_at":    utils.PathSearch("created_at", user, nil),
		"urn":           utils.PathSearch("urn", user, nil),
		"user_id":       utils.PathSearch("user_id", user, nil),
		"tags":          tags,
		"last_login_at": utils.PathSearch("last_login_at", userLastLogin, nil),
	}
	return result
}
