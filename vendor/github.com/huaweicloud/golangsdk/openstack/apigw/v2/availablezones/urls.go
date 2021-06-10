package availablezones

import "github.com/huaweicloud/golangsdk"

const rootPath = "apigw/available-zones"

func rootURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(rootPath)
}
