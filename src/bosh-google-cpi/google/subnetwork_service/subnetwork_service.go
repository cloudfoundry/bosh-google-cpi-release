package subnetwork

type Service interface {
	Find(projectId string, id string, region string) (Subnetwork, error)
}
