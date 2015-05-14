package fakes

import (
	"github.com/frodenas/bosh-google-cpi/google/instance_service"
	"google.golang.org/api/compute/v1"
)

type FakeInstanceService struct {
	AttachedDisksCalled bool
	AttachedDisksErr    error
	AttachedDisksList   ginstance.InstanceAttachedDisks

	AttachDiskCalled     bool
	AttachDiskErr        error
	AttachDiskDeviceName string
	AttachDiskDevicePath string

	DetachDiskCalled bool
	DetachDiskErr    error

	FindCalled   bool
	FindFound    bool
	FindInstance *compute.Instance
	FindErr      error

	RebootCalled bool
	RebootErr    error

	SetMetadataCalled     bool
	SetMetadataErr        error
	SetMetadataVMMetadata ginstance.InstanceMetadata
}

func (i *FakeInstanceService) AttachedDisks(id string) (ginstance.InstanceAttachedDisks, error) {
	i.AttachedDisksCalled = true
	return i.AttachedDisksList, i.AttachedDisksErr
}

func (i *FakeInstanceService) AttachDisk(id string, diskLink string) (string, string, error) {
	i.AttachDiskCalled = true
	return i.AttachDiskDeviceName, i.AttachDiskDevicePath, i.AttachDiskErr
}

func (i *FakeInstanceService) DetachDisk(id string, diskID string) error {
	i.DetachDiskCalled = true
	return i.DetachDiskErr
}

func (i *FakeInstanceService) Find(id string, zone string) (*compute.Instance, bool, error) {
	i.FindCalled = true
	return i.FindInstance, i.FindFound, i.FindErr
}

func (i *FakeInstanceService) Reboot(id string) error {
	i.RebootCalled = true
	return i.RebootErr
}

func (i *FakeInstanceService) SetMetadata(id string, vmMetadata ginstance.InstanceMetadata) error {
	i.SetMetadataCalled = true
	i.SetMetadataVMMetadata = vmMetadata
	return i.SetMetadataErr
}
