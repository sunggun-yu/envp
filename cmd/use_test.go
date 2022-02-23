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

	// setup test env
	var stdout, stderr bytes.Buffer
	var cmd *cobra.Command
	var testConfigFile string

	BeforeEach(func() {
		// use command
		cmd = useCommand()
		cmd.SetOut(&stdout)
		cmd.SetErr(&stderr)
		cmd.SetArgs([]string{})

		// prepare test config file before each test
		testConfigFile = fmt.Sprintf("%v.yaml", GinkgoRandomSeed())
		configFileName = testConfigFile
		oiginal, _ := ioutil.ReadFile("../testdata/config.yaml")
		ioutil.WriteFile(testConfigFile, oiginal, 0644)
	})

	AfterEach(func() {
		// delete test config file before each test - reset
		os.Remove(testConfigFile)
		cmd.SetArgs([]string{})
	})

	When("specify profile name that is existing", func() {
		It("should not return error", func() {
			cmd.SetArgs([]string{"lab.cluster2"})
			err := cmd.Execute()
			Expect(err).ToNot(HaveOccurred())
			Expect(stderr.String()).To(BeEmpty())
		})
	})

	When("specify default profile name as arg", func() {
		It("should return error", func() {
			cmd.SetArgs([]string{"docker"})
			err := cmd.Execute()
			Expect(err).NotTo(HaveOccurred())
			Expect(stdout.String()).NotTo(BeEmpty())
			Expect(stderr.String()).To(BeEmpty())
			fmt.Println(stderr.String())
		})
	})

	When("specify profile name that is not existing", func() {
		It("should return error", func() {
			cmd.SetArgs([]string{fmt.Sprintf("non-existing-profile-name-%v", GinkgoRandomSeed())})
			err := cmd.Execute()
			Expect(err).To(HaveOccurred())
			Expect(stderr.String()).NotTo(BeEmpty())
			fmt.Println(stderr.String())
		})
	})

	When("specify no profile name", func() {
		It("should return error", func() {
			err := cmd.Execute()
			Expect(err).To(HaveOccurred())
			Expect(stderr.String()).NotTo(BeEmpty())
		})
	})

})
