package operation_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestOperationService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Operation Service Suite")
}
