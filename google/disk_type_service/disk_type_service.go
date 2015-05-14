package gdisktype

import (
	"google.golang.org/api/compute/v1"
)

type DiskTypeService interface {
	Find(id string, zone string) (*compute.DiskType, bool, error)
}
