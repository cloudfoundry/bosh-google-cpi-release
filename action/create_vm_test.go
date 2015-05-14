package action_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/frodenas/bosh-google-cpi/action"

	addressfakes "github.com/frodenas/bosh-google-cpi/google/address_service/fakes"
	diskfakes "github.com/frodenas/bosh-google-cpi/google/disk_service/fakes"
	disktypefakes "github.com/frodenas/bosh-google-cpi/google/disk_type_service/fakes"
	imagefakes "github.com/frodenas/bosh-google-cpi/google/image_service/fakes"
	instancefakes "github.com/frodenas/bosh-google-cpi/google/instance_service/fakes"
	machinetypefakes "github.com/frodenas/bosh-google-cpi/google/machine_type_service/fakes"
	networkfakes "github.com/frodenas/bosh-google-cpi/google/network_service/fakes"
	targetpoolfakes "github.com/frodenas/bosh-google-cpi/google/target_pool_service/fakes"
	registryfakes "github.com/frodenas/bosh-registry/client/fakes"

	"github.com/frodenas/bosh-google-cpi/api"
	"github.com/frodenas/bosh-google-cpi/google/instance_service"
	"github.com/frodenas/bosh-registry/client"
	"google.golang.org/api/compute/v1"
)

var _ = Describe("CreateVM", func() {
	var (
		err                      error
		vmCID                    VMCID
		networks                 Networks
		cloudProps               VMCloudProperties
		disks                    []DiskCID
		env                      Environment
		registryOptions          registry.ClientOptions
		agentOptions             registry.AgentOptions
		expectedVMProps          *ginstance.InstanceProperties
		vmNetworks               ginstance.InstanceNetworks
		expectedInstanceNetworks ginstance.GoogleInstanceNetworks
		expectedAgentSettings    registry.AgentSettings

		vmService          *instancefakes.FakeInstanceService
		addressService     *addressfakes.FakeAddressService
		diskService        *diskfakes.FakeDiskService
		diskTypeService    *disktypefakes.FakeDiskTypeService
		machineTypeService *machinetypefakes.FakeMachineTypeService
		networkService     *networkfakes.FakeNetworkService
		stemcellService    *imagefakes.FakeImageService
		targetPoolService  *targetpoolfakes.FakeTargetPoolService
		registryClient     *registryfakes.FakeClient

		createVM CreateVM
	)

	BeforeEach(func() {
		vmService = &instancefakes.FakeInstanceService{}
		addressService = &addressfakes.FakeAddressService{}
		diskService = &diskfakes.FakeDiskService{}
		diskTypeService = &disktypefakes.FakeDiskTypeService{}
		machineTypeService = &machinetypefakes.FakeMachineTypeService{}
		networkService = &networkfakes.FakeNetworkService{}
		stemcellService = &imagefakes.FakeImageService{}
		targetPoolService = &targetpoolfakes.FakeTargetPoolService{}
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
				Type: "fake-blobstore-type",
			},
		}
		createVM = NewCreateVM(
			vmService,
			addressService,
			diskService,
			diskTypeService,
			machineTypeService,
			networkService,
			stemcellService,
			targetPoolService,
			registryClient,
			registryOptions,
			agentOptions,
			"fake-default-zone",
		)
	})

	Describe("Run", func() {
		BeforeEach(func() {
			vmService.CreateID = "fake-vm-id"
			stemcellService.FindFound = true
			machineTypeService.FindFound = true

			diskService.FindDisk = &compute.Disk{Zone: "fake-disk-zone"}
			stemcellService.FindImage = &compute.Image{SelfLink: "fake-image-self-link"}
			machineTypeService.FindMachineType = &compute.MachineType{SelfLink: "fake-machine-type-self-link"}
			diskTypeService.FindDiskType = &compute.DiskType{SelfLink: "fake-disk-type-self-link"}

			cloudProps = VMCloudProperties{
				Zone:              "",
				MachineType:       "fake-machine-type",
				RootDiskSizeGb:    10,
				RootDiskType:      "",
				AutomaticRestart:  true,
				OnHostMaintenance: "TERMINATE",
				ServiceScopes:     []string{},
			}

			networks = Networks{
				"fake-network-name": Network{
					Type:    "dynamic",
					IP:      "fake-network-ip",
					Gateway: "fake-network-gateway",
					Netmask: "fake-network-netmask",
					DNS:     []string{"fake-network-dns"},
					Default: []string{"fake-network-default"},
					CloudProperties: NetworkCloudProperties{
						NetworkName:         "fake-network-cloud-network-name",
						Tags:                NetworkTags{"fake-network-cloud-network-tag"},
						EphemeralExternalIP: true,
						IPForwarding:        false,
						TargetPool:          "fake-network-cloud-target-pool",
					},
				},
			}

			expectedVMProps = &ginstance.InstanceProperties{
				Zone:              "fake-default-zone",
				Stemcell:          "fake-image-self-link",
				MachineType:       "fake-machine-type-self-link",
				RootDiskSizeGb:    10,
				RootDiskType:      "",
				AutomaticRestart:  true,
				OnHostMaintenance: "TERMINATE",
				ServiceScopes:     []string{},
			}

			vmNetworks = networks.AsInstanceServiceNetworks()
			expectedInstanceNetworks = ginstance.NewGoogleInstanceNetworks(
				vmNetworks,
				addressService,
				networkService,
				targetPoolService,
			)

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
						Type:    "dynamic",
						IP:      "fake-network-ip",
						Gateway: "fake-network-gateway",
						Netmask: "fake-network-netmask",
						DNS:     []string{"fake-network-dns"},
						Default: []string{"fake-network-default"},
					},
				},
				VM: registry.VMSettings{
					Name: "fake-vm-id",
				},
			}
		})

		It("creates the vm", func() {
			vmCID, err = createVM.Run("fake-agent-id", "fake-stemcell-id", cloudProps, networks, disks, env)
			Expect(err).NotTo(HaveOccurred())
			Expect(diskService.FindCalled).To(BeFalse())
			Expect(stemcellService.FindCalled).To(BeTrue())
			Expect(machineTypeService.FindCalled).To(BeTrue())
			Expect(diskTypeService.FindCalled).To(BeFalse())
			Expect(vmService.CreateCalled).To(BeTrue())
			Expect(vmService.CleanUpCalled).To(BeFalse())
			Expect(vmService.AddNetworkConfigurationCalled).To(BeTrue())
			Expect(registryClient.UpdateCalled).To(BeTrue())
			Expect(registryClient.UpdateSettings).To(Equal(expectedAgentSettings))
			Expect(vmCID).To(Equal(VMCID("fake-vm-id")))
			Expect(vmService.CreateVMProps).To(Equal(expectedVMProps))
			Expect(vmService.CreateInstanceNetworks).To(Equal(expectedInstanceNetworks))
			Expect(vmService.CreateRegistryEndpoint).To(Equal("http://fake-registry-host:25777"))
		})

		It("returns an error if stemcellService find call returns an error", func() {
			stemcellService.FindErr = errors.New("fake-stemcell-service-error")

			_, err = createVM.Run("fake-agent-id", "fake-stemcell-id", cloudProps, networks, disks, env)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("fake-stemcell-service-error"))
			Expect(diskService.FindCalled).To(BeFalse())
			Expect(stemcellService.FindCalled).To(BeTrue())
			Expect(machineTypeService.FindCalled).To(BeFalse())
			Expect(diskTypeService.FindCalled).To(BeFalse())
			Expect(vmService.CreateCalled).To(BeFalse())
			Expect(vmService.CleanUpCalled).To(BeFalse())
			Expect(vmService.AddNetworkConfigurationCalled).To(BeFalse())
			Expect(registryClient.UpdateCalled).To(BeFalse())
		})

		It("returns an error if stemcell is not found", func() {
			stemcellService.FindFound = false

			_, err = createVM.Run("fake-agent-id", "fake-stemcell-id", cloudProps, networks, disks, env)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Stemcell 'fake-stemcell-id' does not exists"))
			Expect(diskService.FindCalled).To(BeFalse())
			Expect(stemcellService.FindCalled).To(BeTrue())
			Expect(machineTypeService.FindCalled).To(BeFalse())
			Expect(diskTypeService.FindCalled).To(BeFalse())
			Expect(vmService.CreateCalled).To(BeFalse())
			Expect(vmService.CleanUpCalled).To(BeFalse())
			Expect(vmService.AddNetworkConfigurationCalled).To(BeFalse())
			Expect(registryClient.UpdateCalled).To(BeFalse())
		})

		It("returns an error if machine type is not set", func() {
			cloudProps.MachineType = ""

			_, err = createVM.Run("fake-agent-id", "fake-stemcell-id", cloudProps, networks, disks, env)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("'machine_type' must be provided"))
			Expect(diskService.FindCalled).To(BeFalse())
			Expect(stemcellService.FindCalled).To(BeTrue())
			Expect(machineTypeService.FindCalled).To(BeFalse())
			Expect(diskTypeService.FindCalled).To(BeFalse())
			Expect(vmService.CreateCalled).To(BeFalse())
			Expect(vmService.CleanUpCalled).To(BeFalse())
			Expect(vmService.AddNetworkConfigurationCalled).To(BeFalse())
			Expect(registryClient.UpdateCalled).To(BeFalse())
		})

		It("returns an error if machineTypeService find call returns an error", func() {
			machineTypeService.FindErr = errors.New("fake-machine-type-service-error")

			_, err = createVM.Run("fake-agent-id", "fake-stemcell-id", cloudProps, networks, disks, env)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("fake-machine-type-service-error"))
			Expect(diskService.FindCalled).To(BeFalse())
			Expect(stemcellService.FindCalled).To(BeTrue())
			Expect(machineTypeService.FindCalled).To(BeTrue())
			Expect(diskTypeService.FindCalled).To(BeFalse())
			Expect(vmService.CreateCalled).To(BeFalse())
			Expect(vmService.CleanUpCalled).To(BeFalse())
			Expect(vmService.AddNetworkConfigurationCalled).To(BeFalse())
			Expect(registryClient.UpdateCalled).To(BeFalse())
		})

		It("returns an error if machine type is not found", func() {
			machineTypeService.FindFound = false

			_, err = createVM.Run("fake-agent-id", "fake-stemcell-id", cloudProps, networks, disks, env)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Machine Type 'fake-machine-type' does not exists"))
			Expect(diskService.FindCalled).To(BeFalse())
			Expect(stemcellService.FindCalled).To(BeTrue())
			Expect(machineTypeService.FindCalled).To(BeTrue())
			Expect(diskTypeService.FindCalled).To(BeFalse())
			Expect(vmService.CreateCalled).To(BeFalse())
			Expect(vmService.CleanUpCalled).To(BeFalse())
			Expect(vmService.AddNetworkConfigurationCalled).To(BeFalse())
			Expect(registryClient.UpdateCalled).To(BeFalse())
		})

		It("returns an error if vmService create call returns an error", func() {
			vmService.CreateErr = errors.New("fake-vm-service-error")

			_, err = createVM.Run("fake-agent-id", "fake-stemcell-id", cloudProps, networks, disks, env)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("fake-vm-service-error"))
			Expect(diskService.FindCalled).To(BeFalse())
			Expect(stemcellService.FindCalled).To(BeTrue())
			Expect(machineTypeService.FindCalled).To(BeTrue())
			Expect(diskTypeService.FindCalled).To(BeFalse())
			Expect(vmService.CreateCalled).To(BeTrue())
			Expect(vmService.CleanUpCalled).To(BeFalse())
			Expect(vmService.AddNetworkConfigurationCalled).To(BeFalse())
			Expect(registryClient.UpdateCalled).To(BeFalse())
		})

		It("returns an error if vmService add network configuration call returns an error", func() {
			vmService.AddNetworkConfigurationErr = errors.New("fake-vm-service-error")

			_, err = createVM.Run("fake-agent-id", "fake-stemcell-id", cloudProps, networks, disks, env)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("fake-vm-service-error"))
			Expect(diskService.FindCalled).To(BeFalse())
			Expect(stemcellService.FindCalled).To(BeTrue())
			Expect(machineTypeService.FindCalled).To(BeTrue())
			Expect(diskTypeService.FindCalled).To(BeFalse())
			Expect(vmService.CreateCalled).To(BeTrue())
			Expect(vmService.CleanUpCalled).To(BeTrue())
			Expect(vmService.AddNetworkConfigurationCalled).To(BeTrue())
			Expect(registryClient.UpdateCalled).To(BeFalse())
		})

		It("returns an error if registryClient update call returns an error", func() {
			registryClient.UpdateErr = errors.New("fake-registry-client-error")

			_, err = createVM.Run("fake-agent-id", "fake-stemcell-id", cloudProps, networks, disks, env)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("fake-registry-client-error"))
			Expect(diskService.FindCalled).To(BeFalse())
			Expect(stemcellService.FindCalled).To(BeTrue())
			Expect(machineTypeService.FindCalled).To(BeTrue())
			Expect(diskTypeService.FindCalled).To(BeFalse())
			Expect(vmService.CreateCalled).To(BeTrue())
			Expect(vmService.CleanUpCalled).To(BeTrue())
			Expect(vmService.AddNetworkConfigurationCalled).To(BeTrue())
			Expect(registryClient.UpdateCalled).To(BeTrue())
		})

		Context("when root disk type is set", func() {
			BeforeEach(func() {
				diskTypeService.FindFound = true
				cloudProps.RootDiskType = "fake-root-disk-type"
				expectedVMProps.RootDiskType = "fake-disk-type-self-link"
			})

			It("creates the vm with the right properties", func() {
				vmCID, err = createVM.Run("fake-agent-id", "fake-stemcell-id", cloudProps, networks, disks, env)
				Expect(err).NotTo(HaveOccurred())
				Expect(diskService.FindCalled).To(BeFalse())
				Expect(stemcellService.FindCalled).To(BeTrue())
				Expect(machineTypeService.FindCalled).To(BeTrue())
				Expect(diskTypeService.FindCalled).To(BeTrue())
				Expect(vmService.CreateCalled).To(BeTrue())
				Expect(vmService.CleanUpCalled).To(BeFalse())
				Expect(vmService.AddNetworkConfigurationCalled).To(BeTrue())
				Expect(registryClient.UpdateCalled).To(BeTrue())
				Expect(registryClient.UpdateSettings).To(Equal(expectedAgentSettings))
				Expect(vmCID).To(Equal(VMCID("fake-vm-id")))
				Expect(vmService.CreateVMProps).To(Equal(expectedVMProps))
				Expect(vmService.CreateInstanceNetworks).To(Equal(expectedInstanceNetworks))
				Expect(vmService.CreateRegistryEndpoint).To(Equal("http://fake-registry-host:25777"))
			})

			It("returns an error if diskTypeService find call returns an error", func() {
				diskTypeService.FindErr = errors.New("fake-disk-type-service-error")

				_, err = createVM.Run("fake-agent-id", "fake-stemcell-id", cloudProps, networks, disks, env)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("fake-disk-type-service-error"))
				Expect(diskService.FindCalled).To(BeFalse())
				Expect(stemcellService.FindCalled).To(BeTrue())
				Expect(machineTypeService.FindCalled).To(BeTrue())
				Expect(diskTypeService.FindCalled).To(BeTrue())
				Expect(vmService.CreateCalled).To(BeFalse())
				Expect(vmService.CleanUpCalled).To(BeFalse())
				Expect(vmService.AddNetworkConfigurationCalled).To(BeFalse())
				Expect(registryClient.UpdateCalled).To(BeFalse())
			})

			It("returns an error if root disk type is not found", func() {
				diskTypeService.FindFound = false

				_, err = createVM.Run("fake-agent-id", "fake-stemcell-id", cloudProps, networks, disks, env)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Root Disk Type 'fake-root-disk-type' does not exists"))
				Expect(diskService.FindCalled).To(BeFalse())
				Expect(stemcellService.FindCalled).To(BeTrue())
				Expect(machineTypeService.FindCalled).To(BeTrue())
				Expect(diskTypeService.FindCalled).To(BeTrue())
				Expect(vmService.CreateCalled).To(BeFalse())
				Expect(vmService.CleanUpCalled).To(BeFalse())
				Expect(vmService.AddNetworkConfigurationCalled).To(BeFalse())
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
				Expect(stemcellService.FindCalled).To(BeTrue())
				Expect(machineTypeService.FindCalled).To(BeTrue())
				Expect(diskTypeService.FindCalled).To(BeFalse())
				Expect(vmService.CreateCalled).To(BeTrue())
				Expect(vmService.CleanUpCalled).To(BeFalse())
				Expect(vmService.AddNetworkConfigurationCalled).To(BeTrue())
				Expect(registryClient.UpdateCalled).To(BeTrue())
				Expect(registryClient.UpdateSettings).To(Equal(expectedAgentSettings))
				Expect(vmCID).To(Equal(VMCID("fake-vm-id")))
				Expect(vmService.CreateVMProps).To(Equal(expectedVMProps))
				Expect(vmService.CreateInstanceNetworks).To(Equal(expectedInstanceNetworks))
				Expect(vmService.CreateRegistryEndpoint).To(Equal("http://fake-registry-host:25777"))
			})
		})

		Context("when DiskCIDs is set", func() {
			BeforeEach(func() {
				diskService.FindFound = true
				disks = []DiskCID{"fake-disk-1"}
				expectedVMProps.Zone = "fake-disk-zone"
			})

			It("creates the vm at the right zone", func() {
				vmCID, err = createVM.Run("fake-agent-id", "fake-stemcell-id", cloudProps, networks, disks, env)
				Expect(err).NotTo(HaveOccurred())
				Expect(diskService.FindCalled).To(BeTrue())
				Expect(stemcellService.FindCalled).To(BeTrue())
				Expect(machineTypeService.FindCalled).To(BeTrue())
				Expect(diskTypeService.FindCalled).To(BeFalse())
				Expect(vmService.CreateCalled).To(BeTrue())
				Expect(vmService.CleanUpCalled).To(BeFalse())
				Expect(vmService.AddNetworkConfigurationCalled).To(BeTrue())
				Expect(registryClient.UpdateCalled).To(BeTrue())
				Expect(registryClient.UpdateSettings).To(Equal(expectedAgentSettings))
				Expect(vmCID).To(Equal(VMCID("fake-vm-id")))
				Expect(vmService.CreateVMProps).To(Equal(expectedVMProps))
				Expect(vmService.CreateInstanceNetworks).To(Equal(expectedInstanceNetworks))
				Expect(vmService.CreateRegistryEndpoint).To(Equal("http://fake-registry-host:25777"))
			})

			It("returns an error if diskService find call returns an error", func() {
				diskService.FindErr = errors.New("fake-disk-service-error")

				_, err = createVM.Run("fake-agent-id", "fake-stemcell-id", cloudProps, networks, disks, env)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("fake-disk-service-error"))
				Expect(diskService.FindCalled).To(BeTrue())
				Expect(stemcellService.FindCalled).To(BeFalse())
				Expect(machineTypeService.FindCalled).To(BeFalse())
				Expect(diskTypeService.FindCalled).To(BeFalse())
				Expect(vmService.CreateCalled).To(BeFalse())
				Expect(vmService.CleanUpCalled).To(BeFalse())
				Expect(vmService.AddNetworkConfigurationCalled).To(BeFalse())
				Expect(registryClient.UpdateCalled).To(BeFalse())
			})

			It("returns an error if disk is not found", func() {
				diskService.FindFound = false

				_, err = createVM.Run("fake-agent-id", "fake-stemcell-id", cloudProps, networks, disks, env)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal(api.NewDiskNotFoundError("fake-disk-1", false).Error()))
				Expect(diskService.FindCalled).To(BeTrue())
				Expect(stemcellService.FindCalled).To(BeFalse())
				Expect(machineTypeService.FindCalled).To(BeFalse())
				Expect(diskTypeService.FindCalled).To(BeFalse())
				Expect(vmService.CreateCalled).To(BeFalse())
				Expect(vmService.CleanUpCalled).To(BeFalse())
				Expect(vmService.AddNetworkConfigurationCalled).To(BeFalse())
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
					Expect(stemcellService.FindCalled).To(BeFalse())
					Expect(machineTypeService.FindCalled).To(BeFalse())
					Expect(diskTypeService.FindCalled).To(BeFalse())
					Expect(vmService.CreateCalled).To(BeFalse())
					Expect(vmService.CleanUpCalled).To(BeFalse())
					Expect(vmService.AddNetworkConfigurationCalled).To(BeFalse())
					Expect(registryClient.UpdateCalled).To(BeFalse())
				})
			})
		})
	})
})
