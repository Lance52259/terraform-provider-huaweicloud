package signaturekeys

import "github.com/huaweicloud/golangsdk"

const rootPath = "instances"

func rootURL(c *golangsdk.ServiceClient, instanceId string) string {
	return c.ServiceURL(rootPath, instanceId, "signs")
}

func resourceURL(c *golangsdk.ServiceClient, instanceId, signId string) string {
	return c.ServiceURL(rootPath, instanceId, "signs", signId)
}
