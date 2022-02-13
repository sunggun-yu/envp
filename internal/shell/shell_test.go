package shell

import (
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/sunggun-yu/envp/internal/config"
)

var _ = Describe("Shell", func() {

	Describe("run Execute", func() {

		Context("eligible command", func() {
			When("passing non-empty envs", func() {
				It("should not return err", func() {
					cmd := "echo"
					err := Execute([]string{cmd}, []config.Env{
						{Name: "meow", Value: "woof"},
					})
					Expect(err).ToNot(HaveOccurred())
				})
			})

			When("pass wrong arg to command", func() {
				It("should return err", func() {
					cmd := []string{"cat", "/not-existing-dir/not-existing-file-rand-meow"}
					err := Execute(cmd, []config.Env{
						{Name: "meow", Value: "woof"},
					})
					Expect(err).To(HaveOccurred())
				})
			})
		})

		When("run non-existing command", func() {
			It("should not return err", func() {
				cmd := ""
				err := Execute([]string{cmd}, []config.Env{
					{Name: "meow", Value: "woof"},
				})
				Expect(err).To(HaveOccurred())
			})
		})

	})

	Describe("run StartShell", func() {

		// github action has no default SHELL. so set it as /bin/sh before each test case
		JustBeforeEach(func() {
			// make SHELL empty to occur error
			os.Setenv("SHELL", "/bin/sh")
		})

		When("pass not empty envs", func() {
			It("should not return err", func() {
				err := StartShell([]config.Env{
					{Name: "meow", Value: "woof"},
				})
				Expect(err).ToNot(HaveOccurred())
			})
		})

		When("pass nil envs", func() {
			It("should not return err", func() {
				err := StartShell(nil)
				Expect(err).ToNot(HaveOccurred())
			})
		})

		When("error occured", func() {
			// backup original SHELL to set it back after test
			sh := os.Getenv("SHELL")

			JustBeforeEach(func() {
				// make SHELL empty to occur error
				os.Setenv("SHELL", "")
			})
			It("should return err", func() {
				err := StartShell(nil)
				Expect(err).To(HaveOccurred())
			})
			JustAfterEach(func() {
				// revert SHELL to original
				os.Setenv("SHELL", sh)
			})
		})
	})
})
