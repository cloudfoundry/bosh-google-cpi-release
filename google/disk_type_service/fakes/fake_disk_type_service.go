package fakes

import (
	"google.golang.org/api/compute/v1"
)

type FakeDiskTypeService struct {
	FindCalled   bool
	FindFound    bool
	FindDiskType *compute.DiskType
	FindErr      error
}

func (d *FakeDiskTypeService) Find(id string, zone string) (*compute.DiskType, bool, error) {
	d.FindCalled = true
	return d.FindDiskType, d.FindFound, d.FindErr
}
