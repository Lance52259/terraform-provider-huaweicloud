package signaturekeys

import (
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/pagination"
)

type commonResult struct {
	golangsdk.Result
}

// CreateResult represents a result of the Create method.
type CreateResult struct {
	commonResult
}

// UpdateResult represents a result of the Update method.
type UpdateResult struct {
	commonResult
}

type SignatureKey struct {
	// Signature key ID.
	Id string `json:"id"`
	// Signature name.
	Name string `json:"name"`
	// Signature key.
	SignKey string `json:"sign_key"`
	// Signature type, support 'hmac' and 'basic'.
	SignType string `json:"sign_type"`
	// Signature secret.
	SignSecret string `json:"sign_secret"`
	// Creation time.
	CreateTime string `json:"create_time"`
	// Update time.
	UpdateTime string `json:"update_time"`
	// Number of APIs.
	BindNum int `json:"bind_num"`
}

func (r commonResult) Extract() (*SignatureKey, error) {
	var s SignatureKey
	err := r.ExtractInto(&s)
	return &s, err
}

// SignatureKeyPage represents the response pages of the List method.
type SignatureKeyPage struct {
	pagination.SinglePageBase
}

func ExtractSignatureKeys(r pagination.Page) ([]SignatureKey, error) {
	var s []SignatureKey
	err := r.(SignatureKeyPage).Result.ExtractIntoSlicePtr(&s, "signs")
	return s, err
}

// DeleteResult represents a result of the Delete and DeleteVariable method.
type DeleteResult struct {
	golangsdk.ErrResult
}
