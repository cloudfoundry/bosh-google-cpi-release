package action

type Factory interface {
	Create(method string) (Action, error)
}
