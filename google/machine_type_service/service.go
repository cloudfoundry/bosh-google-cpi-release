package machinetype

type Service interface {
	Find(id string, zone string) (MachineType, bool, error)
}
