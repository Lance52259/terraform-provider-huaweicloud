package domains

import (
	"github.com/huaweicloud/golangsdk"
)

// The DomainOpts allows to create a new domain for specifies APIG group or update an existing domain.
type DomainOpts struct {
	//Minimum SSL version. 'TLSv1.1' and 'TLSv1.2' are supported，default to TLSv1.1.
	MinSslVersion string `json:"min_ssl_version,omitempty"`
	// Custom domain name, which can contain a maximum of 255 characters
	// and must comply with domain name specifications.
	// Notes: the update method does not support this parameter.
	UrlDomain string `json:"url_domain" required:"true"`
}

type DomainOptsBuilder interface {
	ToDomainOptsMap() (map[string]interface{}, error)
}

func (opts DomainOpts) ToDomainOptsMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

// Create is a method by which to create function that create a new environment.
func Create(client *golangsdk.ServiceClient, instanceId, groupId string, opts DomainOptsBuilder) (r CreateResult) {
	reqBody, err := opts.ToDomainOptsMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(rootURL(client, instanceId, groupId), reqBody, &r.Body, nil)
	return
}

func Update(client *golangsdk.ServiceClient, instanceId, groupId, domainId string,
	opts DomainOptsBuilder) (r UpdateResult) {
	reqBody, err := opts.ToDomainOptsMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Put(resourceURL(client, instanceId, groupId, domainId), reqBody, &r.Body, nil)
	return
}

func Delete(client *golangsdk.ServiceClient, instanceId, groupId, domainId string) (r DeleteResult) {
	_, r.Err = client.Delete(resourceURL(client, instanceId, groupId, domainId), nil)
	return
}

// The CertificateOpts allows to bind the certificate to a specifies domain.
type CertificateOpts struct {
	// Certificate content.
	CertContent string `json:"cert_content" required:"true"`
	// Certificate name, which can contain 4 to 50 characters, starting with a letter.
	// Only letters, digits, and underscores (_) are allowed.
	Name string `json:"name" required:"ture"`
	// Private key.
	PrivateKey string `json:"private_key" required:"true"`
}

type CertificateOptsBuilder interface {
	ToCertificateOptsMap() (map[string]interface{}, error)
}

func (opts DomainOpts) ToCertificateOptsMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

func BindCertificate(client *golangsdk.ServiceClient, instanceId, groupId, domainId string,
	opts CertificateOptsBuilder) (r BindCertResult) {
	reqBody, err := opts.ToCertificateOptsMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(certRootURL(client, instanceId, groupId, domainId), reqBody, &r.Body, nil)
	return
}

func GetCertificate(client *golangsdk.ServiceClient, instanceId, groupId, domainId, certId string) (r GetCertResult) {
	_, r.Err = client.Get(certResourceURL(client, instanceId, groupId, domainId, certId), &r.Body, nil)
	return
}

func RemoveCertificate(client *golangsdk.ServiceClient, instanceId, groupId, domainId, certId string) (r DeleteResult) {
	_, r.Err = client.Delete(certResourceURL(client, instanceId, groupId, domainId, certId), nil)
	return
}