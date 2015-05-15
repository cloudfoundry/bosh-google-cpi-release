package fakes

import (
	"github.com/frodenas/bosh-google-cpi/google/network_service"
)

type FakeNetworkService struct {
	FindCalled  bool
	FindFound   bool
	FindNetwork gnetwork.Network
	FindErr     error
}

func (n *FakeNetworkService) Find(id string) (gnetwork.Network, bool, error) {
	n.FindCalled = true
	return n.FindNetwork, n.FindFound, n.FindErr
}
