package fakes

import (
	"bosh-google-cpi/google/machinetype"
)

type FakeMachineTypeService struct {
	FindCalled      bool
	FindFound       bool
	FindMachineType machinetype.MachineType
	FindErr         error

	CustomLinkCalled        bool
	CustomLinkMachineSeries string
	CustomLinkLink          string
}

func (d *FakeMachineTypeService) Find(id string, zone string) (machinetype.MachineType, bool, error) {
	d.FindCalled = true
	return d.FindMachineType, d.FindFound, d.FindErr
}

func (d *FakeMachineTypeService) CustomLink(cpu int, ram int, zone string, machineSeries string) string {
	d.CustomLinkCalled = true
	d.CustomLinkMachineSeries = machineSeries
	return d.CustomLinkLink
}
