package domains

import "github.com/huaweicloud/golangsdk"

type commonResult struct {
	golangsdk.Result
}

type CreateResult struct {
	commonResult
}

type UpdateResult struct {
	commonResult
}

type BindCertResult struct {
	commonResult
}

type GetCertResult struct {
	commonResult
}

type Domain struct {
	// Custom domain name.
	UrlDomain string `json:"url_domain"`
	// Domain ID.
	Id string `json:"id"`
	// CNAME resolution status.
	//     1: not resolved
	//     2: resolving
	//     3: resolved
	//     4: resolving failed
	Status int `json:"status"`
	// Minimum SSL version supported.
	MinSslVersion string `json:"min_ssl_version"`
	// Certificate name. Only certificate related API response support.
	SslName string `json:"ssl_name"`
	// Certificate ID. Only certificate related API response support.
	SslId string `json:"ssl_id"`
}

func (r commonResult) Extract() (*Domain, error) {
	var s Domain
	err := r.ExtractInto(&s)
	return &s, err
}

type DeleteResult struct {
	golangsdk.ErrResult
}
