package integration

import (
	"fmt"

	"bosh-google-cpi/util"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	computebeta "google.golang.org/api/compute/v0.beta"
	"google.golang.org/api/compute/v1"

	"testing"
)

var computeService *compute.Service
var computeServiceB *computebeta.Service

func TestIntegration(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Integration Suite")
}

var _ = SynchronizedBeforeSuite(func() []byte {
	// Clean any straggler VMs
	cleanVMs()

	request := fmt.Sprintf(`{
			  "method": "create_stemcell",
			  "arguments": ["%s", {
				  "name": "bosh-google-kvm-ubuntu-trusty",
				  "version": "%s",
				  "infrastructure": "google"
				}]
			}`, stemcellFile, stemcellVersion)
	stemcell := assertSucceedsWithResult(request).(string)

	ips = make(chan string, len(ipAddrs))

	// Parse IP addresses to be used and put on a chan
	for _, addr := range ipAddrs {
		ips <- addr
	}

	return []byte(stemcell)
}, func(data []byte) {
	// Ensure stemcell was initialized
	existingStemcell = string(data)
	Expect(existingStemcell).ToNot(BeEmpty())

	// Required env vars
	Expect(googleProject).ToNot(Equal(""), "GOOGLE_PROJECT must be set")
	Expect(externalStaticIP).ToNot(Equal(""), "EXTERNAL_STATIC_IP must be set")
	Expect(serviceAccount).ToNot(Equal(""), "SERVICE_ACCOUNT must be set")

	// Initialize a compute API client
	ctx := context.Background()
	client, err := google.DefaultClient(ctx, compute.ComputeScope)
	Expect(err).To(BeNil())
	computeService, err = compute.New(client)
	computeServiceB, err = computebeta.New(client)
	Expect(err).To(BeNil())
})

var _ = SynchronizedAfterSuite(func() {}, func() {
	cleanVMs()
	request := fmt.Sprintf(`{
			  "method": "delete_stemcell",
			  "arguments": ["%v"]
			}`, existingStemcell)

	response, err := execCPI(request)
	Expect(err).ToNot(HaveOccurred())
	Expect(response.Error).To(BeNil())
	Expect(response.Result).To(BeNil())
})

func cleanVMs() {
	// Initialize a compute API client
	ctx := context.Background()
	client, err := google.DefaultClient(ctx, compute.ComputeScope)
	Expect(err).To(BeNil())
	computeService, err := compute.New(client)
	Expect(err).To(BeNil())

	// Clean up any VMs left behind from failed tests. Instances with the 'integration-delete' tag will be deleted.
	var pageToken string
	toDelete := make([]*compute.Instance, 0)
	GinkgoWriter.Write([]byte("Looking for VMs with 'integration-delete' tag. Matches will be deleted\n"))
	for {
		// Clean up VMs with 'integration-delete' tag
		listCall := computeService.Instances.AggregatedList(googleProject)
		listCall.PageToken(pageToken)
		aggregatedList, err := listCall.Do()
		Expect(err).To(BeNil())
		for _, list := range aggregatedList.Items {
			for _, instance := range list.Instances {
				for _, tag := range instance.Tags.Items {
					if tag == "integration-delete" {
						toDelete = append(toDelete, instance)
					}
				}
			}
		}
		if aggregatedList.NextPageToken == "" {
			break
		}
		pageToken = aggregatedList.NextPageToken
	}

	for _, vm := range toDelete {
		GinkgoWriter.Write([]byte(fmt.Sprintf("Deleting VM %v\n", vm.Name)))
		_, err := computeService.Instances.Delete(googleProject, util.ResourceSplitter(vm.Zone), vm.Name).Do()
		Expect(err).ToNot(HaveOccurred())
	}
}
