package cmd

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/sunggun-yu/proxy-wrapper/internal/config"
)

type editFlags struct {
	desc    string
	host    string // TODO: deprecate
	noproxy string // TODO: deprecate
	env     []string
}

func init() {
	rootCmd.AddCommand(editCommand())
}

// example of edit command
func cmdExampleEdit() string {
	return `
  # edit exiting proxy server profile
  prw edit my-proxy -p http://192.168.1.10:3128 -d "my proxy server" -n "127.0.0.1,localhost,something-1.com"
  prw edit my-proxy -n "127.0.0.1,localhost,something-1.com,something-2.com"
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

			// update host if input is not empty
			if flags.host != "" {
				profile.Host = flags.host
			}

			// update noproxy if input is not empty
			if flags.noproxy != "" {
				profile.NoProxy = flags.noproxy
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

	// set optional flag "profile". so that user can select proxy server without swithing proxy profile
	// selected profile by `use` command should be the profile if it is omitted
	cmd.Flags().StringVarP(&flags.host, "proxy", "p", "", "proxy host information with port number. e.g. http://<ip or domain>:<port number>")
	cmd.Flags().StringVarP(&flags.desc, "desc", "d", "", "description of proxy host profile")
	cmd.Flags().StringVarP(&flags.noproxy, "noproxy", "n", "127.0.0.1,localhost", "comma seperated string of domains and ip addresses to be applied to no_proxy")
	cmd.Flags().StringSliceVarP(&flags.env, "env", "e", []string{}, "usage string")

	return cmd
}
