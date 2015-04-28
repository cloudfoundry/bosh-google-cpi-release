package registry

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	bosherr "github.com/cloudfoundry/bosh-agent/errors"
	boshlog "github.com/cloudfoundry/bosh-agent/logger"
)

const RegistryClientLogTag = "RegistryClient"
const RegistryClientMaxAttemps = 5
const RegistryClientRetryDelay = 5

type Client struct {
	options Options
	logger  boshlog.Logger
}

func NewClient(
	options Options,
	logger boshlog.Logger,
) Client {
	return Client{
		options: options,
		logger:  logger,
	}
}

func (c Client) Endpoint() string {
	return fmt.Sprintf("%s://%s:%s@%s:%d", c.options.Schema, c.options.Username, c.options.Password, c.options.Host, c.options.Port)
}

func (c Client) PublicEndpoint() string {
	return fmt.Sprintf("%s://%s:%d", c.options.Schema, c.options.Host, c.options.Port)
}

func (c Client) Delete(instanceID string) error {
	endpoint := fmt.Sprintf("%s/instances/%s/settings", c.Endpoint(), instanceID)
	c.logger.Debug(RegistryClientLogTag, "Deleting agent settings from registry endpoint '%s'", endpoint)

	httpClient := http.Client{}
	request, err := http.NewRequest("DELETE", endpoint, nil)
	if err != nil {
		return bosherr.WrapErrorf(err, "Creating DELETE request for registry endpoint '%s'", endpoint)
	}

	httpResponse, err := c.doHTTPRequest(httpClient, request)
	if err != nil {
		return bosherr.WrapErrorf(err, "Deleting agent settings from registry endpoint '%s'", endpoint)
	}

	defer httpResponse.Body.Close()

	if httpResponse.StatusCode != http.StatusOK {
		return bosherr.Errorf("Received status code '%d' when deleting agent settings from registry endpoint '%s'", httpResponse.StatusCode, endpoint)
	}

	return nil
}

func (c Client) Fetch(instanceID string) (AgentSettings, error) {
	endpoint := fmt.Sprintf("%s/instances/%s/settings", c.Endpoint(), instanceID)
	c.logger.Debug(RegistryClientLogTag, "Fetching agent settings from registry endpoint '%s'", endpoint)

	httpClient := http.Client{}
	request, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return AgentSettings{}, bosherr.WrapErrorf(err, "Creating GET request for registry endpoint '%s'", endpoint)
	}

	httpResponse, err := c.doHTTPRequest(httpClient, request)
	if err != nil {
		return AgentSettings{}, bosherr.WrapErrorf(err, "Fetching agent settings from registry endpoint '%s'", endpoint)
	}

	defer httpResponse.Body.Close()

	if httpResponse.StatusCode != http.StatusOK {
		return AgentSettings{}, bosherr.Errorf("Received status code '%d' when fetching agent settings from registry endpoint '%s'", httpResponse.StatusCode, endpoint)
	}

	httpBody, err := ioutil.ReadAll(httpResponse.Body)
	if err != nil {
		return AgentSettings{}, bosherr.WrapErrorf(err, "Reading agent settings response from registry endpoint '%s'", endpoint)
	}

	var settingsResponse AgentSettingsResponse
	err = json.Unmarshal(httpBody, &settingsResponse)
	if err != nil {
		return AgentSettings{}, bosherr.WrapErrorf(err, "Unmarshalling agent settings response from registry endpoint '%s', contents: '%s'", endpoint, httpBody)
	}

	var agentSettings AgentSettings
	err = json.Unmarshal([]byte(settingsResponse.Settings), &agentSettings)
	if err != nil {
		return AgentSettings{}, bosherr.WrapErrorf(err, "Unmarshalling agent settings response from registry endpoint '%s', contents: '%s'", endpoint, httpBody)
	}

	c.logger.Debug(RegistryClientLogTag, "Received agent settings from registry endpoint '%s', contents: '%s'", endpoint, httpBody)
	return agentSettings, nil
}

func (c Client) Update(instanceID string, agentSet AgentSettings) error {
	settingsJSON, err := json.Marshal(agentSet)
	if err != nil {
		return bosherr.WrapError(err, "Marshalling agent settings")
	}

	endpoint := fmt.Sprintf("%s/instances/%s/settings", c.Endpoint(), instanceID)
	c.logger.Debug(RegistryClientLogTag, "Updating registry endpoint '%s' with agent settings '%s'", endpoint, settingsJSON)

	httpClient := http.Client{}
	putPayload := bytes.NewReader(settingsJSON)
	request, err := http.NewRequest("PUT", endpoint, putPayload)
	if err != nil {
		return bosherr.WrapErrorf(err, "Creating PUT request for registry endpoint '%s' with agent settings '%s'", endpoint, settingsJSON)
	}

	httpResponse, err := c.doHTTPRequest(httpClient, request)
	if err != nil {
		return bosherr.WrapErrorf(err, "Updating registry endpoint '%s' with agent settings: '%s'", endpoint, settingsJSON)
	}

	defer httpResponse.Body.Close()

	if httpResponse.StatusCode != http.StatusOK && httpResponse.StatusCode != http.StatusCreated {
		return bosherr.Errorf("Received status code '%d' when updating registry endpoint '%s' with agent settings: '%s'", httpResponse.StatusCode, endpoint, settingsJSON)
	}

	return nil
}

func (c Client) doHTTPRequest(httpClient http.Client, request *http.Request) (httpResponse *http.Response, err error) {
	retryDelay := time.Duration(RegistryClientRetryDelay) * time.Second

	for attempt := 0; attempt < RegistryClientMaxAttemps; attempt++ {
		httpResponse, err = httpClient.Do(request)
		if err == nil {
			return httpResponse, nil
		}
		c.logger.Debug(RegistryClientLogTag, "Performing registry HTTP call #%d got error '%v'", attempt, err)
		time.Sleep(retryDelay)
	}

	return nil, err
}
