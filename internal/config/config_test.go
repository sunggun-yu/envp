package config_test

import (
	"fmt"
	"io/ioutil"
	"reflect"
	"strings"
	"testing"

	"github.com/sunggun-yu/envp/internal/config"
	"gopkg.in/yaml.v2"
)

var (
	testDataEnvs = config.Envs{
		config.Env{Name: "VAR_C", Value: "VAL_C"},
		config.Env{Name: "VAR_A", Value: "VAL_A"},
		config.Env{Name: "VAR_D", Value: "VAL_D"},
		config.Env{Name: "VAR_B", Value: "VAL_B"},
	}

	testDataEnvMap = map[string]string{
		"VAR_A": "VAL_A",
		"VAR_B": "VAL_B",
		"VAR_C": "VAL_C",
		"VAR_D": "VAL_D",
	}

	testDataArrStringFromFlag = []string{
		"VAR_A=VAL_A",
		"something_not_valid", // should be ignored
		"VAR_B=VAL_B",
		"VAR_C=VAL_C",
		"VAR_D=VAL_D",
		"not:valid",      // should be ignored
		" ",              // should be ignored
		"how=about=this", // should be ignored
	}

	testDataConfig = func() config.Config {
		file, _ := ioutil.ReadFile("testdata/config.yaml")

		var cfg config.Config
		err := yaml.Unmarshal(file, &cfg)
		if err != nil {
			panic(err)
		}
		return cfg
	}

	testDataProfiles = func() *config.Profiles {
		cfg := testDataConfig()
		profiles := cfg.Profiles
		if profiles == nil {
			panic("profiles should not be nil")
		}
		return profiles
	}
)

// test String() method and SortEnv
func TestEnv(t *testing.T) {
	envs := testDataEnvs
	// sort
	config.SortEnv(envs)

	// data must be sorted in key
	// Env should return string in VAR=VAL format
	// Envs should return comma separated string
	expected := "VAR_A=VAL_A,VAR_B=VAL_B,VAR_C=VAL_C,VAR_D=VAL_D"
	actual := envs.String()
	if expected != actual {
		t.Error("Not meet expectation", expected, "-", actual)
	}
}

// test ParseEnvFlagToMap func
// ParseEnvFlagToMap should parse string format "env=val" to map "env: val"
func TestParseEnvFlagToMap(t *testing.T) {

	// nil data test
	if config.ParseEnvFlagToMap([]string{}) != nil {
		t.Error("Not meet expectation. empty slice should return nil")
	}

	testData := testDataArrStringFromFlag
	expected := testDataEnvMap
	actual := config.ParseEnvFlagToMap(testData)

	if !reflect.DeepEqual(expected, actual) {
		t.Error("Not meet expectation", expected, "-", actual)
	}
}

// ParseEnvFlagToEnv should parse slice of string "var=val" to []ENV
func TestParseEnvFlagToEnv(t *testing.T) {

	// nil data test
	if config.ParseEnvFlagToEnv([]string{}) != nil {
		t.Error("Not meet expectation. empty slice should return nil")
	}

	testData := testDataArrStringFromFlag
	// invalid format should be ignored without error
	actual := config.ParseEnvFlagToEnv(testData)
	expected := testDataEnvs
	// ParseEnvFlagToEnv sort the result. so expected should be sorted
	config.SortEnv(expected)
	if !reflect.DeepEqual(expected, actual) {
		t.Error("Not meet expectation", expected, "-", actual)
	}

	// nil test
}

// test MapToEnv func
func TestMapToEnv(t *testing.T) {
	testData := testDataEnvMap
	expected := testDataEnvs
	// sort. MapToEnv sort the result. so expected should be sorted
	config.SortEnv(expected)
	actual := config.MapToEnv(testData)
	if !reflect.DeepEqual(expected, actual) {
		t.Error("Not meet expectation", expected, "-", actual)
	}
}

// test Profiles.FindProfile method
func TestFindProfileByDotNotationKey(t *testing.T) {

	testData := testDataConfig()
	profiles := testData.Profiles

	if p, err := profiles.FindProfile("docker"); p == nil && err != nil {
		t.Error("Should not be nil and error")
	} else if p.Desc != "docker" {
		t.Error("Not meet expectation")
		fmt.Println(p)
	}

	// happy path
	if p, err := profiles.FindProfile("org.nprod.argocd.argo2"); p == nil && err != nil {
		t.Error("Should not be nil and error")
	} else if p.Desc != "org.nprod.argocd.argo2" {
		t.Error("Not meet expectation")
		fmt.Println(p)
	}

	// not existing key
	if p, err := profiles.FindProfile("org.nprod.vault"); p != nil && err == nil {
		t.Error("Should be nil and err")
	}

	// empty string
	if p, err := profiles.FindProfile(""); p != nil && err == nil {
		t.Error("Should be nil and err")
	}

	// wonky format
	if p, err := profiles.FindProfile(".aaa..aaa"); p != nil && err == nil {
		t.Error("Should be nil and err")
	}

	// pointer check
	testChangeData := "changed"
	p, _ := profiles.FindProfile("docker")
	p.Desc = testChangeData
	fmt.Println(testData)
	if (*testData.Profiles)["docker"].Desc != testChangeData {
		t.Error("nested item should be pointer")
	}
}

// test case for Profiles.ProfileNames
func TestProfileNames(t *testing.T) {
	profiles := testDataProfiles()
	expected := []string{
		"docker",
		"lab.cluster1",
		"lab.cluster2",
		"lab.cluster3",
		"org.nprod.argocd.argo1",
		"org.nprod.argocd.argo2",
		"org.nprod.vpn.vpn1",
		"org.nprod.vpn.vpn2",
	}

	actual := profiles.ProfileNames()
	if !reflect.DeepEqual(expected, actual) {
		t.Error("Not meet expectation", expected, "-", actual)
	}
}

// testing FindParentProfile
func TestFindParentProfile(t *testing.T) {
	profiles := testDataProfiles()

	var testCaseNormal = func(child, parent string) {
		pp, _ := profiles.FindParentProfile(child)
		p, _ := profiles.FindProfile(parent)
		if pp != p {
			t.Error("supposed to be same instance")
		}
	}

	// normal case
	testCaseNormal("lab.cluster1", "lab")
	// should return parent even child is not exisiting
	testCaseNormal("lab.cluster-not-exising-in-config", "lab")

	// should return nil for non existing profile
	if p, err := profiles.FindParentProfile("non-exising-profile"); p != nil && err != nil {
		t.Error("supposed to be nill and no err")
	}

	// should return nil for non existing profile
	if p, err := profiles.FindParentProfile("non-existing-parent.non-existing-child"); p != nil && err == nil {
		t.Error("supposed to be nil and err")
	}
}

func TestDeleteProfile(t *testing.T) {
	cfg := testDataConfig
	profile := cfg().Profiles

	var testCase = func(key string) {
		// check before
		if p, _ := profile.FindProfile(key); p == nil {
			t.Error("It should not be nil before deleting")
		}
		// delete
		profile.DeleteProfile(key)

		// check after
		if p, _ := profile.FindProfile(key); p != nil {
			t.Error("Profile should be nil after deleting")
		}

		if len(strings.Split(key, ".")) > 1 {
			if p, _ := profile.FindParentProfile(key); p == nil {
				t.Error("Parent should not be nil after deleting")
			}
		}
	}

	var testCaseNonExistingProfile = func(key string) {
		// delete
		err := profile.DeleteProfile(key)
		if err == nil {
			t.Error("It should be error for deleting non existing profile")
		}
	}

	// test case for non-nested profile
	testCase("docker")
	testCase("lab.cluster1")
	testCase("org.nprod.argocd.argo2")
	testCaseNonExistingProfile("non-existing-parent.non-existing-child")
}
