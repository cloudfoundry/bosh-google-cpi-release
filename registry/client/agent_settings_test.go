package registry_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/frodenas/bosh-google-cpi/registry/client"
)

var _ = Describe("AgentSettings", func() {
	Describe("NewAgentSettingsForVM", func() {
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

			agentSettings := NewAgentSettingsForVM(
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
					Persistent: PersistentSettings{},
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
					ID: "fake-vm-id",
				},
			}

			Expect(agentSettings).To(Equal(expectedAgentSettings))
		})
	})

	Describe("AttachPersistentDisk", func() {
		It("sets persistent disk device name for given disk id on an empty agent settings", func() {
			agentSettings := AgentSettings{}

			newAgentSettings := agentSettings.AttachPersistentDisk("fake-disk-id", "fake-disk-device-name")

			Expect(newAgentSettings).To(Equal(AgentSettings{
				Disks: DisksSettings{
					Persistent: PersistentSettings{
						"fake-disk-id": "fake-disk-device-name",
					},
				},
			}))

			// keeps original agent settings not modified
			Expect(agentSettings).To(Equal(AgentSettings{}))
		})

		It("sets persistent disk device name for given disk id", func() {
			agentSettings := AgentSettings{
				Disks: DisksSettings{
					Persistent: PersistentSettings{
						"fake-other-disk-id": "fake-other-disk-device-name",
					},
				},
			}

			newAgentSettings := agentSettings.AttachPersistentDisk("fake-disk-id", "fake-disk-device-name")

			Expect(newAgentSettings).To(Equal(AgentSettings{
				Disks: DisksSettings{
					Persistent: PersistentSettings{
						"fake-other-disk-id": "fake-other-disk-device-name",
						"fake-disk-id":       "fake-disk-device-name",
					},
				},
			}))

			// keeps original agent settings not modified
			Expect(agentSettings).To(Equal(AgentSettings{
				Disks: DisksSettings{
					Persistent: PersistentSettings{
						"fake-other-disk-id": "fake-other-disk-device-name",
					},
				},
			}))
		})

		It("overwrites persistent disk device name for given disk id", func() {
			agentSettings := AgentSettings{
				Disks: DisksSettings{
					Persistent: PersistentSettings{
						"fake-disk-id": "fake-old-disk-device-name",
					},
				},
			}

			newAgentSettings := agentSettings.AttachPersistentDisk("fake-disk-id", "fake-new-disk-device-name")

			Expect(newAgentSettings).To(Equal(AgentSettings{
				Disks: DisksSettings{
					Persistent: PersistentSettings{
						"fake-disk-id": "fake-new-disk-device-name",
					},
				},
			}))

			// keeps original agent settings not modified
			Expect(agentSettings).To(Equal(AgentSettings{
				Disks: DisksSettings{
					Persistent: PersistentSettings{
						"fake-disk-id": "fake-old-disk-device-name",
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

			// keeps original agent settings not modified
			Expect(agentSettings).To(Equal(AgentSettings{}))
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

			// keeps original agent settings not modified
			Expect(agentSettings).To(Equal(AgentSettings{
				Networks: NetworksSettings{
					"fake-dynamic-name": NetworkSettings{
						Type: "dynamic",
						IP:   "1.2.3.4",
					},
				},
			}))
		})
	})

	Describe("DetachPersistentDisk", func() {
		It("unsets persistent disk device name on an empty agent settings", func() {
			agentSettings := AgentSettings{}

			newAgentSettings := agentSettings.DetachPersistentDisk("fake-disk-id")

			Expect(newAgentSettings).To(Equal(AgentSettings{
				Disks: DisksSettings{
					Persistent: PersistentSettings{},
				},
			}))

			// keeps original agent settings not modified
			Expect(agentSettings).To(Equal(AgentSettings{}))
		})

		It("unsets persistent disk device name if previously set", func() {
			agentSettings := AgentSettings{
				Disks: DisksSettings{
					Persistent: PersistentSettings{
						"fake-disk-id": "fake-disk-device-name",
					},
				},
			}

			newAgentSettings := agentSettings.DetachPersistentDisk("fake-disk-id")

			Expect(newAgentSettings).To(Equal(AgentSettings{
				Disks: DisksSettings{
					Persistent: PersistentSettings{},
				},
			}))

			// keeps original agent settings not modified
			Expect(agentSettings).To(Equal(AgentSettings{
				Disks: DisksSettings{
					Persistent: PersistentSettings{
						"fake-disk-id": "fake-disk-device-name",
					},
				},
			}))
		})

		It("does not change anything if persistent disk was not set", func() {
			agentSettings := AgentSettings{
				Disks: DisksSettings{
					Persistent: PersistentSettings{
						"fake-other-disk-id": "fake-other-disk-device-name",
					},
				},
			}

			newAgentSettings := agentSettings.DetachPersistentDisk("fake-disk-id")

			Expect(newAgentSettings).To(Equal(AgentSettings{
				Disks: DisksSettings{
					Persistent: PersistentSettings{
						"fake-other-disk-id": "fake-other-disk-device-name",
					},
				},
			}))

			// keeps original agent settings not modified
			Expect(agentSettings).To(Equal(AgentSettings{
				Disks: DisksSettings{
					Persistent: PersistentSettings{
						"fake-other-disk-id": "fake-other-disk-device-name",
					},
				},
			}))
		})
	})
})
