package action

type Ping struct{}

func NewPing() Ping { return Ping{} }

func (p Ping) Run() (string, error) {
	return "pong", nil
}
