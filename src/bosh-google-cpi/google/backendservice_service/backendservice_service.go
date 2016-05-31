package backendservice

type Service interface {
	AddInstance(id, instanceId string) error
	RemoveInstance(vmLink string) error
}
