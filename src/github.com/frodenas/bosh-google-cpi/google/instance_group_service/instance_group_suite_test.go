package instancegroup_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestInstanceGroupService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Instance Group Service Suite")
}
