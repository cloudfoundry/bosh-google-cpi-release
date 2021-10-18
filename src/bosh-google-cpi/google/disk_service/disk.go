package disk

type Disk struct {
	Name     string
	SelfLink string
	Status   string
	Zone     string
	Users    []string
	SizeGb   int64
}
