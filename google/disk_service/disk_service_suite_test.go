package disk_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestDiskService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Disk Service Suite")
}
