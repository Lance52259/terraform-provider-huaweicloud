package apis

import (
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/pagination"
)

type commonResult struct {
	golangsdk.Result
}

type CreateResult struct {
	commonResult
}

type UpdateResult struct {
	commonResult
}

type GetResult struct {
	commonResult
}

type GetVersionResult struct {
	commonResult
}

type ApiResp struct {
	// API name which can contains of 3 to 64 characters, starting with a letter.
	// Only letters, digits, and underscores (_) are allowed.
	// Chinese characters must be in UTF-8 or Unicode format.
	Name string `json:"name"`
	// API type. The valid types are as following:
	//     1: public API
	//     2: private API
	Type int `json:"type"`
	// API version which can contains maximum of 16 characters.
	Version string `json:"version"`
	// Request protocol.
	//     HTTP
	//     HTTPS (default)
	//     BOTH: The API can be accessed through both HTTP and HTTPS.
	ReqProtocol string `json:"req_protocol"`
	// Request method. The valid values are GET, POST, PUT, DELETE, HEAD, PATCH, OPTIONS and ANY.
	ReqMethod string `json:"req_method"`
	// Request address, which can contain request parameters enclosed with brackets ({}). For example, /getUserInfo/{userId}. The request address can contain special characters, such as asterisks (*), percent signs (%), hyphens (-), and underscores (_). It can contain a maximum of 512 characters and must comply with URI specifications.
	// The request address must comply with URI specifications.
	ReqUri string `json:"req_uri"`
	// Security authentication mode. The valid values are as following:
	//     NONE
	//     APP
	//     IAM
	//     AUTHORIZER
	AuthType string `json:"auth_type"`
	// Security authentication parameter.
	AuthOpt AuthOpt `json:"auth_opt"`
	// Indicates whether CORS is supported.
	// TRUE: supported
	// FALSE: not supported (default).
	Cors bool `json:"cors"`
	// Route matching mode.
	//     SWA: prefix match
	//     NORMAL: exact match (default).
	MatchMode string `json:"match_mode"`
	// Backend type. The valid types are as following:
	//     HTTP: web backend
	//     FUNCTION: FunctionGraph backend
	//     MOCK: Mock backend
	BackendType string `json:"backend_type"`
	// Description of the API, which can contain a maximum of 255 characters.
	// Chinese characters must be in UTF-8 or Unicode format.
	Description string `json:"remark"`
	// ID of the API group to which the API belongs.
	GroupId string `json:"group_id"`
	// API request body, which can be an example request body, media type, or parameters. Ensure that the request body does not exceed 20,480 characters.
	// Chinese characters must be in UTF-8 or Unicode format.
	BodyDescription string `json:"body_remark"`
	// Example response for a successful request. Ensure that the response does not exceed 20,480 characters.
	// Chinese characters must be in UTF-8 or Unicode format.
	ResultNormalSample  string `json:"result_normal_sample"`
	// Example response for a failed request. Ensure that the response does not exceed 20,480 characters.
	// Chinese characters must be in UTF-8 or Unicode format.
	ResultFailureSample string `json:"result_failure_sample"`
	// ID of the frontend custom authorizer.
	AuthorizerId string `json:"authorizer_id"`
	// Tags.
	Tags []string `json:"tags"`
	// Group response ID.
	ResponseId string `json:"response_id"`
	// API ID.
	Id  string `json:"id"`
	// API status. 1: valid
	Status int `json:"status"`
	// Indicates whether to enable orchestration.
	ArrangeNecessary int `json:"arrange_necessary"`
	// Time when the API is registered.
	RegisterTime string `json:"register_time"`
	// Time when the API was last modified.
	UpdateTime string `json:"update_time"`
	// Name of the API group to which the API belongs.
	GroupName  string `json:"group_name"`
	// Version of the API group to which the API belongs.
	// The default value is V1. Other versions are not supported.
	GroupVersion  string `json:"group_version"`
	// ID of the environment in which the API has been published.
	// If there are multiple publication records, separate the environment IDs with vertical bars (|).
	RunEnvId string `json:"run_env_id"`
	// Name of the environment in which the API has been published.
	// If there are multiple publication records, separate the environment names with vertical bars (|).
	RunEnvName string `json:"run_env_name"`
	// Publication record ID.
	// You can separate multiple publication record IDs with vertical bars (|).
	PublishId string `json:"publish_id"`
	// FunctionGraph backend details.
	FuncInfo FuncGraph `json:"func_info"`
	// Mock backend details.
	MockInfo Mock `json:"mock_info"`
	// Web backend details.
	WebInfo Web `json:"backend_api"`
	// Request parameters.
	ReqParams []ReqParamBase  `json:"req_params"`
	// Backend parameters.
	BackendParams BackendParamBase `json:"backend_params"`
	// Mock policy backends.
	PolicyMocks []PolicyMock `json:"policy_mocks"`
	// FunctionGraph policy backends.
	PolicyFunctions []PolicyFuncGraph `json:"policy_functions"`
	// Web policy backends.
	PolicyWebs []PolicyWeb `json:"policy_https"`
	// The following four parameters are only provided by the response of the API Versions.
	//     SlDomain
	//     SlDomains
	//     VersionId
	//     PublishTime.
	// Subdomain name that API Gateway automatically allocates to the API group.
	SlDomain string `json:"sl_domain"`
	// Subdomain names that API Gateway automatically allocates to the API group.
	SlDomains []string `json:"sl_domains"`
	// API version ID.
	VersionId string `json:"version_id"`
	//Time when the API version is published.
	PublishTime string `json:"publish_time"`
}

// Extract is a method to extract an response struct.
func (r commonResult) Extract() (*ApiResp, error) {
	var s ApiResp
	err := r.ExtractInto(&s)
	return &s, err
}

// ApiPage represents the api pages of the List operation.
type ApiPage struct {
	pagination.SinglePageBase
}

// ExtractApis is a method to extract an response struct list.
func ExtractApis(r pagination.Page) ([]ApiResp, error) {
	var s []ApiResp
	err := r.(ApiPage).Result.ExtractIntoSlicePtr(&s, "responses")
	return s, err
}

// DeleteResult represents a result of the Delete method.
type DeleteResult struct {
	golangsdk.ErrResult
}

type ApiVersion struct {
	// API version ID.
	VersionId string `json:"version_id"`
	// API version.
	VersionNo string `json:"version_no"`
	// API ID.
	ApiId string `json:"api_id"`
	// ID of the environment in which the API has been published.
	EnvId string `json:"env_id"`
	// Name of the environment in which the API has been published.
	EnvName string `json:"env_name"`
	// Description about the publication.
	Remark string `json:"remark"`
	// Publication time.
	PublishTime string `json:"publish_time"`
	// Version status.
	//     1: effective
	//     2: not effective
	Status int `json:"status"`
}

type ApiVersionPage struct {
	pagination.SinglePageBase
}

// ExtractApiVersions is a method to extract an ApiVersion struct list.
func ExtractApiVersions(r pagination.Page) ([]ApiVersion, error) {
	var s []ApiVersion
	err := r.(ApiPage).Result.ExtractIntoSlicePtr(&s, "api_versions")
	return s, err
}

type SwitchVersion struct {
	// Publication record ID.
	PublishId string `json:"publish_id"`
	// API ID.
	ApiId string `json:"api_id"`
	// API name.
	ApiName string `json:"api_name"`
	// ID of the environment in which the API has been published.
	EnvId string `json:"env_id"`
	// Description about the publication.
	Description string `json:"remark"`
	// Publication time.
	PublishTime string `json:"publish_time"`
	// API version currently in use.
	VersionId string `json:"version_id"`
}

type SwitchResult struct {
	golangsdk.Result
}

// Extract is a method to extract an SwitchVersion struct.
func (r SwitchResult) Extract() (*SwitchVersion, error) {
	var s SwitchVersion
	err := r.ExtractInto(&s)
	return &s, err
}

type RuntimeResult struct {
	golangsdk.Result
}

type RuntimeDefinition struct {
	// API name. An API name can contain 3 to 64 characters, starting with a letter.
	// Only letters, digits, and underscores (_) are allowed.
	// Chinese characters must be in UTF-8 or Unicode format.
	Name string `json:"name"`
	// API type.
	// 1: public API
	// 2: private API
	Type int `json:"type"`
	// API version.
	// Maximum: 16
	Version string `json:"version"`
	// Request protocol.
	//     HTTP
	//     HTTPS (Default)
	//     BOTH: The API can be accessed through both HTTP and HTTPS.
	ReqProtocol string `json:"req_protocol"`
	// Request method. The support methods are as following:
	//     GET
	//     POST
	//     PUT
	//     DELETE
	//     HEAD
	//     PATCH
	//     OPTIONS
	//     ANY
	ReqMethod string `json:"req_method"`
	// Request address, which can contain request parameters enclosed with brackets ({}). For example, /getUserInfo/{userId}. The request address can contain special characters, such as asterisks (*), percent signs (%), hyphens (-), and underscores (_). It can contain a maximum of 512 characters and must comply with URI specifications.
	// The request address must comply with URI specifications.
	ReqUri string `json:"req_uri"`
	// Security authentication mode. The support mode are as following:
	//     NONE
	//     APP
	//     IAM
	//     AUTHORIZER
	AuthType string `json:"auth_type"`
	// Security authentication parameter.
	AuthOpt AuthOpt `json:"auth_opt"`
	// Indicates whether CORS is supported.
	//     TRUE: supported
	//     FALSE: not supported (default)
	Cors bool `json:"cors"`
	// Route matching mode.
	// SWA: prefix match
	// NORMAL: exact match.
	// The default value is NORMAL.
	// Enumeration values:
	// SWA
	// NORMAL
	MatchMode string `json:"match_mode"`
	// Backend type.
	// HTTP: web backend
	// FUNCTION: FunctionGraph backend
	// MOCK: Mock backend
	// Enumeration values:
	// HTTP
	// FUNCTION
	// MOCK
	BackendType string `json:"backend_type"`
	// Description of the API, which can contain a maximum of 255 characters.
	// Chinese characters must be in UTF-8 or Unicode format.
	Remark string `json:"remark"`
	// ID of the API group to which the API belongs.
	GroupId string `json:"group_id"`
	// API request body, which can be an example request body, media type, or parameters. Ensure that the request body does not exceed 20,480 characters.
	// Chinese characters must be in UTF-8 or Unicode format.
	BodyRemark string `json:"body_remark"`
	// Example response for a successful request. Ensure that the response does not exceed 20,480 characters.
	// Chinese characters must be in UTF-8 or Unicode format.
	ResultNormalSample string `json:"result_normal_sample"`
	// Example response for a failed request. Ensure that the response does not exceed 20,480 characters.
	// Chinese characters must be in UTF-8 or Unicode format.
	ResultFailureSample string `json:"result_failure_sample"`
	// ID of the frontend custom authorizer.
	AuthorizerId string `json:"authorizer_id"`
	// Tags.
	// The value can contain only letters, digits, and underscores (_), and must start with a letter. You can enter multiple tags and separate them with commas (,).
	// Minimum: 1
	// Maximum: 128
	Tags []string `json:"tags"`
	// Group response ID.
	ResponseId string `json:"response_id"`
	// Integration application ID.
	// This parameter is currently not supported.
	RomaAppId string `json:"roma_app_id"`
	// Custom domain name bound to the API.
	// This parameter is currently not supported.
	DomainName string `json:"domain_name"`
	// API ID.
	Id string `json:"id"`
	// Name of the API group to which the API belongs.
	GroupName string `json:"group_name"`
	// Name of the environment in which the API has been published.
	RunEnvName string `json:"run_env_name"`
	// ID of the environment in which the API has been published.
	RunEnvId string `json:"run_env_id"`
	// Publication record ID.
	PublishId string `json:"publish_id"`
	// Subdomain name of the API group.
	SlDomain string `json:"sl_domain"`
	// Subdomain names that API Gateway automatically allocates to the API group.
	SlDomains []string `json:"sl_domains"`
	// Request parameters.
	ReqParams []ReqParamBase `json:"req_params"`
}

// Extract is a method to extract an RuntimeDefinition struct.
func (r RuntimeResult) Extract() (*RuntimeDefinition, error) {
	var s RuntimeDefinition
	err := r.ExtractInto(&s)
	return &s, err
}

type ConnectivityCheckResult struct {
	golangsdk.Result
}

type ConnectivityCheck struct {
	// Connectivity check result.
	// Enumeration values:
	// SUCCESS
	// FAILED
	CheckResult string `json:"check_result"`
	// Backend services that fail the connectivity check.
	Failures []string `json:"failures"`
}

// Extract is a method to extract an ConnectivityCheck struct.
func (r ConnectivityCheckResult) Extract() (*ConnectivityCheck, error) {
	var s ConnectivityCheck
	err := r.ExtractInto(&s)
	return &s, err
}
