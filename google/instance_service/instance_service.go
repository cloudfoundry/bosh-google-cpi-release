package ginstance

import (
	"google.golang.org/api/compute/v1"
)

type InstanceService interface {
	AttachDisk(id string, diskLink string) (string, string, error)
	AttachedDisks(id string) (InstanceAttachedDisks, error)
	Delete(id string) error
	DeleteNetworkConfiguration(id string, instanceNetworks GoogleInstanceNetworks) error
	DetachDisk(id string, diskID string) error
	Find(id string, zone string) (*compute.Instance, bool, error)
	Reboot(id string) error
	SetMetadata(id string, vmMetadata InstanceMetadata) error
	UpdateNetworkConfiguration(id string, instanceNetworks GoogleInstanceNetworks) error
}

type InstanceAttachedDisks []string

type InstanceMetadata map[string]interface{}
