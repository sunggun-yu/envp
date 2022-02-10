package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/sunggun-yu/envp/internal/config"
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
	var profile string
	// unmarshalled object from selected profile in the config file
	var currentProfile config.Profile

	cmd := &cobra.Command{
		Use:               "start profile-name",
		Short:             "Start new shell session with environment variable profile",
		SilenceUsage:      true,
		Example:           cmdExampleStart(),
		ValidArgsFunction: ValidArgsProfileList,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			// TODO: refactoring, cleanup
			switch {
			case len(args) == 0:
				// this case requires default profile.
				if viper.GetString(ConfigKeyDefaultProfile) == "" {
					printExample(cmd)
					return fmt.Errorf("default profile is not set. please set default profile")
				}
				profile = viper.GetString(ConfigKeyDefaultProfile)
			case len(args) > 0:
				profile = args[0]
			}
			// validate if selected profile is existing in the config
			selected := configProfiles.Sub(profile)
			// unmarshal to Profile
			err := selected.Unmarshal(&currentProfile)
			if err != nil {
				return fmt.Errorf("profile '%v' malformed configuration %e", profile, err)
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {

			// TODO: refactoring, cleanup

			// print start of session message
			fmt.Println(color.GreenString("Starting ENVP session..."), color.RedString(profile))
			color.Cyan(currentProfile.Env.String())
			fmt.Println("> press ctrl+d or type 'exit' to close session")

			// set ENVP_PROFILE env var to leverage profile info in the prompt, such as starship.
			os.Setenv(envpEnvVarKey, profile)

			// ignore error message from shell. let shell print out the errors
			shell.StartShell(currentProfile.Env)

			// print end of session message
			fmt.Println(color.GreenString("ENVP session closed..."), color.RedString(profile))
			return nil
		},
	}
	return cmd
}
