package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func Arg0NotExistingProfile() cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		selected := configProfiles.Sub(args[0])
		if selected == nil {
			return fmt.Errorf("%v is not existing in the profile list", args[0])
		}
		return nil
	}
}

func Arg0ExistingProfile() cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		profiles := configProfiles
		if profiles == nil || len(profiles.AllKeys()) == 0 {
			return nil
		}
		selected := profiles.Sub(args[0])
		if selected != nil {
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

var ValidArgsProfileList = func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	return profileList, cobra.ShellCompDirectiveNoFileComp
}
