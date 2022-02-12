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

	testDataConfig = func() config.Config {
		file, _ := ioutil.ReadFile("testdata/config.yaml")

		var cfg config.Config
		err := yaml.Unmarshal(file, &cfg)
		if err != nil {
			panic(err)
		}
		return cfg
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

	t.Run("when set empty data", func(t *testing.T) {
		// nil data test
		if config.ParseEnvFlagToMap([]string{}) != nil {
			t.Error("Not meet expectation. empty slice should return nil")
		}
	})

	t.Run("when data exist", func(t *testing.T) {
		testData := testDataArrStringFromFlag
		expected := testDataEnvMap
		actual := config.ParseEnvFlagToMap(testData)

		if !reflect.DeepEqual(expected, actual) {
			t.Error("Not meet expectation", expected, "-", actual)
		}
	})
}

// ParseEnvFlagToEnv should parse slice of string "var=val" to []ENV
func TestParseEnvFlagToEnv(t *testing.T) {

	t.Run("when set empty data", func(t *testing.T) {
		// nil data test
		if config.ParseEnvFlagToEnv([]string{}) != nil {
			t.Error("Not meet expectation. empty slice should return nil")
		}
	})

	t.Run("when data exist", func(t *testing.T) {
		testData := testDataArrStringFromFlag
		// invalid format should be ignored without error
		actual := config.ParseEnvFlagToEnv(testData)
		expected := testDataEnvs
		// ParseEnvFlagToEnv sort the result. so expected should be sorted
		config.SortEnv(expected)
		if !reflect.DeepEqual(expected, actual) {
			t.Error("Not meet expectation", expected, "-", actual)
		}
	})
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

func TestDefaultProfile(t *testing.T) {
	cfg := testDataConfig()

	t.Run("when default profile exist", func(t *testing.T) {
		if p, err := cfg.DefaultProfile(); err != nil {
			t.Error("Should not be nil")
		} else {
			fmt.Println(p)
		}
	})

	t.Run("when default profile not exist", func(t *testing.T) {
		// make default empty
		cfg.Default = ""
		if _, err := cfg.DefaultProfile(); err == nil {
			t.Error("Should be error")
		} else {
			fmt.Println(err)
		}
	})
}

func TestProfile(t *testing.T) {
	cfg := testDataConfig()

	t.Run("should return current default when set empty string", func(t *testing.T) {
		if p, err := cfg.Profile(""); err != nil {
			t.Error("Should not be nil.")
		} else if p.Desc != "docker" {
			t.Error("Should not be same as default profile")
		}
	})

	t.Run("normal case", func(t *testing.T) {
		if p, err := cfg.Profile("docker"); err != nil {
			t.Error("Should not be nil.")
		} else if p == nil {
			t.Error("Should not be nil.")
		}
	})

	t.Run("when set non-existing profile", func(t *testing.T) {
		// make default empty and find profile that is not existing
		cfg.Default = ""
		if _, err := cfg.Profile("not-existing-profile"); err == nil {
			t.Error("Should be error")
		} else {
			fmt.Println(err)
		}
	})
}
