package fakes

import (
	"github.com/frodenas/bosh-google-cpi/google/instance_service"
	"google.golang.org/api/compute/v1"
)

type FakeInstanceService struct {
	AddAccessConfigCalled bool
	AddAccessConfigErr    error

	AddNetworkConfigurationCalled bool
	AddNetworkConfigurationErr    error

	AttachDiskCalled     bool
	AttachDiskErr        error
	AttachDiskDeviceName string
	AttachDiskDevicePath string

	AttachedDisksCalled bool
	AttachedDisksErr    error
	AttachedDisksList   ginstance.InstanceAttachedDisks

	CleanUpCalled bool

	CreateCalled           bool
	CreateErr              error
	CreateID               string
	CreateVMProps          *ginstance.InstanceProperties
	CreateInstanceNetworks ginstance.GoogleInstanceNetworks
	CreateRegistryEndpoint string

	DeleteCalled bool
	DeleteErr    error

	DeleteAccessConfigCalled bool
	DeleteAccessConfigErr    error

	DeleteNetworkConfigurationCalled bool
	DeleteNetworkConfigurationErr    error

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

	SetTagsCalled bool
	SetTagsErr    error

	UpdateNetworkConfigurationCalled bool
	UpdateNetworkConfigurationErr    error
}

func (i *FakeInstanceService) AddAccessConfig(id string, zone string, networkInterface string, accessConfig *compute.AccessConfig) error {
	i.AddAccessConfigCalled = true
	return i.AddAccessConfigErr
}

func (i *FakeInstanceService) AddNetworkConfiguration(id string, instanceNetworks ginstance.GoogleInstanceNetworks) error {
	i.AddNetworkConfigurationCalled = true
	return i.AddNetworkConfigurationErr
}

func (i *FakeInstanceService) AttachDisk(id string, diskLink string) (string, string, error) {
	i.AttachDiskCalled = true
	return i.AttachDiskDeviceName, i.AttachDiskDevicePath, i.AttachDiskErr
}

func (i *FakeInstanceService) AttachedDisks(id string) (ginstance.InstanceAttachedDisks, error) {
	i.AttachedDisksCalled = true
	return i.AttachedDisksList, i.AttachedDisksErr
}

func (i *FakeInstanceService) CleanUp(id string) {
	i.CleanUpCalled = true
	return
}

func (i *FakeInstanceService) Create(vmProps *ginstance.InstanceProperties, instanceNetworks ginstance.GoogleInstanceNetworks, registryEndpoint string) (string, error) {
	i.CreateCalled = true
	i.CreateVMProps = vmProps
	i.CreateInstanceNetworks = instanceNetworks
	i.CreateRegistryEndpoint = registryEndpoint
	return i.CreateID, i.CreateErr
}

func (i *FakeInstanceService) Delete(id string) error {
	i.DeleteCalled = true
	return i.DeleteErr
}

func (i *FakeInstanceService) DeleteAccessConfig(id string, zone string, networkInterface string, accessConfig string) error {
	i.DeleteCalled = true
	return i.DeleteErr
}

func (i *FakeInstanceService) DeleteNetworkConfiguration(id string, instanceNetworks ginstance.GoogleInstanceNetworks) error {
	i.DeleteNetworkConfigurationCalled = true
	return i.DeleteNetworkConfigurationErr
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

func (i *FakeInstanceService) SetTags(id string, zone string, instanceTags *compute.Tags) error {
	i.SetTagsCalled = true
	return i.SetTagsErr
}

func (i *FakeInstanceService) UpdateNetworkConfiguration(id string, instanceNetworks ginstance.GoogleInstanceNetworks) error {
	i.UpdateNetworkConfigurationCalled = true
	return i.UpdateNetworkConfigurationErr
}
