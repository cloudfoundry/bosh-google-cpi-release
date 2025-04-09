package backendservice

import (
	"google.golang.org/api/compute/v1"

	"bosh-google-cpi/util"
)

type BackendService struct {
	Name     string
	Backends []Backend
	SelfLink string
}

type Backend struct {
	InstanceGroupLink string
	InstanceGroupID   string
	InstanceGroupZone string
}

// FromComputeBackends returns a list of Backends that make working
// with associated Instance Groups simpler. If zone is not nill, only backends
// with matching zones will be returned.
func FromComputeBackends(b []*compute.Backend) []Backend {
	backends := make([]Backend, len(b))
	for _, backend := range b {
		backends = append(backends, Backend{
			InstanceGroupLink: backend.Group,
			InstanceGroupID:   util.ResourceSplitter(backend.Group),
			InstanceGroupZone: util.ZoneFromURL(backend.Group),
		})
	}
	return backends
}
