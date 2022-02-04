package config

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const (
	configRoot     string = ".config/prw"
	configFileName string = "config.yaml"
)

// cli config that includes profile config file path
type Config struct {
	Profile string `mapstructure:"profile"` // profile config file path
}

type Profile struct {
	Desc    string `mapstructure:"desc"`
	Host    string `mapstructure:"host"`    // TODO: deprecate
	NoProxy string `mapstructure:"noproxy"` // TODO: deprecate
	Env     []Env  `mapstructure:"env"`
}

// go yaml doesn't support capitalized key. so follow k8s env format
type Env struct {
	Name  string
	Value string
}

// TODO: working on it
func ConfigFile() (string, error) {
	// read config file
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	configPath := filepath.Join(home, configRoot)
	configFile := filepath.Join(configPath, configFileName)

	// create config file if not exists
	if _, err := os.Stat(configPath); err != nil {
		if os.IsNotExist(err) {
			os.MkdirAll(configPath, 0755)
		} else {
			return "", err
		}
	}
	return configFile, nil
}

// parse string format "env=val" to map "env: val". it can be used fo dup check from slice of Env
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

// Parse string map to slice of Env
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

func SortEnv(e []Env) {
	// sort it by env name
	sort.Slice(e, func(i, j int) bool {
		return e[i].Name < e[j].Name
	})
}
