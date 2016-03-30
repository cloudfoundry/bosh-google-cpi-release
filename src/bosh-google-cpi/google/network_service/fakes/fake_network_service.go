package fakes

import (
	"bosh-google-cpi/google/network_service"
)

type FakeNetworkService struct {
	FindCalled  bool
	FindFound   bool
	FindNetwork network.Network
	FindErr     error
}

func (n *FakeNetworkService) Find(id string) (network.Network, bool, error) {
	n.FindCalled = true
	return n.FindNetwork, n.FindFound, n.FindErr
}
