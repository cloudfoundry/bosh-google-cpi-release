package address

import (
	"fmt"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
)

func (a GoogleAddressService) FindByIP(ipAddress string) (Address, bool, error) {
	a.logger.Debug(googleAddressServiceLogTag, "Finding Google IP Address '%s'", ipAddress)
	filter := fmt.Sprintf("address eq .*%s", ipAddress)
	addresses, err := a.computeService.Addresses.AggregatedList(a.project).Filter(filter).Do()
	if err != nil {
		return Address{}, false, bosherr.WrapErrorf(err, "Failed to find Google IP Address '%s' in any region", ipAddress)
	}

	for _, addressItems := range addresses.Items {
		for _, addressItem := range addressItems.Addresses {
			// Return the first address (it can only be 1 address with the same IP across all regions)
			address := Address{
				Name:     addressItem.Name,
				SelfLink: addressItem.SelfLink,
			}
			return address, true, nil
		}
	}

	return Address{}, false, nil
}
