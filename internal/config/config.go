package config

import (
	"fmt"
	"sort"
	"strings"
)

// cli config that includes profile config file path
type Config struct {
	Profile string `mapstructure:"profile"` // profile config file path
}

// Profiles is struct of profiles
type Profiles map[string]Profile

// Profile is struct of profile
type Profile struct {
	// set it with mapstructure remain to unmashal config file item `profiles` as Profile
	// yaml inline fixed the nested profiles issue
	Profiles `mapstructure:",remain" yaml:",inline"`
	Desc     string `mapstructure:"desc"`
	Env      Envs   `mapstructure:"env"`
}

type Envs []Env

// Env represent environment variable name and value
// go yaml doesn't support capitalized key. so follow k8s env format
type Env struct {
	Name  string
	Value string
}

// Override String() to make it KEY=VAL format
func (e Env) String() string {
	return fmt.Sprint(e.Name, "=", e.Value)
}

func (e Envs) String() string {
	s := []string{}
	for _, i := range e {
		s = append(s, i.String())
	}
	r := strings.Join(s, ",")
	return r
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
			//fmt.Println("WARN: wrong format of env item. it must be var=val.", ev, "will be ignored")
		} else {
			r[ev[0]] = ev[1]
		}
	}

	return r
}

// ParseEnvFlagToEnv parse slice of string "var=val" to []ENV
func ParseEnvFlagToEnv(envs []string) []Env {
	if len(envs) == 0 {
		return nil
	}

	r := []Env{}

	for _, s := range envs {
		ev := strings.Split(s, "=")
		if len(ev) != 2 {
			// TODO: handle unexpected format
			//fmt.Println("WARN: wrong format of env item. it must be var=val.", ev, "will be ignored")
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
func MapToEnv(m map[string]string) []Env {
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
