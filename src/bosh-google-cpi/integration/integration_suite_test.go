package integration

import (
	"fmt"
	"os"
	"slices"

	"bosh-google-cpi/util"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	cloudresourcemanager "google.golang.org/api/cloudresourcemanager/v1"
	computebeta "google.golang.org/api/compute/v0.beta"
	"google.golang.org/api/compute/v1"
	iam "google.golang.org/api/iam/v1"
	"gopkg.in/yaml.v3"

	"testing"
)

var computeService *compute.Service
var computeServiceB *computebeta.Service

func TestIntegration(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Integration Suite")
}

var _ = SynchronizedBeforeSuite(func() []byte {
	// Ensure tests are run using the documented minimum GCP permissions
	validateMinimumPermissions()

	// Clean any straggler VMs
	cleanVMs()

	request := fmt.Sprintf(`{
			  "method": "create_stemcell",
			  "arguments": ["%s", {
				  "name": "bosh-google-kvm-ubuntu-bionic",
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
	Expect(customServiceAccount).ToNot(Equal(""), "CUSTOM_SERVICE_ACCOUNT must be set")
	Expect(jsonKeyServiceAccount).ToNot(Equal(""), "JSON_KEY_SERVICE_ACCOUNT must be set")

	// Initialize a compute API client
	ctx := context.Background()
	client, err := google.DefaultClient(ctx, compute.ComputeScope)
	Expect(err).To(BeNil())
	computeService, err = compute.New(client)
	Expect(err).To(BeNil())
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

type exampleRoleYaml struct {
	IncludedPermissions []string `yaml:"included_permissions"`
}

func validateMinimumPermissions() {
	ctx := context.Background()

	// Get role names attached to current service account
	resourceManager, err := cloudresourcemanager.NewService(ctx)
	Expect(err).ToNot(HaveOccurred())

	policy, err := resourceManager.Projects.GetIamPolicy(googleProject, &cloudresourcemanager.GetIamPolicyRequest{}).Do()
	Expect(err).ToNot(HaveOccurred())

	var roleNames []string
	for _, binding := range policy.Bindings {
		if slices.Contains(binding.Members, fmt.Sprintf("serviceAccount:%s", jsonKeyServiceAccount)) {
			roleNames = append(roleNames, binding.Role)
		}
	}

	// Get permissions for each role
	iamService, err := iam.NewService(ctx)
	Expect(err).ToNot(HaveOccurred())

	var actualPermissions []string

	for _, roleName := range roleNames {
		role, err := iamService.Projects.Roles.Get(roleName).Do()
		if err != nil {
			role, err = iamService.Roles.Get(roleName).Do()
		}
		Expect(err).ToNot(HaveOccurred())

		actualPermissions = append(actualPermissions, role.IncludedPermissions...)
	}

	//  Compare actual permissions against the permissions in the example role file
	var expectedPermissions []string
	var exampleRole *exampleRoleYaml

	yamlContents, err := os.ReadFile("../../../docs/bosh-director-role.yml")
	Expect(err).ToNot(HaveOccurred())
	err = yaml.Unmarshal(yamlContents, &exampleRole)
	Expect(err).ToNot(HaveOccurred())

	// permissions in addition to minimal set of permissions(docs/bosh-director-role.yml) required to run CPI
	// these permissions are for testing specific use cases e.g. "create VM with an accelerator", "create a VM with a static private IP"
	additionalPermissionsForIntegrationTests := []string{
		// Permissions needed to check permissions for this test
		"resourcemanager.projects.getIamPolicy",
		"iam.roles.get",
		// Permission needed when using accelerators property
		"compute.acceleratorTypes.get",
		// Permission needed when using service_account or service_scopes properties
		"compute.instances.setServiceAccount",
		// Permissions associated with the role iam.serviceAccountUser, needed when using the service_account or service_scopes properties.
		"iam.serviceAccounts.actAs",
		"iam.serviceAccounts.get",
		"iam.serviceAccounts.list",
		"resourcemanager.projects.get",
		"resourcemanager.projects.list",
	}

	expectedPermissions = append(exampleRole.IncludedPermissions, additionalPermissionsForIntegrationTests...)

	Expect(actualPermissions).To(ConsistOf(expectedPermissions))
}

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
