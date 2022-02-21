package config

import (
	"fmt"
	"os"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

// perform test whithin single process but multi thread operation
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

			c.SetDefault(strconv.Itoa(n)) // it may not guarentee the order
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
