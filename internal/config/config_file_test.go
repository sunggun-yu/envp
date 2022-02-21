package config

import (
	"fmt"
	"os"
	"strconv"
	"sync"
	"testing"
)

// perform test whithin single process but multi thread operation
func TestConfigFile(t *testing.T) {
	testFile := "./testdata/multi-thread-test.yaml"
	defer os.Remove(testFile) // remove file after testing

	// when create ConfigFile by NewConfigFile
	cf, err := NewConfigFile(testFile)
	if err != nil {
		t.Error("error should not occurred", err)
	}
	if cf == nil {
		t.Error("it should not nil")
	}

	// when perform read/write config concurrently
	cases := 100
	var wg sync.WaitGroup
	for i := 1; i <= cases; i++ {
		wg.Add(1)
		// add profile as much as number of cases concurrently
		go func(n int) {
			c, err := cf.Read()
			if err != nil {
				t.Error("error should not occurred on read operation", err)
			}
			c.SetDefault(strconv.Itoa(n)) // it may not guarentee the order
			c.SetProfile(fmt.Sprintf("hello.world-%v", n), Profile{
				Desc: strconv.Itoa(n),
				Env: Envs{
					{Name: "VAR", Value: strconv.Itoa(n)},
				},
			})
			wg.Done()
		}(i)
	}
	wg.Wait()

	// when perform after update config
	err = cf.Save()
	if err != nil {
		t.Error("error should not occurred on save", err)
	}

	// when read after save the config file
	c, err := cf.Read()
	if err != nil {
		t.Error("error should not occurred on read after save config")
	}

	// when validate the saved result
	ps := c.ProfileNames()
	if len(ps) != cases {
		t.Errorf("saved profile should match with number of cases. actual %v, expected %v", len(ps), cases)
	}

	// when delete profile
	if err := c.DeleteProfile("hello.world-2"); err != nil {
		t.Error("error should not occurred on delete")
	}

	// when perform after update config
	err = cf.Save()
	if err != nil {
		t.Error("error should not occurred on save", err)
	}

	// when read after save the config file
	c, err = cf.Read()
	if err != nil {
		t.Error("error should not occurred on read after save config")
	}

	// when validate the saved result
	ps = c.ProfileNames()
	if len(ps) != cases-1 {
		t.Errorf("saved profile should match with number of cases. actual %v, expected %v", len(ps), cases)
	}
}
