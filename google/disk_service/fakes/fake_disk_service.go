package fakes

import (
	"google.golang.org/api/compute/v1"
)

type FakeDiskService struct {
	CreateCalled   bool
	CreateErr      error
	CreateID       string
	CreateSize     int
	CreateDiskType string
	CreateZone     string

	DeleteCalled bool
	DeleteErr    error

	FindCalled bool
	FindFound  bool
	FindDisk   *compute.Disk
	FindErr    error
}

func (d *FakeDiskService) Create(size int, diskType string, zone string) (string, error) {
	d.CreateCalled = true
	d.CreateSize = size
	d.CreateDiskType = diskType
	d.CreateZone = zone
	return d.CreateID, d.CreateErr
}

func (d *FakeDiskService) Delete(id string) error {
	d.DeleteCalled = true
	return d.DeleteErr
}

func (d *FakeDiskService) Find(id string, zone string) (*compute.Disk, bool, error) {
	d.FindCalled = true
	return d.FindDisk, d.FindFound, d.FindErr
}
