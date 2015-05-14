package gsnapshot

import (
	"google.golang.org/api/compute/v1"
)

type SnapshotService interface {
	Create(diskID string, description string, zone string) (string, error)
	Delete(id string) error
	Find(id string) (*compute.Snapshot, bool, error)
}
