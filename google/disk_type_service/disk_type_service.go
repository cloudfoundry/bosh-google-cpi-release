package gdisktype

type DiskTypeService interface {
	Find(id string, zone string) (DiskType, bool, error)
}
