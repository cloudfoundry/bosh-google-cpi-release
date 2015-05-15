package snapshot

type Service interface {
	Create(diskID string, description string, zone string) (string, error)
	Delete(id string) error
	Find(id string) (Snapshot, bool, error)
}
