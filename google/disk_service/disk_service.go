package gdisk

import (
	"google.golang.org/api/compute/v1"
)

type DiskService interface {
	Create(size int, diskType string, zone string) (string, error)
	Delete(id string) error
	Find(id string, zone string) (*compute.Disk, bool, error)
}
