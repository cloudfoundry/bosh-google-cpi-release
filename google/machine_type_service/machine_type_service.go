package gmachinetype

type MachineTypeService interface {
	Find(id string, zone string) (MachineType, bool, error)
}
