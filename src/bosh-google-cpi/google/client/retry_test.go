package client

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"net/http"
	"net/http/httptest"
)

var _ = Describe("RetryTransport", func() {
	Describe("Validate", func() {
		It("It retries the maximum number of times and then fails", func() {
			maxRetries := 3
			try := 0

			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				try++
				w.WriteHeader(http.StatusServiceUnavailable)
			}))
			defer ts.Close()

			client := http.Client{
				Transport: &RetryTransport{
					Base:       http.DefaultTransport,
					MaxRetries: maxRetries,
				},
			}

			res, err := client.Get(ts.URL)
			Expect(try).To(Equal(maxRetries + 1))
			Expect(err).ToNot(HaveOccurred())
			Expect(res.StatusCode).To(Equal(http.StatusServiceUnavailable))
		})
		It("It retries the maximum number and succeeds on the last try", func() {
			maxRetries := 3
			try := 0

			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if try == maxRetries {
					w.WriteHeader(http.StatusOK)
				} else {
					w.WriteHeader(http.StatusServiceUnavailable)
				}
				try++
			}))
			defer ts.Close()

			client := http.Client{
				Transport: &RetryTransport{
					Base:       http.DefaultTransport,
					MaxRetries: maxRetries,
				},
			}

			res, err := client.Get(ts.URL)
			Expect(try).To(Equal(maxRetries + 1))
			Expect(err).ToNot(HaveOccurred())
			Expect(res.StatusCode).To(Equal(http.StatusOK))
		})
		It("It retries zero times and succeeds", func() {
			maxRetries := 0
			try := 0

			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusServiceUnavailable)
				try++
			}))
			defer ts.Close()

			client := http.Client{
				Transport: &RetryTransport{
					Base:       http.DefaultTransport,
					MaxRetries: maxRetries,
				},
			}
			res, err := client.Get(ts.URL)
			Expect(try).To(Equal(maxRetries + 1))
			Expect(err).ToNot(HaveOccurred())
			Expect(res.StatusCode).To(Equal(http.StatusServiceUnavailable))
		})
		It("It retries zero times and succeeds", func() {
			maxRetries := 0
			try := 0

			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if try == maxRetries {
					w.WriteHeader(http.StatusOK)
				} else {
					w.WriteHeader(http.StatusServiceUnavailable)
				}
				try++
			}))
			defer ts.Close()

			client := http.Client{
				Transport: &RetryTransport{
					Base:       http.DefaultTransport,
					MaxRetries: maxRetries,
				},
			}
			res, err := client.Get(ts.URL)
			Expect(try).To(Equal(maxRetries + 1))
			Expect(err).ToNot(HaveOccurred())
			Expect(res.StatusCode).To(Equal(http.StatusOK))
		})

	})
})
