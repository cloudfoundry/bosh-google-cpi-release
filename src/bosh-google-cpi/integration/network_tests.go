package integration

import (
	"fmt"

	"google.golang.org/api/compute/v1"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Network", func() {
	It("can create a VM with network IP forwarding enabled", func() {
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
					  "network_name": "%v",
					  "ip_forwarding": true
					}
				  }
				},
				[],
				{}
			  ]
			}`, existingStemcell, networkName)
		vmCID = assertSucceedsWithResult(request).(string)
		assertValidVM(vmCID, func(instance *compute.Instance) {
			Expect(instance.CanIpForward).To(BeTrue())
		})

		By("deleting the VM")
		request = fmt.Sprintf(`{
			  "method": "delete_vm",
			  "arguments": ["%v"]
			}`, vmCID)
		assertSucceeds(request)
	})

	It("can create a VM with network tags", func() {
		By("creating a VM")
		var vmCID string
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
					  "network_name": "%v",
					  "tags": ["integration-delete", "tag1", "tag2"]
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

	It("can create a VM with an ephemeral external IP", func() {
		By("creating a VM")
		var vmCID string
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
					  "network_name": "%v",
					  "ephemeral_external_ip": true
					}
				  }
				},
				[],
				{}
			  ]
			}`, existingStemcell, networkName)
		vmCID = assertSucceedsWithResult(request).(string)
		assertValidVM(vmCID, func(instance *compute.Instance) {
			Expect(instance.NetworkInterfaces[0].AccessConfigs[0].NatIP).ToNot(BeEmpty())
		})

		By("deleting the VM")
		request = fmt.Sprintf(`{
			  "method": "delete_vm",
			  "arguments": ["%v"]
			}`, vmCID)
		assertSucceeds(request)
	})

	It("can create a VM with a static external IP", func() {
		By("creating a VM")
		var vmCID string
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
					  "network_name": "%v",
					  "ephemeral_external_ip": true
					}
				  },
				  "vip": {
					"type": "vip",
					"ip": "%v"
				  }
				},
				[],
				{}
			  ]
			}`, existingStemcell, networkName, externalStaticIP)
		vmCID = assertSucceedsWithResult(request).(string)
		assertValidVM(vmCID, func(instance *compute.Instance) {
			Expect(instance.NetworkInterfaces[0].AccessConfigs[0].NatIP).ToNot(BeEmpty())
			Expect(instance.NetworkInterfaces[0].AccessConfigs[0].NatIP).To(Equal(externalStaticIP))
		})

		By("deleting the VM")
		request = fmt.Sprintf(`{
			  "method": "delete_vm",
			  "arguments": ["%v"]
			}`, vmCID)
		toggleAsyncDelete()
		assertSucceeds(request)
		toggleAsyncDelete()
	})

	It("can create a VM in a subnet", func() {
		By("creating a VM")
		var vmCID string
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
					  "network_name": "%v",
					  "subnetwork_name": "%v"
					}
				  }
				},
				[],
				{}
			  ]
			}`, existingStemcell, customNetworkName, customSubnetworkName)
		vmCID = assertSucceedsWithResult(request).(string)
		assertValidVM(vmCID, func(instance *compute.Instance) {
			Expect(instance.NetworkInterfaces[0].Network).To(ContainSubstring(customNetworkName))
			Expect(instance.NetworkInterfaces[0].Subnetwork).To(ContainSubstring(customSubnetworkName))
		})

		By("deleting the VM")
		request = fmt.Sprintf(`{
			  "method": "delete_vm",
			  "arguments": ["%v"]
			}`, vmCID)
		assertSucceeds(request)
	})

	It("can create a VM with a static private IP", func() {
		By("creating a VM")
		var vmCID string
		ip := <-ips
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
					"type": "manual",
					"ip": "%v",
					"cloud_properties": {
					  "tags": ["integration-delete"],
					  "network_name": "%v",
					  "subnetwork_name": "%v"
					}
				  }
				},
				[],
				{}
			  ]
			}`, existingStemcell, ip, customNetworkName, customSubnetworkName)
		vmCID = assertSucceedsWithResult(request).(string)
		assertValidVM(vmCID, func(instance *compute.Instance) {
			Expect(instance.NetworkInterfaces[0].NetworkIP).To(Equal(ip))
			Expect(instance.NetworkInterfaces[0].AccessConfigs).To(BeEmpty())
		})

		By("deleting the VM")
		request = fmt.Sprintf(`{
			  "method": "delete_vm",
			  "arguments": ["%v"]
			}`, vmCID)
		assertSucceeds(request)
		ips <- ip
	})

	It("can create a VM with a static private IP and ephemeral public IP", func() {
		By("creating a VM")
		ip := <-ips
		var vmCID string
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
					"type": "manual",
					"ip": "%v",
					"cloud_properties": {
					  "tags": ["integration-delete"],
					  "network_name": "%v",
					  "subnetwork_name": "%v",
					  "ephemeral_external_ip": true
					}
				  }
				},
				[],
				{}
			  ]
			}`, existingStemcell, ip, customNetworkName, customSubnetworkName)
		vmCID = assertSucceedsWithResult(request).(string)
		assertValidVM(vmCID, func(instance *compute.Instance) {
			Expect(instance.NetworkInterfaces[0].NetworkIP).To(Equal(ip))
			Expect(instance.NetworkInterfaces[0].AccessConfigs[0].NatIP).ToNot(BeEmpty())
		})

		By("deleting the VM")
		request = fmt.Sprintf(`{
			  "method": "delete_vm",
			  "arguments": ["%v"]
			}`, vmCID)
		assertSucceeds(request)
		ips <- ip
	})

	It("can create a VM with a static private IP and static public IP", func() {
		By("creating a VM")
		ip := <-ips
		var vmCID string
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
					"type": "manual",
					"ip": "%v",
					"cloud_properties": {
					  "tags": ["integration-delete"],
					  "network_name": "%v",
					  "subnetwork_name": "%v",
					  "ephemeral_external_ip": true
					}
				  },
				  "vip": {
					"type": "vip",
					"ip": "%v"
				  }
				},
				[],
				{}
			  ]
			}`, existingStemcell, ip, customNetworkName, customSubnetworkName, externalStaticIP)
		vmCID = assertSucceedsWithResult(request).(string)
		assertValidVM(vmCID, func(instance *compute.Instance) {
			Expect(instance.NetworkInterfaces[0].NetworkIP).To(Equal(ip))
			Expect(instance.NetworkInterfaces[0].AccessConfigs[0].NatIP).To(Equal(externalStaticIP))
		})

		By("deleting the VM")
		request = fmt.Sprintf(`{
			  "method": "delete_vm",
			  "arguments": ["%v"]
			}`, vmCID)
		toggleAsyncDelete()
		assertSucceeds(request)
		ips <- ip
		toggleAsyncDelete()
	})

	It("execute the creation and deleting of a VM in a target pool", func() {
		By("creating a VM")
		var vmCID string
		request := fmt.Sprintf(`{
			  "method": "create_vm",
			  "arguments": [
				"agent",
				"%v",
				{
				  "machine_type": "n1-standard-1",
				  "target_pool": "%v"
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
			}`, existingStemcell, targetPool, networkName)
		vmCID = assertSucceedsWithResult(request).(string)
		tp, err := computeService.TargetPools.Get(googleProject, region, targetPool).Do()
		Expect(err).ToNot(HaveOccurred())
		Expect(tp.Instances).ToNot(BeEmpty())
		Expect(tp.Instances).To(ContainElement(ContainSubstring(vmCID)))

		By("deleting the VM and confirming its removal from the target pool")
		request = fmt.Sprintf(`{
			  "method": "delete_vm",
			  "arguments": ["%v"]
			}`, vmCID)
		assertSucceeds(request)

		tp, err = computeService.TargetPools.Get(googleProject, region, targetPool).Do()
		Expect(err).ToNot(HaveOccurred())
		Expect(tp.Instances).ToNot(ContainElement(ContainSubstring(vmCID)))
	})
})
