package acceleratortype_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestAcceleratorTypeService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Accelerator Type Service Suite")
}
