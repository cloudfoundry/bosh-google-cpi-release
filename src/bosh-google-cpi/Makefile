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
testint: check-proj check-destroy get-deps clean fmt lint vet build
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

	@gcloud compute target-pools describe cfintegration --region us-central1 &> /dev/null; if [ $$? -ne 0 ]; then \
		echo "Creating target pool cfintegration"; \
		gcloud --project=$(GOOGLE_PROJECT) compute target-pools create cfintegration --region=us-central1; \
	fi
	
	@gcloud compute target-pools describe cfintegration2 --region us-central1 &> /dev/null; if [ $$? -ne 0 ]; then \
		echo "Creating target pool cfintegration"; \
		gcloud --project=$(GOOGLE_PROJECT) compute target-pools create cfintegration2 --region=us-central1; \
	fi

	@gcloud compute instance-groups unmanaged describe cfintegration-us-central1-a --zone us-central1-a &> /dev/null; if [ $$? -ne 0 ]; then \
		echo "Creating instance group cfintegration-us-central1-a"; \
		gcloud --project=$(GOOGLE_PROJECT) compute instance-groups unmanaged create cfintegration-us-central1-a --zone us-central1-a; \
	fi
	@echo "Your project is ready for integration tests."

# Deletes the resources created by the configint target
cleanint: check-proj 
	@gcloud compute networks describe cfintegration &> /dev/null; if [ $$? -eq 0 ]; then \
		echo "Deleting network cfintegration"; \
		gcloud --project=$(GOOGLE_PROJECT) compute networks delete cfintegration --quiet; \
	fi

	@gcloud compute networks describe cfintegration-custom &> /dev/null; if [ $$? -eq 0 ]; then \
		echo "Deleting network cfintegration-custom"; \
		gcloud --project=$(GOOGLE_PROJECT) compute networks delete cfintegration-custom --quiet; \
	fi

	@gcloud compute networks subnets describe cfintegration-custom-us-central1 --region=us-central1 &> /dev/null; if [ $$? -eq 0 ]; then \
		echo "Deleting subnetwork cfintegration-custom-us-central1"; \
		gcloud --project=$(GOOGLE_PROJECT) compute networks subnets delete cfintegration-custom-us-central1 --region=us-central1 --quiet; \
	fi

	@gcloud compute addresses describe cfintegration --region us-central1 &> /dev/null; if [ $$? -eq 0 ]; then \
		echo "Deleting static IP cfintegration"; \
		gcloud --project=$(GOOGLE_PROJECT) compute addresses delete cfintegration --region=us-central1 --quiet; \
	fi

	@gcloud compute target-pools describe cfintegration --region us-central1 &> /dev/null; if [ $$? -eq 0 ]; then \
		echo "Deleting target pool cfintegration"; \
		gcloud --project=$(GOOGLE_PROJECT) compute target-pools delete cfintegration --region=us-central1 --quiet; \
	fi

	@gcloud compute target-pools describe cfintegration2 --region us-central1 &> /dev/null; if [ $$? -eq 0 ]; then \
		echo "Deleting target pool cfintegration"; \
		gcloud --project=$(GOOGLE_PROJECT) compute target-pools delete cfintegration2 --region=us-central1 --quiet; \
	fi

	@gcloud compute instance-groups unmanaged describe cfintegration-us-central1-a --zone us-central1-a &> /dev/null; if [ $$? -eq 0 ]; then \
		echo "Deleting instance group cfintegration-us-central1-a"; \
		gcloud --project=$(GOOGLE_PROJECT) compute instance-groups unmanaged delete cfintegration-us-central1-a --zone us-central1-a --quiet; \
	fi
	@echo "All integration tests resources have been removed. You will need to run the configint target before doing integration tests again."

check-proj:
	@if test "$(GOOGLE_PROJECT)" = "" ; then \
	  echo "GOOGLE_PROJECT must be set"; \
	  exit 1; \
	fi
	
check-destroy:
	@if test "$(DESTROY_OK)" = "" ; then \
	  echo "This is a potentially destructive operation. Existing instances with the "integration-delete" tag will be deleted. You just set DESTROY_OK to acknowledge."; \
	  exit 1; \
	fi
