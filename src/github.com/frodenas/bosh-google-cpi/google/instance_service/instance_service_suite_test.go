package instance_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestInstanceService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Instance Service Suite")
}
