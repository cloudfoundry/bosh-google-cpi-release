package registry

// Client represents a BOSH Registry Client.
type Client interface {
	Delete(instanceID string) error
	Fetch(instanceID string) (AgentSettings, error)
	Update(instanceID string, agentSettings AgentSettings) error
}
