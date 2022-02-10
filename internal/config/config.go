package config

import (
	"fmt"
	"sort"
	"strings"
)

// cli config that includes profile config file path
type Config struct {
	Profiles Profiles `mapstructure:"profiles"` // profile config file path
}

// Profiles is map of profile. make item pointer to perform add/edit/delete
type Profiles map[string]*Profile

// Profile is struct of profile
type Profile struct {
	// set it with mapstructure remain to unmashal config file item `profiles` as Profile
	// yaml inline fixed the nested profiles issue
	Profiles `mapstructure:",remain" yaml:",inline"`
	Desc     string `mapstructure:"desc"`
	Env      Envs   `mapstructure:"env"`
}

// Envs is slice of Env
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

// Override strings() of Envs go generate comma separated string. this will be used for displaying env vars in list and start command.
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

// FindProfileByDotNotationKey finds profile from dot notation of key such as "a.b.c"
func FindProfileByDotNotationKey(key string, profiles Profiles) *Profile {
	if key == "" {
		return nil
	}
	keys := strings.Split(key, ".")
	current := profiles
	var profile *Profile
	for _, k := range keys {
		if p, ok := current[k]; ok {
			current = p.Profiles
			profile = p
		} else {
			return nil
		}
	}
	return profile
}
