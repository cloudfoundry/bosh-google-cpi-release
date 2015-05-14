package gmachinetype

import (
	"google.golang.org/api/compute/v1"
)

type MachineTypeService interface {
	Find(id string, zone string) (*compute.MachineType, bool, error)
}
