package cmd

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestCmd(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Cmd Suite")
}

// set configFileName as test data
var _ = BeforeSuite(func() {
	configFileName = "../testdata/config.yaml"
})
