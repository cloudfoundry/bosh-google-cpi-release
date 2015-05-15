package fakes

import (
	"github.com/frodenas/bosh-google-cpi/google/machine_type_service"
)

type FakeMachineTypeService struct {
	FindCalled      bool
	FindFound       bool
	FindMachineType machinetype.MachineType
	FindErr         error
}

func (d *FakeMachineTypeService) Find(id string, zone string) (machinetype.MachineType, bool, error) {
	d.FindCalled = true
	return d.FindMachineType, d.FindFound, d.FindErr
}
