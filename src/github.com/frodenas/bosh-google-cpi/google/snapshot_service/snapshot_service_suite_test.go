package snapshot_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestSnapshotService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Snapshot Service Suite")
}
