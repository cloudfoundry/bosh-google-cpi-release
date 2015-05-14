package fakes

import (
	"google.golang.org/api/compute/v1"
)

type FakeNetworkService struct {
	FindCalled  bool
	FindFound   bool
	FindNetwork *compute.Network
	FindErr     error
}

func (n *FakeNetworkService) Find(id string) (*compute.Network, bool, error) {
	n.FindCalled = true
	return n.FindNetwork, n.FindFound, n.FindErr
}
