package integration

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("VM", func() {
	Context("Lifecycle", func() {
		It("can create a simple VM", func() {
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
					  "network_name": "%v"
					}
				  }
				},
				[],
				{}
			  ]
			}`, existingStemcell, networkName)
			vmCID = assertSucceedsWithResult(request).(string)
		})

		It("can confirm existence of an existing VM", func() {
			request := fmt.Sprintf(`{
			  "method": "has_vm",
			  "arguments": ["%v"]
			}`, vmCID)
			exists := assertSucceedsWithResult(request).(bool)
			Expect(exists).To(Equal(true))
		})

		It("can reboot an existing VM", func() {
			request := fmt.Sprintf(`{
			  "method": "reboot_vm",
			  "arguments": ["%v"]
			}`, vmCID)
			assertSucceeds(request)
		})

		It("can delete a VM", func() {
			request := fmt.Sprintf(`{
			  "method": "delete_vm",
			  "arguments": ["%v"]
			}`, vmCID)
			assertSucceeds(request)
		})
	})

	Context("can create a VM with existing disk attachment hints", func() {
		var request, diskCID2 string
		It("creates the disks", func() {
			request = fmt.Sprintf(`{
			  "method": "create_disk",
			  "arguments": [32768, {}, ""]
			}`)
			diskCID = assertSucceedsWithResult(request).(string)
			diskCID2 = assertSucceedsWithResult(request).(string)
		})

		It("creates a VM with the disks", func() {
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
					  "network_name": "%v"
					}
				  }
				},
				["%v", "%v"],
				{}
			  ]
			}`, existingStemcell, networkName, diskCID, diskCID2)
			vmCID = assertSucceedsWithResult(request).(string)
		})

		It("deletes the disks", func() {
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
		})
		It("deletes the VM", func() {
			request = fmt.Sprintf(`{
			  "method": "delete_vm",
			  "arguments": ["%v"]
			}`, vmCID)
			assertSucceeds(request)
		})
	})

	It("can create a VM with custom machine type", func() {
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
					  "network_name": "%v"
					}
				  }
				},
				[],
				{}
			  ]
			}`, existingStemcell, networkName)
		vmCID = assertSucceedsWithResult(request).(string)
		request = fmt.Sprintf(`{
			  "method": "delete_vm",
			  "arguments": ["%v"]
			}`, vmCID)
		assertSucceeds(request)
	})

	It("can create a VM in a particular zone", func() {
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
					  "network_name": "%v"
					}
				  }
				},
				[],
				{}
			  ]
			}`, existingStemcell, networkName)
		vmCID = assertSucceedsWithResult(request).(string)
		request = fmt.Sprintf(`{
			  "method": "delete_vm",
			  "arguments": ["%v"]
			}`, vmCID)
		assertSucceeds(request)
	})
	Context("scheduling params", func() {
		It("can create a VM with automatic restart disabled", func() {
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
					  "network_name": "%v"
					}
				  }
				},
				[],
				{}
			  ]
			}`, existingStemcell, networkName)
			vmCID = assertSucceedsWithResult(request).(string)
			request = fmt.Sprintf(`{
			  "method": "delete_vm",
			  "arguments": ["%v"]
			}`, vmCID)
			assertSucceeds(request)
		})
		It("can create a VM with OnHostMaintenance modified", func() {
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
					  "network_name": "%v"
					}
				  }
				},
				[],
				{}
			  ]
			}`, existingStemcell, networkName)
			vmCID = assertSucceedsWithResult(request).(string)
			request = fmt.Sprintf(`{
			  "method": "delete_vm",
			  "arguments": ["%v"]
			}`, vmCID)
			assertSucceeds(request)
		})
		It("can create a preemtible VM", func() {
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
					  "network_name": "%v"
					}
				  }
				},
				[],
				{}
			  ]
			}`, existingStemcell, networkName)
			vmCID = assertSucceedsWithResult(request).(string)
			request = fmt.Sprintf(`{
			  "method": "delete_vm",
			  "arguments": ["%v"]
			}`, vmCID)
			assertSucceeds(request)
		})
	})
	It("can create a VM with scopes", func() {
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
					  "network_name": "%v"
					}
				  }
				},
				[],
				{}
			  ]
			}`, existingStemcell, networkName)
		vmCID = assertSucceedsWithResult(request).(string)
		request = fmt.Sprintf(`{
			  "method": "delete_vm",
			  "arguments": ["%v"]
			}`, vmCID)
		assertSucceeds(request)
	})

	FContext("Network customization", func() {
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
			request = fmt.Sprintf(`{
			  "method": "delete_vm",
			  "arguments": ["%v"]
			}`, vmCID)
			assertSucceeds(request)
		})
	})
})
