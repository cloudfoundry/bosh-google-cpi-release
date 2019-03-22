default: test

# Builds bosh-google-cpi for linux-amd64
build:
	go build -mod=vendor -ldflags="-X bosh-google-cpi/google/config.cpiRelease=`cat release 2>/dev/null`" -o out/cpi bosh-google-cpi/main

# Build cross-platform binaries
build-all:
	./.bin/gox -output="out/cpi_{{.OS}}_{{.Arch}}" -ldflags="-X bosh-google-cpi/google/config.cpiRelease=`cat release 2>/dev/null`" bosh-google-cpi/main

# Prepration for tests
get-deps:

	test -d .bin || mkdir .bin

	# Go lint tool
	GOBIN=$$PWD/.bin go install -mod=vendor github.com/golang/lint/golint

	# Simplify cross-compiling
	GOBIN=$$PWD/.bin go install -mod=vendor github.com/mitchellh/gox

	# Ginkgo test runner
	GOBIN=$$PWD/.bin go install -mod=vendor github.com/onsi/ginkgo/ginkgo

stemcell.tgz:
	curl https://storage.googleapis.com/bosh-cpi-artifacts/bosh-stemcell-3262.12-google-kvm-ubuntu-trusty-go_agent-raw.tar.gz -o stemcell.tgz

# Cleans up directory and source code with gofmt
clean:
	go clean -mod=vendor ./...

# Run gofmt on all code
fmt:
	gofmt -l -w .

# Run linter with non-stric checking
lint:
	ls -d */ | xargs -L 1 ./.bin/golint

# Vet code
vet:
	go vet -mod=vendor ./...

# Runs the unit tests with coverage
test: get-deps clean fmt lint vet build
	./.bin/ginkgo -r -race -skipPackage=integration -mod vendor .

# Runs the integration tests from Concourse
testintci: get-deps
	./.bin/ginkgo -mod vendor integration -slowSpecThreshold=500 -progress -nodes=3 -randomizeAllSpecs -randomizeSuites $(GINKGO_ARGS) -v

# Runs the integration tests with coverage
testint: check-proj get-deps clean fmt
	CPI_ASYNC_DELETE=true \
	STEMCELL_FILE=$(PWD)/stemcell.tgz \
	STEMCELL_VERSION='3468.20' \
	GOOGLE_PROJECT=$(GOOGLE_PROJECT) \
	NETWORK_NAME=$(shell cd ../../ci/test_infrastructure && terraform output google_auto_network) \
	CUSTOM_NETWORK_NAME=$(shell cd ../../ci/test_infrastructure && terraform output google_network) \
	CUSTOM_SUBNETWORK_NAME=$(shell cd ../../ci/test_infrastructure && terraform output google_subnetwork) \
	PRIVATE_IP=10.0.0.100,10.0.0.101,10.0.0.102 \
	TARGET_POOL=$(shell cd ../../ci/test_infrastructure && terraform output google_target_pool) \
	BACKEND_SERVICE=$(shell cd ../../ci/test_infrastructure && terraform output google_backend_service) \
	REGION_BACKEND_SERVICE=$(shell cd ../../ci/test_infrastructure && terraform output google_region_backend_service) \
	ILB_INSTANCE_GROUP=$(shell cd ../../ci/test_infrastructure && terraform output google_region_backend_service) \
	ZONE=$(shell cd ../../ci/test_infrastructure && terraform output google_zone) \
	REGION=$(shell cd ../../ci/test_infrastructure && terraform output google_region) \
	GOOGLE_PROJECT=$(shell cd ../../ci/test_infrastructure && terraform output google_project) \
	SERVICE_ACCOUNT=$(shell cd ../../ci/test_infrastructure && terraform output google_service_account) \
	EXTERNAL_STATIC_IP=$(shell cd ../../ci/test_infrastructure && terraform output google_address_int_ip) \
    ./bin/ginkgo -mod vendor integration -slowSpecThreshold=500 -progress -nodes=3 -randomizeAllSpecs -randomizeSuites $(GINKGO_ARGS) -v

# Checks and creates, if necessary, resources in a project required to run integration tests.
configint: stemcell.tgz check-proj
	cd ../../ci/test_infrastructure; terraform init; terraform apply -auto-approve \
	  -var 'google_project=$(GOOGLE_PROJECT)' \
	  -var 'google_region=us-central1' \
	  -var 'google_zone=us-central1-a' \
	  -var 'google_subnetwork_range=10.0.0.0/24' \
	  -var 'google_firewall_internal=cpi-integration-internal' \
	  -var 'google_firewall_external=cpi-integration-external' \
	  -var 'google_auto_network=cpi-auto' \
	  -var 'google_network=cpi-manual' \
	  -var 'google_json_key_data=""' \
	  -var 'prefix=cpi-integration'

# Deletes the resources created by the configint target
cleanint: check-proj
	rm -f stemcell.tgz

	cd ../../ci/test_infrastructure; terraform init; terraform destroy -force \
	  -var 'google_project=$(GOOGLE_PROJECT)' \
	  -var 'google_region=us-central1' \
	  -var 'google_zone=us-central1-a' \
	  -var 'google_subnetwork_range=10.0.0.0/24' \
	  -var 'google_firewall_internal=cpi-integration-internal' \
	  -var 'google_firewall_external=cpi-integration-external' \
	  -var 'google_auto_network=cpi-auto' \
	  -var 'google_manual_network=cpi-manual' \
	  -var 'google_auto_network=cpi-auto' \
	  -var 'google_network=cpi-manual' \
	  -var 'google_json_key_data=""' \
	  -var 'prefix=cpi-integration'

	@echo "All integration tests resources have been removed. You will need to run the configint target before doing integration tests again."

check-proj:
	@if test "$(GOOGLE_PROJECT)" = "" ; then \
	  echo "GOOGLE_PROJECT must be set"; \
	  exit 1; \
	fi

check-destroy:
	@if test "$(DESTROY_OK)" = "" ; then \
	  echo "This is a potentially destructive operation. Existing instances with the "integration-delete" tag will be deleted. You must set DESTROY_OK to acknowledge."; \
	  exit 1; \
	fi
