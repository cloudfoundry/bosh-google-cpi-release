package gaddress

import (
	"fmt"

	bosherr "github.com/cloudfoundry/bosh-agent/errors"

	"google.golang.org/api/compute/v1"
)

func (a GoogleAddressService) FindByIP(ipAddress string) (*compute.Address, bool, error) {
	a.logger.Debug(googleAddressServiceLogTag, "Finding Google IP Address '%s'", ipAddress)
	filter := fmt.Sprintf("address eq .*%s", ipAddress)
	addresses, err := a.computeService.Addresses.AggregatedList(a.project).Filter(filter).Do()
	if err != nil {
		return &compute.Address{}, false, bosherr.WrapErrorf(err, "Failed to find Google IP Address '%s' in any region", ipAddress)
	}

	for _, addressItems := range addresses.Items {
		for _, address := range addressItems.Addresses {
			// Return the first address (it can only be 1 address with the same IP across all regions)
			return address, true, nil
		}
	}

	return &compute.Address{}, false, nil
}
