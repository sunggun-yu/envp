package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(useCommand())
}

// example of use command
func cmdExampleUse() string {
	return `
  # set default profile to profile-name
  envp use profile-name
  
  # env vars in the default profile will be set during command execution
  envp -- kubectl get pods
  `
}

// add command
func useCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "use profile-name",
		Short:             "Set default environment variable profile",
		SilenceUsage:      true,
		Example:           cmdExampleUse(),
		ValidArgsFunction: validArgsProfileList,
		RunE: func(cmd *cobra.Command, args []string) error {

			cfg, err := configFile.Read()
			if err != nil {
				return err
			}
			profile, err := currentProfile(cfg, args)
			if err != nil {
				checkErrorAndPrintCommandExample(cmd, err)
				return err
			}
			// just exit if selected profile is already default
			if profile.IsDefault {
				fmt.Println("Profile", profile.Name, "is alreday set as default")
				os.Exit(0)
			}

			// set selected profile as default
			cfg.SetDefault(profile.Name)

			if err := configFile.Save(); err != nil {
				return err
			}

			fmt.Println("Default profile is set to", color.GreenString(cfg.Default))

			return nil
		},
	}
	return cmd
}
