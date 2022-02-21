package config

import (
	"sort"
	"strings"
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

// ParseEnvFlagToMap parse string format "env=val" to map "env: val". it can be used fo dup check from slice of Env
func ParseEnvFlagToMap(envs []string) map[string]string {

	if len(envs) == 0 {
		return nil
	}

	r := map[string]string{}

	for _, s := range envs {
		ev := strings.Split(s, "=")
		if len(ev) != 2 {
			// TODO: handle unexpected format
			continue
		} else {
			r[ev[0]] = ev[1]
		}
	}
	return r
}

// ParseEnvFlagToEnv parse slice of string "var=val" to []ENV
func ParseEnvFlagToEnv(args []string) Envs {
	if len(args) == 0 {
		return nil
	}

	r := []Env{}

	for _, s := range args {
		ev := strings.Split(s, "=")
		if len(ev) != 2 {
			// TODO: handle unexpected format
			//fmt.Println("WARN: wrong format of env item. it must be var=val.", ev, "will be ignored")
			continue
		} else {
			r = append(r, Env{
				Name:  ev[0],
				Value: ev[1],
			})
		}
	}
	SortEnv(r)
	return r
}

// MapToEnv parse string map to slice of Env
func MapToEnv(m map[string]string) Envs {
	r := []Env{}
	for k, v := range m {
		r = append(r, Env{
			Name:  k,
			Value: v,
		})
	}
	// sort it by env name
	SortEnv(r)
	return r
}

// SortEnv sort []Env by name asc
func SortEnv(e []Env) {
	// sort it by env name
	sort.Slice(e, func(i, j int) bool {
		return e[i].Name < e[j].Name
	})
}
