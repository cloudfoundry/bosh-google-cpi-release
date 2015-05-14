package gaddress

import (
	"google.golang.org/api/compute/v1"
)

type AddressService interface {
	Find(id string, region string) (*compute.Address, bool, error)
	FindByIP(ipAddress string) (*compute.Address, bool, error)
}
