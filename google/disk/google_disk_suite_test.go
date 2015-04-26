package gdisk_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestGoogleDiskService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Google Disk Service Suite")
}
