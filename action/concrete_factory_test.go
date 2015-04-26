package action_test

import (
	boshlog "github.com/cloudfoundry/bosh-agent/logger"
	fakeuuid "github.com/cloudfoundry/bosh-agent/uuid/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/frodenas/bosh-google-cpi/action"
	"github.com/frodenas/bosh-google-cpi/google/address"
	"github.com/frodenas/bosh-google-cpi/google/client"
	"github.com/frodenas/bosh-google-cpi/google/disk"
	"github.com/frodenas/bosh-google-cpi/google/disk_type"
	"github.com/frodenas/bosh-google-cpi/google/image"
	"github.com/frodenas/bosh-google-cpi/google/instance"
	"github.com/frodenas/bosh-google-cpi/google/machine_type"
	"github.com/frodenas/bosh-google-cpi/google/network"
	"github.com/frodenas/bosh-google-cpi/google/operation"
	"github.com/frodenas/bosh-google-cpi/google/snapshot"
	"github.com/frodenas/bosh-google-cpi/google/target_pool"
	"github.com/frodenas/bosh-google-cpi/registry"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/storage/v1"
)

var _ = Describe("concreteFactory", func() {
	var (
		project        string
		defaultZone    string
		uuidGen        *fakeuuid.FakeGenerator
		computeService *compute.Service
		storageService *storage.Service
		googleClient   gclient.GoogleClient
		logger         boshlog.Logger

		options = ConcreteFactoryOptions{
			Registry: registry.RegistryOptions{
				Schema:   "http",
				Host:     "fake-host",
				Port:     5555,
				Username: "fake-username",
				Password: "fake-password",
			},
		}

		factory Factory
	)

	var (
		operationService goperation.GoogleOperationService
	)

	BeforeEach(func() {
		//googleClient = fakewrdnclient.New()
		uuidGen = &fakeuuid.FakeGenerator{}
		logger = boshlog.NewLogger(boshlog.LevelNone)

		factory = NewConcreteFactory(
			googleClient,
			uuidGen,
			options,
			logger,
		)
	})

	BeforeEach(func() {
		operationService = goperation.NewGoogleOperationService(
			project,
			computeService,
			logger,
		)
	})

	It("returns error if action cannot be created", func() {
		action, err := factory.Create("fake-unknown-action")
		Expect(err).To(HaveOccurred())
		Expect(action).To(BeNil())
	})

	It("create_disk", func() {
		diskService := gdisk.NewGoogleDiskService(
			project,
			computeService,
			operationService,
			uuidGen,
			logger,
		)

		diskTypeService := gdisktype.NewGoogleDiskTypeService(
			project,
			computeService,
			logger,
		)

		vmService := ginstance.NewGoogleInstanceService(
			project,
			computeService,
			operationService,
			uuidGen,
			logger,
		)

		action, err := factory.Create("create_disk")
		Expect(err).ToNot(HaveOccurred())
		Expect(action).To(Equal(NewCreateDisk(diskService, diskTypeService, vmService, defaultZone)))
	})

	It("delete_disk", func() {
		diskService := gdisk.NewGoogleDiskService(
			project,
			computeService,
			operationService,
			uuidGen,
			logger,
		)

		action, err := factory.Create("delete_disk")
		Expect(err).ToNot(HaveOccurred())
		Expect(action).To(Equal(NewDeleteDisk(diskService)))
	})

	It("attach_disk", func() {
		diskService := gdisk.NewGoogleDiskService(
			project,
			computeService,
			operationService,
			uuidGen,
			logger,
		)

		vmService := ginstance.NewGoogleInstanceService(
			project,
			computeService,
			operationService,
			uuidGen,
			logger,
		)

		registryService := registry.NewRegistryService(
			options.Registry,
			logger,
		)

		action, err := factory.Create("attach_disk")
		Expect(err).ToNot(HaveOccurred())
		Expect(action).To(Equal(NewAttachDisk(diskService, vmService, registryService)))
	})

	It("detach_disk", func() {
		diskService := gdisk.NewGoogleDiskService(
			project,
			computeService,
			operationService,
			uuidGen,
			logger,
		)

		vmService := ginstance.NewGoogleInstanceService(
			project,
			computeService,
			operationService,
			uuidGen,
			logger,
		)

		registryService := registry.NewRegistryService(
			options.Registry,
			logger,
		)

		action, err := factory.Create("detach_disk")
		Expect(err).ToNot(HaveOccurred())
		Expect(action).To(Equal(NewDetachDisk(diskService, vmService, registryService)))
	})

	It("snapshot_disk", func() {
		snapshotService := gsnapshot.NewGoogleSnapshotService(
			project,
			computeService,
			operationService,
			uuidGen,
			logger,
		)

		diskService := gdisk.NewGoogleDiskService(
			project,
			computeService,
			operationService,
			uuidGen,
			logger,
		)

		action, err := factory.Create("snapshot_disk")
		Expect(err).ToNot(HaveOccurred())
		Expect(action).To(Equal(NewSnapshotDisk(snapshotService, diskService)))
	})

	It("delete_snapshot", func() {
		snapshotService := gsnapshot.NewGoogleSnapshotService(
			project,
			computeService,
			operationService,
			uuidGen,
			logger,
		)

		action, err := factory.Create("delete_snapshot")
		Expect(err).ToNot(HaveOccurred())
		Expect(action).To(Equal(NewDeleteSnapshot(snapshotService)))
	})

	It("create_stemcell", func() {
		stemcellService := gimage.NewGoogleImageService(
			project,
			computeService,
			storageService,
			operationService,
			uuidGen,
			logger,
		)

		action, err := factory.Create("create_stemcell")
		Expect(err).ToNot(HaveOccurred())
		Expect(action).To(Equal(NewCreateStemcell(stemcellService)))
	})

	It("delete_stemcell", func() {
		stemcellService := gimage.NewGoogleImageService(
			project,
			computeService,
			storageService,
			operationService,
			uuidGen,
			logger,
		)

		action, err := factory.Create("delete_stemcell")
		Expect(err).ToNot(HaveOccurred())
		Expect(action).To(Equal(NewDeleteStemcell(stemcellService)))
	})

	It("create_vm", func() {
		vmService := ginstance.NewGoogleInstanceService(
			project,
			computeService,
			operationService,
			uuidGen,
			logger,
		)

		diskService := gdisk.NewGoogleDiskService(
			project,
			computeService,
			operationService,
			uuidGen,
			logger,
		)

		machineTypeService := gmachinetype.NewGoogleMachineTypeService(
			project,
			computeService,
			logger,
		)

		networkService := gnetwork.NewGoogleNetworkService(
			project,
			computeService,
			logger,
		)

		stemcellService := gimage.NewGoogleImageService(
			project,
			computeService,
			storageService,
			operationService,
			uuidGen,
			logger,
		)

		registryService := registry.NewRegistryService(
			options.Registry,
			logger,
		)

		action, err := factory.Create("create_vm")
		Expect(err).ToNot(HaveOccurred())
		Expect(action).To(Equal(NewCreateVM(vmService, diskService, machineTypeService, networkService, stemcellService, registryService, options.Agent, defaultZone)))
	})

	It("configure_networks", func() {
		vmService := ginstance.NewGoogleInstanceService(
			project,
			computeService,
			operationService,
			uuidGen,
			logger,
		)

		addressService := gaddress.NewGoogleAddressService(
			project,
			computeService,
			logger,
		)

		networkService := gnetwork.NewGoogleNetworkService(
			project,
			computeService,
			logger,
		)

		targetPoolService := gtargetpool.NewGoogleTargetPoolService(
			project,
			computeService,
			operationService,
			logger,
		)

		registryService := registry.NewRegistryService(
			options.Registry,
			logger,
		)

		action, err := factory.Create("configure_networks")
		Expect(err).ToNot(HaveOccurred())
		Expect(action).To(Equal(NewConfigureNetworks(vmService, addressService, networkService, targetPoolService, registryService)))
	})

	It("delete_vm", func() {
		vmService := ginstance.NewGoogleInstanceService(
			project,
			computeService,
			operationService,
			uuidGen,
			logger,
		)

		registryService := registry.NewRegistryService(
			options.Registry,
			logger,
		)

		action, err := factory.Create("delete_vm")
		Expect(err).ToNot(HaveOccurred())
		Expect(action).To(Equal(NewDeleteVM(vmService, registryService)))
	})

	It("reboot_vm", func() {
		vmService := ginstance.NewGoogleInstanceService(
			project,
			computeService,
			operationService,
			uuidGen,
			logger,
		)

		action, err := factory.Create("reboot_vm")
		Expect(err).ToNot(HaveOccurred())
		Expect(action).To(Equal(NewRebootVM(vmService)))
	})

	It("set_vm_metadata", func() {
		vmService := ginstance.NewGoogleInstanceService(
			project,
			computeService,
			operationService,
			uuidGen,
			logger,
		)

		action, err := factory.Create("set_vm_metadata")
		Expect(err).ToNot(HaveOccurred())
		Expect(action).To(Equal(NewSetVMMetadata(vmService)))
	})

	It("has_vm", func() {
		vmService := ginstance.NewGoogleInstanceService(
			project,
			computeService,
			operationService,
			uuidGen,
			logger,
		)

		action, err := factory.Create("has_vm")
		Expect(err).ToNot(HaveOccurred())
		Expect(action).To(Equal(NewHasVM(vmService)))
	})

	It("get_disks", func() {
		vmService := ginstance.NewGoogleInstanceService(
			project,
			computeService,
			operationService,
			uuidGen,
			logger,
		)

		action, err := factory.Create("get_disks")
		Expect(err).ToNot(HaveOccurred())
		Expect(action).To(Equal(NewGetDisks(vmService)))
	})

	It("returns error because ping is not official CPI method if action is ping", func() {
		action, err := factory.Create("ping")
		Expect(err).To(HaveOccurred())
		Expect(action).To(BeNil())
	})
})
