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

	profile := config.NamedProfile{
		Name:    "my-profile",
		Profile: config.NewProfile(),
	}
	profile.Env = []*config.Env{
		{Name: "meow", Value: "woof"},
	}

	Describe("run Execute", func() {

		var stdout, stderr bytes.Buffer
		sc := NewShellCommand()
		sc.Stdout = &stdout
		sc.Stderr = &stderr

		Context("eligible command", func() {

			When("passing non-empty envs", func() {
				It("should not return err", func() {
					cmd := "echo"
					err := sc.Execute([]string{cmd}, &profile, false)
					Expect(err).ToNot(HaveOccurred())
				})
			})

			When("pass wrong arg to command", func() {
				It("should return err", func() {
					cmd := []string{"cat", "/not-existing-dir/not-existing-file-rand-meow"}
					err := sc.Execute(cmd, &profile, false)
					Expect(err).To(HaveOccurred())
				})
			})
		})

		When("run non-existing command", func() {
			It("should not return err", func() {
				cmd := ""
				err := sc.Execute([]string{cmd}, &profile, false)
				Expect(err).To(HaveOccurred())
			})
		})

		When("skip-init param is true", func() {

			profile := config.NamedProfile{
				Name:    "my-profile",
				Profile: config.NewProfile(),
			}

			// The init script contains 'exit 1' which would normally cause an error,
			// but since skip-init=true, the script will be skipped entirely
			profile.InitScript = "exit 1"

			It("should not error because init-script won't be running", func() {
				cmd := []string{"echo", "hello"}
				err := sc.Execute(cmd, &profile, true)
				Expect(err).NotTo(HaveOccurred())
				Expect(stdout.String()).Should(ContainSubstring("hello"))
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
				err := sc.StartShell(&profile, false)
				Expect(err).ToNot(HaveOccurred())
			})
		})

		When("pass nil envs", func() {
			It("should not return err", func() {
				profile := config.NamedProfile{
					Name:    "",
					Profile: config.NewProfile(),
				}
				err := sc.StartShell(&profile, false)
				Expect(err).ToNot(HaveOccurred())
				Expect(stdout.String()).NotTo(BeEmpty())
				Expect(stderr.String()).To(BeEmpty())
			})
		})

		When("error occurred", func() {
			// backup original SHELL to set it back after test
			sh := os.Getenv("SHELL")

			JustBeforeEach(func() {
				// make SHELL empty for test case
				os.Setenv("SHELL", "")
			})
			It("it should not return err since it uses /bin/sh as default shell even SHELL is empty", func() {
				profile := config.NamedProfile{
					Name:    "my-profile",
					Profile: config.NewProfile(),
				}
				err := sc.StartShell(&profile, false)
				Expect(err).NotTo(HaveOccurred())
			})
			JustAfterEach(func() {
				// revert SHELL to original
				os.Setenv("SHELL", sh)
			})
		})
	})
})

var _ = Describe("parseEnvs function", func() {

	envs := config.Envs{}
	envs.AddEnv("PATH", "~/.config")
	envs.AddEnv("HOME", "$HOME")
	errs := parseEnvs(envs)
	h, _ := os.UserHomeDir()
	pe := appendEnvpProfile(envs.Strings(), "my-profile")

	It("should not occur error", func() {
		Expect(errs).ToNot(HaveOccurred())
	})

	When("has ~ in the value", func() {
		It("should be extracted to abs home dir", func() {
			Expect(pe).To(ContainElement(fmt.Sprintf("PATH=%s/.config", h)))
		})
	})

	When("has $HOME in the value", func() {
		It("should be extracted to abs home dir", func() {
			Expect(pe).To(ContainElement(fmt.Sprintf("HOME=%s", h)))
		})
	})

	When("append profile env var", func() {
		It("should include env var value of profile", func() {
			Expect(pe).To(ContainElement(fmt.Sprintf("%s=my-profile", envpEnvVarKey)))
		})
	})
})

var _ = Describe("parseEnvs with os.ExpandEnv", func() {
	When("var is referencing other vars", func() {
		envs := config.Envs{}
		envs.AddEnv("VAR1", "VAL_1")
		envs.AddEnv("VAR2", "$VAR1")
		envs.AddEnv("VAR3", "my-$VAR2")
		err := parseEnvs(envs)
		pe := envs.Strings()

		It("should not occur error", func() {
			Expect(err).ToNot(HaveOccurred())
		})

		It("expanding value of other vars", func() {
			Expect(pe).To(ContainElement("VAR2=VAL_1"))
			Expect(pe).To(ContainElement("VAR3=my-VAL_1"))
		})
	})
})

var _ = Describe("env shell command substitution", func() {

	envs := config.Envs{}
	envs.AddEnv("TEST_1", "VALUE_1")
	envs.AddEnv("TEST_SUBST_1", "$(echo hello)")
	envs.AddEnv("TEST_SUBST_2", "$(echo $TEST_1)")
	envs.AddEnv("TEST_SUBST_3", "$(this-is-error)")
	errs := parseEnvs(envs)

	profile := config.NamedProfile{
		Name:    "my-profile",
		Profile: config.NewProfile(),
	}
	profile.Env = envs

	When("has $() in the value", func() {
		It("should perform shell command substitution", func() {
			Expect(envs.Strings()).To(ContainElement("TEST_SUBST_1=hello"))
		})
	})

	When("referring another (earlier) env variable with command substitution", func() {
		It("should get the value of other env var value and substitute with it", func() {
			Expect(envs.Strings()).To(ContainElement("TEST_SUBST_2=VALUE_1"))
		})
	})

	When("command in substitution is wrong", func() {
		It("should occur error", func() {
			Expect(errs).To(HaveOccurred())
		})

		It("should not perform substitution. keep original value", func() {
			Expect(envs.Strings()).To(ContainElement("TEST_SUBST_3=$(this-is-error)"))
		})

		var stdout, stderr bytes.Buffer
		sc := NewShellCommand()
		sc.Stdout = &stdout
		sc.Stderr = &stderr

		It("StartShell should show parsing error message in stderr", func() {

			err := sc.StartShell(&profile, false)
			Expect(err).To(HaveOccurred())
			Expect(stderr.String()).NotTo(BeEmpty())
			Expect(stderr.String()).To(ContainSubstring("error processing value of TEST_SUBST_3"))
		})
	})
})

var _ = Describe("init-script", func() {

	When("init-script is defined", func() {

		var stdout, stderr bytes.Buffer
		sc := NewShellCommand()
		sc.Stdout = &stdout
		sc.Stderr = &stderr

		profile := config.NamedProfile{
			Name:    "my-profile",
			Profile: config.NewProfile(),
		}
		profile.InitScript = `echo "hello world"`

		err := sc.executeInitScript(&profile)

		It("should not error", func() {
			Expect(err).NotTo(HaveOccurred())
		})

		It("execute defined script", func() {
			Expect(stdout.String()).To(ContainSubstring("hello world"))
		})
	})

	When("init-script has multi line of shell script", func() {

		var stdout, stderr bytes.Buffer
		sc := NewShellCommand()
		sc.Stdout = &stdout
		sc.Stderr = &stderr

		profile := config.NamedProfile{
			Name:    "my-profile",
			Profile: config.NewProfile(),
		}
		profile.InitScript = `
		if [ 1 -gt 0 ]; then
			echo "hello world"
		else
			echo "wrong~"
		fi
		`

		err := sc.executeInitScript(&profile)

		It("should not error", func() {
			Expect(err).NotTo(HaveOccurred())
		})

		It("execute defined script", func() {
			Expect(stdout.String()).To(ContainSubstring("hello world"))
		})
	})

	When("init-script is wrong", func() {

		var stdout, stderr bytes.Buffer
		sc := NewShellCommand()
		sc.Stdout = &stdout
		sc.Stderr = &stderr

		profile := config.NamedProfile{
			Name:    "my-profile",
			Profile: config.NewProfile(),
		}

		profile.InitScript = "exit 1"
		err := sc.StartShell(&profile, false)

		It("should error", func() {
			Expect(err).To(HaveOccurred())
		})

		It("should show init-script error message in stderr", func() {
			Expect(stderr.String()).To(ContainSubstring("init-script error"))
		})
	})

	When("skip-init param is true", func() {

		var stdout, stderr bytes.Buffer
		sc := NewShellCommand()
		sc.Stdout = &stdout
		sc.Stderr = &stderr

		profile := config.NamedProfile{
			Name:    "my-profile",
			Profile: config.NewProfile(),
		}

		// The init script contains 'exit 1' which would normally cause an error,
		// but since skip-init=true, the script will be skipped entirely
		profile.InitScript = "exit 1"
		err := sc.StartShell(&profile, true)

		It("should not error because init-script won't be running", func() {
			Expect(err).NotTo(HaveOccurred())
		})
	})

	When("init-script use defined Env", func() {

		var stdout, stderr bytes.Buffer
		sc := NewShellCommand()
		sc.Stdout = &stdout
		sc.Stderr = &stderr

		profile := config.NamedProfile{
			Name:    "my-profile",
			Profile: config.NewProfile(),
		}
		profile.Env = []*config.Env{
			{Name: "MY_VAR", Value: "MY_VAL"},
		}
		profile.InitScript = "echo $MY_VAR"

		err := sc.executeInitScript(&profile)

		It("should not error", func() {
			Expect(err).NotTo(HaveOccurred())
		})

		It("execute defined script", func() {
			Expect(stdout.String()).To(ContainSubstring("MY_VAL"))
		})
	})

	When("profile is nil", func() {

		var stdout, stderr bytes.Buffer
		sc := NewShellCommand()
		sc.Stdout = &stdout
		sc.Stderr = &stderr

		err := sc.executeInitScript(nil)

		It("should not error", func() {
			Expect(err).NotTo(HaveOccurred())
		})

		It("does nothing", func() {
			Expect(stdout.String()).To(BeEmpty())
		})
	})

	When("init-script is empty", func() {

		var stdout, stderr bytes.Buffer
		sc := NewShellCommand()
		sc.Stdout = &stdout
		sc.Stderr = &stderr

		profile := config.NamedProfile{
			Name:    "my-profile",
			Profile: config.NewProfile(),
		}
		profile.InitScript = ""

		err := sc.executeInitScript(&profile)

		It("should not error", func() {
			Expect(err).NotTo(HaveOccurred())
		})

		It("does nothing", func() {
			Expect(stdout.String()).To(BeEmpty())
		})
	})
})

var _ = Describe("multiple init-script", func() {
	var stdout, stderr bytes.Buffer
	sc := NewShellCommand()
	sc.Stdout = &stdout
	sc.Stderr = &stderr

	When("multiple init-script is defined", func() {
		profile := config.NamedProfile{
			Name:    "my-profile",
			Profile: config.NewProfile(),
		}

		var initScripts []interface{}
		initScripts = append(initScripts, map[string]interface{}{"run": "echo meow-1"})
		initScripts = append(initScripts, map[string]interface{}{"run": "echo meow-2"})
		initScripts = append(initScripts, map[string]interface{}{"something-else": "echo meow-2"})

		profile.InitScript = initScripts

		err := sc.executeInitScript(&profile)

		It("should not error", func() {
			Expect(err).NotTo(HaveOccurred())
		})

		It("output should only have result of run(s)", func() {
			Expect(stdout.String()).To(Equal("meow-1\nmeow-2\n"))
		})
	})
})

var _ = Describe("nested profiles", func() {

	_ = os.Setenv(envpEnvVarKey, "profile-1")

	When("append another profile into env var", func() {
		envs := os.Environ()
		envs = appendEnvpProfile(envs, "profile-2")
		It("should include previous profile", func() {
			Expect(envs).To(ContainElement(fmt.Sprintf("%s=profile-1 > profile-2", envpEnvVarKey)))
		})
	})
})
