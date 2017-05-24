default: test

# Builds bosh-google-cpi for linux-amd64
build:
	go build -ldflags="-X bosh-google-cpi/google/config.cpiRelease=`cat release 2>/dev/null`" -o out/cpi bosh-google-cpi/main

# Build cross-platform binaries
build-all:
	gox -output="out/cpi_{{.OS}}_{{.Arch}}" -ldflags="-X bosh-google-cpi/google/config.cpiRelease=`cat release 2>/dev/null`" bosh-google-cpi/main

# Prepration for tests
get-deps:
	# Go lint tool
	go get github.com/golang/lint/golint

	# Simplify cross-compiling
	go get github.com/mitchellh/gox

	# Ginkgo and omega test tools
	go get github.com/onsi/ginkgo/ginkgo
	go get github.com/onsi/gomega

# Cleans up directory and source code with gofmt
clean:
	go clean ./...

# Run gofmt on all code
fmt:
	gofmt -l -w .

# Run linter with non-stric checking
lint:
	@echo ls -d */ | grep -v vendor | xargs -L 1 golint
	ls -d */ | grep -v vendor | xargs -L 1 golint

# Vet code
vet:
	go tool vet $$(ls -d */ | grep -v vendor)

# Runs the unit tests with coverage
test: get-deps clean fmt lint vet build
	ginkgo -r -race -skipPackage=integration .

# Runs the integration tests from Concourse
testintci: get-deps
	ginkgo integration -slowSpecThreshold=500 -progress -nodes=3 -randomizeAllSpecs -randomizeSuites $(GINKGO_ARGS) -v

# Runs the integration tests with coverage
testint: check-proj get-deps clean fmt
	$(eval INTEGRATION_ADDRESS = $(shell jq -e --raw-output '.modules[0].outputs.external_ip.value' ../../ci/terraform/terraform.tfstate))

	CPI_ASYNC_DELETE=true STEMCELL_URL=https://storage.googleapis.com/bosh-cpi-artifacts/bosh-stemcell-3262.12-google-kvm-ubuntu-trusty-go_agent-raw.tar.gz SERVICE_ACCOUNT=cfintegration@$(GOOGLE_PROJECT).iam.gserviceaccount.com GOOGLE_PROJECT=$(GOOGLE_PROJECT) EXTERNAL_STATIC_IP=$(INTEGRATION_ADDRESS) ginkgo integration -slowSpecThreshold=500 -progress -nodes=3 -randomizeAllSpecs -randomizeSuites $(GINKGO_ARGS) -v

# Checks and creates, if necessary, resources in a project required to run integration tests.
configint: check-proj
	$(eval HOSTNAME = $(shell hostname))
	@terraform apply \
		-state=../../ci/terraform/terraform.tfstate \
		-var "project=$(GOOGLE_PROJECT)" \
		-var "google_credentials=$(GOOGLE_APPLICATION_CREDENTIALS)" \
		-var "env_name=$(HOSTNAME)-local-cfintegration" \
		../../ci/terraform

	@echo "Your project is ready for integration tests."

# Deletes the resources created by the configint target
cleanint: check-proj
	$(eval HOSTNAME = $(shell hostname))
	@terraform destroy -force \
		-state=../../ci/terraform/terraform.tfstate \
		-var "project=$(GOOGLE_PROJECT)" \
		-var "google_credentials=$(GOOGLE_APPLICATION_CREDENTIALS)" \
		-var "env_name=$(HOSTNAME)-local-cfintegration" \
		../../ci/terraform

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
