package backendservice

import (
	"strings"

	"bosh-google-cpi/google/instance_group_service"
	"bosh-google-cpi/google/operation_service"

	"bosh-google-cpi/util"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/googleapi"
)

const googleBackendServiceServiceLogTag = "GoogleBackendServiceService"

type GoogleBackendServiceService struct {
	project               string
	computeService        *compute.Service
	instanceGroupsService instancegroup.Service
	operationService      operation.Service
	logger                boshlog.Logger
}

func NewGoogleBackendServiceService(
	project string,
	computeService *compute.Service,
	operationService operation.Service,
	logger boshlog.Logger,
) GoogleBackendServiceService {
	igms := instancegroup.NewGoogleInstanceGroupService(project, computeService, operationService, logger)
	return GoogleBackendServiceService{
		project:               project,
		instanceGroupsService: igms,
		computeService:        computeService,
		operationService:      operationService,
		logger:                logger,
	}
}

// AddInstance will add an instance to a Backend Service. A Backend Service may
// have more than one backend/instance group associated with it. In that case,
// the instance will be added to each backend/instance group that is in the
// same zone as the instance.
func (i GoogleBackendServiceService) AddInstance(id, vmLink string) error {
	zone := util.ZoneFromURL(vmLink)
	if zone == "" {
		return bosherr.Errorf("Could not find VM zone in %q", vmLink)
	}
	i.logger.Debug(googleBackendServiceServiceLogTag, "Adding instance %q to all backends for Backend Service %q in zone %q", vmLink, id, zone)
	region := util.RegionFromZone(zone)
	backendService, found, err := i.find(id, region)
	if err != nil {
		return err
	}

	if !found {
		return bosherr.WrapErrorf(err, "Backend Service %q does not exist", id)
	}

	instance, err := i.computeService.Instances.Get(i.project, zone, util.ResourceSplitter(vmLink)).Do()
	if err != nil {
		return bosherr.WrapErrorf(err, "Could not find instance while trying to add to backend service")
	}

	var added bool
	// All backends in a backend service
	for _, b := range backendService.Backends {
		// If backend's instance group is in the instance's zone, this is a candidate
		if b.InstanceGroupZone == zone {
			// TODO(evanbrown): Handle multiple network interfaces
			// Confirm that instance group is in same subnetwork as instance
			if ig, found, _ := i.instanceGroupsService.Find(b.InstanceGroupID, b.InstanceGroupZone); found && (ig.Subnetwork == "" || ig.Subnetwork == instance.NetworkInterfaces[0].Subnetwork) {
				if err = i.instanceGroupsService.AddInstance(b.InstanceGroupID, vmLink); err != nil {
					return bosherr.WrapErrorf(err, "Failed to add instance %q to Backend Service %q's instance group named %q", vmLink, id, b.InstanceGroupID)
				}
				added = true
			}
		}
	}

	if !added {
		return bosherr.Errorf("Backend Service %q does not contain any Unmanaged Instance Groups in zone %q", id, zone)
	}
	return nil
}

// RemoveInstance will remove an instance from all associated Backend Services.
func (i GoogleBackendServiceService) RemoveInstance(vmLink string) error {
	zone := util.ZoneFromURL(vmLink)
	if zone == "" {
		return bosherr.Errorf("Could not find VM zone in %q", vmLink)
	}
	i.logger.Debug(googleBackendServiceServiceLogTag, "Removing instance %q from all backends in zone %q", vmLink, zone)
	backendServices, err := i.findByInstance(vmLink)
	if err != nil {
		return err
	}

	// All backend services
	for _, bs := range backendServices {
		// All backends in a backend service
		for _, b := range bs.Backends {
			// If backend's instance group is in the instance's zone
			if b.InstanceGroupZone == zone {
				if err = i.instanceGroupsService.RemoveInstance(b.InstanceGroupID, vmLink); err != nil {
					return bosherr.WrapErrorf(err, "Failed to remove instance %q from Backend Service %q's instance group named %q", vmLink, bs.Name, b.InstanceGroupID)
				}
			}
		}
	}
	return nil
}

// Find locates a Backend Service by its ID. It returns the Backend Service and
// true if found.
// False and a nil error are returned if the Backend Service
// is not found.
// False and an error value are returned if an error occurred
// while trying to find the Backend Service.
func (i GoogleBackendServiceService) find(id, region string) (BackendService, bool, error) {
	i.logger.Debug(googleBackendServiceServiceLogTag, "Finding Google Backend Service %q", id)

	// Search for a matching backend service amongst an aggregated list containing both global and regional items
	aggregatedBackendServices, err := i.computeService.BackendServices.AggregatedList(i.project).Do()
	// TODO(craigatgoogle): Employ server-side name filtering once the API filter bug is fixed, https://b.corp.google.com/issues/80238913
	if err == nil {
		var backendService *compute.BackendService
		for _, scopedList := range aggregatedBackendServices.Items {
			for _, bs := range scopedList.BackendServices {
				if bs.Name == id && (bs.Region == "" || strings.Contains(bs.Region, region)) {
					// Ensure there doesn't exist a collision in names between global/regional backend services
					if backendService != nil {
						return BackendService{},
							false,
							bosherr.Errorf("Failed to find Google Backend Service %q, given ambiguous name.", id)
					}
					backendService = bs
				}
			}
		}

		if backendService != nil {
			return BackendService{
				Name:     backendService.Name,
				SelfLink: backendService.SelfLink,
				Backends: FromComputeBackends(backendService.Backends),
			}, true, nil
		}
	}

	if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == 404 {
		return BackendService{}, false, nil
	}
	return BackendService{}, false, bosherr.WrapErrorf(err, "Failed to find Google Backend Service %q", id)
}

// FindByInstance returns all Backend Services that an instance in a zone
// belongs to. An instance is not directly associated with a Backend
// Service. Rather, an instance is associated with an Unmanaged Instance
// Group, and that is associated with a Backend Service. Because an
// instance may be associated with more than one Instance Group, and
// an Instance Group may be associated with more than one Backend Service,
// it is possible for an instance to be associated with more than one
// Backend Service.
func (i GoogleBackendServiceService) findByInstance(vmLink string) ([]BackendService, error) {
	zone := util.ZoneFromURL(vmLink)
	if zone == "" {
		return nil, bosherr.Errorf("Could not find VM zone in %q", vmLink)
	}
	id := util.ResourceSplitter(vmLink)
	i.logger.Debug(googleBackendServiceServiceLogTag, "Finding all Google Backend Services that instance %q is associated with", id)

	found := make([]BackendService, 0)
	allBackendServices, err := i.list(zone)
	if err != nil {
		return found, bosherr.WrapErrorf(err, "Failed to list Google Backend Services for instance %q", id)
	}

	backendServices := make([]BackendService, 0)
	for _, bs := range allBackendServices {
		for _, b := range bs.Backends {
			if b.InstanceGroupZone == zone {
				if igs, found, _ := i.instanceGroupsService.Find(b.InstanceGroupLink, b.InstanceGroupZone); found {
					for _, instance := range igs.Instances {
						if instance == vmLink {
							i.logger.Debug(googleBackendServiceServiceLogTag, "Found instance %q in Instance Group %q", vmLink, igs.Name)
							backendServices = append(backendServices, bs)
							continue
						}
					}
				}
			}
		}
	}

	return backendServices, nil
}

// List returns a list of Backend Services that have one or more backend
// instance groups in the provided zone. The returned Backend Services will
// contain all backends/instance groups regardless of zone.
func (i GoogleBackendServiceService) list(zone string) ([]BackendService, error) {
	i.logger.Debug(googleBackendServiceServiceLogTag, "Finding all Google Backend Services with at least one instance group in zone %q", zone)
	backendServiceList, err := i.computeService.BackendServices.List(i.project).Do()
	if err != nil {
		return []BackendService{}, bosherr.WrapErrorf(err, "Failed to list Google Backend Services")
	}

	backendServices := make([]BackendService, 0)
	for _, bs := range backendServiceList.Items {
		for _, b := range bs.Backends {
			if strings.Contains(b.Group, zone) {
				backendService := BackendService{
					Name:     bs.Name,
					SelfLink: bs.SelfLink,
					Backends: FromComputeBackends(bs.Backends),
				}
				backendServices = append(backendServices, backendService)
				break
			}
		}
	}
	return backendServices, nil
}
