package cmd

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(listCommand())
}

// example of edit command
func cmdExampleList() string {
	return `
  envp list
  envp ls
  `
}

// listCommand prints out list of environment variable profiles
func listCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "list",
		Short:        "List all environment variable profiles",
		Aliases:      []string{"ls"},
		SilenceUsage: true,
		Example:      cmdExampleList(),
		RunE: func(cmd *cobra.Command, args []string) error {

			// current default profile name to compare
			defaultProfile := viper.GetString(ConfigKeyDefaultProfile)
			// print profiles. mark default profile with *
			for _, p := range Config.Profiles.ProfileNames() {
				if p == defaultProfile {
					color.Green("* %s", p)
				} else {
					fmt.Println(" ", p)
				}
			}
			return nil
		},
	}
	return cmd
}
