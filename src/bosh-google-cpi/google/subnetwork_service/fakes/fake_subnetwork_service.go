package fakes

import (
	"bosh-google-cpi/google/subnetwork_service"
)

type FakeSubnetworkService struct {
	FindCalled     bool
	FindFound      bool
	FindSubnetwork subnetwork.Subnetwork
	FindErr        error
}

func (s *FakeSubnetworkService) Find(id string, region string) (subnetwork.Subnetwork, bool, error) {
	s.FindCalled = true
	return s.FindSubnetwork, s.FindFound, s.FindErr
}
