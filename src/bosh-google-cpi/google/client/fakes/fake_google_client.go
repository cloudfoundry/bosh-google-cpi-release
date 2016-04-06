package fakes

import (
	"bosh-google-cpi/google/client"
)

func NewFakeGoogleClient() client.GoogleClient { return client.GoogleClient{} }
