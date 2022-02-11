package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func Arg0NotExistingProfile() cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		_, err := Config.Profiles.FindProfile(args[0])
		if err != nil {
			return err
		}
		return nil
	}
}

func Arg0ExistingProfile() cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		//TODO: what was this? lol
		// profiles := configProfiles
		// if profiles == nil || len(profiles.AllKeys()) == 0 {
		// 	return nil
		// }
		if p, _ := Config.Profiles.FindProfile(args[0]); p != nil {
			return fmt.Errorf("%v is existing already", args[0])
		}
		return nil
	}
}

func Arg0AsProfileName() cobra.PositionalArgs {
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

// ValidArgsProfileList is for auto complete
var ValidArgsProfileList = func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	return Config.Profiles.ProfileNames(), cobra.ShellCompDirectiveNoFileComp
}
