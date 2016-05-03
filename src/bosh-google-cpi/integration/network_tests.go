package integration

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"google.golang.org/api/compute/v1"
)

var _ = Describe("Network", func() {
	It("can create a VM with network IP forwarding enabled", func() {
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

		request = fmt.Sprintf(`{
			  "method": "delete_vm",
			  "arguments": ["%v"]
			}`, vmCID)
		assertSucceeds(request)
	})

	It("can create a VM with network tags", func() {
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
					  "tags": ["tag1", "tag2"]
					}
				  }
				},
				[],
				{}
			  ]
			}`, existingStemcell, networkName)
		vmCID = assertSucceedsWithResult(request).(string)
		assertValidVM(vmCID, func(instance *compute.Instance) {
			Expect(instance.Tags.Items).To(ConsistOf("tag1", "tag2"))
		})

		request = fmt.Sprintf(`{
			  "method": "delete_vm",
			  "arguments": ["%v"]
			}`, vmCID)
		assertSucceeds(request)
	})

	It("can create a VM with an ephemeral external IP", func() {
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
		request = fmt.Sprintf(`{
			  "method": "delete_vm",
			  "arguments": ["%v"]
			}`, vmCID)
		assertSucceeds(request)
	})

	It("can create a VM with a static external IP", func() {
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
		request = fmt.Sprintf(`{
			  "method": "delete_vm",
			  "arguments": ["%v"]
			}`, vmCID)
		assertSucceeds(request)
	})

	It("can create a VM in a subnet", func() {
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
		request = fmt.Sprintf(`{
			  "method": "delete_vm",
			  "arguments": ["%v"]
			}`, vmCID)
		assertSucceeds(request)
	})

	It("can create a VM with a static private IP", func() {
		ip := "192.168.100.101"
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
		request = fmt.Sprintf(`{
			  "method": "delete_vm",
			  "arguments": ["%v"]
			}`, vmCID)
		assertSucceeds(request)
	})

	It("can create a VM with a static private IP and ephemeral public IP", func() {
		ip := "192.168.100.102"
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
		request = fmt.Sprintf(`{
			  "method": "delete_vm",
			  "arguments": ["%v"]
			}`, vmCID)
		assertSucceeds(request)
	})

	It("can create a VM with a static private IP and static public IP", func() {
		ip := "192.168.100.103"
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
		request = fmt.Sprintf(`{
			  "method": "delete_vm",
			  "arguments": ["%v"]
			}`, vmCID)
		assertSucceeds(request)
	})

	Context("Target Pools", func() {
		It("can create a VM in a target pool", func() {
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
		})

		It("can delete a VM and remove its target pool association", func() {
			request := fmt.Sprintf(`{
			  "method": "delete_vm",
			  "arguments": ["%v"]
			}`, vmCID)
			assertSucceeds(request)

			tp, err := computeService.TargetPools.Get(googleProject, region, targetPool).Do()
			Expect(err).ToNot(HaveOccurred())
			Expect(tp.Instances).ToNot(ContainElement(ContainSubstring(vmCID)))

		})
	})

	Context("Instance Groups", func() {
		justInstances := func(ig *compute.InstanceGroupsListInstances) []string {
			instances := make([]string, len(ig.Items))
			for _, i := range ig.Items {
				instances = append(instances, i.Instance)
			}
			return instances
		}
		It("can create a VM in an instance group", func() {
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
					  "instance_group": "%v"
					}
				  }
				},
				[],
				{}
			  ]
			}`, existingStemcell, networkName, instanceGroup)
			vmCID = assertSucceedsWithResult(request).(string)
			ig, err := computeService.InstanceGroups.ListInstances(googleProject, zone, instanceGroup, &compute.InstanceGroupsListInstancesRequest{InstanceState: "RUNNING"}).Do()
			Expect(err).ToNot(HaveOccurred())
			Expect(justInstances(ig)).To(ContainElement(ContainSubstring(vmCID)))
		})

		It("can delete a VM and remove its instance group association", func() {
			request := fmt.Sprintf(`{
			  "method": "delete_vm",
			  "arguments": ["%v"]
			}`, vmCID)
			assertSucceeds(request)

			ig, err := computeService.InstanceGroups.ListInstances(googleProject, zone, instanceGroup, &compute.InstanceGroupsListInstancesRequest{InstanceState: "RUNNING"}).Do()
			Expect(err).ToNot(HaveOccurred())
			Expect(justInstances(ig)).ToNot(ContainElement(ContainSubstring(vmCID)))

		})
	})
})
