package server

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"regexp"

	boshlog "github.com/cloudfoundry/bosh-agent/logger"

	"github.com/frodenas/bosh-google-cpi/registry/server/store"
)

const RegistryInstanceHandlerLogTag = "RegistryInstanceHandler"

type InstanceHandler struct {
	config        ListenerConfig
	registryStore store.RegistryStore
	logger        boshlog.Logger
}

func NewInstanceHandler(
	config ListenerConfig,
	registryStore store.RegistryStore,
	logger boshlog.Logger,
) *InstanceHandler {
	return &InstanceHandler{
		config:        config,
		registryStore: registryStore,
		logger:        logger,
	}
}

type SettingsResponse struct {
	Settings string `json:"settings"`
	Status   string `json:"status"`
}

func (ih *InstanceHandler) HandleFunc(w http.ResponseWriter, req *http.Request) {
	ih.logger.Debug(RegistryInstanceHandlerLogTag, "Received %s %s from %s", req.Method, req.URL.Path, req.RemoteAddr)
	instanceID, found := ih.getInstanceID(req)
	if !found {
		ih.logger.Debug(RegistryInstanceHandlerLogTag, "Instance ID not found in request:", req.Method)
		ih.handleNotFound(w)
		return
	}

	ih.logger.Debug(RegistryInstanceHandlerLogTag, "Found instance ID in request: '%s'", instanceID)

	if !ih.isAuthorized(req, instanceID) {
		ih.handleUnauthorized(w)
		return
	}

	switch req.Method {
	case "GET":
		ih.HandleGet(instanceID, w, req)
		return
	case "PUT":
		ih.HandlePut(instanceID, w, req)
		return
	case "DELETE":
		ih.HandleDelete(instanceID, w, req)
		return
	default:
		ih.handleNotFound(w)
		return
	}
}

func (ih *InstanceHandler) HandleGet(instanceID string, w http.ResponseWriter, req *http.Request) {
	settingsJSON, found, err := ih.registryStore.Get(instanceID)
	if err != nil {
		ih.logger.Debug(RegistryInstanceHandlerLogTag, "Failed to read settings for instance '%s': '%v'", err)
		ih.handleBadRequest(w)
		return
	}
	if !found {
		ih.logger.Debug(RegistryInstanceHandlerLogTag, "No settings for instance '%s' found", instanceID)
		ih.handleNotFound(w)
		return
	}

	ih.logger.Debug(RegistryInstanceHandlerLogTag, "Found settings for instance '%s': '%s'", instanceID, string(settingsJSON))

	response := SettingsResponse{
		Settings: string(settingsJSON),
		Status:   "ok",
	}

	responseJSON, err := json.Marshal(response)
	if err != nil {
		ih.handleBadRequest(w)
		return
	}

	w.Write(responseJSON)
}

func (ih *InstanceHandler) HandlePut(instanceID string, w http.ResponseWriter, req *http.Request) {
	reqBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		ih.handleBadRequest(w)
		return
	}

	ih.logger.Debug(RegistryInstanceHandlerLogTag, "Saving settings for instance '%s': '%s'", instanceID, string(reqBody))
	err = ih.registryStore.Save(instanceID, string(reqBody))
	if err != nil {
		ih.logger.Debug(RegistryInstanceHandlerLogTag, "Failed to save settings for instance '%s': '%v'", err)
		ih.handleBadRequest(w)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (ih *InstanceHandler) HandleDelete(instanceID string, w http.ResponseWriter, req *http.Request) {
	ih.logger.Debug(RegistryInstanceHandlerLogTag, "Deleting settings for instance '%s'", instanceID)
	err := ih.registryStore.Delete(instanceID)
	if err != nil {
		ih.logger.Debug(RegistryInstanceHandlerLogTag, "Failed to delete settings for instance '%s': '%v'", err)
		ih.handleBadRequest(w)
		return
	}
}

func (ih *InstanceHandler) getInstanceID(req *http.Request) (string, bool) {
	pattern := regexp.MustCompile("/instances/([^/]+)/settings")
	matches := pattern.FindStringSubmatch(req.URL.Path)

	if len(matches) == 0 {
		return "", false
	}

	return matches[1], true
}

func (ih *InstanceHandler) isAuthorized(req *http.Request, instanceID string) bool {
	auth := ih.config.Username + ":" + ih.config.Password
	expectedAuthorizationHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))

	return expectedAuthorizationHeader == req.Header.Get("Authorization")
}

func (ih *InstanceHandler) handleUnauthorized(w http.ResponseWriter) {
	ih.logger.Debug(RegistryInstanceHandlerLogTag, "Received unauthorized request")
	w.Header().Add("WWW-Authenticate", `Basic realm="Bosh Google Registry"`)
	w.WriteHeader(http.StatusUnauthorized)
}

func (ih *InstanceHandler) handleNotFound(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)

	settingsJSON, err := json.Marshal(SettingsResponse{Status: "not_found"})
	if err != nil {
		ih.logger.Warn(RegistryInstanceHandlerLogTag, "Failed to marshal 'not found' settings response: '%s'", err.Error())
		return
	}
	w.Write(settingsJSON)
}

func (ih *InstanceHandler) handleBadRequest(w http.ResponseWriter) {
	w.WriteHeader(http.StatusBadRequest)

	settingsJSON, err := json.Marshal(SettingsResponse{Status: "error"})
	if err != nil {
		ih.logger.Warn(RegistryInstanceHandlerLogTag, "Failed to marshal 'bad request' settings response: '%s'", err.Error())
		return
	}
	w.Write(settingsJSON)
}
