package ginstance

import (
	"google.golang.org/api/compute/v1"
)

type InstanceService interface {
	AddAccessConfig(id string, zone string, networkInterface string, accessConfig *compute.AccessConfig) error
	AddNetworkConfiguration(id string, instanceNetworks GoogleInstanceNetworks) error
	AttachDisk(id string, diskLink string) (string, string, error)
	AttachedDisks(id string) (InstanceAttachedDisks, error)
	CleanUp(id string)
	Create(vmProps *InstanceProperties, instanceNetworks GoogleInstanceNetworks, registryEndpoint string) (string, error)
	Delete(id string) error
	DeleteAccessConfig(id string, zone string, networkInterface string, accessConfig string) error
	DeleteNetworkConfiguration(id string, instanceNetworks GoogleInstanceNetworks) error
	DetachDisk(id string, diskID string) error
	Find(id string, zone string) (*compute.Instance, bool, error)
	Reboot(id string) error
	SetMetadata(id string, vmMetadata InstanceMetadata) error
	SetTags(id string, zone string, instanceTags *compute.Tags) error
	UpdateNetworkConfiguration(id string, instanceNetworks GoogleInstanceNetworks) error
}

type InstanceAttachedDisks []string

type InstanceMetadata map[string]interface{}

type InstanceProperties struct {
	Zone              string
	Stemcell          string
	MachineType       string
	RootDiskSizeGb    int
	RootDiskType      string
	AutomaticRestart  bool
	OnHostMaintenance string
	ServiceScopes     InstanceServiceScopes
}

type InstanceServiceScopes []string
