package machinetype

type Service interface {
	Find(id string, zone string) (MachineType, bool, error)
	CustomLink(cpu int, ram int, zone string) string
}
