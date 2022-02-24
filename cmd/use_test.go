package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/spf13/cobra"
	"github.com/sunggun-yu/envp/internal/config"
)

var _ = Describe("Use", func() {

	var (
		args           []string       // args to pass to command
		testConfigFile string         // test config file path
		stdout, stderr bytes.Buffer   // stdout and stderr
		cmd            *cobra.Command // command
		err            error          // error
		cfg            *config.Config // config instance
		copy           bool           // whether copy valid config file as test file or not
	)

	// BeforeEach prepare cmd and copy of test config file
	BeforeEach(func() {
		// prepare command
		args = []string{}   // init args
		cmd = useCommand()  // init command
		cmd.SetOut(&stdout) // set stdout
		cmd.SetErr(&stderr) // set stderr

		// prepare test config file before each test
		testConfigFile = fmt.Sprintf("use-%v.yaml", GinkgoRandomSeed()) // set random config file
		configFileName = testConfigFile                                 // set random config file as configFileName. so initConfig will initiate config
		copy = true                                                     // copy valid test config file as default

		// delete test config file
		DeferCleanup(func() {
			os.Remove(testConfigFile) // remove test config file after test case
		})
	})

	// AfterEach reset the stdout and stderr
	AfterEach(func() {
		stdout.Reset() // reset stdout after test case. so the last test case result will be cleared
		stderr.Reset() // reset stderr after test case. so the last test case result will be cleared
	})

	// it runs right before the It
	JustBeforeEach(func() {

		// copy if copy is true. otherwise it will be fresh empty config file
		if copy {
			oiginal, _ := ioutil.ReadFile("../testdata/config.yaml")
			ioutil.WriteFile(testConfigFile, oiginal, 0644)
		}

		cmd.SetArgs(args)          // set the arg for each test case
		err = cmd.Execute()        // execute the command
		cfg, _ = configFile.Read() //  set the config instance after executing command as result
	})

	When("specifying no profile name", func() {
		BeforeEach(func() {
			args = []string{}
		})
		It("should be error and print out error message", func() {
			Expect(err).To(HaveOccurred())
			Expect(stderr.String()).NotTo(BeEmpty())
		})
	})

	When("specify profile name that is existing", func() {

		profileName := "lab.cluster2"

		BeforeEach(func() {
			args = append(args, profileName)
		})
		It("should not error and print out success message", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(stdout.String()).ShouldNot(BeEmpty())
			Expect(stderr.String()).Should(BeEmpty())
		})

		It("should change the default profile of config", func() {
			d, err := cfg.DefaultProfile()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(d.Name).Should(Equal(profileName))
		})
	})

	When("specify default profile name as arg", func() {
		BeforeEach(func() {
			args = append(args, "docker")
		})
		It("should not error and print out success message", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(stdout.String()).ShouldNot(BeEmpty())
			Expect(stderr.String()).Should(BeEmpty())
		})
	})

	When("specify profile name that is not existing", func() {
		profileName := "some-profile-name"
		BeforeEach(func() {
			copy = false
			args = append(args, profileName)
		})
		It("should be error and print out error message", func() {
			Expect(err).Should(HaveOccurred())
			Expect(stderr.String()).ShouldNot(BeEmpty())
			fmt.Println(stderr.String())
		})
	})

	When("run use command with specifying multiple profile name", func() {
		BeforeEach(func() {
			args = append(args, "lab.cluster1", "lab.cluster2")
		})
		It("should be error and print out error message", func() {
			Expect(err).Should(HaveOccurred())
			Expect(stderr.String()).ShouldNot(BeEmpty())
			fmt.Println(stderr.String())
		})
	})
})
