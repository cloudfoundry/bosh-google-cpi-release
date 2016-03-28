package operation

import (
	"google.golang.org/api/compute/v1"
)

type Service interface {
	Waiter(operation *compute.Operation, zone string, region string) (*compute.Operation, error)
}
