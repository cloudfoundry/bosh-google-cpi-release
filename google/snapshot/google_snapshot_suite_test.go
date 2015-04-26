package gsnapshot_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestGoogleSnapshotService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Google Snapshot Service Suite")
}
