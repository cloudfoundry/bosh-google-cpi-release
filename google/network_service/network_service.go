package gnetwork

import (
	"google.golang.org/api/compute/v1"
)

type NetworkService interface {
	Find(id string) (*compute.Network, bool, error)
}
