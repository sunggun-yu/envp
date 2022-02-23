package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/spf13/cobra"
)

var _ = Describe("Use", func() {

	var (
		testConfigFile string
	)

	BeforeEach(func() {
		// prepare test config file before each test
		testConfigFile = fmt.Sprintf("%v.yaml", GinkgoRandomSeed())
		configFileName = testConfigFile
		oiginal, _ := ioutil.ReadFile("../testdata/config.yaml")
		ioutil.WriteFile(testConfigFile, oiginal, 0644)
	})

	AfterEach(func() {
		// delete test config file before each test - reset
		os.Remove(testConfigFile)
	})

	When("run use command without specifying profile name", func() {
		// setup test env
		var (
			stdout, stderr bytes.Buffer
			cmd            *cobra.Command
		)

		BeforeEach(func() {
			// use command
			cmd = useCommand()
			cmd.SetOut(&stdout)
			cmd.SetErr(&stderr)
			cmd.SetArgs(make([]string, 0))
		})

		When("specify no profile name", func() {
			var (
				err error
			)

			JustBeforeEach(func() {
				cmd.SetArgs(make([]string, 0))
				err = cmd.Execute()
			})

			It("should return error", func() {
				Expect(err).To(HaveOccurred())
			})

			It("should print out error message", func() {
				Expect(stderr.String()).NotTo(BeEmpty())
			})
		})
	})

	When("run use command with specifying profile name", func() {
		// setup test env
		var (
			stdout, stderr bytes.Buffer
			cmd            *cobra.Command
		)

		BeforeEach(func() {
			// use command
			cmd = useCommand()
			cmd.SetOut(&stdout)
			cmd.SetErr(&stderr)
			cmd.SetArgs(make([]string, 0))
		})

		When("specify profile name that is existing", func() {
			var (
				err error
			)

			JustBeforeEach(func() {
				cmd.SetArgs([]string{"lab.cluster2"})
				err = cmd.Execute()
			})

			It("should not return error", func() {
				Expect(err).ShouldNot(HaveOccurred())
			})

			It("should print out success message", func() {
				Expect(stdout.String()).ShouldNot(BeEmpty())
			})

			It("should not be error message", func() {
				Expect(stderr.String()).Should(BeEmpty())
			})
		})

		When("specify default profile name as arg", func() {
			var (
				err error
			)

			JustBeforeEach(func() {
				cmd.SetArgs([]string{"docker"})
				err = cmd.Execute()
			})

			It("should not return error", func() {
				Expect(err).ShouldNot(HaveOccurred())
			})

			It("should print out success message", func() {
				Expect(stdout.String()).ShouldNot(BeEmpty())
			})

			It("should not be error message", func() {
				Expect(stderr.String()).Should(BeEmpty())
			})
		})

		When("specify profile name that is not existing", func() {
			var (
				err error
			)

			JustBeforeEach(func() {
				cmd.SetArgs([]string{fmt.Sprintf("non-existing-profile-name-%v", GinkgoRandomSeed())})
				err = cmd.Execute()
			})

			It("should return error", func() {
				Expect(err).Should(HaveOccurred())
			})

			It("should print out error message", func() {
				Expect(stderr.String()).ShouldNot(BeEmpty())
				fmt.Println(stderr.String())
			})
		})
	})

	When("run use command with specifying multiple profile name", func() {
		// setup test env
		var (
			stdout, stderr bytes.Buffer
			cmd            *cobra.Command
			err            error
		)

		BeforeEach(func() {
			// use command
			cmd = useCommand()
			cmd.SetOut(&stdout)
			cmd.SetErr(&stderr)
			cmd.SetArgs([]string{"lab.cluster1", "lab.cluster2"})
		})

		JustBeforeEach(func() {
			cmd.SetArgs([]string{"lab.cluster1", "lab.cluster2"})
			err = cmd.Execute()
		})

		It("should return error", func() {
			Expect(err).Should(HaveOccurred())
		})

		It("should print out error message", func() {
			Expect(stderr.String()).ShouldNot(BeEmpty())
			fmt.Println(stderr.String())
		})
	})
})
