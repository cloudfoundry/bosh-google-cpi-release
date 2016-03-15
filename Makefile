default: test

# Builds bosh-google-cpi for linux-amd64
build:
	env GOOS=linux GOARCH=amd64 go build -o out/cpi github.com/frodenas/bosh-google-cpi/main

build.linux-amd64: build
	env GOOS=linux GOARCH=amd64 go build -o out/cpi github.com/frodenas/bosh-google-cpi/main

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

# Cleans up directory and source code with gofmt
clean:
	go clean ./...

# Prepration for tests
get-deps:
	# Go vet tool
	go get golang.org/x/tools/cmd/vet

	# Ginkgo and omega test tools
	go get github.com/onsi/ginkgo/ginkgo 
	go get github.com/onsi/gomega

# Runs the unit tests with coverage
test: clean fmt vet build
	ginkgo -r -race .
