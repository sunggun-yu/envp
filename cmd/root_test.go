package cmd

import (
	"bytes"
	"os"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestInitConfig(t *testing.T) {
	initConfig()
}

var _ = Describe("run the command with empty args", Ordered, func() {

	var stdout, stderr bytes.Buffer
	testDir := "/tmp/envp/test"

	BeforeAll(func() {
		configPath = testDir
		rootCmd.SetOut(&stdout)
		rootCmd.SetErr(&stderr)
		// this is condition for test case
		rootCmd.SetArgs([]string{})
	})

	It("should exit 0", func() {
		defer func() {
			if r := recover(); r != nil {
				Expect(r).To(ContainSubstring("Exit(0)"))
			}
		}()
		// this should exit 0. so, evaluate expectation in recover
		rootCmd.Execute()
	})

	It("should not error", func() {
		defer func() {
			if r := recover(); r != nil {
				err := stderr.String()
				Expect(err).To(BeEmpty())
			}
		}()
		// this should exit 0. so, evaluate expectation in recover
		rootCmd.Execute()
	})

	It("should print out the help", func() {
		defer func() {
			if r := recover(); r != nil {
				out := stdout.String()
				Expect(out).To(ContainSubstring("Usage:"))
				Expect(out).To(ContainSubstring("Examples:"))
				Expect(out).To(ContainSubstring("Available Commands:"))
				Expect(out).To(ContainSubstring("Flags:"))
			}
		}()
		// this should exit 0. so, evaluate expectation in recover
		rootCmd.Execute()
	})

	AfterAll(func() {
		DeferCleanup(func() {
			os.Remove(testDir)
		})
	})
})
