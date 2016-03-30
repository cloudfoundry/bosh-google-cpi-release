package machinetype

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"

	"bosh-google-cpi/util"
	"google.golang.org/api/googleapi"
)

func (m GoogleMachineTypeService) Find(id string, zone string) (MachineType, bool, error) {
	m.logger.Debug(googleMachineTypeServiceLogTag, "Finding Google Machine Type '%s' in zone '%s'", id, zone)
	machineTypeItem, err := m.computeService.MachineTypes.Get(m.project, util.ResourceSplitter(zone), id).Do()
	if err != nil {
		if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == 404 {
			return MachineType{}, false, nil
		}

		return MachineType{}, false, bosherr.WrapErrorf(err, "Failed to find Google Machine Type '%s' in zone '%s'", id, zone)
	}

	machineType := MachineType{
		Name:     machineTypeItem.Name,
		SelfLink: machineTypeItem.SelfLink,
		Zone:     machineTypeItem.Zone,
	}
	return machineType, true, nil
}
