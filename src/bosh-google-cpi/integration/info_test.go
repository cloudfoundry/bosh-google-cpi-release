package integration

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Misc", func() {
	Describe("info", func() {
		It("provides stemcell_formats", func() {
			result := assertSucceedsWithResult(`{
				"method": "info",
				"arguments": []
			}`).(map[string]interface{})

			Expect(result).To(HaveKey("stemcell_formats"))
		})
	})
})
