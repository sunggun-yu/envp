package cmd

import (
	"bytes"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Version", func() {

	// setup test env
	var stdout, stderr bytes.Buffer
	version := "my version"
	cmd := versionCommand()
	cmd.SetOut(&stdout)
	cmd.SetErr(&stderr)
	cmd.SetArgs([]string{})
	SetVersion(version)

	When("execute the version command", func() {

		err := cmd.Execute()

		It("should not return error", func() {
			Expect(err).ToNot(HaveOccurred())
			Expect(stderr.String()).To(BeEmpty())
		})

		It("should return same version that was set to root command", func() {
			Expect(stdout.String()).ToNot(BeEmpty())
			Expect(stdout.String()).To(ContainSubstring(version))
		})
	})
})
