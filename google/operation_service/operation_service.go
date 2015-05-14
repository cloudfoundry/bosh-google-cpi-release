package goperation

import (
	"google.golang.org/api/compute/v1"
)

type OperationService interface {
	Waiter(operation *compute.Operation, zone string, region string) (*compute.Operation, error)
}
