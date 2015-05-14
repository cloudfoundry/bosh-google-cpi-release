package gtargetpool

import (
	"google.golang.org/api/compute/v1"
)

type TargetPoolService interface {
	AddInstance(id string, vmLink string) error
	Find(id string, region string) (*compute.TargetPool, bool, error)
	FindByInstance(vmLink string, region string) (string, bool, error)
	List(region string) ([]*compute.TargetPool, error)
	RemoveInstance(id string, vmLink string) error
}
