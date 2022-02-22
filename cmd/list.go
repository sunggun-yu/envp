package cmd

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
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
			cfg, err := configFile.Read()
			if err != nil {
				return err
			}

			// print profiles.
			for _, p := range cfg.ProfileNames() {
				if p == cfg.Default {
					// mark default profile with * and green
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
