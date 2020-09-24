module bosh-google-cpi

go 1.12

require (
	cloud.google.com/go v0.66.0 // indirect
	github.com/bmatcuk/doublestar v1.1.1 // indirect
	github.com/charlievieth/fs v0.0.0-20170613215519-7dc373669fa1 // indirect
	github.com/cloudfoundry/bosh-utils v0.0.0-20180413212538-2c869a1a0cce
	github.com/golang/lint v0.0.0-20181217174547-8f45f776aaf1
	github.com/mitchellh/gox v0.4.0
	github.com/mitchellh/iochan v1.0.0 // indirect
	github.com/nu7hatch/gouuid v0.0.0-20131221200532-179d4d0c4d8d // indirect
	github.com/onsi/ginkgo v1.8.0
	github.com/onsi/gomega v1.5.0
	golang.org/x/net v0.0.0-20200904194848-62affa334b73
	golang.org/x/oauth2 v0.0.0-20200902213428-5d25da1a8d43
	golang.org/x/sys v0.0.0-20200922070232-aee5d888a860 // indirect
	google.golang.org/api v0.32.0
	google.golang.org/genproto v0.0.0-20200921165018-b9da36f5f452 // indirect
	google.golang.org/grpc v1.32.0 // indirect
)

replace bosh-google-cpi => ../
