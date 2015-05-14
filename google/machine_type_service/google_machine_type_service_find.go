package gmachinetype

import (
	bosherr "github.com/cloudfoundry/bosh-agent/errors"

	"github.com/frodenas/bosh-google-cpi/util"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/googleapi"
)

func (m GoogleMachineTypeService) Find(id string, zone string) (*compute.MachineType, bool, error) {
	m.logger.Debug(googleMachineTypeServiceLogTag, "Finding Google Machine Type '%s' in zone '%s'", id, zone)
	machineType, err := m.computeService.MachineTypes.Get(m.project, util.ResourceSplitter(zone), id).Do()
	if err != nil {
		if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == 404 {
			return &compute.MachineType{}, false, nil
		}

		return &compute.MachineType{}, false, bosherr.WrapErrorf(err, "Failed to find Google Machine Type '%s' in zone '%s'", id, zone)
	}

	return machineType, true, nil
}
