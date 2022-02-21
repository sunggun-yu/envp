package config

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/sunggun-yu/envp/internal/util"
	"gopkg.in/yaml.v2"
)

var (
	// default config to be used when initiate the empty config file
	defaultConfig = Config{
		Default:  "",
		Profiles: &Profiles{},
	}
)

// ConfigFile is struct that representing the envp config file
// It contains Config and perform read and save operation
type ConfigFile struct {
	mu     sync.RWMutex
	name   string
	config *Config
	isInit bool
}

// NewConfigFile returns ConfigFile. it create the config file directory and file if not exists
func NewConfigFile(name string) (*ConfigFile, error) {

	// ensure if file is exsiting. if not create directory and file
	// expand and replace file path if it is refering home dir, `~`, `$HOME`
	p, err := util.EnsureConfigFilePath(filepath.Dir(name))
	if err != nil {
		return nil, err
	}
	filePath := filepath.Join(p, filepath.Base(name))

	cf := &ConfigFile{
		name: filePath,
	}
	// init config file if not exist
	if err := cf.initConfigFile(); err != nil {
		return nil, err
	}

	return cf, nil
}

// initConfigFile initiate the config file
// create config file directory if not exist
// create config file with minimal content if not exist
// set config file permission as 0600. so only owner can R/W
func (c *ConfigFile) initConfigFile() error {

	c.mu.Lock()
	defer c.mu.Unlock()

	// return error if config file name is not seet
	if c.name == "" {
		return fmt.Errorf("Config file is not set")
	}

	// open config file. create if not exist. and set file permission as 0600
	f, err := os.OpenFile(c.name, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	// close file once handling is done
	defer f.Close()

	// create default empty config file if file is empty (check by size)
	fs, _ := f.Stat()
	if fs.Size() == 0 {
		b, _ := yaml.Marshal(&defaultConfig)
		if _, err := f.Write(b); err != nil {
			return err
		}
	}

	// set init flag true
	c.isInit = true
	return nil
}

// Read read config file return Config
func (c *ConfigFile) Read() (*Config, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// unmarshal
	b, _ := os.ReadFile(c.name)
	if err := yaml.Unmarshal(b, &c.config); err != nil {
		return nil, err
	}
	if c.config == nil {
		return nil, fmt.Errorf("cannot read config")
	}
	// set mutex to Config to syncronize object along with file operation
	c.config.SetMutex(&c.mu)
	return c.config, nil
}

// Save saves config as a file
func (c *ConfigFile) Save() error {

	c.mu.Lock()
	defer c.mu.Unlock()

	// Marshal config
	b, err := yaml.Marshal(c.config)
	if err != nil {
		return err
	}

	// open file with read/write and trunc flag
	f, err := os.OpenFile(c.name, os.O_RDWR|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(b)
	if err != nil {
		return err
	}

	return nil
}
