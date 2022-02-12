package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
		ValidArgsFunction: ValidArgsProfileList,
		RunE: func(cmd *cobra.Command, args []string) error {

			name, _, isDefault, err := CurrentProfile(args)
			if err != nil {
				checkErrorAndPrintCommandExample(cmd, err)
				return err
			}

			// set default="" if default profile is being deleted
			if isDefault {
				Config.Default = ""
				color.Yellow("WARN: You are deleting default profile '%s'. please set default profile once it is deleted", name)
			}
			// ask y/n decision before proceed delete
			if prompt.PromptConfirm(fmt.Sprintf("Delete profile %s", color.RedString(name))) {
				// delete profile
				Config.Profiles.DeleteProfile(name)
			} else {
				fmt.Println("Cancelled")
				os.Exit(0)
			}

			// wait for the config file update and verify profile is added or not
			rc := make(chan error, 1)
			// I think underlying of viper.OnConfiChange is goroutine. but just run it as goroutine just in case
			// it's being watched in root initConfig - viper.WatchConfig()
			go viper.OnConfigChange(func(e fsnotify.Event) {
				if p, _ := Config.Profiles.FindProfile(name); p != nil {
					rc <- fmt.Errorf("profile %v not deleted", name)
					return
				}
				fmt.Println("Profile", name, "deleted successfully:", e.Name)
				rc <- nil
			})

			// update config and save
			if err := updateAndSaveConfigFile(&Config, viper.GetViper()); err != nil {
				return err
			}

			// wait for profile validation channel
			errOnChange := <-rc
			if errOnChange != nil {
				return errOnChange
			}
			return nil
		},
	}
	return cmd
}
