package action_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "bosh-google-cpi/action"

	snapshotfakes "bosh-google-cpi/google/snapshot_service/fakes"
)

var _ = Describe("DeleteSnapshot", func() {
	var (
		err error

		snapshotService *snapshotfakes.FakeSnapshotService

		deleteSnapshot DeleteSnapshot
	)

	BeforeEach(func() {
		snapshotService = &snapshotfakes.FakeSnapshotService{}
		deleteSnapshot = NewDeleteSnapshot(snapshotService)
	})

	Describe("Run", func() {
		It("deletes the snapshot", func() {
			_, err = deleteSnapshot.Run("fake-snapshot-id")
			Expect(err).NotTo(HaveOccurred())
			Expect(snapshotService.DeleteCalled).To(BeTrue())
		})

		It("returns an error if snapshotService delete call returns an error", func() {
			snapshotService.DeleteErr = errors.New("fake-snapshot-service-error")

			_, err = deleteSnapshot.Run("fake-snapshot-id")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("fake-snapshot-service-error"))
			Expect(snapshotService.DeleteCalled).To(BeTrue())
		})
	})
})
