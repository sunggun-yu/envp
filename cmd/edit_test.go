package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/spf13/cobra"
	"github.com/sunggun-yu/envp/internal/config"
)

var _ = Describe("Edit Command", func() {
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
		cmd = editCommand() // init command
		cmd.SetOut(&stdout) // set stdout
		cmd.SetErr(&stderr) // set stderr

		// prepare test config file before each test
		testConfigFile = fmt.Sprintf("edit-%v.yaml", GinkgoRandomSeed()) // set random config file
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

	When("edit profile", func() {
		profileName := "lab.cluster1"
		desc := "lab.cluster1" + "-updated"
		envs := []string{"env1=var1", "env2=var2"}
		BeforeEach(func() {
			args = append(args, profileName, "-d", desc, "-e", envs[0], "-e", envs[1])
		})

		It("should not be error", func() {
			Expect(err).ShouldNot(HaveOccurred())
			fmt.Println(stdout.String(), err)
		})

		It("should be added correctly with given args and flags", func() {
			p, err := cfg.Profile(profileName)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(p).ShouldNot(BeNil())
			Expect(p.Name).Should(Equal(profileName))
			Expect(p.Env.String()).Should(ContainSubstring(strings.Join(envs, ",")))
			Expect(p.Desc).Should(Equal(desc))
		})
	})

	When("edit profile that is not exisiting", func() {
		profileName := "lab.cluster1"
		envs := []string{"env1=var1", "env2=var2"}
		BeforeEach(func() {
			copy = false // empty default config file. so nothing exist
			args = append(args, profileName, "-d", profileName, "-e", envs[0], "-e", envs[1])
		})

		It("should be error", func() {
			Expect(err).Should(HaveOccurred())
			Expect(stderr.String()).ShouldNot(BeEmpty())
			fmt.Println(stderr.String(), err)
		})
	})

	When("pass multiple profile names", func() {
		BeforeEach(func() {
			args = append(args, "lab.cluster1", "lab.cluster2", "-d", "some-desc", "-e", "env=var")
		})
		It("should be error", func() {
			Expect(err).Should(HaveOccurred())
			fmt.Println(stderr.String(), err)
		})
	})

	When("pass no args and flags at all", func() {
		BeforeEach(func() {
			args = []string{}
		})
		It("should be error", func() {
			Expect(err).Should(HaveOccurred())
			fmt.Println(stderr.String(), err)
		})
	})

	When("pass no env flag but desc flag only", func() {
		profileName := "lab.cluster1"
		desc := "some-updated-desc"
		BeforeEach(func() {
			args = []string{profileName, "-d", desc}
		})
		It("should not be error", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(stdout.String()).ShouldNot(BeEmpty())
			fmt.Println(stdout.String(), err)
		})
		It("should update desc", func() {
			p, err := cfg.Profile(profileName)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(stdout.String()).ShouldNot(BeEmpty())
			Expect(p.Desc).Should(Equal(desc))
			fmt.Println(stdout.String(), err)
		})
	})
})
