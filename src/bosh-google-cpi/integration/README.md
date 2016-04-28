# Integration Testing
Running integration tests requires a Google Compute Engine project with billing enabled. It creates real resources and costs real money to run.

## Configure your project
Run the `configint` target with your project specified to provision the resources required to run integration tests. For example:

```
$ make configint GOOGLE_PROJECT=my-test-project
```

This will provision several networks, a subnetwork, and a static external IP address. It is **important** to note that the IP address is a billable resource. To minimize charges, you can use `gcloud` to manually delete the address after tests have run.

## Run integration
Use the `testint` target to run integration tests.

### All tests
This command will run all possible integration tests, and may take close to an hour to complete:

```
$ GOOGLE_PROJECT=evandbrown17 KEEP_REUSABLE_VM=true GINKGO_ARGS=-ginkgo.focus=Stemcell make testint
```

Setting KEEP_REUSABLE_VM will prevent the VM created by tests from being deleted, allowing it to be reused in future runs. This is helpful when writing tests and iterating quickly as it eliminates VM start time on subsequent test runs:

```
$ GOOGLE_PROJECT=evandbrown17 KEEP_REUSABLE_VM=true make testint
```

Finally, you can pass arguments to Ginkgo. To run only the stemcell tests, set the focus:

```
$ GOOGLE_PROJECT=evandbrown17 GINKGO_ARGS=-ginkgo.focus=Stemcell make testint
```

