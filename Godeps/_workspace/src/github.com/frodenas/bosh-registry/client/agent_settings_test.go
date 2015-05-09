package registry_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/frodenas/bosh-registry/client"
)

var _ = Describe("AgentSettings", func() {
	Describe("NewAgentSettings", func() {
		It("returns agent settings", func() {
			networks := NetworksSettings{
				"fake-net-name": NetworkSettings{
					Type:    "fake-type",
					IP:      "fake-ip",
					Gateway: "fake-gateway",
					Netmask: "fake-netmask",
					DNS:     []string{"fake-dns"},
					Default: []string{"fake-default"},
					CloudProperties: map[string]interface{}{
						"fake-cp-key": "fake-cp-value",
					},
				},
			}

			env := EnvSettings{"fake-env-key": "fake-env-value"}

			agentOptions := AgentOptions{
				Mbus: "fake-mbus",
				Ntp:  []string{"fake-ntp"},
				Blobstore: BlobstoreOptions{
					Type: "fake-blobstore-type",
					Options: map[string]interface{}{
						"fake-blobstore-key": "fake-blobstore-value",
					},
				},
			}

			agentSettings := NewAgentSettings(
				"fake-agent-id",
				"fake-vm-id",
				networks,
				env,
				agentOptions,
			)

			expectedAgentSettings := AgentSettings{
				AgentID: "fake-agent-id",

				Blobstore: BlobstoreSettings{
					Provider: "fake-blobstore-type",
					Options: map[string]interface{}{
						"fake-blobstore-key": "fake-blobstore-value",
					},
				},

				Disks: DisksSettings{
					System:     "/dev/sda",
					Persistent: map[string]PersistentSettings{},
				},

				Env: EnvSettings{
					"fake-env-key": "fake-env-value",
				},

				Mbus: "fake-mbus",

				Networks: NetworksSettings{
					"fake-net-name": NetworkSettings{
						Type:    "fake-type",
						IP:      "fake-ip",
						Gateway: "fake-gateway",
						Netmask: "fake-netmask",
						DNS:     []string{"fake-dns"},
						Default: []string{"fake-default"},
						CloudProperties: map[string]interface{}{
							"fake-cp-key": "fake-cp-value",
						},
					},
				},

				Ntp: []string{"fake-ntp"},

				VM: VMSettings{
					Name: "fake-vm-id",
				},
			}

			Expect(agentSettings).To(Equal(expectedAgentSettings))
		})
	})

	Describe("AttachPersistentDisk", func() {
		It("sets persistent disk device name and path for given disk id on an empty agent settings", func() {
			agentSettings := AgentSettings{}

			newAgentSettings := agentSettings.AttachPersistentDisk("fake-disk-id", "fake-disk-device-name", "fake-disk-path")

			Expect(newAgentSettings).To(Equal(AgentSettings{
				Disks: DisksSettings{
					Persistent: map[string]PersistentSettings{
						"fake-disk-id": PersistentSettings{
							ID:       "fake-disk-id",
							VolumeID: "fake-disk-device-name",
							Path:     "fake-disk-path",
						},
					},
				},
			}))
		})

		It("sets persistent disk device name and path for given disk id", func() {
			agentSettings := AgentSettings{
				Disks: DisksSettings{
					Persistent: map[string]PersistentSettings{
						"fake-other-disk-id": PersistentSettings{
							ID:       "fake-other-disk-id",
							VolumeID: "fake-other-disk-device-name",
							Path:     "fake-other-disk-path",
						},
					},
				},
			}

			newAgentSettings := agentSettings.AttachPersistentDisk("fake-disk-id", "fake-disk-device-name", "fake-disk-path")

			Expect(newAgentSettings).To(Equal(AgentSettings{
				Disks: DisksSettings{
					Persistent: map[string]PersistentSettings{
						"fake-other-disk-id": PersistentSettings{
							ID:       "fake-other-disk-id",
							VolumeID: "fake-other-disk-device-name",
							Path:     "fake-other-disk-path",
						},
						"fake-disk-id": PersistentSettings{
							ID:       "fake-disk-id",
							VolumeID: "fake-disk-device-name",
							Path:     "fake-disk-path",
						},
					},
				},
			}))
		})

		It("overwrites persistent disk device name and path for given disk id", func() {
			agentSettings := AgentSettings{
				Disks: DisksSettings{
					Persistent: map[string]PersistentSettings{
						"fake-disk-id": PersistentSettings{
							ID:       "fake-disk-id",
							VolumeID: "fake-disk-device-name",
							Path:     "fake-disk-path",
						},
					},
				},
			}

			newAgentSettings := agentSettings.AttachPersistentDisk("fake-disk-id", "fake-new-disk-device-name", "fake-new-disk-path")

			Expect(newAgentSettings).To(Equal(AgentSettings{
				Disks: DisksSettings{
					Persistent: map[string]PersistentSettings{
						"fake-disk-id": PersistentSettings{
							ID:       "fake-disk-id",
							VolumeID: "fake-new-disk-device-name",
							Path:     "fake-new-disk-path",
						},
					},
				},
			}))
		})
	})

	Describe("ConfigureNetworks", func() {
		It("sets networks on an empty agent settings", func() {
			agentSettings := AgentSettings{}

			newAgentSettings := agentSettings.ConfigureNetworks(NetworksSettings{
				"fake-dynamic-name": NetworkSettings{
					Type: "dynamic",
					IP:   "1.2.3.4",
				},
			})

			Expect(newAgentSettings).To(Equal(AgentSettings{
				Networks: NetworksSettings{
					"fake-dynamic-name": NetworkSettings{
						Type: "dynamic",
						IP:   "1.2.3.4",
					},
				},
			}))
		})

		It("overwrites networks", func() {
			agentSettings := AgentSettings{
				Networks: NetworksSettings{
					"fake-dynamic-name": NetworkSettings{
						Type: "dynamic",
						IP:   "1.2.3.4",
					},
				},
			}

			newAgentSettings := agentSettings.ConfigureNetworks(NetworksSettings{
				"fake-vip-name": NetworkSettings{
					Type: "vip",
					IP:   "5.6.7.8",
				},
			})

			Expect(newAgentSettings).To(Equal(AgentSettings{
				Networks: NetworksSettings{
					"fake-vip-name": NetworkSettings{
						Type: "vip",
						IP:   "5.6.7.8",
					},
				},
			}))
		})
	})

	Describe("DetachPersistentDisk", func() {
		It("unsets persistent disk device name on an empty agent settings", func() {
			agentSettings := AgentSettings{}

			newAgentSettings := agentSettings.DetachPersistentDisk("fake-disk-id")

			Expect(newAgentSettings).To(Equal(AgentSettings{}))
		})

		It("unsets persistent disk device name if previously set", func() {
			agentSettings := AgentSettings{
				Disks: DisksSettings{
					Persistent: map[string]PersistentSettings{
						"fake-disk-id": PersistentSettings{
							ID:       "fake-disk-id",
							VolumeID: "fake-disk-device-name",
							Path:     "fake-disk-path",
						},
					},
				},
			}

			newAgentSettings := agentSettings.DetachPersistentDisk("fake-disk-id")

			Expect(newAgentSettings).To(Equal(AgentSettings{
				Disks: DisksSettings{
					Persistent: map[string]PersistentSettings{},
				},
			}))
		})

		It("does not change anything if persistent disk was not set", func() {
			agentSettings := AgentSettings{
				Disks: DisksSettings{
					Persistent: map[string]PersistentSettings{
						"fake-other-disk-id": PersistentSettings{
							ID:       "fake-other-disk-id",
							VolumeID: "fake-other-disk-device-name",
							Path:     "fake-other-disk-path",
						},
					},
				},
			}

			newAgentSettings := agentSettings.DetachPersistentDisk("fake-disk-id")

			Expect(newAgentSettings).To(Equal(AgentSettings{
				Disks: DisksSettings{
					Persistent: map[string]PersistentSettings{
						"fake-other-disk-id": PersistentSettings{
							ID:       "fake-other-disk-id",
							VolumeID: "fake-other-disk-device-name",
							Path:     "fake-other-disk-path",
						},
					},
				},
			}))
		})
	})
})
