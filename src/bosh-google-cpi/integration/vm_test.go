package integration

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"google.golang.org/api/compute/v1"
)

var _ = Describe("VM", func() {
	It("creates a VM with an invalid configuration and receives an error message with logs", func() {
		request := fmt.Sprintf(`{
			  "method": "create_vm",
			  "arguments": [
				"agent",
				"%v",
				{
				  "machine_type": "n1-standard-error"
				},
				{
				  "default": {
					"type": "dynamic",
					"cloud_properties": {
					  "tags": ["integration-delete"],
					  "network_name": "%v"
					}
				  }
				},
				[],
				{}
			  ]
			}`, existingStemcell, networkName)
		resp, err := execCPI(request)
		Expect(err).ToNot(HaveOccurred())
		Expect(resp.Error.Message).ToNot(BeEmpty())
		Expect(resp.Log).ToNot(BeEmpty())
	})

	It("executes the VM lifecycle", func() {
		var vmCID string
		By("creating a VM")
		request := fmt.Sprintf(`{
			  "method": "create_vm",
			  "arguments": [
				"agent",
				"%v",
				{
				  "machine_type": "n1-standard-1"
				},
				{
				  "default": {
					"type": "dynamic",
					"cloud_properties": {
					  "tags": ["integration-delete"],
					  "network_name": "%v"
					}
				  }
				},
				[],
				{}
			  ]
			}`, existingStemcell, networkName)
		vmCID = assertSucceedsWithResult(request).(string)

		By("locating the VM")
		request = fmt.Sprintf(`{
			  "method": "has_vm",
			  "arguments": ["%v"]
			}`, vmCID)
		exists := assertSucceedsWithResult(request).(bool)
		Expect(exists).To(Equal(true))

		By("rebooting the VM")
		request = fmt.Sprintf(`{
			  "method": "reboot_vm",
			  "arguments": ["%v"]
			}`, vmCID)
		assertSucceeds(request)

		By("deleting the VM")
		request = fmt.Sprintf(`{
			  "method": "delete_vm",
			  "arguments": ["%v"]
			}`, vmCID)
		assertSucceeds(request)

	})

	It("can create a VM with tags", func() {
		By("creating a VM")
		var vmCID string
		request := fmt.Sprintf(`{
			  "method": "create_vm",
			  "arguments": [
				"agent",
				"%v",
				{
				   "machine_type": "n1-standard-1",
				   "tags": ["tag1", "tag2"]
				},
				{
				  "default": {
					"type": "dynamic",
					"cloud_properties": {
					  "tags": ["integration-delete"],
					  "network_name": "%v"
					}
				  }
				},
				[],
				{}
			  ]
			}`, existingStemcell, networkName)
		vmCID = assertSucceedsWithResult(request).(string)
		assertValidVM(vmCID, func(instance *compute.Instance) {
			Expect(instance.Tags.Items).To(ConsistOf("integration-delete", "tag1", "tag2"))
		})

		By("deleting the VM")
		request = fmt.Sprintf(`{
			  "method": "delete_vm",
			  "arguments": ["%v"]
			}`, vmCID)
		assertSucceeds(request)
	})

	It("can create a VM with overlapping VM and network tags", func() {
		By("creating a VM")
		var vmCID string
		request := fmt.Sprintf(`{
			  "method": "create_vm",
			  "arguments": [
				"agent",
				"%v",
				{
				   "machine_type": "n1-standard-1",
				   "tags": ["tag1", "tag2", "integration-delete"]
				},
				{
				  "default": {
					"type": "dynamic",
					"cloud_properties": {
					  "tags": ["integration-delete"],
					  "network_name": "%v"
					}
				  }
				},
				[],
				{}
			  ]
			}`, existingStemcell, networkName)
		vmCID = assertSucceedsWithResult(request).(string)
		assertValidVM(vmCID, func(instance *compute.Instance) {
			Expect(instance.Tags.Items).To(ConsistOf("integration-delete", "tag1", "tag2"))
		})

		By("deleting the VM")
		request = fmt.Sprintf(`{
			  "method": "delete_vm",
			  "arguments": ["%v"]
			}`, vmCID)
		assertSucceeds(request)
	})
	It("executes the VM lifecycle with disk attachment hints", func() {
		By("creating two disks")
		var request, diskCID, diskCID2, vmCID string
		request = fmt.Sprintf(`{
			  "method": "create_disk",
			  "arguments": [32768, {}, ""]
			}`)
		diskCID = assertSucceedsWithResult(request).(string)
		diskCID2 = assertSucceedsWithResult(request).(string)

		By("creating a VM with the disk attachment hints")
		request = fmt.Sprintf(`{
			  "method": "create_vm",
			  "arguments": [
				"agent",
				"%v",
				{
				  "machine_type": "n1-standard-1"
				},
				{
				  "default": {
					"type": "dynamic",
					"cloud_properties": {
					  "tags": ["integration-delete"],
					  "network_name": "%v"
					}
				  }
				},
				["%v", "%v"],
				{}
			  ]
			}`, existingStemcell, networkName, diskCID, diskCID2)
		vmCID = assertSucceedsWithResult(request).(string)

		By("deleting the disks")
		request = fmt.Sprintf(`{
			  "method": "delete_disk",
			  "arguments": ["%v"]
			}`, diskCID)
		assertSucceeds(request)

		request = fmt.Sprintf(`{
			  "method": "delete_disk",
			  "arguments": ["%v"]
			}`, diskCID2)
		assertSucceeds(request)

		By("deleting the VM")
		request = fmt.Sprintf(`{
			  "method": "delete_vm",
			  "arguments": ["%v"]
			}`, vmCID)
		assertSucceeds(request)
	})

	It("executes the VM lifecycle with custom machine type", func() {
		By("creating a VM")
		var vmCID string
		request := fmt.Sprintf(`{
			  "method": "create_vm",
			  "arguments": [
				"agent",
				"%v",
				{
				  "cpu": 2,
				  "ram": 5120
				},
				{
				  "default": {
					"type": "dynamic",
					"cloud_properties": {
					  "tags": ["integration-delete"],
					  "network_name": "%v"
					}
				  }
				},
				[],
				{}
			  ]
			}`, existingStemcell, networkName)
		vmCID = assertSucceedsWithResult(request).(string)

		By("deleting the VM")
		request = fmt.Sprintf(`{
			  "method": "delete_vm",
			  "arguments": ["%v"]
			}`, vmCID)
		assertSucceeds(request)
	})

	It("executes the VM lifecycle in a specific zone", func() {
		By("creating a VM")
		var vmCID string

		request := fmt.Sprintf(`{
			  "method": "create_vm",
			  "arguments": [
				"agent",
				"%v",
				{
				  "machine_type": "n1-standard-1",
				  "zone": "us-central1-b"
				},
				{
				  "default": {
					"type": "dynamic",
					"cloud_properties": {
					  "tags": ["integration-delete"],
					  "network_name": "%v"
					}
				  }
				},
				[],
				{}
			  ]
			}`, existingStemcell, networkName)
		vmCID = assertSucceedsWithResult(request).(string)

		By("deleting the VM")
		request = fmt.Sprintf(`{
			  "method": "delete_vm",
			  "arguments": ["%v"]
			}`, vmCID)
	})

	It("executes the VM lifecycle with automatic restart disabled", func() {
		By("creating a VM")
		var vmCID string
		request := fmt.Sprintf(`{
			  "method": "create_vm",
			  "arguments": [
				"agent",
				"%v",
				{
				  "machine_type": "n1-standard-1",
				  "automatic_restart": false
				},
				{
				  "default": {
					"type": "dynamic",
					"cloud_properties": {
					  "tags": ["integration-delete"],
					  "network_name": "%v"
					}
				  }
				},
				[],
				{}
			  ]
			}`, existingStemcell, networkName)
		vmCID = assertSucceedsWithResult(request).(string)

		By("deleting the VM")
		request = fmt.Sprintf(`{
			  "method": "delete_vm",
			  "arguments": ["%v"]
			}`, vmCID)
		assertSucceeds(request)
	})

	It("execute the VM lifecycle with OnHostMaintenance modified", func() {
		By("creating a VM")
		var vmCID string

		request := fmt.Sprintf(`{
			  "method": "create_vm",
			  "arguments": [
				"agent",
				"%v",
				{
				  "machine_type": "n1-standard-1",
				  "on_host_maintenance": "TERMINATE"
				},
				{
				  "default": {
					"type": "dynamic",
					"cloud_properties": {
					  "tags": ["integration-delete"],
					  "network_name": "%v"
					}
				  }
				},
				[],
				{}
			  ]
			}`, existingStemcell, networkName)

		By("deleting the VM")
		vmCID = assertSucceedsWithResult(request).(string)
		request = fmt.Sprintf(`{
			  "method": "delete_vm",
			  "arguments": ["%v"]
			}`, vmCID)
		assertSucceeds(request)
	})

	It("can execute the VM lifecycle with a preemptible VM", func() {
		By("creating a VM")
		var vmCID string

		request := fmt.Sprintf(`{
			  "method": "create_vm",
			  "arguments": [
				"agent",
				"%v",
				{
				  "machine_type": "n1-standard-1",
				  "preemtible": true
				},
				{
				  "default": {
					"type": "dynamic",
					"cloud_properties": {
					  "tags": ["integration-delete"],
					  "network_name": "%v"
					}
				  }
				},
				[],
				{}
			  ]
			}`, existingStemcell, networkName)
		vmCID = assertSucceedsWithResult(request).(string)

		By("deleting the VM")
		request = fmt.Sprintf(`{
			  "method": "delete_vm",
			  "arguments": ["%v"]
			}`, vmCID)
		assertSucceeds(request)
	})

	It("executes the VM lifecycle with scopes", func() {
		By("creating a VM")
		var vmCID string
		request := fmt.Sprintf(`{
			  "method": "create_vm",
			  "arguments": [
				"agent",
				"%v",
				{
				  "machine_type": "n1-standard-1",
				  "service_scopes": ["devstorage.read_write"]
				},
				{
				  "default": {
					"type": "dynamic",
					"cloud_properties": {
					  "tags": ["integration-delete"],
					  "network_name": "%v"
					}
				  }
				},
				[],
				{}
			  ]
			}`, existingStemcell, networkName)
		vmCID = assertSucceedsWithResult(request).(string)

		By("deleting the VM")
		request = fmt.Sprintf(`{
			  "method": "delete_vm",
			  "arguments": ["%v"]
			}`, vmCID)
		assertSucceeds(request)
	})

	It("executes the VM lifecycle with a backend service", func() {
		justInstances := func(ig *compute.InstanceGroupsListInstances) []string {
			instances := make([]string, len(ig.Items))
			for _, i := range ig.Items {
				instances = append(instances, i.Instance)
			}
			return instances
		}
		By("creating a VM")
		var vmCID string
		request := fmt.Sprintf(`{
			  "method": "create_vm",
			  "arguments": [
				"agent",
				"%v",
				{
				  "machine_type": "n1-standard-1",
				  "backend_service": "%v"
				},
				{
				  "default": {
					"type": "dynamic",
					"cloud_properties": {
					  "tags": ["integration-delete"],
					  "network_name": "%v"
					}
				  }
				},
				[],
				{}
			  ]
			}`, existingStemcell, backendService, networkName)
		vmCID = assertSucceedsWithResult(request).(string)

		ig, err := computeService.InstanceGroups.ListInstances(googleProject, zone, instanceGroup, &compute.InstanceGroupsListInstancesRequest{InstanceState: "RUNNING"}).Do()
		Expect(err).ToNot(HaveOccurred())
		Expect(justInstances(ig)).To(ContainElement(ContainSubstring(vmCID)))

		By("deleting the VM and confirming its removal from backend service instance group")
		toggleAsyncDelete()
		request = fmt.Sprintf(`{
			  "method": "delete_vm",
			  "arguments": ["%v"]
			}`, vmCID)
		assertSucceeds(request)
		toggleAsyncDelete()
		ig, err = computeService.InstanceGroups.ListInstances(googleProject, zone, instanceGroup, &compute.InstanceGroupsListInstancesRequest{InstanceState: "RUNNING"}).Do()
		Expect(err).ToNot(HaveOccurred())
		Expect(justInstances(ig)).ToNot(ContainElement(ContainSubstring(vmCID)))

	})

})
