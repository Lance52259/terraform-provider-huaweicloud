package vpcchannels

import (
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/pagination"
)

// ChannelOpts allows to create a new vpc channel or update the existing vpc channel using given parameters.
type ChannelOpts struct {
	// VPC channel name, which can contain 3 to 64 characters, starting with a letter.
	// Only letters, digits, hyphens (-), and underscores (_) are allowed.
	// Chinese characters must be in UTF-8 or Unicode format.
	Name string `json:"name" required:"true"`
	// VPC channel type.
	//     1: private network ELB channel (to be deprecated)
	//     2: fast channel with the load balancing function
	Type string `json:"type" required:"true"`
	// Host port of the VPC channel. The value range is 1–65535.
	// This parameter is required and valid if the VPC channel type is set to 2.
	Port int `json:"port,omitempty"`
	// Distribution algorithm. The valid algorithms are as following:
	//     1: WRR
	//     2: WLC
	//     3: SH
	//     4: URI hashing
	// Default to 1 (WRR). This parameter is mandatory if the VPC channel type is set to 2.
	BalanceStrategy int `json:"balance_strategy,omitempty"`
	// Member type of the VPC channel. The valid types are 'ip' and 'ecs', default to ecs.
	// This parameter is required if the VPC channel type is set to 2.
	MemberType string `json:"member_type,omitempty"`
	// Backend server list. Only one backend server is included if the VPC channel type is set to 1.
	Members []MemberInfo `json:"members" required:"true"`
	// Health config of vpc channel.
	HealthConfig HealthConfig `json:"vpc_health_config"`
}

type MemberInfo struct {
	// Backend server address.
	// This parameter is valid when the member type is IP address.
	Host string `json:"host,omitempty"`
	// Backend server weight. The valid value is range from 0 to 100.
	// The higher the weight is, the more requests a cloud server will receive.
	// The weight is only available for the WRR and WLC algorithms and it is valid only when the VPC channel type is set
	// to 2.
	Weight int `json:"weight,omitempty"`
	// Backend server ID. This parameter is valid when the member type is instance.
	// The value can contain 1 to 64 characters, including letters, digits, hyphens (-), and underscores (_).
	EcsId string `json:"ecs_id,omitempty"`
	// Backend server name which can contain 1 to 64 characters, including letters, digits, periods (.), hyphens (-),
	// and underscores (_).
	// This parameter is valid when the member type is instance.
	EcsName string `json:"ecs_name,omitempty"`
}

type HealthConfig struct {
	// Protocol for performing health checks on backend servers in the VPC channel.
	// The supported protocols are as follows: TCP, HTTP and HTTPS.
	Protocol string `json:"protocol" required:"true"`
	// Destination path for health checks.
	// This parameter is required if protocol is set to http.
	Path string `json:"path,omitempty"`
	// Request method for health checks. The valid options are GET and HEAD.
	Method string `json:"method"`
	// Destination port for health checks. By default, the host port of the VPC channel is used.
	// The valid value is range from 1 to 65535.
	Port int `json:"port"`
	// Healthy threshold, which refers to the number of consecutive successful checks required for a backend server to
	// be considered healthy.
	// The valid value is range from 2 to 10.
	ThresholdNormal int `json:"threshold_normal"`
	// Unhealthy threshold, which refers to the number of consecutive failed checks required for a backend server to be
	// considered unhealthy.
	// The valid value is range from 2 to 10.
	ThresholdAbnormal int `json:"threshold_abnormal"`
	// Interval between consecutive checks, in second.
	// The value must be greater than the value of timeout.
	// The valid range is from 5 to 300.
	TimeInterval int `json:"time_interval"`
	// Response codes for determining a successful HTTP response.
	// The value can be any integer within 100–599 in one of the following formats:
	//     Multiple values: 200,201,202
	//     Range: 200-299
	//     Multiple values and ranges: 201,202,210-299.
	// This parameter is required if protocol is set to http.
	HttpCode string `json:"http_code"`
	// Indicates whether to enable two-way authentication. If this function is enabled, the certificate specified in the backend_client_certificate configuration item of the gateway is used.
	// Default: false
	EnableClientSsl bool `json:"enable_client_ssl"`
	// Timeout for determining whether a health check fails, in second.
	// The value must be less than the value of time_interval.
	// The valid value is range from 2 to 30.
	Timeout int `json:"timeout"`
}

type ChannelOptsBuilder interface {
	ToChannelOptsMap() (map[string]interface{}, error)
}

func (opts ChannelOpts) ToChannelOptsMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

// Create is a method by which to create function that create a new vpc channel.
func Create(client *golangsdk.ServiceClient, instanceId string, opts ChannelOptsBuilder) (r CreateResult) {
	reqBody, err := opts.ToChannelOptsMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(rootURL(client, instanceId), reqBody, &r.Body, nil)
	return
}

// Update is a method by which to create function that update the existing .
func Update(client *golangsdk.ServiceClient, instanceId, chanId string,
	opts ChannelOptsBuilder) (r UpdateResult) {
	reqBody, err := opts.ToChannelOptsMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Put(resourceURL(client, instanceId, chanId), reqBody, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})
	return
}

// Get is a method to obtain the specified custom response according to the instanceId, appId and respId.
func Get(client *golangsdk.ServiceClient, instanceId, chanId string) (r GetResult) {
	_, r.Err = client.Get(resourceURL(client, instanceId, chanId), &r.Body, nil)
	return
}

// ListOpts allows to filter list data using given parameters.
type ListOpts struct {
	// Channel ID.
	Id string `q:"id"`
	// Channel name.
	Name string `q:"name"`
	// Channel type.
	VpcType string `q:"vpc_type"`
	// Offset from which the query starts.
	// If the offset is less than 0, the value is automatically converted to 0. Default to 0.
	Offset int `q:"offset"`
	// Number of items displayed on each page. The valid values are range form 1 to 500, default to 20.
	Limit int `q:"limit"`
	//
	PreciseSearch string `q:"precise_search"`
}

type ListOptsBuilder interface {
	ToListQuery() (string, error)
}

func (opts ListOpts) ToListQuery() (string, error) {
	q, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return "", err
	}
	return q.String(), err
}

// List is a method to obtain an array of one or more vpc channel according to the query parameters.
func List(client *golangsdk.ServiceClient, instanceId string, opts ListOptsBuilder) pagination.Pager {
	url := rootURL(client, instanceId)
	if opts != nil {
		query, err := opts.ToListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}

	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return ChannelPage{pagination.SinglePageBase(r)}
	})
}

// Delete is a method to delete the existing custom response.
func Delete(client *golangsdk.ServiceClient, instanceId, chanId string) (r DeleteResult) {
	_, r.Err = client.Delete(resourceURL(client, instanceId, chanId), nil)
	return
}
