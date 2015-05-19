package registry_test

import (
	"encoding/json"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/frodenas/bosh-registry/client"

	boshlog "github.com/cloudfoundry/bosh-utils/logger"

	"github.com/frodenas/bosh-registry/server/fakes"
)

var _ = Describe("HTTPClient", func() {
	var (
		err error

		instanceHandler *fakes.FakeInstanceHandler
		mux             *http.ServeMux
		server          *httptest.Server

		options ClientOptions
		logger  boshlog.Logger
		client  HTTPClient

		instanceID           string
		expectedAgentSet     AgentSettings
		expectedAgentSetJSON []byte
	)

	BeforeEach(func() {
		logger = boshlog.NewLogger(boshlog.LevelNone)
		instanceHandler = fakes.NewFakeInstanceHandler("fake-username", "fake-password")
		mux = http.NewServeMux()
		mux.HandleFunc("/", instanceHandler.HandleFunc)

		instanceID = "fake-instance-id"
		expectedAgentSet = AgentSettings{AgentID: "fake-agent-id"}
		expectedAgentSetJSON, err = json.Marshal(expectedAgentSet)
		Expect(err).ToNot(HaveOccurred())
	})

	Context("when using http", func() {
		BeforeEach(func() {
			server = httptest.NewServer(mux)
			serverURL, err := url.Parse(server.URL)
			Expect(err).ToNot(HaveOccurred())
			serverHost, serverPortString, err := net.SplitHostPort(serverURL.Host)
			Expect(err).ToNot(HaveOccurred())
			serverPort, err := strconv.ParseInt(serverPortString, 10, 64)
			Expect(err).ToNot(HaveOccurred())

			options = ClientOptions{
				Protocol: serverURL.Scheme,
				Host:     serverHost,
				Port:     int(serverPort),
				Username: "fake-username",
				Password: "fake-password",
			}
			client = NewHTTPClient(options, logger)
		})

		AfterEach(func() {
			server.Close()
		})

		Describe("Delete", func() {
			Context("when settings for the instance exist in the registry", func() {
				BeforeEach(func() {
					instanceHandler.InstanceSettings = expectedAgentSetJSON
				})

				It("deletes settings in the registry", func() {
					err = client.Delete(instanceID)
					Expect(err).ToNot(HaveOccurred())
					Expect(instanceHandler.InstanceSettings).To(Equal([]byte{}))
				})
			})

			Context("when settings for instance does not exist", func() {
				It("should not return an error", func() {
					Expect(instanceHandler.InstanceSettings).To(Equal([]byte{}))
					err = client.Delete(instanceID)
					Expect(err).ToNot(HaveOccurred())
				})
			})
		})

		Describe("Fetch", func() {
			Context("when settings for the instance exist in the registry", func() {
				BeforeEach(func() {
					instanceHandler.InstanceSettings = expectedAgentSetJSON
				})

				It("fetches settings from the registry", func() {
					agentSet, err := client.Fetch(instanceID)
					Expect(err).ToNot(HaveOccurred())
					Expect(agentSet).To(Equal(expectedAgentSet))
				})
			})

			Context("when settings for instance does not exist", func() {
				It("returns an error", func() {
					Expect(instanceHandler.InstanceSettings).To(Equal([]byte{}))
					agentSet, err := client.Fetch(instanceID)
					Expect(err).To(HaveOccurred())
					Expect(agentSet).To(Equal(AgentSettings{}))
				})
			})
		})

		Describe("Update", func() {
			It("updates settings in the registry", func() {
				Expect(instanceHandler.InstanceSettings).To(Equal([]byte{}))
				err := client.Update(instanceID, expectedAgentSet)
				Expect(err).ToNot(HaveOccurred())
				Expect(instanceHandler.InstanceSettings).To(Equal(expectedAgentSetJSON))
			})
		})
	})

	Context("when using https", func() {
		BeforeEach(func() {
			server = httptest.NewTLSServer(mux)
			serverURL, err := url.Parse(server.URL)
			Expect(err).ToNot(HaveOccurred())
			serverHost, serverPortString, err := net.SplitHostPort(serverURL.Host)
			Expect(err).ToNot(HaveOccurred())
			serverPort, err := strconv.ParseInt(serverPortString, 10, 64)
			Expect(err).ToNot(HaveOccurred())

			options = ClientOptions{
				Protocol: serverURL.Scheme,
				Host:     serverHost,
				Port:     int(serverPort),
				Username: "fake-username",
				Password: "fake-password",
				TLS: ClientTLSOptions{
					InsecureSkipVerify: true,
					CertFile:           "../test/assets/public.pem",
					KeyFile:            "../test/assets/private.pem",
					CACertFile:         "../test/assets/ca.pem",
				},
			}
			client = NewHTTPClient(options, logger)
		})

		AfterEach(func() {
			server.Close()
		})

		Describe("Delete", func() {
			It("deletes settings in the registry", func() {
				instanceHandler.InstanceSettings = expectedAgentSetJSON
				err = client.Delete(instanceID)
				Expect(err).ToNot(HaveOccurred())
				Expect(instanceHandler.InstanceSettings).To(Equal([]byte{}))
			})
		})

		Describe("Fetch", func() {
			It("fetches settings from the registry", func() {
				instanceHandler.InstanceSettings = expectedAgentSetJSON
				agentSet, err := client.Fetch(instanceID)
				Expect(err).ToNot(HaveOccurred())
				Expect(agentSet).To(Equal(expectedAgentSet))
			})
		})

		Describe("Update", func() {
			It("updates settings in the registry", func() {
				Expect(instanceHandler.InstanceSettings).To(Equal([]byte{}))
				err := client.Update(instanceID, expectedAgentSet)
				Expect(err).ToNot(HaveOccurred())
				Expect(instanceHandler.InstanceSettings).To(Equal(expectedAgentSetJSON))
			})
		})
	})
})
