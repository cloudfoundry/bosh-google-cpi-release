package acceleratortype

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	"google.golang.org/api/googleapi"

	"bosh-google-cpi/util"
)

func (m GoogleAcceleratorTypeService) Find(id string, zone string) (AcceleratorType, bool, error) {
	m.logger.Debug(googleAcceleratorTypeServiceLogTag, "Finding Google Accelerator Type '%s' in zone '%s'", id, zone)
	acceleratorTypeItem, err := m.computeService.AcceleratorTypes.Get(m.project, util.ResourceSplitter(zone), id).Do()
	if err != nil {
		if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == 404 {
			return AcceleratorType{}, false, nil
		}

		return AcceleratorType{}, false, bosherr.WrapErrorf(err, "Failed to find Google Accelerator Type '%s' in zone '%s'", id, zone)
	}

	acceleratorType := AcceleratorType{
		Name:     acceleratorTypeItem.Name,
		SelfLink: acceleratorTypeItem.SelfLink,
		Zone:     acceleratorTypeItem.Zone,
	}
	return acceleratorType, true, nil
}
