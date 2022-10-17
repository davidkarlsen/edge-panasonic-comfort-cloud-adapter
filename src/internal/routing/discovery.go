package routing

import (
	"github.com/futurehomeno/cliffhanger/discovery"
)

// GetDiscoveryResource returns a service discovery configuration.
func GetDiscoveryResource() *discovery.Resource {
	return &discovery.Resource{
		ResourceName:           ResourceName,
		ResourceType:           discovery.ResourceTypeAd,
		ResourceFullName:       "Panasonic Comfort Cloud",
		Author:                 "support@futurehome.no",
		IsInstanceConfigurable: false,
		InstanceID:             "1",
		Version:                "1",
		AdapterInfo: discovery.AdapterInfo{
			Technology:            "panasonic_comfort_cloud",
			FwVersion:             "all",
			NetworkManagementType: "inclusion_exclusion",
		},
	}
}
