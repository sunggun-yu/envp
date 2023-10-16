package config

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testDataEnvs = Envs{
		&Env{Name: "VAR_C", Value: "VAL_C"},
		&Env{Name: "VAR_A", Value: "VAL_A"},
		&Env{Name: "VAR_D", Value: "VAL_D"},
		&Env{Name: "VAR_B", Value: "VAL_B"},
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
	SortEnv(envs)

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
		if ParseEnvFlagToMap([]string{}) != nil {
			t.Error("Not meet expectation. empty slice should return nil")
		}
	})

	t.Run("when data exist", func(t *testing.T) {
		testData := testDataArrStringFromFlag
		expected := testDataEnvMap
		actual := ParseEnvFlagToMap(testData)

		if !reflect.DeepEqual(expected, actual) {
			t.Error("Not meet expectation", expected, "-", actual)
		}
	})
}

// ParseEnvFlagToEnv should parse slice of string "var=val" to []ENV
func TestParseEnvFlagToEnv(t *testing.T) {

	t.Run("when set empty data", func(t *testing.T) {
		// nil data test
		if ParseEnvFlagToEnv([]string{}) != nil {
			t.Error("Not meet expectation. empty slice should return nil")
		}
	})

	t.Run("when data exist", func(t *testing.T) {
		testData := testDataArrStringFromFlag
		// invalid format should be ignored without error
		actual := ParseEnvFlagToEnv(testData)
		expected := testDataEnvs
		// ParseEnvFlagToEnv sort the result. so expected should be sorted
		SortEnv(expected)
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
	SortEnv(expected)
	actual := MapToEnv(testData)
	if !reflect.DeepEqual(expected, actual) {
		t.Error("Not meet expectation", expected, "-", actual)
	}
}

// TestEnvsStrings tests Strings func in Envs
func TestEnvsStrings(t *testing.T) {

	expected := []string{"VAR_C=VAL_C", "VAR_A=VAL_A", "VAR_D=VAL_D", "VAR_B=VAL_B"}
	assert.ElementsMatch(t, expected, testDataEnvs.Strings())
}

// TestEnvsAddEnv tests AddEnv func in Envs
func TestEnvsAddEnv(t *testing.T) {

	expected := []string{"VAR_1=VAL_1", "VAR_2=VAL_2", "VAR_3=VAL_3"}

	actual := Envs{}
	actual.AddEnv("VAR_1", "VAL_1")
	actual.AddEnv("VAR_2", "VAL_2")
	actual.AddEnv("VAR_3", "VAL_3")

	assert.ElementsMatch(t, expected, actual.Strings())
}
