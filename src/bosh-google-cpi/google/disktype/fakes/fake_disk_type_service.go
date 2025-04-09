package fakes

import (
	"bosh-google-cpi/google/disktype"
)

type FakeDiskTypeService struct {
	FindCalled   bool
	FindFound    bool
	FindDiskType disktype.DiskType
	FindErr      error
}

func (d *FakeDiskTypeService) Find(id string, zone string) (disktype.DiskType, bool, error) {
	d.FindCalled = true
	return d.FindDiskType, d.FindFound, d.FindErr
}
