package availablezones

import "github.com/huaweicloud/golangsdk"

// List is a method to obtain an array containing one or more available zones.
func List(client *golangsdk.ServiceClient) (r AvailableZoneResult) {
	_, r.Err = client.Get(rootURL(client), &r.Body, nil)
	return
}
