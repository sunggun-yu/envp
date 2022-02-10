package config

import (
	"fmt"
	"sort"
	"strings"
)

// Profiles is map of profile. make item pointer to perform add/edit/delete
type Profiles map[string]*Profile

// Profile is struct of profile
// TODO: linked list might be better. but unmarshal may not support. so need to rebuild structure after reading the config.
type Profile struct {
	// set it with mapstructure remain to unmashal config file item `profiles` as Profile
	// yaml inline fixed the nested profiles issue
	Profiles `mapstructure:",remain" yaml:",inline"`
	Desc     string `mapstructure:"desc" yaml:"desc"`
	Env      Envs   `mapstructure:"env" yaml:"env"`
}

// Envs is slice of Env
type Envs []Env

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

// Override strings() of Envs go generate comma separated string. this will be used for displaying env vars in list and start command.
func (e Envs) String() string {
	s := []string{}
	for _, i := range e {
		s = append(s, i.String())
	}
	r := strings.Join(s, ",")
	return r
}

// FindProfile finds profile from dot notation of profile name such as "a.b.c"
func (p *Profiles) FindProfile(key string) (*Profile, error) {
	keys := strings.Split(key, ".")
	result := findProfileByDotNotationKey(keys, p)
	if result == nil {
		return nil, fmt.Errorf("profile %v is not exising", key)
	}
	return result, nil
}

// FindParentProfile ...
func (p *Profiles) FindParentProfile(key string) (*Profile, error) {
	keys := strings.Split(key, ".")
	if len(keys) == 1 {
		return nil, nil
	}
	result := findProfileByDotNotationKey(keys[:len(keys)-1], p)
	if result == nil {
		return nil, fmt.Errorf("parent profile of %v is not existing", key)
	}
	return result, nil
}

// ProfileNames list up all profile names in the Profiles
// nested profile's name will be formatted as "my-grand-parent.my-parent.me"
func (p *Profiles) ProfileNames() []string {
	// generate profile list
	profileList := *listProfileKeys("", *p, &[]string{})
	sort.Strings(profileList)
	return profileList
}

// DeleteProfile delete profile
func (p *Profiles) DeleteProfile(key string) error {

	keys := strings.Split(key, ".")

	var profile string
	var parent Profiles

	switch {
	case len(keys) == 1:
		profile = key
		parent = *p
	default:
		// last one should be the final profile name
		profile = keys[len(keys)-1]
		// get parent
		pp, err := p.FindParentProfile(key)
		if err != nil {
			// no parent means profile is not existing
			return fmt.Errorf("profile %v is not existing", key)
		}
		parent = pp.Profiles
	}

	delete(parent, profile)
	return nil
}

// FindProfileByDotNotationKey finds profile from dot notation of key such as "a.b.c"
// keys is array of string that in-order by nested profile. finding parent profile will be possible by keys[:len(keys)-1]
func findProfileByDotNotationKey(keys []string, profiles *Profiles) *Profile {
	current := *profiles
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

// list all the profiles in dot "." format. e.g. mygroup.my-subgroup.my-profile
// Do DFS to build viper keys for profiles
func listProfileKeys(key string, profiles Profiles, arr *[]string) *[]string {
	for k, v := range profiles {
		var s string
		if key == "" {
			s = k
		} else {
			s = fmt.Sprint(key, ".", k)
		}
		// only Profile item has env items will be considered as profile
		// even group(parent Profile that has children Profiles) will be considered as Profile if it has env items.
		if len(v.Env) > 0 {
			*arr = append(*arr, s)
		}
		// recursion
		listProfileKeys(s, v.Profiles, arr)
	}
	return arr
}
