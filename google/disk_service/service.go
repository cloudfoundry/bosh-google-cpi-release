package disk

type Service interface {
	Create(size int, diskType string, zone string) (string, error)
	Delete(id string) error
	Find(id string, zone string) (Disk, bool, error)
}
