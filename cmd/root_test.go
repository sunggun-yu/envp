package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/sunggun-yu/envp/internal/config"
	"github.com/sunggun-yu/envp/internal/shell"
)

func TestInitConfig(t *testing.T) {
	initConfig()
}

func TestExecute(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			assert.Contains(t, r, "Exit(0)")
		}
	}()
	Execute()
}

var _ = Describe("Root command with empty args", Ordered, func() {

	var stdout, stderr bytes.Buffer
	testDir := "/tmp/envp/test"

	BeforeAll(func() {
		configFileName = testDir
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

var _ = Describe("Root Command", func() {
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

		// should default shell for ci testing
		os.Setenv("SHELL", "/bin/sh")

		// prepare the shell command
		sc := shell.NewShellCommand()
		sc.Stdout = &stdout
		sc.Stderr = &stderr

		// prepare command
		args = []string{}     // init args
		cmd = rootCommand(sc) // init command
		cmd.SetOut(&stdout)   // set stdout
		cmd.SetErr(&stderr)   // set stderr

		// prepare test config file before each test
		testConfigFile = fmt.Sprintf("/tmp/root-%v.yaml", GinkgoRandomSeed()) // set random config file
		configFileName = testConfigFile                                       // set random config file as configFileName. so initConfig will initiate config
		copy = true                                                           // copy valid test config file as default

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

	When("execute command with valid inputs", func() {
		profileName := "lab.cluster2"
		BeforeEach(func() {
			args = append(args, profileName, "--", "env")
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

	When("execute command with valid inputs but no specify profile name", func() {
		BeforeEach(func() {
			args = append(args, "--", "env")
		})

		It("should not be error", func() {
			Expect(err).ShouldNot(HaveOccurred())
		})

		It("should include environment variable of default profile", func() {
			p, err := cfg.DefaultProfile()
			Expect(err).ShouldNot(HaveOccurred())
			for _, e := range p.Env {
				Expect(stdout.String()).Should(ContainSubstring(e.String()))
			}
		})

		It("should include ENVP_PROFILE environment variable that is match to default profile name", func() {
			p, err := cfg.DefaultProfile()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(err).ShouldNot(HaveOccurred())
			Expect(stdout.String()).Should(ContainSubstring(fmt.Sprintf("ENVP_PROFILE=%s", p.Name)))
		})
	})

	When("execute command but no specifying command", func() {
		profileName := "lab.cluster2"
		BeforeEach(func() {
			args = append(args, profileName)
		})

		It("should be error", func() {
			Expect(err).Should(HaveOccurred())
		})
	})

	When("execute command with profile that is not exisiting", func() {
		profileName := fmt.Sprintf("not-exisiting-profile-%v", GinkgoRandomSeed())
		BeforeEach(func() {
			args = append(args, profileName, "--", "echo", "hello")
		})

		It("should be error", func() {
			Expect(err).Should(HaveOccurred())
		})
	})

	When("execute command but shell command that is not exisiting", func() {
		profileName := "lab.cluster2"
		BeforeEach(func() {
			args = append(args, profileName, "--", "1293471029384701298374019872498-aslkaslkjasdfjklasdfjklasdf-202020202")
		})

		It("should be error", func() {
			Expect(err).Should(HaveOccurred())
		})
	})

	// When("execute start command with empty string of profile", func() {
	// 	profileName := ""
	// 	BeforeEach(func() {
	// 		args = append(args, profileName)
	// 		input = "echo hello"
	// 	})

	// 	It("should be error", func() {
	// 		Expect(err).Should(HaveOccurred())
	// 	})
	// })

	// When("execute start command with valid inputs but omit profile name", func() {
	// 	BeforeEach(func() {
	// 		args = []string{}
	// 		input = "env"
	// 	})

	// 	It("should not be error", func() {
	// 		Expect(err).ShouldNot(HaveOccurred())
	// 	})

	// 	It("should execute for default profile and include environment variable of default profile", func() {
	// 		p, err := cfg.DefaultProfile()
	// 		Expect(err).ShouldNot(HaveOccurred())
	// 		for _, e := range p.Env {
	// 			Expect(stdout.String()).Should(ContainSubstring(e.String()))
	// 		}
	// 	})

	// 	It("should include ENVP_PROFILE environment variable that is match to default profile", func() {
	// 		p, err := cfg.DefaultProfile()
	// 		Expect(err).ShouldNot(HaveOccurred())
	// 		Expect(stdout.String()).Should(ContainSubstring(fmt.Sprintf("ENVP_PROFILE=%s", p.Name)))
	// 	})
	// })
})
