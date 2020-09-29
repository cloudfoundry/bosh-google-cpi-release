package action_test

import (
	"encoding/json"
	"fmt"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	fakeuuid "github.com/cloudfoundry/bosh-utils/uuid/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	. "bosh-google-cpi/action"

	clientfakes "bosh-google-cpi/google/client/fakes"

	"bosh-google-cpi/config"
	"bosh-google-cpi/google/address_service"
	"bosh-google-cpi/google/backendservice_service"
	"bosh-google-cpi/google/client"
	"bosh-google-cpi/google/disk_service"
	"bosh-google-cpi/google/disk_type_service"
	"bosh-google-cpi/google/image_service"
	"bosh-google-cpi/google/instance_service"
	"bosh-google-cpi/google/machine_type_service"
	"bosh-google-cpi/google/network_service"
	"bosh-google-cpi/google/operation_service"
	"bosh-google-cpi/google/project_service"
	"bosh-google-cpi/google/snapshot_service"
	"bosh-google-cpi/google/subnetwork_service"
	"bosh-google-cpi/google/target_pool_service"

	"bosh-google-cpi/google/accelerator_type_service"
	"bosh-google-cpi/registry"
)

var _ = Describe("ConcreteFactory", func() {
	var (
		uuidGen      *fakeuuid.FakeGenerator
		googleClient client.GoogleClient
		logger       boshlog.Logger
		ctx          map[string]interface{}

		cfg = config.Config{
			Cloud: config.Cloud{
				Properties: config.CPIProperties{
					Registry: registry.ClientOptions{
						Protocol: "http",
						Host:     "fake-host",
						Port:     5555,
						Username: "fake-username",
						Password: "fake-password",
					},
				},
			},
		}

		factory Factory
	)

	var (
		operationService       operation.GoogleOperationService
		addressService         address.Service
		diskService            disk.Service
		diskTypeService        disktype.Service
		imageService           image.Service
		backendServiceService  backendservice.Service
		machineTypeService     machinetype.Service
		acceleratorTypeService acceleratortype.Service
		networkService         network.Service
		snapshotService        snapshot.Service
		subnetworkService      subnetwork.Service
		registryClient         registry.Client
		targetPoolService      targetpool.Service
		vmService              instance.Service
	)

	BeforeEach(func() {
		GoogleClientFunc = clientfakes.NewFakeGoogleClient
		uuidGen = &fakeuuid.FakeGenerator{}
		logger = boshlog.NewLogger(boshlog.LevelNone)

		ctx = map[string]interface{}{
			"project":                   "fake-project",
			"user_agent_prefix":         "fake-user-agent-prefix",
			"json_key":                  "{}",
			"default_root_disk_size_gb": 10,
			"default_root_disk_type":    "fake-root-disk-type",
		}

		ctxBytes, err := json.Marshal(ctx)
		Expect(err).ToNot(HaveOccurred())
		err = json.Unmarshal(ctxBytes, &cfg.Cloud.Properties.Google)
		Expect(err).ToNot(HaveOccurred())

		googleClient, _ = GoogleClientFunc(cfg.Cloud.Properties.Google, logger)

		factory = NewConcreteFactory(
			uuidGen,
			cfg,
			logger,
		)

		operationService = operation.NewGoogleOperationService(
			ctx["project"].(string),
			googleClient.ComputeService(),
			googleClient.ComputeBetaService(),
			logger,
		)

		addressService = address.NewGoogleAddressService(
			ctx["project"].(string),
			googleClient.ComputeService(),
			logger,
		)

		diskService = disk.NewGoogleDiskService(
			ctx["project"].(string),
			googleClient.ComputeService(),
			operationService,
			uuidGen,
			logger,
		)

		diskTypeService = disktype.NewGoogleDiskTypeService(
			ctx["project"].(string),
			googleClient.ComputeService(),
			logger,
		)

		imageService = image.NewGoogleImageService(
			ctx["project"].(string),
			googleClient.ComputeService(),
			googleClient.StorageService(),
			operationService,
			uuidGen,
			logger,
		)

		backendServiceService = backendservice.NewGoogleBackendServiceService(
			ctx["project"].(string),
			googleClient.ComputeService(),
			operationService,
			logger,
		)

		machineTypeService = machinetype.NewGoogleMachineTypeService(
			ctx["project"].(string),
			googleClient.ComputeService(),
			logger,
		)

		acceleratorTypeService = acceleratortype.NewGoogleAcceleratorTypeService(
			ctx["project"].(string),
			googleClient.ComputeService(),
			logger,
		)

		projectService := project.NewGoogleProjectService(
			ctx["project"].(string),
		)

		networkService = network.NewGoogleNetworkService(
			projectService,
			googleClient.ComputeService(),
			logger,
		)

		registryClient = registry.NewMetadataClient(
			googleClient,
			cfg.Cloud.Properties.Registry,
			logger,
		)

		snapshotService = snapshot.NewGoogleSnapshotService(
			ctx["project"].(string),
			googleClient.ComputeService(),
			operationService,
			uuidGen,
			logger,
		)

		subnetworkService = subnetwork.NewGoogleSubnetworkService(
			projectService,
			googleClient.ComputeService(),
			logger,
		)

		targetPoolService = targetpool.NewGoogleTargetPoolService(
			ctx["project"].(string),
			googleClient.ComputeService(),
			operationService,
			logger,
		)

		vmService = instance.NewGoogleInstanceService(
			ctx["project"].(string),
			googleClient.ComputeService(),
			googleClient.ComputeBetaService(),
			addressService,
			backendServiceService,
			networkService,
			operationService,
			subnetworkService,
			targetPoolService,
			diskTypeService,
			uuidGen,
			logger,
		)
	})

	apiVersions := []int{1, 2}
	for _, val := range apiVersions {
		apiVersion := val

		Context(fmt.Sprintf("Api Version %d", apiVersion), func() {

			It("returns error if action cannot be created", func() {
				action, err := factory.Create("fake-unknown-action", ctx, apiVersion)
				Expect(err).To(HaveOccurred())
				Expect(action).To(BeNil())
			})

			It("create_disk", func() {
				action, err := factory.Create("create_disk", ctx, apiVersion)
				Expect(err).ToNot(HaveOccurred())
				Expect(action).To(Equal(NewCreateDisk(
					diskService,
					diskTypeService,
					vmService,
				)))
			})

			It("delete_disk", func() {
				action, err := factory.Create("delete_disk", ctx, apiVersion)
				Expect(err).ToNot(HaveOccurred())
				Expect(action).To(Equal(NewDeleteDisk(diskService)))
			})

			It("detach_disk", func() {
				action, err := factory.Create("detach_disk", ctx, apiVersion)
				Expect(err).ToNot(HaveOccurred())
				Expect(action).To(Equal(NewDetachDisk(vmService, registryClient)))
			})

			It("snapshot_disk", func() {
				action, err := factory.Create("snapshot_disk", ctx, apiVersion)
				Expect(err).ToNot(HaveOccurred())
				Expect(action).To(Equal(NewSnapshotDisk(snapshotService, diskService)))
			})

			It("delete_snapshot", func() {
				action, err := factory.Create("delete_snapshot", ctx, apiVersion)
				Expect(err).ToNot(HaveOccurred())
				Expect(action).To(Equal(NewDeleteSnapshot(snapshotService)))
			})

			It("create_stemcell", func() {
				action, err := factory.Create("create_stemcell", ctx, apiVersion)
				Expect(err).ToNot(HaveOccurred())
				Expect(action).To(Equal(NewCreateStemcell(imageService)))
			})

			It("delete_stemcell", func() {
				action, err := factory.Create("delete_stemcell", ctx, apiVersion)
				Expect(err).ToNot(HaveOccurred())
				Expect(action).To(Equal(NewDeleteStemcell(imageService)))
			})

			It("configure_networks", func() {
				action, err := factory.Create("configure_networks", ctx, apiVersion)
				Expect(err).ToNot(HaveOccurred())
				Expect(action).To(Equal(NewConfigureNetworks(vmService, registryClient)))
			})

			It("delete_vm", func() {
				action, err := factory.Create("delete_vm", ctx, apiVersion)
				Expect(err).ToNot(HaveOccurred())
				Expect(action).To(Equal(NewDeleteVM(vmService, registryClient)))
			})

			It("reboot_vm", func() {
				action, err := factory.Create("reboot_vm", ctx, apiVersion)
				Expect(err).ToNot(HaveOccurred())
				Expect(action).To(Equal(NewRebootVM(vmService)))
			})

			It("set_vm_metadata", func() {
				action, err := factory.Create("set_vm_metadata", ctx, apiVersion)
				Expect(err).ToNot(HaveOccurred())
				Expect(action).To(Equal(NewSetVMMetadata(vmService)))
			})

			It("has_vm", func() {
				action, err := factory.Create("has_vm", ctx, apiVersion)
				Expect(err).ToNot(HaveOccurred())
				Expect(action).To(Equal(NewHasVM(vmService)))
			})

			It("get_disks", func() {
				action, err := factory.Create("get_disks", ctx, apiVersion)
				Expect(err).ToNot(HaveOccurred())
				Expect(action).To(Equal(NewGetDisks(vmService)))
			})

			It("ping", func() {
				action, err := factory.Create("ping", ctx, apiVersion)
				Expect(err).ToNot(HaveOccurred())
				Expect(action).To(Equal(NewPing()))
			})

			It("info", func() {
				action, err := factory.Create("info", ctx, apiVersion)
				Expect(err).ToNot(HaveOccurred())
				Expect(action).To(Equal(NewInfo()))
			})

			It("when action is current_vm_id returns an error because this CPI does not implement the method", func() {
				action, err := factory.Create("current_vm_id", ctx, apiVersion)
				Expect(err).To(HaveOccurred())
				Expect(action).To(BeNil())
			})

			It("when action is wrong returns an error because it is not an official CPI method", func() {
				action, err := factory.Create("wrong", ctx, apiVersion)
				Expect(err).To(HaveOccurred())
				Expect(action).To(BeNil())
			})
		})
	}

	DescribeTable("attach_disk",
		func(apiVersion int, f func() interface{}) {
			action, err := factory.Create("attach_disk", ctx, apiVersion)
			Expect(err).ToNot(HaveOccurred())
			Expect(action).To(Equal(f()))
		},
		Entry("apiVersion 1", 1, func() interface{} { return NewAttachDiskV1(diskService, vmService, registryClient) }),
		Entry("apiVersion 2", 2, func() interface{} { return NewAttachDiskV2(diskService, vmService, registryClient) }),
	)

	DescribeTable("create_vm",
		func(apiVersion int, f func() interface{}) {
			action, err := factory.Create("create_vm", ctx, apiVersion)
			Expect(err).ToNot(HaveOccurred())
			Expect(action).To(Equal(f()))
		},

		Entry("apiVersion 1", 1,
			func() interface{} {
				return NewCreateVMV1(
					vmService,
					diskService,
					diskTypeService,
					imageService,
					machineTypeService,
					acceleratorTypeService,
					registryClient,
					cfg.Cloud.Properties.Registry,
					cfg.Cloud.Properties.Agent,
					ctx["default_root_disk_size_gb"].(int),
					ctx["default_root_disk_type"].(string),
				)
			},
		),

		Entry("apiVersion 2", 2,
			func() interface{} {
				return NewCreateVMV2(
					vmService,
					diskService,
					diskTypeService,
					imageService,
					machineTypeService,
					acceleratorTypeService,
					registryClient,
					cfg.Cloud.Properties.Registry,
					cfg.Cloud.Properties.Agent,
					ctx["default_root_disk_size_gb"].(int),
					ctx["default_root_disk_type"].(string),
				)
			},
		),
	)
})
