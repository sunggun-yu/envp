package cmd

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/sunggun-yu/envp/internal/config"
)

type editFlags struct {
	desc string
	env  []string
}

func init() {
	rootCmd.AddCommand(editCommand())
}

// example of edit command
func cmdExampleEdit() string {
	return `
  envp edit my-proxy \
    -d "updated profile desc" \
    -e "NO_PROXY=127.0.0.1,localhost"
  `
}

func editCommand() *cobra.Command {
	var flags editFlags

	cmd := &cobra.Command{
		Use:          "edit [profile-name-with-no-space] [flags]",
		Short:        "edit profile",
		SilenceUsage: true,
		Example:      cmdExampleEdit(),
		Args: cobra.MatchAll(
			Arg0AsProfileName(),
			Arg0NotExistingProfile(),
		),
		RunE: func(cmd *cobra.Command, args []string) error {

			profileName := args[0]
			var profile config.Profile

			// get entire profiles to update it properly
			sub := viper.Sub(ConfigKeyProfile)

			// validate selected profile
			selected := sub.Sub(profileName)
			// unmarshal into Profile
			err := selected.Unmarshal(&profile)
			if err != nil {
				return fmt.Errorf("profile '%v' malformed configuration %e", profile, err)
			}

			// update desc if input is not empty
			if flags.desc != "" {
				profile.Desc = flags.desc
			}

			// update env
			// parse flag.env into a map for easy checking
			menv := config.ParseEnvFlagToMap(flags.env)
			if menv != nil {
				// loop profile.Env and check if flag.env has updated value(exist)
				for _, e := range profile.Env {
					if _, exist := menv[e.Name]; !exist {
						menv[e.Name] = e.Value
					}
				}
				profile.Env = config.MapToEnv(menv)
			}

			// set updated profile to sub - profiles
			sub.Set(profileName, profile)

			// overwrite the profile
			viper.Set(ConfigKeyProfile, sub.AllSettings())
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
				fmt.Println("profile", profileName, "updated successfully:", e.Name)
				rc <- nil
			})

			if err := viper.WriteConfig(); err != nil {
				return err
			}
			// wait for profile validation channel
			err = <-rc
			if err != nil {
				return err
			}
			return nil
		},
	}

	// set optional flag "profile". so that user can select profile without swithing profile
	// selected profile by `use` command should be the profile if it is omitted
	cmd.Flags().StringVarP(&flags.desc, "desc", "d", "", "description of profile")
	cmd.Flags().StringSliceVarP(&flags.env, "env", "e", []string{}, "usage string")

	return cmd
}
