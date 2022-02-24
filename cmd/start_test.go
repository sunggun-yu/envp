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
	"github.com/sunggun-yu/envp/internal/shell"
)

var _ = Describe("Start Command", func() {
	var (
		args                  []string       // args to pass to command
		testConfigFile        string         // test config file path
		stdout, stderr, stdin bytes.Buffer   // stdout and stderr
		cmd                   *cobra.Command // command
		err                   error          // error
		cfg                   *config.Config // config instance
		copy                  bool           // whether copy valid config file as test file or not
		input                 string         // user input for prompt
	)

	// BeforeEach prepare cmd and copy of test config file
	BeforeEach(func() {

		// should default shell for ci testing
		os.Setenv("SHELL", "/bin/sh")

		// prepare the shell command
		sc := shell.NewShellCommand()
		sc.Stdin = &stdin
		sc.Stdout = &stdout
		sc.Stderr = &stderr

		// prepare command
		args = []string{}      // init args
		cmd = startCommand(sc) // init command
		cmd.SetOut(&stdout)    // set stdout
		cmd.SetErr(&stderr)    // set stderr
		cmd.SetIn(&stdin)      // set stderr

		// prepare test config file before each test
		testConfigFile = fmt.Sprintf("/tmp/start-%v.yaml", GinkgoRandomSeed()) // set random config file
		configFileName = testConfigFile                                        // set random config file as configFileName. so initConfig will initiate config
		copy = true                                                            // copy valid test config file as default

		// delete test config file
		DeferCleanup(func() {
			os.Remove(testConfigFile) // remove test config file after test case
		})
	})

	// AfterEach reset the stdout and stderr
	AfterEach(func() {
		stdout.Reset() // reset stdout after test case. so the last test case result will be cleared
		stderr.Reset() // reset stderr after test case. so the last test case result will be cleared
		stdin.Reset()  // reset stdin after test case. so the last test case result will be cleared
	})

	// it runs right before the It
	JustBeforeEach(func() {

		// copy if copy is true. otherwise it will be fresh empty config file
		if copy {
			oiginal, _ := ioutil.ReadFile("../testdata/config.yaml")
			ioutil.WriteFile(testConfigFile, oiginal, 0644)
		}

		cmd.SetArgs(args)                               // set the arg for each test case
		stdin.Write([]byte(fmt.Sprintf("%s\n", input))) // write input
		err = cmd.Execute()                             // execute the command
		cfg, _ = configFile.Read()                      //  set the config instance after executing command as result
	})

	When("execute start command with valid inputs", func() {
		profileName := "lab.cluster1"
		BeforeEach(func() {
			args = append(args, profileName)
			input = "env"
		})

		It("should not be error", func() {
			Expect(err).ShouldNot(HaveOccurred())
		})

		It("should include environment variable of profile", func() {
			p, err := cfg.Profile(profileName)
			Expect(err).ShouldNot(HaveOccurred())
			for _, e := range p.Env {
				Expect(stdout.String()).Should(ContainSubstring(e.String()))
			}
		})

		It("should include ENVP_PROFILE environment variable that is match to selected profile", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(stdout.String()).Should(ContainSubstring(fmt.Sprintf("ENVP_PROFILE=%s", profileName)))
		})
	})

	When("execute start command with profile that is not exisiting", func() {
		profileName := fmt.Sprintf("not-exisiting-profile-%v", GinkgoRandomSeed())
		BeforeEach(func() {
			args = append(args, profileName)
			input = "echo hello"
		})

		It("should be error", func() {
			Expect(err).Should(HaveOccurred())
		})
	})

	When("execute start command with empty string of profile", func() {
		profileName := ""
		BeforeEach(func() {
			args = append(args, profileName)
			input = "echo hello"
		})

		It("should be error", func() {
			Expect(err).Should(HaveOccurred())
		})
	})

	When("execute start command with valid inputs but omit profile name", func() {
		BeforeEach(func() {
			args = []string{}
			input = "env"
		})

		It("should not be error", func() {
			Expect(err).ShouldNot(HaveOccurred())
		})

		It("should execute for default profile and include environment variable of default profile", func() {
			p, err := cfg.DefaultProfile()
			Expect(err).ShouldNot(HaveOccurred())
			for _, e := range p.Env {
				Expect(stdout.String()).Should(ContainSubstring(e.String()))
			}
		})

		It("should include ENVP_PROFILE environment variable that is match to default profile", func() {
			p, err := cfg.DefaultProfile()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(stdout.String()).Should(ContainSubstring(fmt.Sprintf("ENVP_PROFILE=%s", p.Name)))
		})
	})
})
