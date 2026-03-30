package fakes

import (
	"bosh-google-cpi/google/disk"
)

type FakeDiskService struct {
	CreateCalled   bool
	CreateErr      error
	CreateID       string
	CreateSize     int
	CreateDiskType string
	CreateZone     string

	CreateFromSnapshotCalled           bool
	CreateFromSnapshotErr              error
	CreateFromSnapshotID               string
	CreateFromSnapshotSnapshotSelfLink string
	CreateFromSnapshotSize             int
	CreateFromSnapshotDiskType         string
	CreateFromSnapshotZone             string

	DeleteCalled bool
	DeleteErr    error

	ResizeCalled bool
	ResizeErr    error

	FindCalled bool
	FindFound  bool
	FindDisk   disk.Disk
	FindErr    error
}

func (d *FakeDiskService) Create(size int, diskType string, zone string) (string, error) {
	d.CreateCalled = true
	d.CreateSize = size
	d.CreateDiskType = diskType
	d.CreateZone = zone
	return d.CreateID, d.CreateErr
}

func (d *FakeDiskService) CreateFromSnapshot(snapshotSelfLink string, size int, diskType string, zone string) (string, error) {
	d.CreateFromSnapshotCalled = true
	d.CreateFromSnapshotSnapshotSelfLink = snapshotSelfLink
	d.CreateFromSnapshotSize = size
	d.CreateFromSnapshotDiskType = diskType
	d.CreateFromSnapshotZone = zone
	return d.CreateFromSnapshotID, d.CreateFromSnapshotErr
}

func (d *FakeDiskService) Delete(id string) error {
	d.DeleteCalled = true
	return d.DeleteErr
}

func (d *FakeDiskService) Resize(id string, newSize int) error {
	d.ResizeCalled = true
	return d.ResizeErr
}

func (d *FakeDiskService) Find(id string, zone string) (disk.Disk, bool, error) {
	d.FindCalled = true
	return d.FindDisk, d.FindFound, d.FindErr
}
