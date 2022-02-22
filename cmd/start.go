package cmd

import (
	"github.com/spf13/cobra"
	"github.com/sunggun-yu/envp/internal/shell"
)

func init() {
	rootCmd.AddCommand(startCommand())
}

// example of delete command
func cmdExampleStart() string {
	return `
  # start new shell session with default profile
  envp start

  # start new shell session with specific profile
  envp start <profile-name>
  `
}

// deleteCommand delete/remove environment variable profile and it's envionment variables from the config file
func startCommand() *cobra.Command {

	cmd := &cobra.Command{
		Use:               "start profile-name",
		Short:             "Start new shell session with environment variable profile",
		SilenceUsage:      true,
		Example:           cmdExampleStart(),
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

			// ignore error message from shell. let shell print out the errors
			sc := shell.NewShellCommand()
			sc.StartShell(profile.Env, profile.Name)

			return nil
		},
	}
	return cmd
}
