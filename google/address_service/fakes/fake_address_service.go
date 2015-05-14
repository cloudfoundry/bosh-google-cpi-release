package fakes

import (
	"google.golang.org/api/compute/v1"
)

type FakeAddressService struct {
	FindCalled  bool
	FindFound   bool
	FindAddress *compute.Address
	FindErr     error

	FindByIPCalled  bool
	FindByIPFound   bool
	FindByIPAddress *compute.Address
	FindByIPErr     error
}

func (n *FakeAddressService) Find(id string, region string) (*compute.Address, bool, error) {
	n.FindCalled = true
	return n.FindAddress, n.FindFound, n.FindErr
}

func (n *FakeAddressService) FindByIP(ipAddress string) (*compute.Address, bool, error) {
	n.FindByIPCalled = true
	return n.FindByIPAddress, n.FindByIPFound, n.FindByIPErr
}
