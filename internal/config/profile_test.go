package config_test

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/sunggun-yu/envp/internal/config"
	"gopkg.in/yaml.v3"
)

var (
	testDataProfiles = func() *config.Profiles {
		cfg := testDataConfig()
		profiles := cfg.Profiles
		if profiles == nil {
			panic("profiles should not be nil")
		}
		return profiles
	}
)

// test Profiles.FindProfile method
func TestFindProfileByDotNotationKey(t *testing.T) {

	testData := testDataConfig()
	profiles := testData.Profiles

	t.Run("find non-nested profile", func(t *testing.T) {
		if p, err := profiles.FindProfile("docker"); p == nil && err != nil {
			t.Error("Should not be nil and error")
		} else if p.Desc != "docker" {
			t.Error("Not meet expectation")
			fmt.Println(p)
		}
	})

	t.Run("find nested profile", func(t *testing.T) {
		if p, err := profiles.FindProfile("org.nprod.argocd.argo2"); p == nil && err != nil {
			t.Error("Should not be nil and error")
		} else if p.Desc != "org.nprod.argocd.argo2" {
			t.Error("Not meet expectation")
			fmt.Println(p)
		}
	})

	t.Run("find non-existing profile", func(t *testing.T) {
		if p, err := profiles.FindProfile("org.nprod.vault"); p != nil && err == nil {
			t.Error("Should be nil and err")
		}
	})

	t.Run("find with empty string", func(t *testing.T) {
		if p, err := profiles.FindProfile(""); p != nil && err == nil {
			t.Error("Should be nil and err")
		}
	})

	t.Run("find with wonky format string", func(t *testing.T) {
		if p, err := profiles.FindProfile(".aaa..aaa"); p != nil && err == nil {
			t.Error("Should be nil and err")
		}
	})

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
		"parent-has-env",
		"profile-with-init-script",
		"profile-with-multi-init-script",
		"profile-with-multi-init-script-but-no-run",
		"profile-with-no-init-script",
		"profile-with-single-init-script-but-array",
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

	t.Run("find parent of existing child profile", func(t *testing.T) {
		// normal case
		testCaseNormal("lab.cluster1", "lab")
	})

	t.Run("find existing parent of non-existing child profile", func(t *testing.T) {
		// should return parent even child is not existing
		testCaseNormal("lab.cluster-not-existing-in-config", "lab")
	})

	t.Run("find parent of non-existing profile", func(t *testing.T) {
		// should return nil for non existing profile
		if p, err := profiles.FindParentProfile("non-existing-profile"); p != nil && err != nil {
			t.Error("supposed to be nill and no err")
		}
	})

	t.Run("find non-existing parent of non-existing child profile", func(t *testing.T) {
		// should return nil for non existing profile
		if p, err := profiles.FindParentProfile("non-existing-parent.non-existing-child"); p != nil && err == nil {
			t.Error("supposed to be nil and err")
		}
	})
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

		// error case
		err := profile.DeleteProfile("")
		if err == nil {
			t.Error("deleting empty string profile name should be error")
		} else if err.Error() == "" {
			t.Error("should return some error message")
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
	t.Run("delete non-nested profile", func(t *testing.T) {
		testCase("docker")
	})

	t.Run("delete nested profile", func(t *testing.T) {
		testCase("lab.cluster1")
		testCase("org.nprod.argocd.argo2")
	})

	t.Run("delete non-existing nested profile", func(t *testing.T) {
		testCaseNonExistingProfile("non-existing-parent.non-existing-child")
	})
}

// test SetProfile method
func TestSetProfile(t *testing.T) {

	profiles := testDataProfiles()

	var testCaseNormal = func(n, d string) {
		p := config.Profile{
			Desc: d,
		}
		err := profiles.SetProfile(n, p)
		if err != nil {
			t.Errorf("It shouldn't be err: %v", err)
		}

		// after set
		s, err := profiles.FindProfile(n)
		if err != nil {
			t.Errorf("Not updated: %v", err)
		}

		if s != nil && s.Desc != d {
			t.Errorf("Not updated: %v", s)
		}
	}

	// adding non-existing 1st level
	t.Run("adding non-existing 1st level", func(t *testing.T) {
		testCaseNormal("something", "something")
	})

	// adding into non-existing nested profile
	t.Run("adding into non-existing nested profile", func(t *testing.T) {
		testCaseNormal("some.thing", "hello")
	})

	// adding into non-existing nested profile - deeper
	t.Run("adding into non-existing nested profile - deeper", func(t *testing.T) {
		testCaseNormal("how.about.this.deep.case.meow.meow.woof.woof", "meow")
	})

	t.Run("adding into existing nested profile", func(t *testing.T) {
		// adding into existing nested profile
		testCaseNormal("org.nprod.argocd.argo100", "argocd")
		// sibling that was existing before appending should exist after
		s, _ := profiles.FindProfile("org.nprod.argocd.argo2")
		if s == nil {
			t.Error("sibling is not exist after appending")
		}
	})

	t.Run("adding into profile that has env info but no children", func(t *testing.T) {
		// adding into existing nested profile
		testCaseNormal("parent-has-env.new-child", "new-child")
	})

	t.Run("overwriting existing profile", func(t *testing.T) {
		// overwriting existing profile
		testCaseNormal("lab.cluster3", "updated lab.cluster3")
		// sibling that was existing before appending should exist after
		s, _ := profiles.FindProfile("lab.cluster2")
		if s == nil {
			t.Error("sibling is not exist after appending")
		}
	})

	t.Run("set empty string profile", func(t *testing.T) {
		err := profiles.SetProfile("", *config.NewProfile())
		if err == nil {
			t.Errorf("It supposed to be error")
		}
	})

	t.Run("marshalling", func(t *testing.T) {
		out, err := yaml.Marshal(profiles)
		if err != nil {
			t.Errorf("marshalling is failed after adding profiles")
		} else {
			fmt.Println(string(out))
		}
	})
}

func TestProfileInitScript(t *testing.T) {

	cfg := testDataConfig
	profile := cfg().Profiles

	t.Run("profile with no init-script", func(t *testing.T) {
		p, err := profile.FindProfile("profile-with-no-init-script")
		assert.NoError(t, err, "error should not occurred")
		assert.NotEmpty(t, p, "profile is found")
		expect := 0
		assert.Len(t, p.InitScripts(), expect, fmt.Sprintf("should be %v init-script", expect))
	})

	t.Run("profile with single init-script", func(t *testing.T) {
		p, err := profile.FindProfile("profile-with-init-script")
		assert.NoError(t, err, "error should not occurred")
		assert.NotEmpty(t, p, "profile is found")
		expect := 1
		assert.Len(t, p.InitScripts(), expect, fmt.Sprintf("should be %v init-script", expect))
	})

	t.Run("profile with single init-script but array type", func(t *testing.T) {
		p, err := profile.FindProfile("profile-with-single-init-script-but-array")
		assert.NoError(t, err, "error should not occurred")
		assert.NotEmpty(t, p, "profile is found")
		expect := 1
		assert.Len(t, p.InitScripts(), expect, fmt.Sprintf("should be %v init-script", expect))
	})

	t.Run("profile with multiple init-script", func(t *testing.T) {
		p, err := profile.FindProfile("profile-with-multi-init-script")
		assert.NoError(t, err, "error should not occurred")
		assert.NotEmpty(t, p, "profile is found")
		expect := 2
		assert.Len(t, p.InitScripts(), expect, fmt.Sprintf("should be %v init-script", expect))
	})

	t.Run("profile with multiple init-script but has no map of run keyword", func(t *testing.T) {
		p, err := profile.FindProfile("profile-with-multi-init-script-but-no-run")
		assert.NoError(t, err, "error should not occurred")
		assert.NotEmpty(t, p, "profile is found")
		expect := 0
		assert.Len(t, p.InitScripts(), expect, fmt.Sprintf("should be %v init-script", expect))
	})
}
