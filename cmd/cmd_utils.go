package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/sunggun-yu/envp/internal/config"
)

// CurrentProfile is function that returns config.Profile
// it checks args and return default profile if args has no profile name
// otherwise speicified profile will be return
func CurrentProfile(args []string) (name string, profile *config.Profile, isDefault bool, err error) {
	switch {
	case len(args) > 0:
		name = args[0]
		profile, err = Config.Profiles.FindProfile(name)
		if Config.Default == name {
			isDefault = true
		}
	default:
		name = Config.Default
		isDefault = true
		profile, err = Config.DefaultProfile()
	}
	return name, profile, isDefault, err
}

// print command example
func printExample(cmd *cobra.Command) {
	fmt.Println("Example:")
	fmt.Println(cmd.Example)
}

// check the error type and print out command help for specific error types
func checkErrorAndPrintCommandExample(cmd *cobra.Command, err error) {
	switch err.(type) {
	case *config.DefaultProfileNotSetError:
		printExample(cmd)
	case *config.ProfileNameInputEmptyError:
		printExample(cmd)
	}
}

// set current status of Config into viper and save it to config file
func updateAndSaveConfigFile(cfg *config.Config, v *viper.Viper) error {

	v.Set("default", cfg.Default)
	v.Set("profiles", cfg.Profiles)

	if err := v.WriteConfig(); err != nil {
		return err
	}
	return nil
}
