package machinetype_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestMachineTypeService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Machine Type Service Suite")
}
