package util

import (
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("File", func() {
	Describe("has prefix of home dir", func() {
		When("with value ~", func() {
			It("should return home dir", func() {
				p, err := ExpandHomeDir("~")
				Expect(err).ToNot(HaveOccurred())
				Expect(p).ToNot(BeEmpty())
			})
		})

		When("with value $HOME", func() {
			It("should return home dir", func() {
				p, err := ExpandHomeDir("$HOME")
				Expect(err).ToNot(HaveOccurred())
				Expect(p).ToNot(BeEmpty())
			})
		})

		When("error occured getting $HOME", func() {
			// backup original home path to set it back after test
			var originalHome string
			originalHome, _ = os.UserHomeDir()

			JustBeforeEach(func() {
				// make env Home empty to make error
				os.Setenv("HOME", "")
			})

			It("should return error and same path as input", func() {
				path := "$HOME/something/something"
				p, err := ExpandHomeDir(path)
				Expect(err).To(HaveOccurred())
				Expect(p).To(Equal(path))
			})

			JustAfterEach(func() {
				// revert HOME to original
				os.Setenv("HOME", originalHome)
			})
		})
	})

	Describe("has no prefix of home dir", func() {
		It("should return same path as original", func() {
			path := "/tmp/something/something"
			p, err := ExpandHomeDir(path)
			Expect(err).ToNot(HaveOccurred())
			Expect(p).To(Equal(path))
		})
	})
})
