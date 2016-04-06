package fakes

import (
	"bosh-google-cpi/google/address_service"
)

type FakeAddressService struct {
	FindCalled  bool
	FindFound   bool
	FindAddress address.Address
	FindErr     error

	FindByIPCalled  bool
	FindByIPFound   bool
	FindByIPAddress address.Address
	FindByIPErr     error
}

func (n *FakeAddressService) Find(id string, region string) (address.Address, bool, error) {
	n.FindCalled = true
	return n.FindAddress, n.FindFound, n.FindErr
}

func (n *FakeAddressService) FindByIP(ipAddress string) (address.Address, bool, error) {
	n.FindByIPCalled = true
	return n.FindByIPAddress, n.FindByIPFound, n.FindByIPErr
}
