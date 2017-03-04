package client

import (
	"bytes"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

// A function that will modify the request before it is made
type RequestModifier func(req *http.Request)

type RetryTransport struct {
	MaxRetries      int
	Base            http.RoundTripper
	RequestModifier RequestModifier
}

func (rt *RetryTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return rt.try(req)
}

func (rt *RetryTransport) try(req *http.Request) (resp *http.Response, err error) {
	var body []byte

	if rt.RequestModifier != nil {
		rt.RequestModifier(req)
	}

	// Save the req body for future retries as it will be read and closed
	// by Base.RoundTrip.
	if req.Body != nil {
		body, err = ioutil.ReadAll(req.Body)
		if err != nil {
			return
		}
	}

	for try := 0; try <= rt.MaxRetries; try++ {
		r := bytes.NewReader(body)
		req.Body = ioutil.NopCloser(r)
		resp, err = rt.Base.RoundTrip(req)

		sleep := func() {
			time.Sleep(200 * time.Millisecond << uint64(try))
		}

		// Retry on net.Error
		switch err.(type) {
		case net.Error:
			if !err.(net.Error).Temporary() {
				return
			}
			sleep()
			continue
		case error:
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
