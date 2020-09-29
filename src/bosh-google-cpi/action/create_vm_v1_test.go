package action_test

import (
	. "bosh-google-cpi/action"
	"bosh-google-cpi/api"
	acceleratortypefakes "bosh-google-cpi/google/accelerator_type_service/fakes"
	"bosh-google-cpi/google/disk_service"
	diskfakes "bosh-google-cpi/google/disk_service/fakes"
	"bosh-google-cpi/google/disk_type_service"
	disktypefakes "bosh-google-cpi/google/disk_type_service/fakes"
	"bosh-google-cpi/google/image_service"
	imagefakes "bosh-google-cpi/google/image_service/fakes"
	"bosh-google-cpi/google/instance_service"
	instancefakes "bosh-google-cpi/google/instance_service/fakes"
	"bosh-google-cpi/google/machine_type_service"
	machinetypefakes "bosh-google-cpi/google/machine_type_service/fakes"
	"bosh-google-cpi/registry"
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	registryfakes "bosh-google-cpi/registry/fakes"

	"bosh-google-cpi/google/accelerator_type_service"
)

var _ = Describe("CreateVM", func() {
	var (
		err                      error
		vmCID                    interface{}
		networks                 Networks
		cloudProps               VMCloudProperties
		disks                    []DiskCID
		env                      Environment
		defaultRootDiskSizeGb    int
		defaultRootDiskType      string
		registryOptions          registry.ClientOptions
		agentOptions             registry.AgentOptions
		expectedVMProps          *instance.Properties
		expectedInstanceNetworks instance.Networks
		expectedAgentSettings    registry.AgentSettings

		vmService              *instancefakes.FakeInstanceService
		diskService            *diskfakes.FakeDiskService
		diskTypeService        *disktypefakes.FakeDiskTypeService
		machineTypeService     *machinetypefakes.FakeMachineTypeService
		imageService           *imagefakes.FakeImageService
		registryClient         *registryfakes.FakeClient
		acceleratorTypeService *acceleratortypefakes.FakeAcceleratorTypeService

		createVM CreateVMV1
	)

	BeforeEach(func() {
		vmService = &instancefakes.FakeInstanceService{}
		diskService = &diskfakes.FakeDiskService{}
		diskTypeService = &disktypefakes.FakeDiskTypeService{}
		machineTypeService = &machinetypefakes.FakeMachineTypeService{}
		acceleratorTypeService = &acceleratortypefakes.FakeAcceleratorTypeService{}
		imageService = &imagefakes.FakeImageService{}
		registryClient = &registryfakes.FakeClient{}
		registryOptions = registry.ClientOptions{
			Protocol: "http",
			Host:     "fake-registry-host",
			Port:     25777,
			Username: "fake-registry-username",
			Password: "fake-registry-password",
		}
		agentOptions = registry.AgentOptions{
			Mbus: "http://fake-mbus",
			Blobstore: registry.BlobstoreOptions{
				Provider: "fake-blobstore-type",
			},
		}
		defaultRootDiskSizeGb = 0
		defaultRootDiskType = ""
	})

	JustBeforeEach(func() {
		createVM = NewCreateVMV1(
			vmService,
			diskService,
			diskTypeService,
			imageService,
			machineTypeService,
			acceleratorTypeService,
			registryClient,
			registryOptions,
			agentOptions,
			defaultRootDiskSizeGb,
			defaultRootDiskType,
		)
	})

	Describe("Run", func() {
		BeforeEach(func() {
			vmService.CreateID = "fake-vm-id"
			imageService.FindFound = true
			machineTypeService.FindFound = true

			diskService.FindDisk = disk.Disk{Zone: "fake-default-zone"}
			imageService.FindImage = image.Image{SelfLink: "fake-image-self-link"}
			machineTypeService.FindMachineType = machinetype.MachineType{SelfLink: "fake-machine-type-self-link"}
			diskTypeService.FindDiskType = disktype.DiskType{SelfLink: "fake-disk-type-self-link"}

			cloudProps = VMCloudProperties{
				Zone:              "fake-default-zone",
				MachineType:       "fake-machine-type",
				RootDiskSizeGb:    0,
				RootDiskType:      "",
				AutomaticRestart:  true,
				OnHostMaintenance: "TERMINATE",
				Preemptible:       true,
				ServiceScopes:     []string{},
				BackendService:    "fake-backend-service",
			}

			networks = Networks{
				"fake-network-name": &Network{
					Type:    "manual",
					IP:      "fake-network-ip",
					Gateway: "fake-network-gateway",
					Netmask: "fake-network-netmask",
					DNS:     []string{"fake-network-dns"},
					DHCP:    true,
					Default: []string{"fake-network-default"},
					CloudProperties: NetworkCloudProperties{
						NetworkName:         "fake-network-cloud-network-name",
						Tags:                instance.Tags([]string{"fake-network-cloud-network-tag"}),
						EphemeralExternalIP: true,
						IPForwarding:        false,
					},
				},
			}

			expectedVMProps = &instance.Properties{
				Zone:              "fake-default-zone",
				Stemcell:          "fake-image-self-link",
				MachineType:       "fake-machine-type-self-link",
				RootDiskSizeGb:    0,
				RootDiskType:      "",
				AutomaticRestart:  true,
				OnHostMaintenance: "TERMINATE",
				Preemptible:       true,
				ServiceScopes:     []string{},
				BackendService:    instance.BackendService{Name: "fake-backend-service"},
			}

			expectedInstanceNetworks = networks.AsInstanceServiceNetworks()

			expectedAgentSettings = registry.AgentSettings{
				AgentID: "fake-agent-id",
				Blobstore: registry.BlobstoreSettings{
					Provider: "fake-blobstore-type",
				},
				Disks: registry.DisksSettings{
					System:     "/dev/sda",
					Persistent: map[string]registry.PersistentSettings{},
				},
				Mbus: "http://fake-mbus",
				Networks: registry.NetworksSettings{
					"fake-network-name": registry.NetworkSettings{
						Type:    "manual",
						IP:      "fake-network-ip",
						Gateway: "fake-network-gateway",
						Netmask: "fake-network-netmask",
						DNS:     []string{"fake-network-dns"},
						DHCP:    true,
						Default: []string{"fake-network-default"},
					},
				},
				VM: registry.VMSettings{
					Name: "fake-vm-id",
				},
			}
		})
		Context("when BackendService is set as a hash", func() {
			var backendService map[string]string
			BeforeEach(func() {
				backendService = make(map[string]string)
				backendService["name"] = "foobar"
				expectedVMProps.BackendService = instance.BackendService{Name: "foobar"}

				cloudProps.BackendService = backendService
			})

			It("defaults to external", func() {
				vmCID, err = createVM.Run("fake-agent-id", "fake-stemcell-id", cloudProps, networks, disks, env)
				Expect(err).NotTo(HaveOccurred())
				Expect(vmService.CreateVMProps).To(Equal(expectedVMProps))
			})

			It("supports internal", func() {
				backendService["scheme"] = "INTERNAL"
				expectedVMProps.BackendService = instance.BackendService{Name: "foobar"}

				vmCID, err = createVM.Run("fake-agent-id", "fake-stemcell-id", cloudProps, networks, disks, env)
				Expect(err).NotTo(HaveOccurred())
				Expect(vmService.CreateVMProps).To(Equal(expectedVMProps))
			})

			It("requires name", func() {
				delete(backendService, "name")
				vmCID, err = createVM.Run("fake-agent-id", "fake-stemcell-id", cloudProps, networks, disks, env)
				Expect(err).To(HaveOccurred())
			})
		})

		It("creates the vm", func() {
			vmCID, err = createVM.Run("fake-agent-id", "fake-stemcell-id", cloudProps, networks, disks, env)
			Expect(err).NotTo(HaveOccurred())
			Expect(diskService.FindCalled).To(BeFalse())
			Expect(imageService.FindCalled).To(BeTrue())
			Expect(machineTypeService.FindCalled).To(BeTrue())
			Expect(machineTypeService.CustomLinkCalled).To(BeFalse())
			Expect(acceleratorTypeService.FindCalled).To(BeFalse())
			Expect(diskTypeService.FindCalled).To(BeFalse())
			Expect(vmService.CreateCalled).To(BeTrue())
			Expect(vmService.CleanUpCalled).To(BeFalse())
			Expect(registryClient.UpdateCalled).To(BeTrue())
			Expect(registryClient.UpdateSettings).To(Equal(expectedAgentSettings))
			Expect(vmCID).To(Equal(VMCID("fake-vm-id")))
			Expect(vmService.CreateVMProps).To(Equal(expectedVMProps))
			Expect(vmService.CreateNetworks).To(Equal(expectedInstanceNetworks))
			Expect(vmService.CreateRegistryEndpoint).To(Equal("http://fake-registry-username:fake-registry-password@fake-registry-host:25777"))
		})

		It("returns an error if imageService find call returns an error", func() {
			imageService.FindErr = errors.New("fake-image-service-error")

			_, err = createVM.Run("fake-agent-id", "fake-stemcell-id", cloudProps, networks, disks, env)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("fake-image-service-error"))
			Expect(diskService.FindCalled).To(BeFalse())
			Expect(imageService.FindCalled).To(BeTrue())
			Expect(machineTypeService.FindCalled).To(BeFalse())
			Expect(diskTypeService.FindCalled).To(BeFalse())
			Expect(vmService.CreateCalled).To(BeFalse())
			Expect(vmService.CleanUpCalled).To(BeFalse())
			Expect(registryClient.UpdateCalled).To(BeFalse())
		})

		It("uses the gcp image url directly if provided", func() {
			stemcellLink := "https://www.googleapis.com/compute/v1/projects/fake/stemcell/path"
			_, err = createVM.Run("fake-agent-id", StemcellCID(stemcellLink), cloudProps, networks, disks, env)
			Expect(err).ToNot(HaveOccurred())
			Expect(imageService.FindCalled).To(BeFalse())
			Expect(vmService.CreateVMProps.Stemcell).To(Equal(stemcellLink))
		})

		It("returns an error if stemcell is not found", func() {
			imageService.FindFound = false

			_, err = createVM.Run("fake-agent-id", "fake-stemcell-id", cloudProps, networks, disks, env)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Stemcell 'fake-stemcell-id' does not exists"))
			Expect(diskService.FindCalled).To(BeFalse())
			Expect(imageService.FindCalled).To(BeTrue())
			Expect(machineTypeService.FindCalled).To(BeFalse())
			Expect(diskTypeService.FindCalled).To(BeFalse())
			Expect(vmService.CreateCalled).To(BeFalse())
			Expect(vmService.CleanUpCalled).To(BeFalse())
			Expect(registryClient.UpdateCalled).To(BeFalse())
		})

		It("returns an error if machine type and cpu are set", func() {
			cloudProps.CPU = 2

			_, err = createVM.Run("fake-agent-id", "fake-stemcell-id", cloudProps, networks, disks, env)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("'machine_type' and 'cpu' or 'ram' cannot be provided together"))
			Expect(diskService.FindCalled).To(BeFalse())
			Expect(imageService.FindCalled).To(BeTrue())
			Expect(machineTypeService.FindCalled).To(BeFalse())
			Expect(diskTypeService.FindCalled).To(BeFalse())
			Expect(vmService.CreateCalled).To(BeFalse())
			Expect(vmService.CleanUpCalled).To(BeFalse())
			Expect(registryClient.UpdateCalled).To(BeFalse())
		})

		It("returns an error if machine type and ram are set", func() {
			cloudProps.RAM = 5120

			_, err = createVM.Run("fake-agent-id", "fake-stemcell-id", cloudProps, networks, disks, env)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("'machine_type' and 'cpu' or 'ram' cannot be provided together"))
			Expect(diskService.FindCalled).To(BeFalse())
			Expect(imageService.FindCalled).To(BeTrue())
			Expect(machineTypeService.FindCalled).To(BeFalse())
			Expect(diskTypeService.FindCalled).To(BeFalse())
			Expect(vmService.CreateCalled).To(BeFalse())
			Expect(vmService.CleanUpCalled).To(BeFalse())
			Expect(registryClient.UpdateCalled).To(BeFalse())
		})

		It("returns an error if machineTypeService find call returns an error", func() {
			machineTypeService.FindErr = errors.New("fake-machine-type-service-error")

			_, err = createVM.Run("fake-agent-id", "fake-stemcell-id", cloudProps, networks, disks, env)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("fake-machine-type-service-error"))
			Expect(diskService.FindCalled).To(BeFalse())
			Expect(imageService.FindCalled).To(BeTrue())
			Expect(machineTypeService.FindCalled).To(BeTrue())
			Expect(diskTypeService.FindCalled).To(BeFalse())
			Expect(vmService.CreateCalled).To(BeFalse())
			Expect(vmService.CleanUpCalled).To(BeFalse())
			Expect(registryClient.UpdateCalled).To(BeFalse())
		})

		It("returns an error if machine type is not found", func() {
			machineTypeService.FindFound = false

			_, err = createVM.Run("fake-agent-id", "fake-stemcell-id", cloudProps, networks, disks, env)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Machine Type 'fake-machine-type' does not exists"))
			Expect(diskService.FindCalled).To(BeFalse())
			Expect(imageService.FindCalled).To(BeTrue())
			Expect(machineTypeService.FindCalled).To(BeTrue())
			Expect(diskTypeService.FindCalled).To(BeFalse())
			Expect(vmService.CreateCalled).To(BeFalse())
			Expect(vmService.CleanUpCalled).To(BeFalse())
			Expect(registryClient.UpdateCalled).To(BeFalse())
		})

		Context("when VM props override network props", func() {
			BeforeEach(func() {
				var t bool = true
				cloudProps.EphemeralExternalIP = &t
				cloudProps.IPForwarding = &t

				networks["fake-network-name"].CloudProperties.IPForwarding = false
				networks["fake-network-name"].CloudProperties.EphemeralExternalIP = false
				expectedInstanceNetworks.Network().IPForwarding = true
				expectedInstanceNetworks.Network().EphemeralExternalIP = true
			})

			It("creates the vm with the right properties", func() {
				vmCID, err = createVM.Run("fake-agent-id", "fake-stemcell-id", cloudProps, networks, disks, env)
				Expect(err).NotTo(HaveOccurred())
				Expect(vmService.CreateNetworks).To(Equal(expectedInstanceNetworks))
			})
		})

		Context("when custom service account and service scopes are provided", func() {
			BeforeEach(func() {
				cloudProps.ServiceAccount = "fake-service-account"
				cloudProps.ServiceScopes = []string{"fake-service-scope"}

				expectedVMProps.ServiceScopes = instance.ServiceScopes([]string{"fake-service-scope"})
				expectedVMProps.ServiceAccount = instance.ServiceAccount("fake-service-account")
			})

			It("creates the vm with a default service account and single scope", func() {
				vmCID, err = createVM.Run("fake-agent-id", "fake-stemcell-id", cloudProps, networks, disks, env)
				Expect(err).NotTo(HaveOccurred())
				Expect(vmService.CreateVMProps).To(Equal(expectedVMProps))
			})
		})

		Context("when custom machine type is set", func() {
			BeforeEach(func() {
				cloudProps.MachineType = ""
				cloudProps.CPU = 2
				cloudProps.RAM = 5120

				machineTypeService.CustomLinkLink = "custom-machine-type-link"
				expectedVMProps.MachineType = "custom-machine-type-link"
			})

			It("creates the vm with the right properties", func() {
				vmCID, err = createVM.Run("fake-agent-id", "fake-stemcell-id", cloudProps, networks, disks, env)
				Expect(err).NotTo(HaveOccurred())
				Expect(diskService.FindCalled).To(BeFalse())
				Expect(imageService.FindCalled).To(BeTrue())
				Expect(machineTypeService.FindCalled).To(BeFalse())
				Expect(machineTypeService.CustomLinkCalled).To(BeTrue())
				Expect(diskTypeService.FindCalled).To(BeFalse())
				Expect(vmService.CreateCalled).To(BeTrue())
				Expect(vmService.CleanUpCalled).To(BeFalse())
				Expect(registryClient.UpdateCalled).To(BeTrue())
				Expect(registryClient.UpdateSettings).To(Equal(expectedAgentSettings))
				Expect(vmCID).To(Equal(VMCID("fake-vm-id")))
				Expect(vmService.CreateVMProps).To(Equal(expectedVMProps))
				Expect(vmService.CreateNetworks).To(Equal(expectedInstanceNetworks))
				Expect(vmService.CreateRegistryEndpoint).To(Equal("http://fake-registry-username:fake-registry-password@fake-registry-host:25777"))
			})

			It("returns an error if cpu is not set", func() {
				cloudProps.CPU = 0

				_, err = createVM.Run("fake-agent-id", "fake-stemcell-id", cloudProps, networks, disks, env)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("'machine_type' or 'cpu' and 'ram' must be provided"))
				Expect(diskService.FindCalled).To(BeFalse())
				Expect(imageService.FindCalled).To(BeTrue())
				Expect(machineTypeService.FindCalled).To(BeFalse())
				Expect(diskTypeService.FindCalled).To(BeFalse())
				Expect(vmService.CreateCalled).To(BeFalse())
				Expect(vmService.CleanUpCalled).To(BeFalse())
				Expect(registryClient.UpdateCalled).To(BeFalse())
			})

			It("returns an error if ram is not set", func() {
				cloudProps.RAM = 0

				_, err = createVM.Run("fake-agent-id", "fake-stemcell-id", cloudProps, networks, disks, env)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("'machine_type' or 'cpu' and 'ram' must be provided"))
				Expect(diskService.FindCalled).To(BeFalse())
				Expect(imageService.FindCalled).To(BeTrue())
				Expect(machineTypeService.FindCalled).To(BeFalse())
				Expect(diskTypeService.FindCalled).To(BeFalse())
				Expect(vmService.CreateCalled).To(BeFalse())
				Expect(vmService.CleanUpCalled).To(BeFalse())
				Expect(registryClient.UpdateCalled).To(BeFalse())
			})
		})

		It("returns an error if vmService create call returns an error", func() {
			vmService.CreateErr = errors.New("fake-vm-service-error")

			_, err = createVM.Run("fake-agent-id", "fake-stemcell-id", cloudProps, networks, disks, env)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("fake-vm-service-error"))
			Expect(diskService.FindCalled).To(BeFalse())
			Expect(imageService.FindCalled).To(BeTrue())
			Expect(machineTypeService.FindCalled).To(BeTrue())
			Expect(diskTypeService.FindCalled).To(BeFalse())
			Expect(vmService.CreateCalled).To(BeTrue())
			Expect(vmService.CleanUpCalled).To(BeFalse())
			Expect(registryClient.UpdateCalled).To(BeFalse())
		})

		It("returns an error if registryClient update call returns an error", func() {
			registryClient.UpdateErr = errors.New("fake-registry-client-error")

			_, err = createVM.Run("fake-agent-id", "fake-stemcell-id", cloudProps, networks, disks, env)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("fake-registry-client-error"))
			Expect(diskService.FindCalled).To(BeFalse())
			Expect(imageService.FindCalled).To(BeTrue())
			Expect(machineTypeService.FindCalled).To(BeTrue())
			Expect(diskTypeService.FindCalled).To(BeFalse())
			Expect(vmService.CreateCalled).To(BeTrue())
			Expect(vmService.CleanUpCalled).To(BeTrue())
			Expect(registryClient.UpdateCalled).To(BeTrue())
		})

		Context("when default root disk size is set", func() {
			BeforeEach(func() {
				defaultRootDiskSizeGb = 20
				expectedVMProps.RootDiskSizeGb = 20
			})

			It("creates the vm with the right properties", func() {
				vmCID, err = createVM.Run("fake-agent-id", "fake-stemcell-id", cloudProps, networks, disks, env)
				Expect(err).NotTo(HaveOccurred())
				Expect(diskService.FindCalled).To(BeFalse())
				Expect(imageService.FindCalled).To(BeTrue())
				Expect(machineTypeService.FindCalled).To(BeTrue())
				Expect(diskTypeService.FindCalled).To(BeFalse())
				Expect(vmService.CreateCalled).To(BeTrue())
				Expect(vmService.CleanUpCalled).To(BeFalse())
				Expect(registryClient.UpdateCalled).To(BeTrue())
				Expect(registryClient.UpdateSettings).To(Equal(expectedAgentSettings))
				Expect(vmCID).To(Equal(VMCID("fake-vm-id")))
				Expect(vmService.CreateVMProps).To(Equal(expectedVMProps))
				Expect(vmService.CreateNetworks).To(Equal(expectedInstanceNetworks))
				Expect(vmService.CreateRegistryEndpoint).To(Equal("http://fake-registry-username:fake-registry-password@fake-registry-host:25777"))
			})
		})

		Context("when cloud properties root disk size is set", func() {
			BeforeEach(func() {
				cloudProps.RootDiskSizeGb = 20
				expectedVMProps.RootDiskSizeGb = 20
			})

			It("creates the vm with the right properties", func() {
				vmCID, err = createVM.Run("fake-agent-id", "fake-stemcell-id", cloudProps, networks, disks, env)
				Expect(err).NotTo(HaveOccurred())
				Expect(diskService.FindCalled).To(BeFalse())
				Expect(imageService.FindCalled).To(BeTrue())
				Expect(machineTypeService.FindCalled).To(BeTrue())
				Expect(diskTypeService.FindCalled).To(BeFalse())
				Expect(vmService.CreateCalled).To(BeTrue())
				Expect(vmService.CleanUpCalled).To(BeFalse())
				Expect(registryClient.UpdateCalled).To(BeTrue())
				Expect(registryClient.UpdateSettings).To(Equal(expectedAgentSettings))
				Expect(vmCID).To(Equal(VMCID("fake-vm-id")))
				Expect(vmService.CreateVMProps).To(Equal(expectedVMProps))
				Expect(vmService.CreateNetworks).To(Equal(expectedInstanceNetworks))
				Expect(vmService.CreateRegistryEndpoint).To(Equal("http://fake-registry-username:fake-registry-password@fake-registry-host:25777"))
			})
		})

		Context("when default root disk type is set", func() {
			BeforeEach(func() {
				diskTypeService.FindFound = true
				defaultRootDiskType = "fake-default-root-disk-type"
				diskTypeService.FindDiskType = disktype.DiskType{SelfLink: "fake-default-root-disk-type-self-link"}
				expectedVMProps.RootDiskType = "fake-default-root-disk-type-self-link"
			})

			It("creates the vm with the right properties", func() {
				vmCID, err = createVM.Run("fake-agent-id", "fake-stemcell-id", cloudProps, networks, disks, env)
				Expect(err).NotTo(HaveOccurred())
				Expect(diskService.FindCalled).To(BeFalse())
				Expect(imageService.FindCalled).To(BeTrue())
				Expect(machineTypeService.FindCalled).To(BeTrue())
				Expect(diskTypeService.FindCalled).To(BeTrue())
				Expect(vmService.CreateCalled).To(BeTrue())
				Expect(vmService.CleanUpCalled).To(BeFalse())
				Expect(registryClient.UpdateCalled).To(BeTrue())
				Expect(registryClient.UpdateSettings).To(Equal(expectedAgentSettings))
				Expect(vmCID).To(Equal(VMCID("fake-vm-id")))
				Expect(vmService.CreateVMProps).To(Equal(expectedVMProps))
				Expect(vmService.CreateNetworks).To(Equal(expectedInstanceNetworks))
				Expect(vmService.CreateRegistryEndpoint).To(Equal("http://fake-registry-username:fake-registry-password@fake-registry-host:25777"))
			})
		})

		Context("when cloud properties root disk type is set", func() {
			BeforeEach(func() {
				diskTypeService.FindFound = true
				cloudProps.RootDiskType = "fake-root-disk-type"
				expectedVMProps.RootDiskType = "fake-disk-type-self-link"
			})

			It("creates the vm with the right properties", func() {
				vmCID, err = createVM.Run("fake-agent-id", "fake-stemcell-id", cloudProps, networks, disks, env)
				Expect(err).NotTo(HaveOccurred())
				Expect(diskService.FindCalled).To(BeFalse())
				Expect(imageService.FindCalled).To(BeTrue())
				Expect(machineTypeService.FindCalled).To(BeTrue())
				Expect(diskTypeService.FindCalled).To(BeTrue())
				Expect(vmService.CreateCalled).To(BeTrue())
				Expect(vmService.CleanUpCalled).To(BeFalse())
				Expect(registryClient.UpdateCalled).To(BeTrue())
				Expect(registryClient.UpdateSettings).To(Equal(expectedAgentSettings))
				Expect(vmCID).To(Equal(VMCID("fake-vm-id")))
				Expect(vmService.CreateVMProps).To(Equal(expectedVMProps))
				Expect(vmService.CreateNetworks).To(Equal(expectedInstanceNetworks))
				Expect(vmService.CreateRegistryEndpoint).To(Equal("http://fake-registry-username:fake-registry-password@fake-registry-host:25777"))
			})

			It("returns an error if diskTypeService find call returns an error", func() {
				diskTypeService.FindErr = errors.New("fake-disk-type-service-error")

				_, err = createVM.Run("fake-agent-id", "fake-stemcell-id", cloudProps, networks, disks, env)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("fake-disk-type-service-error"))
				Expect(diskService.FindCalled).To(BeFalse())
				Expect(imageService.FindCalled).To(BeTrue())
				Expect(machineTypeService.FindCalled).To(BeTrue())
				Expect(diskTypeService.FindCalled).To(BeTrue())
				Expect(vmService.CreateCalled).To(BeFalse())
				Expect(vmService.CleanUpCalled).To(BeFalse())
				Expect(registryClient.UpdateCalled).To(BeFalse())
			})

			It("returns an error if root disk type is not found", func() {
				diskTypeService.FindFound = false

				_, err = createVM.Run("fake-agent-id", "fake-stemcell-id", cloudProps, networks, disks, env)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Root Disk Type 'fake-root-disk-type' does not exists"))
				Expect(diskService.FindCalled).To(BeFalse())
				Expect(imageService.FindCalled).To(BeTrue())
				Expect(machineTypeService.FindCalled).To(BeTrue())
				Expect(diskTypeService.FindCalled).To(BeTrue())
				Expect(vmService.CreateCalled).To(BeFalse())
				Expect(vmService.CleanUpCalled).To(BeFalse())
				Expect(registryClient.UpdateCalled).To(BeFalse())
			})
		})

		Context("when zone is set", func() {
			BeforeEach(func() {
				cloudProps.Zone = "fake-zone"
				expectedVMProps.Zone = "fake-zone"
			})

			It("creates the vm at the right zone", func() {
				vmCID, err = createVM.Run("fake-agent-id", "fake-stemcell-id", cloudProps, networks, disks, env)
				Expect(err).NotTo(HaveOccurred())
				Expect(diskService.FindCalled).To(BeFalse())
				Expect(imageService.FindCalled).To(BeTrue())
				Expect(machineTypeService.FindCalled).To(BeTrue())
				Expect(diskTypeService.FindCalled).To(BeFalse())
				Expect(vmService.CreateCalled).To(BeTrue())
				Expect(vmService.CleanUpCalled).To(BeFalse())
				Expect(registryClient.UpdateCalled).To(BeTrue())
				Expect(registryClient.UpdateSettings).To(Equal(expectedAgentSettings))
				Expect(vmCID).To(Equal(VMCID("fake-vm-id")))
				Expect(vmService.CreateVMProps).To(Equal(expectedVMProps))
				Expect(vmService.CreateNetworks).To(Equal(expectedInstanceNetworks))
				Expect(vmService.CreateRegistryEndpoint).To(Equal("http://fake-registry-username:fake-registry-password@fake-registry-host:25777"))
			})
		})

		Context("when local-ssd has been set", func() {
			BeforeEach(func() {
				cloudProps.EphemeralDiskType = "local-ssd"

				expectedVMProps.EphemeralDiskType = "local-ssd"
				expectedAgentSettings.Disks.Ephemeral = "/dev/nvme0n1"
			})

			It("creates the vm with the right properties", func() {
				vmCID, err = createVM.Run("fake-agent-id", "fake-stemcell-id", cloudProps, networks, disks, env)
				Expect(err).NotTo(HaveOccurred())
				Expect(diskService.FindCalled).To(BeFalse())
				Expect(imageService.FindCalled).To(BeTrue())
				Expect(machineTypeService.FindCalled).To(BeTrue())
				Expect(machineTypeService.CustomLinkCalled).To(BeFalse())
				Expect(diskTypeService.FindCalled).To(BeFalse())
				Expect(vmService.CreateCalled).To(BeTrue())
				Expect(vmService.CleanUpCalled).To(BeFalse())
				Expect(registryClient.UpdateCalled).To(BeTrue())
				Expect(registryClient.UpdateSettings).To(Equal(expectedAgentSettings))
				Expect(vmCID).To(Equal(VMCID("fake-vm-id")))
				Expect(vmService.CreateVMProps).To(Equal(expectedVMProps))
				Expect(vmService.CreateNetworks).To(Equal(expectedInstanceNetworks))
				Expect(vmService.CreateRegistryEndpoint).To(Equal("http://fake-registry-username:fake-registry-password@fake-registry-host:25777"))
			})
		})

		Context("when accelerator is set", func() {
			BeforeEach(func() {
				acceleratorTypeService.FindFound = true
				acceleratorTypeService.FindAcceleratorType = acceleratortype.AcceleratorType{SelfLink: "fake-accelerator-type-self-link"}
				acc := Accelerator{
					AcceleratorType: "fake-accelerator-type",
					Count:           1,
				}
				expectedAcc := instance.Accelerator{
					AcceleratorType: "fake-accelerator-type-self-link",
					Count:           1,
				}
				cloudProps.Accelerators = []Accelerator{acc}
				expectedVMProps.Accelerators = []instance.Accelerator{expectedAcc}
				expectedVMProps.OnHostMaintenance = "TERMINATE"
			})

			It("creates the vm with the accelerator", func() {
				_, err = createVM.Run("fake-agent-id", "fake-stemcell-id", cloudProps, networks, disks, env)
				Expect(err).NotTo(HaveOccurred())
				Expect(acceleratorTypeService.FindCalled).To(BeTrue())
				Expect(vmService.CreateCalled).To(BeTrue())
				Expect(vmService.CreateVMProps).To(Equal(expectedVMProps))

			})

			It("returns an error if acceleratorTypeService find call returns an error", func() {
				acceleratorTypeService.FindErr = errors.New("fake-accelerator-type-service-error")

				_, err = createVM.Run("fake-agent-id", "fake-stemcell-id", cloudProps, networks, disks, env)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("fake-accelerator-type-service-error"))
				Expect(diskService.FindCalled).To(BeFalse())
				Expect(imageService.FindCalled).To(BeTrue())
				Expect(machineTypeService.FindCalled).To(BeTrue())
				Expect(acceleratorTypeService.FindCalled).To(BeTrue())
				Expect(diskTypeService.FindCalled).To(BeFalse())
				Expect(vmService.CreateCalled).To(BeFalse())
				Expect(vmService.CleanUpCalled).To(BeFalse())
				Expect(registryClient.UpdateCalled).To(BeFalse())
			})

			It("returns an error if accelerator type is not found", func() {
				acceleratorTypeService.FindFound = false

				_, err = createVM.Run("fake-agent-id", "fake-stemcell-id", cloudProps, networks, disks, env)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Accelerator Type 'fake-accelerator-type' does not exists"))
				Expect(diskService.FindCalled).To(BeFalse())
				Expect(imageService.FindCalled).To(BeTrue())
				Expect(machineTypeService.FindCalled).To(BeTrue())
				Expect(acceleratorTypeService.FindCalled).To(BeTrue())
				Expect(diskTypeService.FindCalled).To(BeFalse())
				Expect(vmService.CreateCalled).To(BeFalse())
				Expect(vmService.CleanUpCalled).To(BeFalse())
				Expect(registryClient.UpdateCalled).To(BeFalse())
			})
		})

		Context("when DiskCIDs is set", func() {
			BeforeEach(func() {
				diskService.FindFound = true
				disks = []DiskCID{"fake-disk-1"}
				expectedVMProps.Zone = "fake-default-zone"
			})

			It("creates the vm at the right zone", func() {
				vmCID, err = createVM.Run("fake-agent-id", "fake-stemcell-id", cloudProps, networks, disks, env)
				Expect(err).NotTo(HaveOccurred())
				Expect(diskService.FindCalled).To(BeTrue())
				Expect(imageService.FindCalled).To(BeTrue())
				Expect(machineTypeService.FindCalled).To(BeTrue())
				Expect(diskTypeService.FindCalled).To(BeFalse())
				Expect(vmService.CreateCalled).To(BeTrue())
				Expect(vmService.CleanUpCalled).To(BeFalse())
				Expect(registryClient.UpdateCalled).To(BeTrue())
				Expect(registryClient.UpdateSettings).To(Equal(expectedAgentSettings))
				Expect(vmCID).To(Equal(VMCID("fake-vm-id")))
				Expect(vmService.CreateVMProps).To(Equal(expectedVMProps))
				Expect(vmService.CreateNetworks).To(Equal(expectedInstanceNetworks))
				Expect(vmService.CreateRegistryEndpoint).To(Equal("http://fake-registry-username:fake-registry-password@fake-registry-host:25777"))
			})

			It("returns an error if diskService find call returns an error", func() {
				diskService.FindErr = errors.New("fake-disk-service-error")

				_, err = createVM.Run("fake-agent-id", "fake-stemcell-id", cloudProps, networks, disks, env)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("fake-disk-service-error"))
				Expect(diskService.FindCalled).To(BeTrue())
				Expect(imageService.FindCalled).To(BeFalse())
				Expect(machineTypeService.FindCalled).To(BeFalse())
				Expect(diskTypeService.FindCalled).To(BeFalse())
				Expect(vmService.CreateCalled).To(BeFalse())
				Expect(vmService.CleanUpCalled).To(BeFalse())
				Expect(registryClient.UpdateCalled).To(BeFalse())
			})

			It("returns an error if disk is not found", func() {
				diskService.FindFound = false

				_, err = createVM.Run("fake-agent-id", "fake-stemcell-id", cloudProps, networks, disks, env)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal(api.NewDiskNotFoundError("fake-disk-1", false).Error()))
				Expect(diskService.FindCalled).To(BeTrue())
				Expect(imageService.FindCalled).To(BeFalse())
				Expect(machineTypeService.FindCalled).To(BeFalse())
				Expect(diskTypeService.FindCalled).To(BeFalse())
				Expect(vmService.CreateCalled).To(BeFalse())
				Expect(vmService.CleanUpCalled).To(BeFalse())
				Expect(registryClient.UpdateCalled).To(BeFalse())
			})

			Context("and zone is set", func() {
				BeforeEach(func() {
					cloudProps.Zone = "fake-zone"
				})

				It("returns an error if zone and disk are different", func() {
					_, err = createVM.Run("fake-agent-id", "fake-stemcell-id", cloudProps, networks, disks, env)
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("can't use multiple zones:"))
					Expect(diskService.FindCalled).To(BeTrue())
					Expect(imageService.FindCalled).To(BeFalse())
					Expect(machineTypeService.FindCalled).To(BeFalse())
					Expect(diskTypeService.FindCalled).To(BeFalse())
					Expect(vmService.CreateCalled).To(BeFalse())
					Expect(vmService.CleanUpCalled).To(BeFalse())
					Expect(registryClient.UpdateCalled).To(BeFalse())
				})
			})
		})

		Context("when network is of type 'dynamic'", func() {
			BeforeEach(func() {
				networks["fake-network-name"].Type = "dynamic"
				expectedInstanceNetworks.Network().Type = "dynamic"
				expectedInstanceNetworks.Network().IP = ""
				diskService.FindFound = true
				disks = []DiskCID{"fake-disk-1"}
			})

			It("clear out the ip part of the network settings", func() {
				vmCID, err = createVM.Run("fake-agent-id", "fake-stemcell-id", cloudProps, networks, disks, env)
				Expect(err).NotTo(HaveOccurred())
				Expect(vmService.CreateNetworks).To(Equal(expectedInstanceNetworks))
			})
		})
	})
})
