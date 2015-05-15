package fakes

import (
	"github.com/frodenas/bosh-google-cpi/google/client"
)

func NewFakeGoogleClient() client.GoogleClient { return client.GoogleClient{} }
