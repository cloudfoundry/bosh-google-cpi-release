package redactor_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestDispatcher(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Redactor Suite")
}
