package gclient_fakes

import (
	"google.golang.org/api/compute/v1"
)

type FakeGoogleClient struct {
	project        string
	defaultZone    string
	computeService FakeComputeService
}

func NewFakeGoogleClient(
	project string,
	jsonKey string,
	defaultZone string,
	accessKeyId string,
	secretAccessKey string,
) (FakeGoogleClient, error) {
	var fakeComputeService FakeComputeService
	return FakeGoogleClient{
		project:        project,
		defaultZone:    defaultZone,
		computeService: fakeComputeService,
	}, nil
}

func (c FakeGoogleClient) Project() string {
	return c.project
}

func (c FakeGoogleClient) DefaultZone() string {
	return c.defaultZone
}

func (c FakeGoogleClient) ComputeService() FakeComputeService {
	return c.computeService
}

type FakeComputeService compute.Service

func (c FakeComputeService) DiskTypes() FakeDiskTypes {
	return FakeDiskTypes{}
}

type FakeDiskTypes compute.DiskTypesService

func (c FakeDiskTypes) Get(project string, zone string, id string) FakeDiskTypes {
	return FakeDiskTypes{}
}
