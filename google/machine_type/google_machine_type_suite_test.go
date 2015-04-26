package gmachinetype_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestGoogleMachineTypeService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Google Machine Type Service Suite")
}
