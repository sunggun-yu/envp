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

			var c map[string]config.Profile
			ps := viper.Sub(ConfigKeyProfile)
			if err := ps.Unmarshal(&c); err != nil {
				return err
			}

			// corresponding profile object in the profiles
			// existence of profile has been guaranteed by arg validation. but nil checking might be good to have.
			profile := c[profileName]

			// update desc if input is not empty
			if flags.desc != "" {
				profile.Desc = flags.desc
			}

			menv := config.ParseEnvFlagToMap(flags.env)
			if menv != nil {
				for _, e := range profile.Env {
					if _, exist := menv[e.Name]; !exist {
						menv[e.Name] = e.Value
					}
				}
				profile.Env = config.MapToEnv(menv)
			}

			// overwrite the profile
			viper.Set(fmt.Sprintf("%v.%v", ConfigKeyProfile, profileName), profile)
			// watch config changes
			viper.WatchConfig()
			// wait for the config file update and verify profile is added or not
			rc := make(chan error, 1)
			// I think underlying of viper.OnConfiChange is goroutine. but just run it as goroutine just in case
			go viper.OnConfigChange(func(e fsnotify.Event) {
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
			err := <-rc
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
