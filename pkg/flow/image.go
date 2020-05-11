package flow

type Image struct {
	Id                 Id        `json:"id"`
	OperatingSystem    string    `json:"os"`
	Version            string    `json:"version"`
	Key                string    `json:"key"`
	Category           string    `json:"category"`
	Type               string    `json:"type"`
	MinRootDiskSize    int       `json:"min_root_disk_size"`
	Sorting            int       `json:"sorting"`
	RequiredLicenses   []Product `json:"required_licenses"`
	AvailableLocations []int     `json:"available_locations"`
}
