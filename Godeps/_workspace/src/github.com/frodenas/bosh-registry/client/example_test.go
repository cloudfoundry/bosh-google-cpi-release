package registry_test

import (
	"fmt"

	"github.com/frodenas/bosh-registry/client"

	boshlog "github.com/cloudfoundry/bosh-utils/logger"
)

func Example() {
	var err error

	clientOptions := registry.ClientOptions{
		Protocol: "http",
		Host:     "127.0.0.1",
		Port:     25777,
		Username: "username",
		Password: "password",
	}
	logger := boshlog.NewLogger(boshlog.LevelDebug)
	registryClient := registry.NewHTTPClient(clientOptions, logger)

	instanceID := "instance-id"

	networksSettings := registry.NetworksSettings{}
	envSettings := registry.EnvSettings{}
	agentOptions := registry.AgentOptions{}
	settings := registry.NewAgentSettings("agent-id", "vm-id", networksSettings, envSettings, agentOptions)

	// Set the agent settings for a VM
	fmt.Printf("Updating settings for instance '%s' with '%#v'", instanceID, settings)
	err = registryClient.Update(instanceID, settings)
	if err != nil {
		fmt.Printf("Update call returned an error: %s", err)
	}

	// Get the agent settings for a VM
	settings, err = registryClient.Fetch(instanceID)
	if err != nil {
		fmt.Printf("Fetch call returned an error: %s", err)
	}
	fmt.Printf("Settings for instance '%s' are '%#v'", instanceID, settings)

	// Delete the agent settings for a VM
	err = registryClient.Delete(instanceID)
	if err != nil {
		fmt.Printf("Delete call returned an error: %s", err)
	}
}
