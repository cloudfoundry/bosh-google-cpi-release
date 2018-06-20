package action

type Factory interface {
	Create(method string, ctx map[string]interface{}) (Action, error)
}
