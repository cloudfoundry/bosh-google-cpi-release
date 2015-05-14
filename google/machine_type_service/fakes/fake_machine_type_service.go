package fakes

import (
	"google.golang.org/api/compute/v1"
)

type FakeMachineTypeService struct {
	FindCalled      bool
	FindFound       bool
	FindMachineType *compute.MachineType
	FindErr         error
}

func (d *FakeMachineTypeService) Find(id string, zone string) (*compute.MachineType, bool, error) {
	d.FindCalled = true
	return d.FindMachineType, d.FindFound, d.FindErr
}
