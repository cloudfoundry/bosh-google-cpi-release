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

var _ = Describe("RegistryClient", func() {
	var (
		logger               boshlog.Logger
		registryClient       Client
		registryServer       *RegistryServer
		instanceID           string
		expectedAgentSet     AgentSettings
		expectedAgentSetJSON []byte
	)

	BeforeEach(func() {
		options := ClientOptions{
			Schema:   "http",
			Host:     "127.0.0.1",
			Port:     6307,
			Username: "fake-username",
			Password: "fake-password",
		}
		registryServer = NewRegistryServer(options)
		readyCh := make(chan struct{})
		go registryServer.Start(readyCh)
		<-readyCh

		instanceID = "fake-instance-id"
		logger = boshlog.NewLogger(boshlog.LevelNone)
		registryClient = NewClient(options, logger)

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
				err := registryClient.Delete(instanceID)
				Expect(err).ToNot(HaveOccurred())
				Expect(registryServer.InstanceSettings).To(Equal([]byte{}))
			})
		})

		Context("when settings for instance do not exist", func() {
			It("returns an error", func() {
				Expect(registryServer.InstanceSettings).To(Equal([]byte{}))
				err := registryClient.Delete(instanceID)
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
				agentSet, err := registryClient.Fetch(instanceID)
				Expect(err).ToNot(HaveOccurred())
				Expect(agentSet).To(Equal(expectedAgentSet))
			})
		})

		Context("when settings for instance do not exist", func() {
			It("returns an error", func() {
				Expect(registryServer.InstanceSettings).To(Equal([]byte{}))
				agentSet, err := registryClient.Fetch(instanceID)
				Expect(err).To(HaveOccurred())
				Expect(agentSet).To(Equal(AgentSettings{}))
			})
		})
	})

	Describe("Update", func() {
		It("updates settings in the registry", func() {
			Expect(registryServer.InstanceSettings).To(Equal([]byte{}))
			err := registryClient.Update(instanceID, expectedAgentSet)
			Expect(err).ToNot(HaveOccurred())
			Expect(registryServer.InstanceSettings).To(Equal(expectedAgentSetJSON))
		})
	})
})

type RegistryServer struct {
	InstanceSettings []byte
	options          ClientOptions
	listener         net.Listener
}

func NewRegistryServer(options ClientOptions) *RegistryServer {
	return &RegistryServer{
		InstanceSettings: []byte{},
		options:          options,
	}
}

func (s *RegistryServer) Start(readyCh chan struct{}) error {
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

func (s *RegistryServer) Stop() error {
	// if client keeps connection alive, server will still be running
	s.InstanceSettings = nil

	err := s.listener.Close()
	if err != nil {
		return err
	}

	return nil
}

func (s *RegistryServer) instanceHandler(w http.ResponseWriter, req *http.Request) {
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

func (s *RegistryServer) isAuthorized(req *http.Request) bool {
	auth := s.options.Username + ":" + s.options.Password
	expectedAuthorizationHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))

	return expectedAuthorizationHeader == req.Header.Get("Authorization")
}
