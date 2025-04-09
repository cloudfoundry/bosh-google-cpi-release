package network

type Service interface {
	Find(proectId, id string) (Network, bool, error)
}
