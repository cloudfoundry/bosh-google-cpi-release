package client

import (
	"bytes"
	"io"
	"net"
	"net/http"
	"time"

	boshlog "github.com/cloudfoundry/bosh-utils/logger"
)

const (
	defaultFirstRetrySleep = 50 * time.Millisecond
	retryLogTag            = "RetryTransport"
)

// A function that will modify the request before it is made
type RequestModifier func(req *http.Request)

type RetryTransport struct {
	MaxRetries      int
	FirstRetrySleep time.Duration
	Base            http.RoundTripper
	RequestModifier RequestModifier
	logger          boshlog.Logger
}

func (rt *RetryTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return rt.try(req)
}

func (rt *RetryTransport) try(req *http.Request) (resp *http.Response, err error) {
	if rt.FirstRetrySleep == 0 {
		rt.FirstRetrySleep = defaultFirstRetrySleep
	}

	var body []byte

	if rt.RequestModifier != nil {
		rt.RequestModifier(req)
	}

	// Save the req body for future retries as it will be read and closed
	// by Base.RoundTrip.
	if req.Body != nil {
		body, err = io.ReadAll(req.Body)
		if err != nil {
			return
		}
	}
	var sleep func()
	for try := 0; try <= rt.MaxRetries; try++ {
		r := bytes.NewReader(body)
		req.Body = io.NopCloser(r)
		resp, err = rt.Base.RoundTrip(req)

		sleep = func() {
			d := rt.FirstRetrySleep << uint64(try)
			rt.logger.Info(retryLogTag, "Retrying request (%d/%d) after %s", try, rt.MaxRetries, d)
			time.Sleep(d)
		}

		// Retry on net.Error
		switch err.(type) { //nolint:staticcheck
		case net.Error:
			if err.(net.Error).Temporary() || err.(net.Error).Timeout() { //nolint:staticcheck
				rt.logger.Info(retryLogTag, "net.Error is retryable: %s", err.Error())
				sleep()
				continue
			}
			rt.logger.Info(retryLogTag, "net.Error was not retryable: %s", err.Error())
			return
		case error:
			rt.logger.Info(retryLogTag, "Error was not retryable: %s", err.Error())
			return
		}

		// Retry on status code >= 500
		if resp.StatusCode >= 500 {
			sleep()
			continue
		}
		return
	}
	return
}
