package gdisktype_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestGoogleDiskTypeService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Google Disk Type Service Suite")
}
