package vpcchannels

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

// GetResult represents a result of the Get operation.
type GetResult struct {
	commonResult
}

// UpdateResult represents a result of the Update operation.
type UpdateResult struct {
	commonResult
}

type Channel struct {
	// Channel ID.
	Id string `json:"id"`
	// Channel name.
	Name string `json:"name"`
	// VPC channel type.
	//     1: private network ELB channel (to be deprecated)
	//     2: fast channel with the load balancing function
	Type int `json:"type"`
	// Host port of the VPC channel. The value range is 1–65535.
	// This parameter is valid only when the VPC channel type is set to 2.
	Port int `json:"port"`
	// Distribution algorithm.
	//     1: WRR
	//     2: WLC
	//     3: SH
	//     4: URI hashing
	BalanceStrategy int `json:"balance_strategy"`
	// Member type of the VPC channel, contains 'ip' and 'ecs'.
	MemberType string `json:"member_type"`
	// Time when the VPC channel is created.
	CreateTime string `json:"create_time"`
	// VPC channel status.
	//     1: normal
	//     2: abnormal
	Status int `json:"status"`
	// ID of a private network ELB channel.
	ElbId string `json:"elb_id"`
}

// Extract is a method to extract an response struct.
func (r commonResult) Extract() (*Channel, error) {
	var s Channel
	err := r.ExtractInto(&s)
	return &s, err
}

// ChannelPage represents the response pages of the List operation.
type ChannelPage struct {
	pagination.SinglePageBase
}

// ExtractChannels is a method to extract an response struct list.
func ExtractChannels(r pagination.Page) ([]Channel, error) {
	var s []Channel
	err := r.(ChannelPage).Result.ExtractIntoSlicePtr(&s, "vpc_channels")
	return s, err
}

// DeleteResult represents a result of the Delete method.
type DeleteResult struct {
	golangsdk.ErrResult
}
