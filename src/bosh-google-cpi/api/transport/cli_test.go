package transport_test

import (
	"errors"
	"io"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "bosh-google-cpi/api/transport"

	boshlog "github.com/cloudfoundry/bosh-utils/logger"

	fakedisp "bosh-google-cpi/api/dispatcher/fakes"
)

type FakeReader struct {
	ReadBytes []byte
	ReadErr   error
}

func (r *FakeReader) Read(bytes []byte) (int, error) {
	copy(bytes, r.ReadBytes)

	if r.ReadErr != nil {
		return len(r.ReadBytes), r.ReadErr
	}

	return len(r.ReadBytes), io.EOF
}

type FakeWriter struct {
	WriteBytes []byte
	WriteErr   error
}

func (w *FakeWriter) Write(b []byte) (int, error) {
	w.WriteBytes = b
	return len(b), w.WriteErr
}

var _ = Describe("CLI", func() {
	var (
		in         *FakeReader // io.Reader
		out        *FakeWriter // io.Writer
		dispatcher *fakedisp.FakeDispatcher
		logger     boshlog.Logger
		cli        CLI
	)

	BeforeEach(func() {
		in = &FakeReader{}
		out = &FakeWriter{}
		dispatcher = &fakedisp.FakeDispatcher{}
		logger = boshlog.NewLogger(boshlog.LevelNone)
		cli = NewCLI(in, out, dispatcher, logger)
	})

	Describe("ServeOnce", func() {
		It("reads request from in and writes response to out", func() {
			in.ReadBytes = []byte("fake-bytes-in")

			dispatcher.DispatchRespBytes = []byte("fake-bytes-out")

			err := cli.ServeOnce()
			Expect(err).ToNot(HaveOccurred())

			Expect(dispatcher.DispatchReqBytes).To(Equal([]byte("fake-bytes-in")))

			Expect(out.WriteBytes).To(Equal([]byte("fake-bytes-out")))
		})

		It("returns error if reading request from in fails", func() {
			in.ReadErr = errors.New("fake-read-err")

			err := cli.ServeOnce()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("fake-read-err"))
		})

		It("returns error if writing response to out fails", func() {
			out.WriteErr = errors.New("fake-write-err")

			err := cli.ServeOnce()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("fake-write-err"))
		})
	})
})
