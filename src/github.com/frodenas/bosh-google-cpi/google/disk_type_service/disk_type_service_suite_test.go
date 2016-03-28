package disktype_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestDiskTypeService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Disk Type Service Suite")
}
