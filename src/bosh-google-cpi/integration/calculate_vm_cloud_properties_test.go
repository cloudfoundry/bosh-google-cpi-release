package integration

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Misc", func() {
	Describe("calculate_vm_cloud_properties", func() {
		It("provides a basic match", func() {
			result := assertSucceedsWithResult(`{
				"method": "calculate_vm_cloud_properties",
				"arguments": [{"cpu":1,"ram":1024,"ephemeral_disk_size":1024}]
			}`).(map[string]interface{})

			Expect(result).To(HaveKey("cpu"))
			Expect(result).To(HaveKey("ram"))
			Expect(result).To(HaveKey("root_disk_size_gb"))
		})
	})
})
