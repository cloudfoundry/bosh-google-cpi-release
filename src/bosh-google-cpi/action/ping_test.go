package action_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "bosh-google-cpi/action"
)

var _ = Describe("Ping", func() {
	var (
		err  error
		ping Ping
		pong string
	)

	BeforeEach(func() {
		ping = NewPing()
	})

	Describe("Run", func() {
		It("returns pong", func() {
			pong, err = ping.Run()
			Expect(err).NotTo(HaveOccurred())
			Expect(pong).To(Equal("pong"))
		})
	})
})
