package registry

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
)

const httpClientLogTag = "RegistryHTTPClient"
const httpClientMaxAttemps = 5
const httpClientRetryDelay = 5

// HTTPClient represents a BOSH Registry Client.
type HTTPClient struct {
	options ClientOptions
	logger  boshlog.Logger
}

// NewHTTPClient creates a new BOSH Registry Client.
func NewHTTPClient(
	options ClientOptions,
	logger boshlog.Logger,
) HTTPClient {
	return HTTPClient{
		options: options,
		logger:  logger,
	}
}

// Delete deletes the instance settings for a given instance ID.
func (c HTTPClient) Delete(instanceID string) error {
	endpoint := fmt.Sprintf("%s/instances/%s/settings", c.options.EndpointWithCredentials(), instanceID)
	c.logger.Debug(httpClientLogTag, "Deleting agent settings from registry endpoint '%s'", endpoint)

	request, err := http.NewRequest("DELETE", endpoint, nil)
	if err != nil {
		return bosherr.WrapErrorf(err, "Creating DELETE request for registry endpoint '%s'", endpoint)
	}

	httpResponse, err := c.doRequest(request)
	if err != nil {
		return bosherr.WrapErrorf(err, "Deleting agent settings from registry endpoint '%s'", endpoint)
	}

	defer httpResponse.Body.Close()

	if httpResponse.StatusCode != http.StatusOK {
		return bosherr.Errorf("Received status code '%d' when deleting agent settings from registry endpoint '%s'", httpResponse.StatusCode, endpoint)
	}

	c.logger.Debug(httpClientLogTag, "Deleted agent settings from registry endpoint '%s'", endpoint)
	return nil
}

// Fetch gets the agent settings for a given instance ID.
func (c HTTPClient) Fetch(instanceID string) (AgentSettings, error) {
	endpoint := fmt.Sprintf("%s/instances/%s/settings", c.options.EndpointWithCredentials(), instanceID)
	c.logger.Debug(httpClientLogTag, "Fetching agent settings from registry endpoint '%s'", endpoint)

	request, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return AgentSettings{}, bosherr.WrapErrorf(err, "Creating GET request for registry endpoint '%s'", endpoint)
	}

	httpResponse, err := c.doRequest(request)
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

	var settingsResponse agentSettingsResponse
	err = json.Unmarshal(httpBody, &settingsResponse)
	if err != nil {
		return AgentSettings{}, bosherr.WrapErrorf(err, "Unmarshalling agent settings response from registry endpoint '%s', contents: '%s'", endpoint, httpBody)
	}

	var agentSettings AgentSettings
	err = json.Unmarshal([]byte(settingsResponse.Settings), &agentSettings)
	if err != nil {
		return AgentSettings{}, bosherr.WrapErrorf(err, "Unmarshalling agent settings response from registry endpoint '%s', contents: '%s'", endpoint, httpBody)
	}

	c.logger.Debug(httpClientLogTag, "Received agent settings from registry endpoint '%s', contents: '%s'", endpoint, httpBody)
	return agentSettings, nil
}

// Update updates the agent settings for a given instance ID. If there are not already agent settings for the instance, it will create ones.
func (c HTTPClient) Update(instanceID string, agentSettings AgentSettings) error {
	settingsJSON, err := json.Marshal(agentSettings)
	if err != nil {
		return bosherr.WrapErrorf(err, "Marshalling agent settings, contents: '%#v", agentSettings)
	}

	endpoint := fmt.Sprintf("%s/instances/%s/settings", c.options.EndpointWithCredentials(), instanceID)
	c.logger.Debug(httpClientLogTag, "Updating registry endpoint '%s' with agent settings '%s'", endpoint, settingsJSON)

	putPayload := bytes.NewReader(settingsJSON)
	request, err := http.NewRequest("PUT", endpoint, putPayload)
	if err != nil {
		return bosherr.WrapErrorf(err, "Creating PUT request for registry endpoint '%s' with agent settings '%s'", endpoint, settingsJSON)
	}

	httpResponse, err := c.doRequest(request)
	if err != nil {
		return bosherr.WrapErrorf(err, "Updating registry endpoint '%s' with agent settings: '%s'", endpoint, settingsJSON)
	}

	defer httpResponse.Body.Close()

	if httpResponse.StatusCode != http.StatusOK && httpResponse.StatusCode != http.StatusCreated {
		return bosherr.Errorf("Received status code '%d' when updating registry endpoint '%s' with agent settings: '%s'", httpResponse.StatusCode, endpoint, settingsJSON)
	}

	c.logger.Debug(httpClientLogTag, "Updated registry endpoint '%s' with agent settings '%s'", endpoint, settingsJSON)
	return nil
}

func (c HTTPClient) doRequest(request *http.Request) (httpResponse *http.Response, err error) {
	httpClient, err := c.httpClient()
	if err != nil {
		return nil, bosherr.WrapErrorf(err, "Creating HTTP Client")
	}

	retryDelay := time.Duration(httpClientRetryDelay) * time.Second
	for attempt := 0; attempt < httpClientMaxAttemps; attempt++ {
		httpResponse, err = httpClient.Do(request)
		if err == nil {
			return httpResponse, nil
		}
		c.logger.Debug(httpClientLogTag, "Performing registry HTTP call #%d got error '%v'", attempt, err)
		time.Sleep(retryDelay)
	}

	return nil, err
}

func (c HTTPClient) httpClient() (http.Client, error) {
	httpClient := http.Client{}

	if c.options.Protocol == "https" {
		certificates, err := tls.LoadX509KeyPair(c.options.TLS.CertFile, c.options.TLS.KeyFile)
		if err != nil {
			return httpClient, bosherr.WrapError(err, "Loading X509 Key Pair")
		}

		certPool := x509.NewCertPool()
		if c.options.TLS.CACertFile != "" {
			caCert, err := ioutil.ReadFile(c.options.TLS.CACertFile)
			if err != nil {
				return httpClient, bosherr.WrapError(err, "Loading CA certificate")
			}

			if !certPool.AppendCertsFromPEM(caCert) {
				return httpClient, bosherr.WrapError(err, "Invalid CA Certificate")
			}
		}

		tlsConfig := &tls.Config{
			Certificates:       []tls.Certificate{certificates},
			InsecureSkipVerify: c.options.TLS.InsecureSkipVerify,
			RootCAs:            certPool,
		}

		httpClient.Transport = &http.Transport{TLSClientConfig: tlsConfig}
	}

	return httpClient, nil
}
