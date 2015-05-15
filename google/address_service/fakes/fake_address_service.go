package fakes

import (
	"github.com/frodenas/bosh-google-cpi/google/address_service"
)

type FakeAddressService struct {
	FindCalled  bool
	FindFound   bool
	FindAddress gaddress.Address
	FindErr     error

	FindByIPCalled  bool
	FindByIPFound   bool
	FindByIPAddress gaddress.Address
	FindByIPErr     error
}

func (n *FakeAddressService) Find(id string, region string) (gaddress.Address, bool, error) {
	n.FindCalled = true
	return n.FindAddress, n.FindFound, n.FindErr
}

func (n *FakeAddressService) FindByIP(ipAddress string) (gaddress.Address, bool, error) {
	n.FindByIPCalled = true
	return n.FindByIPAddress, n.FindByIPFound, n.FindByIPErr
}
