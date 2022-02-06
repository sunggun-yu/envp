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
    -d "profile desc" \
    -e HTTPS_PROXY=http://some-proxy:3128 \
    -e "NO_PROXY=127.0.0.1,localhost" \
    -e "DOCKER_HOST=ssh://myuser@some-server"
  `
}

// add command
func addCommand() *cobra.Command {
	// add flags
	var flags addFlags

	// add command
	cmd := &cobra.Command{
		Use:          "add [profile-name-with-no-space] [flags]",
		Short:        "add profile",
		SilenceUsage: true,
		Example:      cmdExampleAdd(),
		Args: cobra.MatchAll(
			Arg0AsProfileName(),
			Arg0ExistingProfile(),
		),
		RunE: func(cmd *cobra.Command, args []string) error {

			profileName := args[0]
			profile := config.Profile{
				Desc: flags.desc,
				Env:  []config.Env{},
			}
			profile.Env = config.ParseEnvFlagToEnv(flags.env)

			// get new viper instance to add item properly
			sub := viper.Sub(ConfigKeyProfile)
			sub.Set(profileName, profile)
			// overwrite the entire profiles
			viper.Set(ConfigKeyProfile, sub.AllSettings())

			// TODO: study viper more. watch may not needed if viper.WriteConfig() reloads config after writing file.
			// watch config changes
			viper.WatchConfig()
			// wait for the config file update and verify profile is added or not
			rc := make(chan error, 1)
			viper.OnConfigChange(func(e fsnotify.Event) {
				// assuming
				if viper.Sub(ConfigKeyProfile).Get(profileName) == nil {
					rc <- fmt.Errorf("profile %v not added", profileName)
					return
				}
				fmt.Println("profile", profileName, "added successfully:", e.Name)
				rc <- nil
			})

			if err := viper.WriteConfig(); err != nil {
				return err
			}
			// wait for profile validation channel
			err := <-rc
			if err != nil {
				return err
			}
			return nil
		},
	}

	// set optional flag "profile". so that user can select profile without swithing it
	// selected profile by `use` command should be the profile if it is omitted
	cmd.Flags().StringVarP(&flags.desc, "desc", "d", "", "description of profile")
	cmd.Flags().StringSliceVarP(&flags.env, "env", "e", []string{}, "usage string")

	return cmd
}
