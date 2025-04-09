package operation

import (
	computebeta "google.golang.org/api/compute/v0.beta"
	"google.golang.org/api/compute/v1"
)

type Service interface {
	Waiter(operation *compute.Operation, zone string, region string) (*compute.Operation, error)
	WaiterB(operation *computebeta.Operation, zone string, region string) (*computebeta.Operation, error)
}
