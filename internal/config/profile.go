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
	Profiles Profiles `mapstructure:",remain" yaml:",inline"`
	Desc     string   `mapstructure:"desc" yaml:"desc,omitempty"`
	Env      Envs     `mapstructure:"env" yaml:"env,omitempty"`
}

// NewProfile creates the Profile
func NewProfile() *Profile {
	return &Profile{
		Profiles: Profiles{},
	}
}

// Envs is slice of Env
type Envs []Env

// Env represent environment variable name and value
// go yaml doesn't support capitalized key. so follow k8s env format
type Env struct {
	Name  string `mapstructure:"name" yaml:"name"`
	Value string `mapstructure:"value" yaml:"value"`
}

// ProfileNotExistingError is an error when expected profile is not existing
type ProfileNotExistingError struct {
	profile string
}

// NewProfileNotExistingError create new ProfileNotExistingError
func NewProfileNotExistingError(profile string) *ProfileNotExistingError {
	return &ProfileNotExistingError{
		profile: profile,
	}
}

// Error is to make ProfileNotExistingError errors
func (e *ProfileNotExistingError) Error() string {
	return fmt.Sprintf("profile %s is not existing", e.profile)
}

// ProfileNameInputEmptyError is an error when mandatory profile input is empty
type ProfileNameInputEmptyError struct{}

// Error is to make ProfileNotExistingError errors
func (e *ProfileNameInputEmptyError) Error() string {
	return "input profile name is empty"
}

// NewProfileNameInputEmptyError create new ProfileNameInputEmptyError
func NewProfileNameInputEmptyError() *ProfileNameInputEmptyError {
	return &ProfileNameInputEmptyError{}
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

// SetProfile sets profile into the Profiles
// key is dot "." delimetered or plain string without no space.
// if it is dot delimeterd, considering it as nested profile
func (p *Profiles) SetProfile(key string, profile Profile) error {
	if key == "" {
		return NewProfileNameInputEmptyError()
	}
	keys := strings.Split(key, ".")
	if len(keys) == 1 {
		(*p)[keys[0]] = &profile
		return nil
	}

	// in case it's nested profile
	// build/get nested parents
	var parent *Profile
	// loop until last parent
	for _, k := range keys[:len(keys)-1] {
		if parent == nil {
			if (*p)[k] != nil {
				parent = (*p)[k]
			} else {
				parent = NewProfile()
				(*p)[k] = parent
			}
			continue
		}
		if parent.Profiles[k] == nil {
			parent.Profiles[k] = NewProfile()
		}
		parent = parent.Profiles[k]
	}
	// add last profile into last parent
	pname := keys[len(keys)-1]
	parent.Profiles[pname] = &profile
	return nil
}

// FindProfile finds profile from dot notation of profile name such as "a.b.c"
func (p *Profiles) FindProfile(key string) (*Profile, error) {
	if key == "" {
		return nil, NewProfileNameInputEmptyError()
	}
	keys := strings.Split(key, ".")
	result := findProfileByDotNotationKey(keys, p)
	if result == nil {
		return nil, NewProfileNotExistingError(key)
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
	if key == "" {
		return NewProfileNameInputEmptyError()
	}

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
			return NewProfileNotExistingError(key)
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
