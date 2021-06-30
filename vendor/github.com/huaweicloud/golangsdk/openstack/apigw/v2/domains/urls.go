package domains

import (
	"fmt"
	"github.com/huaweicloud/golangsdk"
)

func buildDomainRootPath(instanceId, groupId string) string {
	return fmt.Sprintf("instances/%s/api-groups/%s/domains", instanceId, groupId)
}

func rootURL(c *golangsdk.ServiceClient, instanceId, groupId string) string {
	return c.ServiceURL(buildDomainRootPath(instanceId, groupId))
}

func resourceURL(c *golangsdk.ServiceClient, instanceId, groupId, domainId string) string {
	return c.ServiceURL(buildDomainRootPath(instanceId, groupId), domainId)
}

func certRootURL(c *golangsdk.ServiceClient, instanceId, groupId, domainId string) string {
	return c.ServiceURL(buildDomainRootPath(instanceId, groupId), domainId, "certificate")
}

func certResourceURL(c *golangsdk.ServiceClient, instanceId, groupId, domainId, certId string) string {
	return c.ServiceURL(buildDomainRootPath(instanceId, groupId), domainId, "certificate", certId)
}