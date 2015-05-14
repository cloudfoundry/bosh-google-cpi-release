package ginstance

import (
	"google.golang.org/api/compute/v1"
)

type InstanceService interface {
	AttachedDisks(id string) (InstanceAttachedDisks, error)
	AttachDisk(id string, diskLink string) (string, string, error)
	DetachDisk(id string, diskID string) error
	Find(id string, zone string) (*compute.Instance, bool, error)
	Reboot(id string) error
	SetMetadata(id string, vmMetadata InstanceMetadata) error
}

type InstanceAttachedDisks []string

type InstanceMetadata map[string]interface{}
