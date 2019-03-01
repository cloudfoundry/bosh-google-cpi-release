package fakes

import (
	"bosh-google-cpi/google/instance_service"
	"google.golang.org/api/compute/v1"
)

type FakeInstanceService struct {
	AddAccessConfigCalled bool
	AddAccessConfigErr    error

	AttachDiskCalled     bool
	AttachDiskErr        error
	AttachDiskDeviceName string
	AttachDiskDevicePath string

	AttachedDisksCalled bool
	AttachedDisksErr    error
	AttachedDisksList   instance.AttachedDisks

	CleanUpCalled bool

	CreateCalled           bool
	CreateErr              error
	CreateID               string
	CreateVMProps          *instance.Properties
	CreateNetworks         instance.Networks
	CreateRegistryEndpoint string

	DeleteCalled bool
	DeleteErr    error

	DeleteAccessConfigCalled bool
	DeleteAccessConfigErr    error

	DetachDiskCalled bool
	DetachDiskErr    error

	DiskDetailCalled bool

	FindCalled   bool
	FindFound    bool
	FindInstance *compute.Instance
	FindErr      error

	RebootCalled bool
	RebootErr    error

	SetMetadataCalled     bool
	SetMetadataErr        error
	SetMetadataVMMetadata instance.Metadata

	SetTagsCalled bool
	SetTagsErr    error

	UpdateNetworkConfigurationCalled bool
	UpdateNetworkConfigurationErr    error
}

func (i *FakeInstanceService) AddAccessConfig(id string, zone string, networkInterface string, accessConfig *compute.AccessConfig) error {
	i.AddAccessConfigCalled = true
	return i.AddAccessConfigErr
}

func (i *FakeInstanceService) AttachDisk(id string, diskLink string) (*instance.DiskAttachmentDetail, error) {
	i.AttachDiskCalled = true
	return i.DiskDetail(id, diskLink)
}

func (i *FakeInstanceService) DiskDetail(vmID string, diskLink string) (*instance.DiskAttachmentDetail, error) {
	i.DiskDetailCalled = true
	return &instance.DiskAttachmentDetail{
		Name: i.AttachDiskDeviceName,
		Path: i.AttachDiskDevicePath,
	}, i.AttachDiskErr
}

func (i *FakeInstanceService) AttachedDisks(id string) (instance.AttachedDisks, error) {
	i.AttachedDisksCalled = true
	return i.AttachedDisksList, i.AttachedDisksErr
}

func (i *FakeInstanceService) CleanUp(id string) {
	i.CleanUpCalled = true
	return
}

func (i *FakeInstanceService) Create(vmProps *instance.Properties, networks instance.Networks, registryEndpoint string) (string, error) {
	i.CreateCalled = true
	i.CreateVMProps = vmProps
	i.CreateNetworks = networks
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

func (i *FakeInstanceService) SetMetadata(id string, vmMetadata instance.Metadata) error {
	i.SetMetadataCalled = true
	i.SetMetadataVMMetadata = vmMetadata
	return i.SetMetadataErr
}

func (i *FakeInstanceService) SetTags(id string, zone string, instanceTags *compute.Tags) error {
	i.SetTagsCalled = true
	return i.SetTagsErr
}

func (i *FakeInstanceService) UpdateNetworkConfiguration(id string, networks instance.Networks) error {
	i.UpdateNetworkConfigurationCalled = true
	return i.UpdateNetworkConfigurationErr
}
