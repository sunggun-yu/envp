package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
		ValidArgsFunction: ValidArgsProfileList,
		RunE: func(cmd *cobra.Command, args []string) error {

			profile, err := currentProfile(args)
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
			Config.Default = profile.Name

			// wait for the config file update and verify profile is added or not
			rc := make(chan error, 1)
			// it's being watched in root initConfig - viper.WatchConfig()
			go viper.OnConfigChange(func(e fsnotify.Event) {
				if Config.Default != profile.Name {
					rc <- fmt.Errorf("default profile is not updated")
					return
				}
				fmt.Println("Default profile is set to", color.GreenString(Config.Default))
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
