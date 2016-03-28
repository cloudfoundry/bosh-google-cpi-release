package targetpool

type Service interface {
	AddInstance(id string, vmLink string) error
	Find(id string, region string) (TargetPool, bool, error)
	FindByInstance(vmLink string, region string) (string, bool, error)
	List(region string) ([]TargetPool, error)
	RemoveInstance(id string, vmLink string) error
}
