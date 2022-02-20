package shell

import (
	"bytes"
	"fmt"
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/sunggun-yu/envp/internal/config"
)

var _ = Describe("Shell", func() {

	Describe("run Execute", func() {

		sc := NewShellCommand()

		Context("eligible command", func() {

			When("passing non-empty envs", func() {
				It("should not return err", func() {
					cmd := "echo"
					err := sc.Execute([]string{cmd}, []config.Env{
						{Name: "meow", Value: "woof"},
					}, "my-profile")
					Expect(err).ToNot(HaveOccurred())
				})
			})

			When("pass wrong arg to command", func() {
				It("should return err", func() {
					cmd := []string{"cat", "/not-existing-dir/not-existing-file-rand-meow"}
					err := sc.Execute(cmd, []config.Env{
						{Name: "meow", Value: "woof"},
					}, "my-profile")
					Expect(err).To(HaveOccurred())
				})
			})
		})

		When("run non-existing command", func() {
			It("should not return err", func() {
				cmd := ""
				err := sc.Execute([]string{cmd}, []config.Env{
					{Name: "meow", Value: "woof"},
				}, "my-profile")
				Expect(err).To(HaveOccurred())
			})
		})

	})

	Describe("run StartShell", func() {

		var stdout, stderr bytes.Buffer
		sc := NewShellCommand()
		sc.Stdout = &stdout
		sc.Stderr = &stderr

		// github action has no default SHELL. so set it as /bin/sh before each test case
		JustBeforeEach(func() {
			// make SHELL empty to occur error
			os.Setenv("SHELL", "/bin/sh")
		})

		When("pass not empty envs", func() {
			It("should not return err", func() {
				err := sc.StartShell([]config.Env{
					{Name: "meow", Value: "woof"},
				}, "my-profile")
				Expect(err).ToNot(HaveOccurred())
			})
		})

		When("pass nil envs", func() {
			It("should not return err", func() {
				err := sc.StartShell(nil, "")
				Expect(err).ToNot(HaveOccurred())
				Expect(stdout.String()).NotTo(BeEmpty())
				Expect(stderr.String()).To(BeEmpty())
			})
		})

		When("error occurred", func() {
			// backup original SHELL to set it back after test
			sh := os.Getenv("SHELL")

			JustBeforeEach(func() {
				// make SHELL empty to occur error
				os.Setenv("SHELL", "")
			})
			It("should return err", func() {
				err := sc.StartShell(nil, "my-profile")
				Expect(err).To(HaveOccurred())
			})
			JustAfterEach(func() {
				// revert SHELL to original
				os.Setenv("SHELL", sh)
			})
		})
	})
})

var _ = Describe("env functions", func() {

	envs := config.Envs{
		config.Env{
			Name:  "PATH",
			Value: "~/.config",
		},
		config.Env{
			Name:  "HOME",
			Value: "$HOME",
		},
	}
	h, _ := os.UserHomeDir()
	pe := appendEnvpProfile(parseEnvs(envs), "my-profile")

	When("has ~ in the value", func() {
		It("should extraced to abs home dir", func() {
			Expect(pe).To(ContainElement(fmt.Sprintf("PATH=%s/.config", h)))
		})
	})

	When("has $HOME in the value", func() {
		It("should extraced to abs home dir", func() {
			Expect(pe).To(ContainElement(fmt.Sprintf("HOME=%s", h)))
		})
	})

	When("append profile env var", func() {
		It("should include env var value of profile", func() {
			Expect(pe).To(ContainElement(fmt.Sprintf("%s=my-profile", envpEnvVarKey)))
		})
	})
})
