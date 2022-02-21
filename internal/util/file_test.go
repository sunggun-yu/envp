package util

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ExpandHomeDir", func() {
	Context("has prefix of home dir", func() {
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

		When("error occurred getting $HOME", func() {
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

	When("has no prefix of home dir", func() {
		It("should return same path as original", func() {
			path := "/tmp/something/something"
			p, err := ExpandHomeDir(path)
			Expect(err).ToNot(HaveOccurred())
			Expect(p).To(Equal(path))
		})
	})
})

var _ = Describe("EnsureConfigFilePath", func() {

	When("working with existing directory", func() {
		Describe("set existing directory but non-existing dir", func() {
			It("should return same path as original", func() {
				file := "/tmp/some-file-not-existing.yaml"
				path, err := EnsureConfigFilePath(file)
				Expect(err).ToNot(HaveOccurred())
				Expect(path).To(Equal(file))
			})
		})

		When("set existing file as directory", func() {
			It("should return error", func() {

				r, _ := generateRandomString(7)
				file := fmt.Sprintf("/tmp/%s/some-file-not-existing.yaml", r)
				testDir := filepath.Dir(file)
				os.Create(testDir)
				// delete dir after test
				defer os.Remove(testDir)

				_, err := EnsureConfigFilePath(file)
				Expect(err).To(HaveOccurred())
				fmt.Println(err)
			})
		})
	})

	When("error occurred getting $HOME", func() {
		// backup original home path to set it back after test
		var originalHome string
		originalHome, _ = os.UserHomeDir()

		JustBeforeEach(func() {
			// make env Home empty to make error
			os.Setenv("HOME", "")
		})

		It("should return error and same path as input", func() {
			path := "$HOME/something/something"
			p, err := EnsureConfigFilePath(path)
			Expect(err).To(HaveOccurred())
			Expect(p).To(Equal(path))
		})

		JustAfterEach(func() {
			// revert HOME to original
			os.Setenv("HOME", originalHome)
		})
	})

	Context("working with non-existing directory", func() {
		When("set non-existing sub directory", func() {
			It("should return same path as original", func() {
				r, _ := generateRandomString(7)
				file := fmt.Sprintf("/tmp/%s/some-file-not-existing", r)
				testDir := filepath.Dir(file)
				// delete dir after test
				defer os.Remove(testDir)

				path, err := EnsureConfigFilePath(file)
				_, errCheck := os.Stat(testDir)

				Expect(err).ToNot(HaveOccurred())
				Expect(path).To(Equal(file))
				Expect(errCheck).ToNot(HaveOccurred())
			})
		})

		When("parent directory is existing but no permission", func() {
			It("should return error", func() {
				r, _ := generateRandomString(7)
				parent := fmt.Sprintf("/tmp/%s/", r)
				file := filepath.Join(parent, "child-dir/another-non-existing-child-dir")
				testDir := filepath.Dir(parent)
				// delete dir after test
				defer os.Remove(testDir)

				// no write permission
				os.Mkdir(parent, 0555)
				_, err := EnsureConfigFilePath(file)
				Expect(err).To(HaveOccurred())
			})
		})
	})
})

// generate random string
func generateRandomString(n int) (string, error) {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-"
	r := make([]byte, n)
	for i := 0; i < n; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return "", err
		}
		r[i] = letters[num.Int64()]
	}
	return string(r), nil
}
