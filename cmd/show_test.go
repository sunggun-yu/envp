package cmd

import (
	"bytes"
	"fmt"
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/spf13/cobra"
)

var _ = Describe("Show", func() {

	// setup test env
	var stdout, stderr bytes.Buffer
	var cmd *cobra.Command
	var err error

	BeforeEach(func() {
		cmd = showCommand()
		cmd.SetOut(&stdout)
		cmd.SetErr(&stderr)
	})

	When("execute the show command for default profile", func() {

		JustBeforeEach(func() {
			configFileName = "../testdata/config.yaml"
			cmd.SetArgs([]string{})
			err = cmd.Execute()
			fmt.Println(stdout.String())
		})

		It("should not return error", func() {
			Expect(err).ToNot(HaveOccurred())
			Expect(stderr.String()).To(BeEmpty())
		})

		It("should print out environment variable sets", func() {
			Expect(stdout.String()).ToNot(BeEmpty())
			Expect(stdout.String()).To(ContainSubstring("DOCKER_HOST=ssh://meow@192.168.1.40"))
		})
	})

	When("execute the show command for specific profile", func() {

		JustBeforeEach(func() {
			configFileName = "../testdata/config.yaml"
			cmd.SetArgs([]string{"lab.cluster1"})
			err = cmd.Execute()
			fmt.Println(stdout.String())
		})

		It("should not return error", func() {
			Expect(err).ToNot(HaveOccurred())
			Expect(stderr.String()).To(BeEmpty())
		})

		It("should print out environment variable sets", func() {
			Expect(stdout.String()).ToNot(BeEmpty())
			Expect(stdout.String()).To(ContainSubstring("lab-cluster1"))
		})
	})

	When("execute the show command with export flag", func() {

		JustBeforeEach(func() {
			configFileName = "../testdata/config.yaml"
			cmd.SetArgs([]string{"--export"})
			err = cmd.Execute()
			fmt.Println(stdout.String())
		})

		It("should not return error", func() {
			Expect(err).ToNot(HaveOccurred())
			Expect(stderr.String()).To(BeEmpty())
		})

		It("should print out environment variable sets", func() {
			Expect(stdout.String()).ToNot(BeEmpty())
			Expect(stdout.String()).To(ContainSubstring("DOCKER_HOST=ssh://meow@192.168.1.40"))
		})

		It("should print out environment variable sets with export command", func() {
			Expect(stdout.String()).ToNot(BeEmpty())
			Expect(stdout.String()).To(ContainSubstring("export DOCKER_HOST=ssh://meow@192.168.1.40"))
		})
	})

	When("execute the show command with empty default config", func() {

		var testFile string

		JustBeforeEach(func() {
			testFile = fmt.Sprintf("../testdata/%v.yaml", GinkgoRandomSeed())
			configFileName = testFile
			cmd.SetArgs([]string{})
			err = cmd.Execute()
			fmt.Println(stdout.String())
		})

		It("should not return error and print out example", func() {
			Expect(err).To(HaveOccurred())
			Expect(stderr.String()).NotTo(BeEmpty())
		})

		JustAfterEach(func() {
			os.Remove(testFile)
		})
	})
})
