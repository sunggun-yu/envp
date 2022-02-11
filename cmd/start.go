package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/sunggun-yu/envp/internal/shell"
)

func init() {
	rootCmd.AddCommand(startCommand())
}

const envpEnvVarKey = "ENVP_PROFILE"

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
		ValidArgsFunction: ValidArgsProfileList,
		RunE: func(cmd *cobra.Command, args []string) error {

			name, profile, _, err := CurrentProfile(args)
			if err != nil {
				checkErrorAndPrintCommandExample(cmd, err)
				return err
			}

			// print start of session message
			fmt.Println(color.GreenString("Starting ENVP session..."), color.RedString(name))
			color.Cyan(profile.Env.String())
			fmt.Println("> press ctrl+d or type 'exit' to close session")

			// set ENVP_PROFILE env var to leverage profile info in the prompt, such as starship.
			os.Setenv(envpEnvVarKey, name)

			// ignore error message from shell. let shell print out the errors
			shell.StartShell(profile.Env)

			// print end of session message
			fmt.Println(color.GreenString("ENVP session closed..."), color.RedString(name))
			return nil
		},
	}
	return cmd
}
