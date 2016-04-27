package integration

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func assertSucceeds(request string) {
	response, err := execCPI(request)
	Expect(err).ToNot(HaveOccurred())
	Expect(response.Error).To(BeNil())
}

func assertSucceedsWithResult(request string) interface{} {
	GinkgoWriter.Write([]byte(fmt.Sprintf("CPI request: %v\n", request)))
	response, err := execCPI(request)
	GinkgoWriter.Write([]byte(fmt.Sprintf("CPI response: %v\n", response)))
	GinkgoWriter.Write([]byte(fmt.Sprintf("CPI error: %v\n", err)))
	Expect(err).ToNot(HaveOccurred())
	Expect(response.Error).To(BeNil())
	Expect(response.Result).ToNot(BeNil())
	return response.Result
}

func toStringArray(raw []interface{}) []string {
	strings := make([]string, len(raw), len(raw))
	for i := range raw {
		strings[i] = raw[i].(string)
	}
	return strings
}
