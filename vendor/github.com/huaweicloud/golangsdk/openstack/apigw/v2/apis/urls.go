package apis

import (
	"fmt"

	"github.com/huaweicloud/golangsdk"
)

const rootPath = "instances"

func buildRootPath(instanceId string) string {
	return fmt.Sprintf("instances/%s/apis", instanceId)
}

func rootURL(c *golangsdk.ServiceClient, instanceId string) string {
	return c.ServiceURL(buildRootPath(instanceId))
}

func resourceURL(c *golangsdk.ServiceClient, instanceId, apiId string) string {
	return c.ServiceURL(buildRootPath(instanceId), apiId)
}

func actionURL(c *golangsdk.ServiceClient, instanceId string) string {
	return c.ServiceURL(buildRootPath(instanceId), "action")
}

func releaseURL(c *golangsdk.ServiceClient, instanceId string) string {
	return c.ServiceURL(buildRootPath(instanceId), "publish")
}

func historyURL(c *golangsdk.ServiceClient, instanceId, apiId string) string {
	return c.ServiceURL(buildRootPath(instanceId), "publish", apiId)
}

func runtimeURL(c *golangsdk.ServiceClient, instanceId, apiId string) string {
	return c.ServiceURL(buildRootPath(instanceId), "runtime", apiId)
}

func versionURL(c *golangsdk.ServiceClient, instanceId, apiId string) string {
	return c.ServiceURL(buildRootPath(instanceId), "versions", apiId)
}

func debugURL(c *golangsdk.ServiceClient, instanceId, apiId string) string {
	return c.ServiceURL(buildRootPath(instanceId), "debug", apiId)
}

func connectivityURL(c *golangsdk.ServiceClient, instanceId string) string {
	return c.ServiceURL("instances", instanceId, "backend/connectivity/check")
}
