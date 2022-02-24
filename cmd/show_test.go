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

var _ = Describe("Show", func() {
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
		cmd = showCommand() // init command
		cmd.SetOut(&stdout) // set stdout
		cmd.SetErr(&stderr) // set stderr

		// prepare test config file before each test
		testConfigFile = fmt.Sprintf("show-%v.yaml", GinkgoRandomSeed()) // set random config file
		configFileName = testConfigFile                                  // set random config file as configFileName. so initConfig will initiate config
		copy = true                                                      // copy valid test config file as default

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

	When("execute the show command for default profile", func() {
		BeforeEach(func() {
			args = []string{}
		})
		It("should not return error", func() {
			Expect(err).ToNot(HaveOccurred())
			Expect(stderr.String()).To(BeEmpty())
		})

		It("should print out environment variable sets of default profile", func() {
			p, _ := cfg.DefaultProfile()
			for _, e := range p.Env {
				Expect(stdout.String()).To(ContainSubstring(e.String()))
			}
		})
	})

	When("specify profile that is existing", func() {
		profileName := "lab.cluster1"
		BeforeEach(func() {
			args = append(args, profileName)
		})
		It("should not return error", func() {
			Expect(err).ToNot(HaveOccurred())
			Expect(stderr.String()).To(BeEmpty())
		})
		It("should print out environment variable sets", func() {
			p, _ := cfg.Profile(profileName)
			for _, e := range p.Env {
				Expect(stdout.String()).To(ContainSubstring(e.String()))
			}
		})
	})

	When("execute the show command with export flag", func() {
		profileName := "lab.cluster1"
		BeforeEach(func() {
			args = append(args, profileName, "--export")
		})
		It("should not return error", func() {
			Expect(err).ToNot(HaveOccurred())
			Expect(stderr.String()).To(BeEmpty())
		})
		It("should print out environment variable sets", func() {
			p, _ := cfg.Profile(profileName)
			for _, e := range p.Env {
				Expect(stdout.String()).To(ContainSubstring(fmt.Sprintf("export %s", e.String())))
			}
		})
	})

	When("default profile is empty", func() {
		BeforeEach(func() {
			copy = false
			args = []string{}
		})
		It("should be error", func() {
			Expect(err).To(HaveOccurred())
			Expect(stderr.String()).NotTo(BeEmpty())
		})
	})

	When("profile name is empty string", func() {
		BeforeEach(func() {
			copy = false
			args = []string{""}
		})
		It("should be error", func() {
			Expect(err).To(HaveOccurred())
			Expect(stderr.String()).NotTo(BeEmpty())
		})
	})
})
