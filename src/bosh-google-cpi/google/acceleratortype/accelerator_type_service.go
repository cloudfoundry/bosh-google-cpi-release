package acceleratortype

type Service interface {
	Find(id string, zone string) (AcceleratorType, bool, error)
}
