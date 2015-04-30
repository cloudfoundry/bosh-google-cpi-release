package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestRegistryMain(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Registry Main Suite")
}
