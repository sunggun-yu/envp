package cmd

import (
	"bytes"
	"fmt"
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/spf13/cobra"
	"github.com/sunggun-yu/envp/internal/config"
)

var _ = Describe("Edit Command", func() {
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
		// prepare command
		args = []string{}     // init args
		cmd = deleteCommand() // init command
		cmd.SetOut(&stdout)   // set stdout
		cmd.SetErr(&stderr)   // set stderr
		cmd.SetIn(&stdin)     // set stderr

		// prepare test config file before each test
		testConfigFile = fmt.Sprintf("/tmp/delete-%v.yaml", GinkgoRandomSeed()) // set random config file
		configFileName = testConfigFile                                         // set random config file as configFileName. so initConfig will initiate config
		copy = true                                                             // copy valid test config file as default

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
			original, _ := os.ReadFile("../testdata/config.yaml")
			os.WriteFile(testConfigFile, original, 0644)
		}

		cmd.SetArgs(args)                               // set the arg for each test case
		stdin.Write([]byte(fmt.Sprintf("%s\n", input))) // write input
		err = cmd.Execute()                             // execute the command
		cfg, _ = configFile.Read()                      //  set the config instance after executing command as result
	})

	When("delete profile", func() {
		profileName := "lab.cluster1"
		BeforeEach(func() {
			args = append(args, profileName)
			input = "y"
		})

		It("should not be error", func() {
			Expect(err).ShouldNot(HaveOccurred())
		})

		It("should be deleted", func() {
			p, err := cfg.Profile(profileName)
			Expect(err).Should(HaveOccurred())
			Expect(p).Should(BeNil())
		})
	})

	When("cancel delete profile", func() {
		profileName := "lab.cluster1"
		BeforeEach(func() {
			args = append(args, profileName)
			input = "N"
		})

		It("should not be error", func() {
			Expect(err).ShouldNot(HaveOccurred())
			fmt.Println(stdout.String(), err)
		})

		It("should not be deleted", func() {
			p, err := cfg.Profile(profileName)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(p).ShouldNot(BeNil())
			Expect(p.Name).Should(Equal(profileName))
		})
	})

	When("delete profile that is not existing", func() {
		profileName := fmt.Sprintf("not-exisiting-profile-%v", GinkgoRandomSeed())
		BeforeEach(func() {
			args = append(args, profileName)
			input = "y"
		})
		It("should not be error", func() {
			Expect(err).Should(HaveOccurred())
		})
	})

	When("set empty string of profile name", func() {
		profileName := ""
		BeforeEach(func() {
			args = append(args, profileName)
			input = "y"
		})
		It("should not be error", func() {
			Expect(err).Should(HaveOccurred())
		})
	})

	When("delete default profile", func() {
		profileName := "docker"
		BeforeEach(func() {
			args = append(args, profileName)
			input = "y"
		})

		It("should display warning message of deleting default profile", func() {
			Expect(stdout.String()).Should(ContainSubstring("WARN: "))
		})

		It("should not be error", func() {
			Expect(err).ShouldNot(HaveOccurred())
		})

		It("should be deleted", func() {
			p, err := cfg.Profile(profileName)
			Expect(err).Should(HaveOccurred())
			Expect(p).Should(BeNil())
		})

		It("should be updated default profile as empty string", func() {
			d, err := cfg.DefaultProfile()
			Expect(err).Should(HaveOccurred())
			Expect(d).Should(BeNil())
		})
	})
})
