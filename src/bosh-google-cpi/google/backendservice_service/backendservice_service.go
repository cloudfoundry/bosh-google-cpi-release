package backendservice

type Service interface {
	AddInstance(id, scheme, instanceId string) error
	RemoveInstance(vmLink string) error
}
