package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/assert"
)

// perform test within single process but multi thread operation
func TestConfigFile(t *testing.T) {

	// assert
	assert := assert.New(t)

	testFile := "./testdata/multi-thread-test.yaml"
	defer os.Remove(testFile) // remove file after testing

	// when create ConfigFile by NewConfigFile
	cf, err := NewConfigFile(testFile)
	assert.NoError(err, "error should not occurred")
	assert.NotNil(cf, "it should not be nil")

	// when perform read/write config concurrently
	cases := 100
	var wg sync.WaitGroup
	for i := 1; i <= cases; i++ {
		wg.Add(1)
		// add profile as much as number of cases concurrently
		go func(n int) {
			c, err := cf.Read()

			assert.NoError(err, "error should not occurred on read operation")

			c.SetDefault(strconv.Itoa(n)) // it may not guarantee the order
			c.SetProfile(fmt.Sprintf("hello.world-%v", n), Profile{
				Desc: strconv.Itoa(n),
				Env: Envs{
					{Name: "VAR", Value: strconv.Itoa(n)},
				},
			})
			// when perform right after update config
			err = cf.Save()
			assert.NoError(err, "error should not occurred on save")
			wg.Done()
		}(i)
	}
	wg.Wait()

	// when read after save the config file
	c, err := cf.Read()
	assert.NoError(err, "error should not occurred on read after save config")

	// when validate the saved result
	ps := c.ProfileNames()
	assert.Equal(cases, len(ps), "saved profile should match with number of cases")

	// when delete profile
	err = c.DeleteProfile("hello.world-2")
	assert.NoError(err, "error should not occurred on delete")

	// when perform after update config
	err = cf.Save()
	assert.NoError(err, "error should not occurred on save")

	// when read after save the config file
	c, err = cf.Read()
	assert.NoError(err, "error should not occurred on read after save config")

	// when validate the saved result
	ps = c.ProfileNames()
	assert.Equal(cases-1, len(ps), "saved profile should match with number of cases")
}

func TestInitConfig(t *testing.T) {
	// assert
	assert := assert.New(t)

	t.Run("when pass empty string of config file path", func(t *testing.T) {
		_, err := NewConfigFile("")
		assert.Error(err, "config file should not empty string")
		fmt.Println(err)
	})

	t.Run("when create empty ConfigFile instance directly", func(t *testing.T) {
		cf := ConfigFile{}
		err := cf.initConfigFile()
		assert.Error(err, "config file should not empty")
	})
}

func TestRead(t *testing.T) {
	// assert
	assert := assert.New(t)
	testFile := fmt.Sprintf("%v", GinkgoRandomSeed())
	defer os.Remove(testFile) // remove file after testing

	t.Run("when create empty ConfigFile instance directly", func(t *testing.T) {
		cf, _ := NewConfigFile(testFile)

		// inject wrong format of yaml data into file
		wrongData := `default: {}
		profiles:
      - wrong
      - 1
		`
		os.WriteFile(testFile, []byte(wrongData), 0600)

		cf.config = nil
		_, err := cf.Read()
		assert.Error(err, "should occur error when have wrong format of config file")
	})
}

func TestWrite(t *testing.T) {
	// assert
	assert := assert.New(t)

	t.Run("when write without read - nil config", func(t *testing.T) {
		testFile := fmt.Sprintf("%v", GinkgoRandomSeed())
		defer os.Remove(testFile) // remove file after testing
		cf, _ := NewConfigFile(testFile)
		err := cf.Save()
		assert.Error(err, "should occur error when have wrong format of config file")
	})

	t.Run("when have no permission on config file", func(t *testing.T) {
		testFile := fmt.Sprintf("%v", GinkgoRandomSeed())
		defer os.Remove(testFile) // remove file after testing
		cf, _ := NewConfigFile(testFile)
		cf.Read()
		// make it read-only
		os.Chmod(testFile, 0400)
		err := cf.Save()
		assert.Error(err, "should occur error when have wrong format of config file")
	})
}

// ---------------------------------------------------------------------------
// Ginkgo test suite
// ---------------------------------------------------------------------------
var _ = Describe("NewConfigFile", func() {
	When("set existing directory as config file", func() {
		testFile := fmt.Sprintf("/tmp/%v/%v", GinkgoRandomSeed(), GinkgoRandomSeed())
		testDir := filepath.Dir(testFile)
		os.Create(testDir)
		defer os.Remove(testFile) // remove file after testing

		It("should return error", func() {
			_, err := NewConfigFile(testFile)
			Expect(err).To(HaveOccurred())
		})
	})

	When("error occurred getting $HOME", func() {
		// backup original home path to set it back after test
		var originalHome string
		originalHome, _ = os.UserHomeDir()
		testFile := fmt.Sprintf("$HOME/%v/%v.yaml", GinkgoRandomSeed(), GinkgoRandomSeed())
		defer os.Remove(testFile) // remove file after testing

		JustBeforeEach(func() {
			// make env Home empty to make error
			os.Setenv("HOME", "")
		})

		It("should return error", func() {
			_, err := NewConfigFile(testFile)
			Expect(err).To(HaveOccurred())
		})

		JustAfterEach(func() {
			// revert HOME to original
			os.Setenv("HOME", originalHome)
		})
	})

	When("config file is already existing", func() {
		testFile := "../../testdata/config.yaml"

		// copy test config file
		original, _ := os.ReadFile("../../testdata/config.yaml")

		_, err := NewConfigFile(testFile)

		It("should not return error", func() {
			Expect(err).NotTo(HaveOccurred())
		})

		It("should not change existing config content", func() {
			actual, err := os.ReadFile(testFile)
			Expect(err).NotTo(HaveOccurred())
			Expect(actual).To(Equal(original))
		})
	})
})
