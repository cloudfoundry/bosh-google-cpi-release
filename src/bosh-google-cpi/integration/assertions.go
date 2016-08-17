package integration

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"google.golang.org/api/compute/v1"
)

func assertSucceeds(request string) {
	response, err := execCPI(request)
	Expect(err).ToNot(HaveOccurred())
	Expect(response.Error).To(BeNil())
}

func assertFails(request string) error {
	response, _ := execCPI(request)
	Expect(response.Error).ToNot(BeNil())
	return response.Error
}

func assertSucceedsWithResult(request string) interface{} {
	response, err := execCPI(request)
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

func assertValidVM(id string, valFunc func(*compute.Instance)) {
	listCall := computeService.Instances.AggregatedList(googleProject)
	listCall.Filter(fmt.Sprintf("name eq %v", id))
	aggregatedList, err := listCall.Do()
	Expect(err).To(BeNil())
	for _, list := range aggregatedList.Items {
		if len(list.Instances) == 1 {
			valFunc(list.Instances[0])
			return
		}
	}
	Fail(fmt.Sprintf("Instance %q not found\n", id))
}
