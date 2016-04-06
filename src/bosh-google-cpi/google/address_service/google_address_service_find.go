package address

import (
	"fmt"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"

	"bosh-google-cpi/util"
	"google.golang.org/api/googleapi"
)

func (a GoogleAddressService) Find(id string, region string) (Address, bool, error) {
	if region == "" {
		a.logger.Debug(googleAddressServiceLogTag, "Finding Google Address '%s'", id)
		filter := fmt.Sprintf("name eq .*%s", id)
		addresses, err := a.computeService.Addresses.AggregatedList(a.project).Filter(filter).Do()
		if err != nil {
			return Address{}, false, bosherr.WrapErrorf(err, "Failed to find Google Address '%s'", id)
		}

		for _, addressItems := range addresses.Items {
			for _, addressItem := range addressItems.Addresses {
				// Return the first address (it can only be 1 address with the same name across all regions)
				address := Address{
					Name:     addressItem.Name,
					SelfLink: addressItem.SelfLink,
				}
				return address, true, nil
			}
		}

		return Address{}, false, nil
	}

	a.logger.Debug(googleAddressServiceLogTag, "Finding Google Address '%s' in region '%s'", id, region)
	addressItem, err := a.computeService.Addresses.Get(a.project, util.ResourceSplitter(region), id).Do()
	if err != nil {
		if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == 404 {
			return Address{}, false, nil
		}

		return Address{}, false, bosherr.WrapErrorf(err, "Failed to find Google Address '%s' in region '%s'", id, region)
	}

	address := Address{
		Name:     addressItem.Name,
		SelfLink: addressItem.SelfLink,
	}
	return address, true, nil
}
