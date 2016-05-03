package integration

import (
	"fmt"

	"bosh-google-cpi/util"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/compute/v1"

	"testing"
)

var computeService *compute.Service

func TestIntegration(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Integration Suite")
}

var _ = BeforeSuite(func() {
	// Required env vars
	Expect(googleProject).ToNot(Equal(""), "GOOGLE_PROJECT must be set")
	Expect(externalStaticIP).ToNot(Equal(""), "EXTERNAL_STATIC_IP must be set")

	// Provision an instance for update tests if it doesn't
	// already exist.
	request := fmt.Sprintf(`{
			  "method": "has_vm",
			  "arguments": ["%v"]
			}`, reusableVMName)
	response, err := execCPI(request)
	Expect(err).To(BeNil())
	exists, ok := response.Result.(bool)
	if ok {
		if !exists {
			GinkgoWriter.Write([]byte("Creating VM that will be reused for several tests.\n"))
			request = fmt.Sprintf(`{
						"method": "create_vm",
						"arguments": [
						  "agent",
						  "%v",
						  {
							"machine_type": "n1-standard-1",
							"name": "%v"
						  },
						  {
							"default": {
							  "type": "dynamic",
							  "cloud_properties": {
								"network_name": "%v"
							  }
							}
						  },
						  [],
						  {}
						]
					  }`, existingStemcell, reusableVMName, networkName)
			assertSucceedsWithResult(request)
		} else {
			GinkgoWriter.Write([]byte("Reusing existing VM\n"))
		}
	}

	// Initialize a compute API client
	ctx := context.Background()
	client, err := google.DefaultClient(ctx, compute.ComputeScope)
	Expect(err).To(BeNil())
	computeService, err = compute.New(client)
	Expect(err).To(BeNil())

	// Clean any straggler VMs
	cleanVMs()

})

var _ = AfterSuite(func() {
	if keepResuableVM == "" {
		GinkgoWriter.Write([]byte("Deleting reusable VM. Set KEEP_REUSABLE_VM to any value to prevent this VM from being terminated after a test run, speeding up future tests.\n"))
		request := fmt.Sprintf(`{
			  "method": "delete_vm",
			  "arguments": ["%v"]
			}`, reusableVMName)
		assertSucceeds(request)
	} else {
		GinkgoWriter.Write([]byte(fmt.Sprintf("The reusable VM named %v will not be deleted because KEEP_REUSABLE_VM is set. This will incur Google Compute Engine usage charges. Manually delete the vm with `gcloud compute instances delete %v` if you do not want the VM to continue running.\n", reusableVMName, reusableVMName)))
	}
	cleanVMs()
})

func cleanVMs() {
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
