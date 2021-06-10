package availablezones

import "github.com/huaweicloud/golangsdk"

type commonResult struct {
	golangsdk.Result
}

type AvailableZoneResult struct {
	commonResult
}

type AvailableZone struct {
	// Available zone name.
	Name string `json:"name"`
	// Available zone ID.
	ID string `json:"id"`
	// Available zone code.
	Code string `json:"code"`
	// Available zone port.
	Port string `json:"port"`
	// Available zone names, contains english name and chinese name.
	Names LocalNames `json:"local_name"`
	// Gateway edititons available in the available zone.
	Editions map[string]bool `json:"specs"`
}

type LocalNames struct {
	EnglishName string `json:"en_us"`
	ChineseName string `json:"zh_cn"`
}

func (r AvailableZoneResult) Extract() ([]AvailableZone, error) {
	var s struct {
		AvailableZones []AvailableZone `json:"available_zones"`
	}
	err := r.ExtractInto(&s)
	return s.AvailableZones, err
}
