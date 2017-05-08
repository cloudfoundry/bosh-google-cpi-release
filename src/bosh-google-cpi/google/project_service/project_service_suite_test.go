package project_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestProjectService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Project Service Suite")
}
