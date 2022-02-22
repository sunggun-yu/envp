package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/sunggun-yu/envp/internal/prompt"
)

func init() {
	rootCmd.AddCommand(deleteCommand())
}

// example of delete command
func cmdExampleDelete() string {
	return `
  envp delete profile
  envp del another-profile
  `
}

// deleteCommand delete/remove environment variable profile and it's envionment variables from the config file
func deleteCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "delete profile-name",
		Short:             "Delete environment variable profile",
		Aliases:           []string{"del"},
		SilenceUsage:      true,
		Example:           cmdExampleDelete(),
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

			// set default="" if default profile is being deleted
			if profile.IsDefault {
				cfg.SetDefault("")
				color.Yellow("WARN: You are deleting default profile '%s'. please set default profile once it is deleted", profile.Name)
			}
			// ask y/n decision before proceed delete
			if prompt.PromptConfirm(fmt.Sprintf("Delete profile %s", color.RedString(profile.Name))) {
				// delete profile
				cfg.DeleteProfile(profile.Name)
			} else {
				fmt.Println("Cancelled")
				os.Exit(0)
			}
			if err := configFile.Save(); err != nil {
				return err
			}
			fmt.Println("Profile", profile.Name, "deleted successfully")
			return nil
		},
	}
	return cmd
}
