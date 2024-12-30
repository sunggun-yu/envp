package cmd

import (
	"github.com/spf13/cobra"
	"github.com/sunggun-yu/envp/internal/shell"
)

func init() {
	// run with default shell command
	rootCmd.AddCommand(startCommand(shell.NewShellCommand()))
}

// flags struct for start command
type startFlags struct {
	skipInitScript bool
}

// example of delete command
func cmdExampleStart() string {
	return `
  # start new shell session with default profile
  envp start

  # start new shell session with specific profile
  envp start <profile-name>

  # skip "init-script" of profile when start new shell session with specific profile
  envp start <profile-name> --skip-init
  `
}

// deleteCommand delete/remove environment variable profile and it's envionment variables from the config file
func startCommand(sh *shell.ShellCommand) *cobra.Command {

	// add flags
	var flags startFlags

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
			sh.StartShell(profile, flags.skipInitScript)

			return nil
		},
	}

	cmd.Flags().BoolVarP(&flags.skipInitScript, "skip-init", "s", false, `Skip running initialization scripts from the profile's "init-script" during shell startup`)

	return cmd
}
