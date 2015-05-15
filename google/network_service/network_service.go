package gnetwork

type NetworkService interface {
	Find(id string) (Network, bool, error)
}
