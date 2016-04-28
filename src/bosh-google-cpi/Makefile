default: test

# Builds bosh-google-cpi for linux-amd64
build:
	go build -o out/cpi bosh-google-cpi/main

# Build cross-platform binaries
build-all:
	gox -output="out/cpi_{{.OS}}_{{.Arch}}" bosh-google-cpi/main

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
	ginkgo -r -race .
	
# Runs the integration tests with coverage
testint: check-proj
	$(eval INTEGRATION_ADDRESS = $(shell gcloud compute addresses describe cfintegration --region=us-central1 | head  -n1 | cut -f2 -d' '))
	GOOGLE_PROJECT=$(GOOGLE_PROJECT) EXTERNAL_STATIC_IP=$(INTEGRATION_ADDRESS) go test -v -timeout=60m ./integration/... -ginkgo.slowSpecThreshold=500 -ginkgo.progress -ginkgo.v $(GINKGO_ARGS)

# Checks and creates, if necessary, resources in a project required to run integration tests.
configint: check-proj
	@gcloud compute networks describe cfintegration &> /dev/null; if [ $$? -ne 0 ]; then \
		echo "Creating network cfintegration"; \
		gcloud --project=$(GOOGLE_PROJECT) compute networks create cfintegration; \
	fi

	@gcloud compute networks describe cfintegration-custom &> /dev/null; if [ $$? -ne 0 ]; then \
		echo "Creating network cfintegration-custom"; \
		gcloud --project=$(GOOGLE_PROJECT) compute networks create cfintegration-custom --mode=custom; \
	fi

	@gcloud compute networks subnets describe cfintegration-custom-us-central1 --region=us-central1 &> /dev/null; if [ $$? -ne 0 ]; then \
		echo "Creating subnetwork cfintegration-custom-us-central1"; \
		gcloud --project=$(GOOGLE_PROJECT) compute networks subnets create cfintegration-custom-us-central1 --network=cfintegration-custom --range=192.168.0.0/16 --region=us-central1; \
	fi

	@gcloud compute addresses describe cfintegration --region us-central1 &> /dev/null; if [ $$? -ne 0 ]; then \
		echo "Creating static IP cfintegration"; \
		gcloud --project=$(GOOGLE_PROJECT) compute addresses create cfintegration --region=us-central1; \
	fi

	@echo "Your project is ready for integration tests."

check-proj:
	@if test "$(GOOGLE_PROJECT)" = "" ; then \
	  echo "GOOGLE_PROJECT must be set"; \
	  exit 1; \
	fi
