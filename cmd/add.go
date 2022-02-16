package cmd

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/sunggun-yu/envp/internal/config"
)

// flags struct for add command
type addFlags struct {
	desc string
	env  []string
}

func init() {
	rootCmd.AddCommand(addCommand())
}

// example of add command
func cmdExampleAdd() string {
	return `
  envp add my-proxy \
    -d 'profile desc' \
    -e 'VAR=VAL' \
    -e HTTPS_PROXY=http://some-proxy:3128 \
    -e 'NO_PROXY=127.0.0.1,localhost' \
    -e 'DOCKER_HOST=ssh://myuser@some-server'
  `
}

// addCommand add/create environment variable profile and it's envionment variables in the config file
func addCommand() *cobra.Command {
	// add flags
	var flags addFlags

	// add command
	cmd := &cobra.Command{
		Use:          "add profile-name",
		Short:        "Add environment variable profile",
		SilenceUsage: true,
		Example:      cmdExampleAdd(),
		Args: cobra.MatchAll(
			Arg0AsProfileName(),
			Arg0ExistingProfile(),
		),
		RunE: func(cmd *cobra.Command, args []string) error {

			name := args[0]
			profile := config.Profile{
				Desc: flags.desc,
				Env:  []config.Env{},
			}
			profile.Env = config.ParseEnvFlagToEnv(flags.env)

			// set profile
			Config.Profiles.SetProfile(name, profile)

			// set profile as default profile if default is empty and no profile is existing
			if Config.Default == "" {
				Config.Default = name
			}

			// wait for the config file update and verify profile is added or not
			rc := make(chan error, 1)
			// it's being watched in root initConfig - viper.WatchConfig()
			go viper.OnConfigChange(func(e fsnotify.Event) {
				if p, _ := Config.Profile(name); p == nil {
					rc <- fmt.Errorf("profile %v not added", name)
					return
				}
				fmt.Println("Profile", name, "added successfully:", e.Name)
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

	cmd.Flags().StringVarP(&flags.desc, "desc", "d", "", "description of profile")
	cmd.Flags().StringArrayVarP(&flags.env, "env", "e", []string{}, "'VAR=VAL' format of string")
	cmd.MarkFlagRequired("env")
	return cmd
}
