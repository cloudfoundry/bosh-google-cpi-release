package registry

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	boshlog "github.com/cloudfoundry/bosh-agent/logger"
)

var _ = Describe("RegistryService", func() {
	var (
		logger               boshlog.Logger
		registryService      RegistryService
		registryServer       *registryServer
		instanceID           string
		expectedAgentSet     AgentSettings
		expectedAgentSetJSON []byte
	)

	BeforeEach(func() {
		registryOptions := RegistryOptions{
			Schema:   "http",
			Host:     "127.0.0.1",
			Port:     6307,
			Username: "fake-username",
			Password: "fake-password",
		}
		registryServer = NewRegistryServer(registryOptions)
		readyCh := make(chan struct{})
		go registryServer.Start(readyCh)
		<-readyCh

		instanceID = "fake-instance-id"
		logger = boshlog.NewLogger(boshlog.LevelNone)
		registryService = NewRegistryService(registryOptions, logger)

		expectedAgentSet = AgentSettings{AgentID: "fake-agent-id"}
		var err error
		expectedAgentSetJSON, err = json.Marshal(expectedAgentSet)
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		registryServer.Stop()
	})

	Describe("Delete", func() {
		Context("when settings for the instance exist in the registry", func() {
			BeforeEach(func() {
				registryServer.InstanceSettings = expectedAgentSetJSON
			})

			It("deletes settings in the registry", func() {
				err := registryService.Delete(instanceID)
				Expect(err).ToNot(HaveOccurred())
				Expect(registryServer.InstanceSettings).To(Equal([]byte{}))
			})
		})

		Context("when settings for instance do not exist", func() {
			It("returns an error", func() {
				err := registryService.Delete(instanceID)
				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("Fetch", func() {
		Context("when settings for the instance exist in the registry", func() {
			BeforeEach(func() {
				registryServer.InstanceSettings = expectedAgentSetJSON
			})

			It("fetches settings from the registry", func() {
				agentSet, err := registryService.Fetch(instanceID)
				Expect(err).ToNot(HaveOccurred())
				Expect(agentSet).To(Equal(expectedAgentSet))
			})
		})

		Context("when settings for instance do not exist", func() {
			It("returns an error", func() {
				agentSet, err := registryService.Fetch(instanceID)
				Expect(err).To(HaveOccurred())
				Expect(agentSet).To(Equal(AgentSettings{}))
			})
		})
	})

	Describe("Update", func() {
		It("updates settings in the registry", func() {
			Expect(registryServer.InstanceSettings).To(Equal([]byte{}))
			err := registryService.Update(instanceID, expectedAgentSet)
			Expect(err).ToNot(HaveOccurred())
			Expect(registryServer.InstanceSettings).To(Equal(expectedAgentSetJSON))
		})
	})

})

type registryServer struct {
	InstanceSettings []byte
	options          RegistryOptions
	listener         net.Listener
}

func NewRegistryServer(options RegistryOptions) *registryServer {
	return &registryServer{
		InstanceSettings: []byte{},
		options:          options,
	}
}

func (s *registryServer) Start(readyCh chan struct{}) error {
	var err error
	s.listener, err = net.Listen("tcp", fmt.Sprintf("%s:%d", s.options.Host, s.options.Port))
	if err != nil {
		return err
	}

	readyCh <- struct{}{}

	httpServer := http.Server{}
	mux := http.NewServeMux()
	httpServer.Handler = mux
	mux.HandleFunc("/instances/fake-instance-id/settings", s.instanceHandler)

	return httpServer.Serve(s.listener)
}

func (s *registryServer) Stop() error {
	// if client keeps connection alive, server will still be running
	s.InstanceSettings = nil

	err := s.listener.Close()
	if err != nil {
		return err
	}

	return nil
}

func (s *registryServer) instanceHandler(w http.ResponseWriter, req *http.Request) {
	if !s.isAuthorized(req) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if req.Method == "GET" {
		if s.InstanceSettings != nil {
			response := AgentSettingsResponse{
				Settings: string(s.InstanceSettings),
				Status:   "ok",
			}
			responseJSON, err := json.Marshal(response)
			if err != nil {
				http.Error(w, "Error marshalling response", http.StatusBadRequest)
				return
			}
			w.Write(responseJSON)
			return
		}
		http.NotFound(w, req)
		return
	}

	if req.Method == "PUT" {
		reqBody, _ := ioutil.ReadAll(req.Body)
		s.InstanceSettings = reqBody

		w.WriteHeader(http.StatusOK)
		return
	}

	if req.Method == "DELETE" {
		if s.InstanceSettings != nil {
			s.InstanceSettings = []byte{}
			w.WriteHeader(http.StatusOK)
			return
		}
		http.NotFound(w, req)
		return
	}
}

func (s *registryServer) isAuthorized(req *http.Request) bool {
	auth := s.options.Username + ":" + s.options.Password
	expectedAuthorizationHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))

	return expectedAuthorizationHeader == req.Header.Get("Authorization")
}
