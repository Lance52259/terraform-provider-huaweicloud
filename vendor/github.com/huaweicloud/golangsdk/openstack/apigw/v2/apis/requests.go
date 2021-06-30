package apis

import (
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/pagination"
)

type ApiOpts struct {
	// ID of the API group to which the API belongs.
	GroupId string `json:"group_id" required:"true"`
	// API name, which can contain 3 to 64 characters, starting with a letter.
	// Only letters, digits, and underscores (_) are allowed.
	// Chinese characters must be in UTF-8 or Unicode format.
	Name string `json:"name" required:"true"`
	// API type. The valid types are as following:
	//     1: public API
	//     2: private API
	Type int `json:"type" required:"true"`
	// Request protocol. The valid protocols are as following:
	//     HTTP.
	//     HTTPS (default).
	//     BOTH: The API can be accessed through both HTTP and HTTPS.
	ReqProtocol string `json:"req_protocol" required:"true"`
	// Request method. The valid values are GET,  POST,  PUT,  DELETE, HEAD, PATCH, OPTIONS and ANY.
	ReqMethod string `json:"req_method" required:"true"`
	// Request address, which can contain a maximum of 512 characters request parameters enclosed with brackets ({}).
	// For example, /getUserInfo/{userId}.
	// The request address can contain special characters, such as asterisks (), percent signs (%), hyphens (-), and
	// underscores (_) and must comply with URI specifications.
	// The address can contain environment variables, each starting with a letter and consisting of 3 to 32 characters.
	// Only letters, digits, hyphens (-), and underscores (_) are allowed in environment variables.
	ReqUri string `json:"req_uri" required:"true"`
	// Security authentication mode. The valid modes are as following:
	//     NONE
	//     APP
	//     IAM
	//     AUTHORIZER
	AuthType string `json:"auth_type" required:"true"`
	// Backend type. The valid types are as following:
	//     HTTP: web backend.
	//     FUNCTION: FunctionGraph backend.
	//     MOCK: Mock backend.
	BackendType string `json:"backend_type" required:"true"`
	// API version. The maximum length of version string is 16.
	Version string `json:"version,omitempty"`
	// Security authentication parameter.
	AuthOpt AuthOpt `json:"auth_opt,omitempty"`
	// Indicates whether CORS is supported. The valid values are as following:
	//     TRUE: supported.
	//     FALSE: not supported (default).
	Cors bool `json:"cors,omitempty"`
	// Route matching mode.  The valid modes are as following:
	//     SWA: prefix match
	//     NORMAL: exact match (default).
	MatchMode string `json:"match_mode,omitempty"`
	// Description of the API, which can contain a maximum of 255 characters.
	// Chinese characters must be in UTF-8 or Unicode format.
	Description string `json:"remark,omitempty"`
	// API request body, which can be an example request body, media type, or parameters.
	// Ensure that the request body does not exceed 20,480 characters.
	// Chinese characters must be in UTF-8 or Unicode format.
	BodyDescription string `json:"body_remark,omitempty"`
	// Example response for a successful request. Ensure that the response does not exceed 20,480 characters.
	// Chinese characters must be in UTF-8 or Unicode format.
	ResultNormalSample string `json:"result_normal_sample,omitempty"`
	// Example response for a failed request. Ensure that the response does not exceed 20,480 characters.
	// Chinese characters must be in UTF-8 or Unicode format.
	ResultFailureSample string `json:"result_failure_sample,omitempty"`
	// ID of the frontend custom authorizer.
	AuthorizerId string `json:"authorizer_id,omitempty"`
	// List of tags. The length of the tags list is range from 1 to 128.
	// The value can contain only letters, digits, and underscores (_), and must start with a letter.
	Tags []string `json:"tags,omitempty"`
	// Group response ID.
	ResponseId string `json:"response_id,omitempty"`
	// Request parameters.
	ReqParams []ReqParamBase `json:"req_params,omitempty"`
	// Backend parameters.
	BackendParams []BackendParamBase `json:"backend_params,omitempty"`
	// Mock backend details.
	MockInfo Mock `json:"mock_info,omitempty"`
	// FunctionGraph backend details.
	FuncInfo FuncGraph `json:"func_info,omitempty"`
	// Web backend details.
	WebInfo Web `json:"backend_api,omitempty"`
	// Mock policy backends.
	PolicyMocks []PolicyMock `json:"policy_mocks,omitempty"`
	// FunctionGraph policy backends.
	PolicyFunctions []PolicyFuncGraph `json:"policy_functions,omitempty"`
	// Web policy backends.
	PolicyWebs []PolicyWeb `json:"policy_https,omitempty"`
}

type AuthOpt struct {
	// Indicates whether AppCode authentication is enabled. The valid types are as following:
	//     DISABLE: AppCode authentication is disabled (default).
	//     HEADER: AppCode authentication is enabled and the AppCode is located in the header.
	// This parameter is valid only if auth_type is set to App.
	AppCodeAuthType string `json:"app_code_auth_type,omitempty"`
}

type Mock struct {
	// Description about the backend, which can contain a maximum of 255 characters.
	// Chinese characters must be in UTF-8 or Unicode format.
	Description string `json:"remark,omitempty"`
	// Response.
	ResultContent string `json:"result_content,omitempty"`
	// Function version Ensure that the version does not exceed 64 characters.
	Version string `json:"version,omitempty"`
	// Backend custom authorizer ID.
	AuthorizerId string `json:"authorizer_id,omitempty"`
}

type FuncGraph struct {
	// Function URN.
	FunctionUrn string `json:"function_urn" required:"true"`
	// Invocation mode. The valid modes are as following:
	//     async: asynchronous
	//     sync: synchronous
	InvocationType string `json:"invocation_type" required:"true"`
	// Timeout, in ms, which allowed for API Gateway to request the backend service.
	// The valid value is range from 1 to 600,000.
	Timeout int `json:"timeout" required:"true"`
	// Backend custom authorizer ID.
	AuthorizerId string `json:"authorizer_id,omitempty"`
	// Description about the backend, which can contain a maximum of 255 characters.
	// Chinese characters must be in UTF-8 or Unicode format.
	Description string `json:"remark,omitempty"`
	// Function version.
	// Maximum: 64
	Version string `json:"version,omitempty"`
}

type Web struct {
	// Request method. The valid methods are GET, POST, PUT, DELETE, HEAD, PATCH, OPTIONS and ANY.
	ReqMethod string `json:"req_method" required:"true"`
	// Request protocol. The valid protocols are HTTP and HTTPS
	ReqProtocol string `json:"req_protocol" required:"true"`
	// Request address, which can contain a maximum of 512 characters request parameters enclosed with brackets ({}).
	// For example, /getUserInfo/{userId}.
	// The request address can contain special characters, such as asterisks (), percent signs (%), hyphens (-), and
	// underscores (_) and must comply with URI specifications.
	// The address can contain environment variables, each starting with a letter and consisting of 3 to 32 characters.
	// Only letters, digits, hyphens (-), and underscores (_) are allowed in environment variables.
	ReqUri string `json:"req_uri" required:"true"`
	// Timeout, in ms, which allowed for API Gateway to request the backend service.
	// The valid value is range from 1 to 600,000.
	Timeout int `json:"timeout" required:"true"`
	// Backend custom authorizer ID.
	AuthorizerId string `json:"authorizer_id,omitempty"`
	// Backend service address which consists of a domain name or IP address and a port number, with not more than 255
	// characters. It must be in the format "Host name:Port number", for example, apig.example.com:7443.
	// If the port number is not specified, the default HTTPS port 443 or the default HTTP port 80 is used.
	// The backend service address can contain environment variables, each starting with a letter and consisting of
	// 3 to 32 characters. Only letters, digits, hyphens (-), and underscores (_) are allowed.
	UrlDomain string `json:"url_domain,omitempty"`
	// Description, which can contain a maximum of 255 characters.
	// Chinese characters must be in UTF-8 or Unicode format.
	Description string `json:"remark,omitempty"`
	// Web backend version, which can contain a maximum of 16 characters.
	Version string `json:"version,omitempty"`
	// Indicates whether to enable two-way authentication.
	EnableClientSsl bool `json:"enable_client_ssl,omitempty"`
	// VPC channel details. This parameter is required if vpc_channel_status is set to 1.
	VpcChannelInfo VpcChannel `json:"vpc_channel_info,omitempty"`
	// Indicates whether to use a VPC channel. The valid values are as following:
	//     1: (A VPC channel is used).
	//     2: (No VPC channel is used).
	VpcChannelEnable int `json:"vpc_channel_status,omitempty"`
}

type VpcChannel struct {
	// Proxy host.
	VpcChannelProxyHost string `json:"vpc_channel_proxy_host,omitempty"`
	// VPC channel ID.
	VpcChannelId string `json:"vpc_channel_id" required:"true"`
}

type ReqParamBase struct {
	// The parameter name, which contain of 1 to 32 characters and start with a letter.
	// Only letters, digits, hyphens (-), underscores (_), and periods (.) are allowed.
	Name string `json:"name" required:"true"`
	// Parameter type. The valid types are as following:
	//     STRING
	//     NUMBER
	Type string `json:"type" required:"true"`
	// Parameter location. The valid modes are as following:
	//     PATH
	//     QUERY
	//     HEADER
	Location string `json:"location" required:"true"`
	// Default value.
	DefaultValue string `json:"default_value,omitempty"`
	// Example value.
	SampleValue string `json:"sample_value,omitempty"`
	// Indicates whether the parameter is required. The valid values are 1 (yes) and 2 (no).
	// The value of this parameter is 1 if Location is set to PATH, and 2 if Location is set to another value.
	Required int `json:"required,omitempty"`
	// Indicates whether validity check is enabled. The valid modes are as following:
	// 1: enabled.
	// 2: disabled (default).
	ValidEnable int `json:"valid_enable,omitempty"`
	// Description about the backend, which can contain a maximum of 255 characters.
	// Chinese characters must be in UTF-8 or Unicode format.
	Description string `json:"remark,omitempty"`
	// Enumerated value.
	Enumerations string `json:"enumerations,omitempty"`
	// Minimum value.
	// This parameter is valid when type is set to NUMBER.
	MinNum int `json:"min_num,omitempty"`
	// Maximum value.
	// This parameter is valid when type is set to NUMBER.
	MaxNum int `json:"max_num,omitempty"`
	// Minimum length.
	// This parameter is valid when type is set to STRING.
	MinSize int `json:"min_size,omitempty"`
	// Maximum length.
	// This parameter is valid when type is set to STRING.
	MaxSize int `json:"max_size,omitempty"`
	// Indicates whether to transparently transfer the parameter. The valid values are 1 (yes) and 2 (no).
	PassThrough string `json:"pass_through,omitempty"`
}

type PolicyMock struct {
	// Policy conditions.
	Conditions []ApiConditionBase `json:"conditions" required:"true"`
	// Effective mode of the backend policy. The valid modes are as following:
	//     ALL: All conditions are met.
	//     ANY: Any condition is met.
	EffectMode string `json:"effect_mode" required:"true"`
	// Backend name, which consists of 3 to 64 characters and must start with a letter and can contain letters, digits,
	// and underscores (_).
	Name string `json:"name" required:"true"`
	// Authorizer ID.
	AuthorizerId string `json:"authorizer_id,omitempty"`
	// Backend parameters.
	BackendParams []BackendParamBase `json:"backend_params,omitempty"`
	// Response.
	ResultContent string `json:"result_content,omitempty"`
}

type PolicyFuncGraph struct {
	// Policy conditions.
	Conditions []ApiConditionBase `json:"conditions" required:"true"`
	// Effective mode of the backend policy.
	//     ALL: All conditions are met.
	//     ANY: Any condition is met.
	EffectMode string `json:"effect_mode" required:"true"`
	// Function URN.
	FunctionUrn string `json:"function_urn" required:"true"`
	// Invocation mode. The valid modes are as following:
	//     async: asynchronous
	//     sync: synchronous
	InvocationType string `json:"invocation_type" required:"true"`
	// The backend name consists of 3 to 64 characters, which must start with a letter and can contain letters, digits,
	// and underscores (_).
	Name string `json:"name" required:"true"`
	// Authorizer ID.
	AuthorizerId string `json:"authorizer_id,omitempty"`
	// Backend parameters.
	BackendParams []BackendParamBase `json:"backend_params,omitempty"`
	// Timeout, in ms, which allowed for API Gateway to request the backend service.
	// The valid value is range from 1 to 600,000.
	Timeout int `json:"timeout,omitempty"`
	// Function version Ensure that the version does not exceed 64 characters.
	Version string `json:"version,omitempty"`
}

type PolicyWeb struct {
	// Request protocol. The value can be HTTP or HTTPS.
	ReqProtocol string `json:"req_protocol" required:"true"`
	// Request method. The valid methods are GET, POST, PUT, DELETE, HEAD, PATCH, OPTIONS and ANY.
	ReqMethod string `json:"req_method" required:"true"`
	// Request address, which can contain request parameters enclosed with brackets ({}). For example, /getUserInfo/{userId}. The request address can contain special characters, such as asterisks (), percent signs (%), hyphens (-), and underscores (_). It can contain a maximum of 512 characters and must comply with URI specifications.
	// The request address can contain environment variables, each starting with a letter and consisting of 3 to 32 characters. Only letters, digits, hyphens (-), and underscores (_) are allowed in environment variables.
	// The request address must comply with URI specifications.
	ReqUri string `json:"req_uri" required:"true"`
	// Endpoint of the policy backend.
	// An endpoint consists of a domain name or IP address and a port number, with not more than 255 characters. It must be in the format "Domain name:Port number", for example, apig.example.com:7443. If the port number is not specified, the default HTTPS port 443 or the default HTTP port 80 is used.
	// The endpoint can contain environment variables, each starting with a letter and consisting of 3 to 32 characters. Only letters, digits, hyphens (-), and underscores (_) are allowed.
	UrlDomain string `json:"url_domain,omitempty"`
	// Timeout, in ms, which allowed for API Gateway to request the backend service.
	// The valid value is range from 1 to 600,000.
	Timeout int `json:"timeout,omitempty"`
	// Effective mode of the backend policy. The valid modes are as following:
	//     ALL: All conditions are met.
	//     ANY: Any condition is met.
	EffectMode string `json:"effect_mode" required:"true"`
	// Backend name, which contains of 3 to 64, must start with a letter and can contain letters, digits, and underscores (_).
	// Minimum: 3
	// Maximum: 64
	Name string `json:"name" required:"true"`
	// Backend parameters.
	BackendParams []BackendParamBase `json:"backend_params,omitempty"`
	// Policy conditions.
	Conditions []ApiConditionBase `json:"conditions" required:"true"`
	// Authorizer ID.
	AuthorizerId string `json:"authorizer_id,omitempty"`
	// VPC channel details. This parameter is required if vpc_channel_status is set to 1.
	VpcChannelInfo VpcChannel `json:"vpc_channel_info,omitempty"`
	// Indicates whether to use a VPC channel. The valid value are as following:
	//     1: A VPC channel is used.
	//     2: No VPC channel is used.
	VpcChannelEnable int `json:"vpc_channel_status,omitempty"`
}

type BackendParamBase struct {
	// Parameter type. The valid types are as following:
	//     REQUEST: Backend parameter.
	//     CONSTANT: Constant parameter.
	//     SYSTEM: System parameter.
	Origin string `json:"origin" required:"true"`
	// Parameter name.
	// The parameter name must start with a letter and can only contain letters, digits, hyphens (-), underscores (_), and periods (.).
	// Minimum: 1
	// Maximum: 32
	Name string `json:"name" required:"true"`
	// Description, which can contain a maximum of 255 characters.
	// Chinese characters must be in UTF-8 or Unicode format.
	Description string `json:"remark,omitempty"`
	// Parameter location. The valid values are PATH, QUERY and HEADER.
	Location string `json:"location" required:"true"`
	// Parameter value, which can contain a maximum of 255 characters. If the origin type is REQUEST, the value of this parameter is the parameter name in req_params.
	// If the origin type is CONSTANT, the value is a constant.
	// If the origin type is SYSTEM, the value is a system parameter name. System parameters include gateway parameters, frontend authentication parameters, and backend authentication parameters. You can set the frontend or backend authentication parameters after enabling custom frontend or backend authentication.
	// The gateway parameters are as follows:
	//     $context.sourceIp: source IP address of the API caller.
	//     $context.stage: deployment environment in which the API is called.
	//     $context.apiId: API ID.
	//     $context.appId: ID of the app used by the API caller.
	//     $context.requestId: request ID generated when the API is called.
	//     $context.serverAddr: address of the gateway server.
	//     $context.serverName: name of the gateway server.
	//     $context.handleTime: time when the API request is processed.
	//     $context.providerAppId: ID of the app used by the API owner. This parameter is currently not supported.
	// Frontend authentication parameter: prefixed with "$context.authorizer.frontend.". For example, to return "aaa" upon successful custom authentication, set this parameter to "$context.authorizer.frontend.aaa".
	// Backend authentication parameter: prefixed with "$context.authorizer.backend.". For example, to return "aaa" upon successful custom authentication, set this parameter to "$context.authorizer.backend.aaa".
	Value string `json:"value" required:"true"`
}

type ApiConditionBase struct {
	// Input parameter name. This parameter is required if the policy type is param.
	ReqParamName string `json:"req_param_name,omitempty"`
	// Policy condition.
	// exact: exact match
	// enum: enumeration
	// pattern: regular expression
	// This parameter is required if the policy type is param.
	// Enumeration values:
	// exact
	// enum
	// pattern
	ConditionType string `json:"condition_type,omitempty"`
	// Policy type.
	// param: input parameter
	// source: source IP address
	// Enumeration values:
	// param
	// source
	ConditionOrigin string `json:"condition_origin" required:"true"`
	// Condition value.
	ConditionValue string `json:"condition_value" required:"true"`
}

type ApiOptsBuilder interface {
	ToApiOptsMap() (map[string]interface{}, error)
}

func (opts ApiOpts) ToApiOptsMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

// Create is a method by which to create function that create a new custom api.
func Create(client *golangsdk.ServiceClient, instanceId string, opts ApiOptsBuilder) (r CreateResult) {
	reqBody, err := opts.ToApiOptsMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(rootURL(client, instanceId), reqBody, &r.Body, nil)
	return
}

// Update is a method to update an existing custom api.
func Update(client *golangsdk.ServiceClient, instanceId, appId string, opts ApiOptsBuilder) (r UpdateResult) {
	reqBody, err := opts.ToApiOptsMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Put(resourceURL(client, instanceId, appId), reqBody, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})
	return
}

// Get is a method to obtain the specified api according to the instanceId and apiId.
func Get(client *golangsdk.ServiceClient, instanceId, apiId string) (r GetResult) {
	_, r.Err = client.Get(resourceURL(client, instanceId, apiId), &r.Body, nil)
	return
}

type ListOpts struct {
	// API ID.
	Id string `q:"id"`
	// API name.
	Name string `q:"name"`
	// API group ID.
	GroupId string `q:"group_id"`
	// Request protocol.
	ReqProtocol string `q:"req_protocol"`
	// Request method.
	ReqMethod string `q:"req_method"`
	// Request path.
	ReqUri string `q:"req_uri"`
	// Security authentication mode.
	AuthType string `q:"auth_type"`
	// ID of the environment in which the API has been published.
	EnvId string `q:"env_id"`
	// API type.
	Type int `q:"type"`
	// Offset from which the query starts.
	// If the offset is less than 0, the value is automatically converted to 0. Default to 0.
	Offset int `q:"offset"`
	// Number of items displayed on each page. The range of number is form 1 to 500, default to 20.
	Limit int `q:"limit"`
	// Parameter name (name or req_uri) for exact matching.
	PreciseSearch string `q:"precise_search"`
}

type ListOptsBuilder interface {
	ToListOptsQuery() (string, error)
}

func (opts ListOpts) ToListOptsQuery() (string, error) {
	q, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return "", err
	}
	return q.String(), err
}

// List is a method to obtain an array of one or more APIG api according to the query parameters.
func List(client *golangsdk.ServiceClient, instanceId string, opts ListOptsBuilder) pagination.Pager {
	url := rootURL(client, instanceId)
	if opts != nil {
		query, err := opts.ToListOptsQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}

	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return ApiPage{pagination.SinglePageBase(r)}
	})
}

// Delete is a method to delete an existing custom api.
func Delete(client *golangsdk.ServiceClient, instanceId, apiId string) (r DeleteResult) {
	_, r.Err = client.Delete(resourceURL(client, instanceId, apiId), nil)
	return
}

type PublishActionOpts struct {
	// Operation to perform. The valid options are as following:
	//     online: publishing the APIs
	//     offline: taking the APIs offline
	// This parameter is automatically generated by method.
	Action string `json:"action" required:"true"`
	// ID of the environment in which the API will be published.
	EnvId string `json:"env_id" required:"true"`
	// ID of the API to be published or taken offline.
	ApiId string `json:"api_id" required:"true"`
	// Description about the operation, which can contain a maximum of 255 characters.
	// Chinese characters must be in UTF-8 or Unicode format.
	Description string `json:"remark,omitempty"`
}

type PublishActionOptsBuilder interface {
	ToOnlineOptsMap() (map[string]interface{}, error)
	ToOfflineOptsMap() (map[string]interface{}, error)
}

func (opts PublishActionOpts) ToOnlineOptsMap() (map[string]interface{}, error) {
	opts.Action = "online"
	return golangsdk.BuildRequestBody(opts, "")
}

func (opts PublishActionOpts) ToOfflineOptsMap() (map[string]interface{}, error) {
	opts.Action = "offline"
	return golangsdk.BuildRequestBody(opts, "")
}

// Release is a method to release an existing api to an environment.
func Release(client *golangsdk.ServiceClient, instanceId string, opts PublishActionOptsBuilder) (r CreateResult) {
	reqBody, err := opts.ToOnlineOptsMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(actionURL(client, instanceId), reqBody, &r.Body, nil)
	return
}

// Offline is a method to offline an existing api form an environment.
func Offline(client *golangsdk.ServiceClient, instanceId string, opts PublishActionOptsBuilder) (r CreateResult) {
	reqBody, err := opts.ToOfflineOptsMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(actionURL(client, instanceId), reqBody, &r.Body, nil)
	return
}

type BatchPublicActionOpts struct {
	// Operation to perform. The valid options are as following:
	//     online: publishing the APIs
	//     offline: taking the APIs offline
	Action string `q:"action"`
	// ID of the environment in which the API will be published.
	EnvId string `json:"env_id" required:"true"`
	// ID of the API to be published or taken offline.
	Apis []string `json:"api_id" required:"true"`
	// Description about the operation, which can contain a maximum of 255 characters.
	// Chinese characters must be in UTF-8 or Unicode format.
	Description string `json:"remark"`
}

type BatchPublishOptsBuilder interface {
	ToBatchPublicActionOptsMap() (map[string]interface{}, error)
	ToBatchPublishQuery() (string, error)
	ToBatchOfflineQuery() (string, error)
}

func (opts BatchPublicActionOpts) ToBatchPublishQuery() (string, error) {
	opts.Action = "online"
	q, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return "", err
	}
	return q.String(), err
}

func (opts BatchPublicActionOpts) ToBatchOfflineQuery() (string, error) {
	opts.Action = "offline"
	q, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return "", err
	}
	return q.String(), err
}

func (opts BatchPublicActionOpts) ToBatchPublicActionOptsMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

// BatchRelease is a method to release one or more existing apis to an environment.
func BatchRelease(client *golangsdk.ServiceClient, instanceId string,
	opts BatchPublishOptsBuilder) (r CreateResult) {
	url := releaseURL(client, instanceId)
	if opts != nil {
		query, err := opts.ToBatchPublishQuery()
		if err != nil {
			r.Err = err
			return
		}
		url += query
	}
	reqBody, err := opts.ToBatchPublicActionOptsMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(actionURL(client, instanceId), reqBody, &r.Body, nil)
	return
}

// BatchOffline is a method to offline one or more existing apis from an environment.
func BatchOffline(client *golangsdk.ServiceClient, instanceId string,
	opts BatchPublishOptsBuilder) (r CreateResult) {
	url := releaseURL(client, instanceId)
	if opts != nil {
		query, err := opts.ToBatchPublishQuery()
		if err != nil {
			r.Err = err
			return
		}
		url += query
	}
	reqBody, err := opts.ToBatchPublicActionOptsMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(actionURL(client, instanceId), reqBody, &r.Body, nil)
	return
}

type ListApiHistoryOpts struct {
	// Environment ID.
	EnvId string `json:"env_id"`
	// Environment name.
	EnvName string `json:"env_name"`
	// Offset from which the query starts.
	// If the offset is less than 0, the value is automatically converted to 0. Default to 0.
	Offset int `q:"offset"`
	// Number of items displayed on each page. The range of number is form 1 to 500, default to 20.
	Limit int `q:"limit"`
}

type ListHistoryOptsBuilder interface {
	ToListOptsQuery() (string, error)
}

func (opts ListApiHistoryOpts) ToListOptsQuery() (string, error) {
	q, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return "", err
	}
	return q.String(), err
}

// ListApiHistoryVersions is a method to obtain an array of one or more history versions of API according to the query
// parameters.
func ListApiHistoryVersions(client *golangsdk.ServiceClient, instanceId, apiId string,
	opts ListOptsBuilder) pagination.Pager {
	url := historyURL(client, instanceId, apiId)
	if opts != nil {
		query, err := opts.ToListOptsQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}

	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return ApiVersionPage{pagination.SinglePageBase(r)}
	})
}

type VersionSwitchOpts struct {
	// Api version ID.
	VersionId string `json:"version_id"`
}

type ApiVersionOptsBuilder interface {
	ToVersionSwitchMap() (map[string]interface{}, error)
}

func (opts VersionSwitchOpts) ToVersionSwitchMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

// SwitchSpecVersion is a method to switch an specifies version between multiple versions of an API.
func SwitchSpecVersion(client *golangsdk.ServiceClient, instanceId, apiId string,
	opts ApiVersionOptsBuilder) (r SwitchResult) {
	reqBody, err := opts.ToVersionSwitchMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Put(historyURL(client, instanceId, apiId), reqBody, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})
	return
}

// GetVersion is a method to obtain the runtime information for an existing API.
func GetVersion(client *golangsdk.ServiceClient, instanceId, apiId string) (r RuntimeResult) {
	_, r.Err = client.Get(versionURL(client, instanceId, apiId), &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})
	return
}

type RuntimeDefinitionOpts struct {
	// ID of the environment in which the API is published.
	EnvId string `q:"env_id"`
}

type RuntimeDefinitionOptsBuilder interface {
	ToRuntimeDefinitionOptsQuery() (string, error)
}

func (opts RuntimeDefinitionOpts) ToRuntimeDefinitionOptsQuery() (string, error) {
	q, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return "", err
	}
	return q.String(), err
}

// GetRuntimeDefinition is a method to obtain the runtime information for an existing API.
func GetRuntimeDefinition(client *golangsdk.ServiceClient, instanceId, apiId string,
	opts RuntimeDefinitionOptsBuilder) (r RuntimeResult) {
	url := runtimeURL(client, instanceId, apiId)
	if opts != nil {
		query, err := opts.ToRuntimeDefinitionOptsQuery()
		if err != nil {
			r.Err = err
			return
		}
		url += query
	}
	_, r.Err = client.Get(url, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})
	return
}

// OfflineApiById is a method to offline an published custom api by id.
func OfflineApiById(client *golangsdk.ServiceClient, instanceId, apiId string) (r DeleteResult) {
	_, r.Err = client.Delete(versionURL(client, instanceId, apiId), nil)
	return
}

type DebugOpts struct {
	// Message body, with a maximum length of 2,097,152 bytes.
	Body string `json:"body,omitempty"`
	// Header parameters, with each value being a character string array. Each parameter name must meet the following requirements:
	// Contains letters, digits, periods (.), or hyphens (-).
	// Starts with a letter, with a maximum of 32 bytes.
	// Case-insensitive and cannot start with X-Apig- or X-Sdk-.
	// Case-insensitive and cannot be X-Stage.
	// Case-insensitive and cannot be X-Auth-Token or Authorization when mode is set to MARKET or CONSUMER.
	// Each header name is normalized before use. For example, x-MY-hEaDer is normalized as X-My-Header.
	Header map[string]string `json:"header,omitempty"`
	// Request method.
	// Enumeration values:
	// GET
	// POST
	// PUT
	// DELETE
	// HEAD
	// PATCH
	// OPTIONS
	Method string `json:"method" required:"true"`
	// Debugging mode.
	// DEVELOPER: Debugging the definition of an API that has not been published.
	// MARKET: Debugging the definition of an API purchased from the marketplace.
	// CONSUMER: Debugging the definition of an API that has been published in a specified environment.
	// In DEVELOPER mode, the API caller must be the API provider.
	// In MARKET mode, the API caller must be the API purchaser or owner.
	// In CONSUMER mode, the API caller must be the API provider or has been authorized to access the API in a specific environment.
	Mode string `json:"mode" required:"true"`
	// Request path of the API, starting with a slash (/) and containing up to 1024 characters.
	// The request path must meet path requirements so that it can be correctly decoded after percent-encoding.
	Path string `json:"path" required:"true"`
	// Query strings, with each value being a character string array. Each parameter name must meet the following requirements:
	// Contains letters, digits, periods (.), hyphens (-), or underscores (_).
	// Starts with a letter, with a maximum of 32 bytes.
	// Case-insensitive and cannot start with X-Apig- or X-Sdk-.
	// Case-insensitive and cannot be X-Stage.
	Query map[string]string `json:"query,omitempty"`
	// Request protocol.
	// HTTP
	// HTTPS
	Scheme string `json:"scheme" required:"true"`
	// AppKey used in the debugging request.
	AppKey string `json:"app_key,omitempty"`
	// AppSecret used in the debugging request.
	AppSecret string `json:"app_secret,omitempty"`
	// Access domain name of the API. If no value is specified, one of the following default values will be used based on the setting of mode:
	// For DEVELOPER, the subdomain name of the API group will be used.
	// For MARKET, the domain name of the API group allocated by the marketplace will be used.
	// For CONSUMER, the subdomain name of the API group will be used.
	Domain string `json:"domain,omitempty"`
	// Running environment specified by the debugging request. This parameter is valid only when mode is set to CONSUMER. If this parameter is not specified, the following default value is used:
	// CONSUMER RELEASE
	Stage string `json:"stage,omitempty"`
}

type DebugOptsBuilder interface {
	ToApiDebugMap() (map[string]interface{}, error)
}

func (opts VersionSwitchOpts) ToApiDebugMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

// Debug is a method to test an custom API use Header, Request Body, etc.
func Debug(client *golangsdk.ServiceClient, instanceId, apiId string,
	opts DebugOptsBuilder) (r RuntimeResult) {
	reqBody, err := opts.ToApiDebugMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(debugURL(client, instanceId, apiId), reqBody, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})
	return
}

type BackendConnectionOpts struct {
	// Backend service configuration type. backend_address indicates the backend service address (no VPC channel is used), and vpc_channel indicates a VPC channel.
	// Enumeration values:
	// backend_address
	// vpc_channel
	BackendType string `json:"backend_type" required:"true"`
	// Backend service address. This parameter is required if backend_type is set to backend_address.
	BackendAddress string `json:"backend_address,omitempty"`
	// VPC channel ID. This parameter is required if backend_type is set to vpc_channel.
	VpcChannelId string `json:"vpc_channel_id,omitempty"`
}

type BackendConnectionBuilder interface {
	ToBackendConnectionMap() (map[string]interface{}, error)
}

func (opts BackendConnectionOpts) ToBackendConnectionMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

func CheckBackendConnectivity(client *golangsdk.ServiceClient, instanceId string,
	opts BackendConnectionBuilder) (r RuntimeResult) {
	reqBody, err := opts.ToBackendConnectionMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(connectivityURL(client, instanceId), reqBody, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})
	return
}
