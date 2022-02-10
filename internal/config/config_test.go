package config_test

import (
	"fmt"
	"io/ioutil"
	"reflect"
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

func TestFindProfileByDotNotationKey(t *testing.T) {

	cfg, _ := ioutil.ReadFile("testdata/config.yaml")

	var testData config.Config
	err := yaml.Unmarshal(cfg, &testData)
	if err != nil {
		t.Error(err)
	}

	if p := config.FindProfileByDotNotationKey("docker", testData.Profiles); p == nil {
		t.Error("Should not be nil")
	} else if p.Desc != "docker" {
		t.Error("Not meet expectation")
		fmt.Println(p)
	}

	// happy path
	if p := config.FindProfileByDotNotationKey("org.nprod.argocd.argo2", testData.Profiles); p == nil {
		t.Error("Should not be nil")
	} else if p.Desc != "org.nprod.argocd.argo2" {
		t.Error("Not meet expectation")
		fmt.Println(p)
	}

	// not existing key
	if p := config.FindProfileByDotNotationKey("org.nprod.vault", testData.Profiles); p != nil {
		t.Error("Should be nil")
	}

	// empty string
	if p := config.FindProfileByDotNotationKey("", testData.Profiles); p != nil {
		t.Error("Should be nil")
	}

	// messed format
	if p := config.FindProfileByDotNotationKey(".aaa..aaa", testData.Profiles); p != nil {
		t.Error("Should be nil")
	}

	// pointer check
	testChangeData := "changed"
	p := config.FindProfileByDotNotationKey("docker", testData.Profiles)
	p.Desc = testChangeData
	if testData.Profiles["docker"].Desc != testChangeData {
		t.Error("nested item should be pointer")
	}
}
