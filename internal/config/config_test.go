package config_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/sunggun-yu/envp/internal/config"
	"gopkg.in/yaml.v2"
)

var (
	testDataConfig = func() config.Config {
		file, _ := os.ReadFile("../../testdata/config.yaml")

		var cfg config.Config
		err := yaml.Unmarshal(file, &cfg)
		if err != nil {
			panic(err)
		}
		return cfg
	}
)

func TestDefaultProfile(t *testing.T) {
	cfg := testDataConfig()

	t.Run("when default profile exist", func(t *testing.T) {
		if p, err := cfg.DefaultProfile(); err != nil {
			t.Error("Should not be nil")
		} else {
			fmt.Println(p)
		}
	})

	t.Run("when default profile is not set", func(t *testing.T) {
		// make default empty
		cfg.Default = ""
		if _, err := cfg.DefaultProfile(); err == nil {
			t.Error("Should be error")
		} else {
			fmt.Println(err)
		}
	})

	t.Run("when default profile not exist", func(t *testing.T) {
		// make default empty
		cfg.Default = "non-existing-profile-name"
		if _, err := cfg.DefaultProfile(); err == nil {
			t.Error("Should be error")
		} else {
			fmt.Println(err)
		}
	})
}

func TestProfile(t *testing.T) {
	cfg := testDataConfig()

	t.Run("should return nil and error when pass empty string", func(t *testing.T) {
		if p, err := cfg.Profile(""); err == nil {
			t.Error("Should occur error")
		} else if p != nil {
			t.Error("Should be nil.")
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
		if _, err := cfg.Profile("not-existing-profile"); err == nil {
			t.Error("Should be error")
		} else {
			fmt.Println(err)
		}
	})
}
