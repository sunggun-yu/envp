package config

import (
	"fmt"
	"sort"
	"strings"
)

// Env represent environment variable name and value
// go yaml doesn't support capitalized key. so follow k8s env format
type Env struct {
	Name  string `mapstructure:"name" yaml:"name"`
	Value string `mapstructure:"value" yaml:"value"`
}

// Override String() to make it KEY=VAL format
func (e Env) String() string {
	return fmt.Sprint(e.Name, "=", e.Value)
}

// Envs is slice of Env
type Envs []*Env

func (e *Envs) AddEnv(name, value string) {
	env := &Env{Name: name, Value: value}
	*e = append(*e, env)
}

// Strings returns KEY=VAL array of Env
func (e Envs) Strings() []string {
	s := []string{}
	for _, i := range e {
		s = append(s, i.String())
	}
	return s
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

	r := []*Env{}

	for _, s := range args {
		ev := strings.Split(s, "=")
		if len(ev) != 2 {
			// TODO: handle unexpected format
			//fmt.Println("WARN: wrong format of env item. it must be var=val.", ev, "will be ignored")
			continue
		} else {
			r = append(r, &Env{
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
	r := []*Env{}
	for k, v := range m {
		r = append(r, &Env{
			Name:  k,
			Value: v,
		})
	}
	// sort it by env name
	SortEnv(r)
	return r
}

// SortEnv sort []Env by name asc
func SortEnv(e []*Env) {
	// sort it by env name
	sort.Slice(e, func(i, j int) bool {
		return e[i].Name < e[j].Name
	})
}
