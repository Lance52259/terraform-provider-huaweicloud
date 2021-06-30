package vpcchannels

import (
	"fmt"
	"github.com/huaweicloud/golangsdk"
)

func buildRootPath(instanceId string) string {
	return fmt.Sprintf("instances/%s/vpc-channels", instanceId)
}

func rootURL(c *golangsdk.ServiceClient, instanceId string) string {
	return c.ServiceURL(buildRootPath(instanceId), "vpc-channels")
}

func resourceURL(c *golangsdk.ServiceClient, instanceId, chanId string) string {
	return c.ServiceURL(buildRootPath(instanceId), "vpc-channels", chanId)
}

func rootMembersURL(c *golangsdk.ServiceClient, instanceId, chanId string) string {
	return c.ServiceURL(buildRootPath(instanceId), chanId, "members")
}

func resourceMembersURL(c *golangsdk.ServiceClient, instanceId, chanId, memberId string) string {
	return c.ServiceURL(buildRootPath(instanceId), chanId, "members", memberId)
}
