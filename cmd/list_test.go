package cmd

import (
	"bytes"
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("List", func() {

	// setup test env
	var stdout, stderr bytes.Buffer

	cmd := listCommand()
	cmd.SetOut(&stdout)
	cmd.SetErr(&stderr)
	cmd.SetArgs([]string{})

	When("execute the list command", func() {
		configFileName = "../testdata/config.yaml"

		err := cmd.Execute()
		fmt.Println(stdout.String())

		It("should not return error", func() {
			Expect(err).ToNot(HaveOccurred())
			Expect(stderr.String()).To(BeEmpty())
		})

		It("should return list of profiles", func() {
			Expect(stdout.String()).ToNot(BeEmpty())
		})
	})
})
