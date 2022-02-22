package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// arg0NotExistingProfile ensure arg for profile name is existing in the profiles
func arg0NotExistingProfile() cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		// TODO: must check the length of args

		cfg, err := configFile.Read()
		if err != nil {
			return err
		}

		_, err = cfg.Profile(args[0])
		if err != nil {
			return err
		}
		return nil
	}
}

// arg0ExistingProfile ensure arg for profile name is not existing in the profiles
func arg0ExistingProfile() cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		cfg, err := configFile.Read()
		if err != nil {
			return err
		}
		if p, _ := cfg.Profile(args[0]); p != nil {
			return fmt.Errorf("%v is existing already", args[0])
		}
		return nil
	}
}

// arg0AsProfileName ensure to receive only 1 args as profile name
func arg0AsProfileName() cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("please specify the profile name")
		}
		if len(args) > 1 {
			return fmt.Errorf("space is not allowed for the profile name")
		}
		return nil
	}
}

// validArgsProfileList is for auto complete
var validArgsProfileList = func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	cfg, _ := configFile.Read()
	return cfg.ProfileNames(), cobra.ShellCompDirectiveNoFileComp
}
