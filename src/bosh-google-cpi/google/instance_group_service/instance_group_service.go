package instancegroup

type Service interface {
	AddInstance(id string, vmLink string) error
	Find(id string, region string) (InstanceGroup, bool, error)
	FindByInstance(vmLink string, zone string) (string, bool, error)
	List(zone string) ([]InstanceGroup, error)
	ListInstances(id string, zone string) ([]string, error)
	RemoveInstance(id string, vmLink string) error
}
