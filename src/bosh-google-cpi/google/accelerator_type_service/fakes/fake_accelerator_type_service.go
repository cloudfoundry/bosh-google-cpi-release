package fakes

import (
	"bosh-google-cpi/google/accelerator_type_service"
)

type FakeAcceleratorTypeService struct {
	FindCalled          bool
	FindFound           bool
	FindAcceleratorType acceleratortype.AcceleratorType
	FindErr             error
}

func (d *FakeAcceleratorTypeService) Find(id string, zone string) (acceleratortype.AcceleratorType, bool, error) {
	d.FindCalled = true
	return d.FindAcceleratorType, d.FindFound, d.FindErr
}
