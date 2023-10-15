package config

import (
	"sync"
)

// Config is struct that represents configuration of config file
type Config struct {
	mu       *sync.RWMutex
	Default  string    `mapstructure:"default" yaml:"default"`
	Profiles *Profiles `mapstructure:"profiles" yaml:"profiles"`
}

// DefaultProfileNotSetError is error when default profile is not set
type DefaultProfileNotSetError struct{}

// NewDefaultProfileNotSetError create new DefaultProfileNotSetError
func NewDefaultProfileNotSetError() *DefaultProfileNotSetError {
	return &DefaultProfileNotSetError{}
}

// Error is to make ProfileNotExistingError errors
func (e *DefaultProfileNotSetError) Error() string {
	return "default profile is not set"
}

// SetMutex set the pointer of RWMutex
func (c *Config) SetMutex(m *sync.RWMutex) {
	c.mu = m
}

// initMutex init RWMutext for Config
// mutex of Config instance is pointer to make it possible to set mutext from ConfigFile for sync operation with file handliing
func (c *Config) initMutex() {
	if c.mu == nil {
		c.SetMutex(new(sync.RWMutex))
	}
}

// DefaultProfile returns default profile of config. it returns DefaultProfileNotSetError when default file is not set
func (c *Config) DefaultProfile() (*NamedProfile, error) {
	if c.Default == "" {
		return nil, NewDefaultProfileNotSetError()
	}
	name := c.Default
	p, err := c.Profiles.FindProfile(name)
	if err != nil {
		return nil, err
	}
	profile := NamedProfile{
		Profile:   p,
		Name:      name,
		IsDefault: true,
	}
	return &profile, nil
}

// Profile find and return NamedProfile of name
// IsDefault will be true if the profile is same as default
func (c *Config) Profile(name string) (*NamedProfile, error) {
	c.initMutex()
	c.mu.Lock()
	defer c.mu.Unlock()
	p, err := c.Profiles.FindProfile(name)
	if err != nil {
		return nil, err
	}
	return &NamedProfile{
		Profile:   p,
		Name:      name,
		IsDefault: c.Default == name,
	}, nil
}

// SetDefault set the default profile
func (c *Config) SetDefault(key string) {
	c.initMutex()
	c.mu.Lock()
	defer c.mu.Unlock()
	c.Default = key
}

// SetProfile set the profile in the Config
func (c *Config) SetProfile(key string, profile Profile) error {
	c.initMutex()
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.Profiles.SetProfile(key, profile)
}

// DeleteProfile delete profile
func (c *Config) DeleteProfile(key string) error {
	c.initMutex()
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.Profiles.DeleteProfile(key)
}

// ProfileNames list up all profile names in the Config Profiles
func (c *Config) ProfileNames() []string {
	c.initMutex()
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.Profiles.ProfileNames()
}
