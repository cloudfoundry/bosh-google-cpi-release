package disktype

type Service interface {
	Find(id string, zone string) (DiskType, bool, error)
}
