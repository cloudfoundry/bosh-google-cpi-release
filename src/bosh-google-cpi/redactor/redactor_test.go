package redactor_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"bosh-google-cpi/redactor"
)

var _ = Describe("Redactor", func() {
	Describe("RedactSecrets", func() {
		It("Redacts", func() {
			secretString := `{
  "method": "fake-action",
  "arguments": [
    {
      "Password": "secret_data",
      "private_key": "more\n_secret_data",
      "public_key": "public_data",
      "account_key": "secret_data",
      "json_key": "secret_data",
      "secret_access_key": "secret_data"
    }
  ],
  "context": "with agent settings \"{\\\"env\\\":{\\\"bosh\\\":{\\\"blobstores\\\":[{\\\"options\\\":{\\\"Password\\\": \\\"escaped_secret_data\\\",\\\"private_key\\\": \\\"escaped_more\n_secret_data\\\",\\\"public_key\\\": \\\"escaped_public_data\\\",\\\"account_key\\\": \\\"escaped_secret_data\\\",\\\"json_key\\\": \\\"escaped_secret_data\\\",\\\"secret_access_key\\\": \\\"escaped_secret_data\\\""
}`
			Expect(redactor.RedactSecrets(secretString)).NotTo(ContainSubstring("secret_data"))
		})
	})
})
