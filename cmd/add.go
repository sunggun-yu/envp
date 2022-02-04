package cmd

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/sunggun-yu/proxy-wrapper/internal/config"
)

// flags struct for add command
type addFlags struct {
	desc    string
	host    string // TODO: deprecate
	noproxy string // TODO: deprecate
	env     []string
}

func init() {
	rootCmd.AddCommand(addCommand())
}

// example of add command
func cmdExampleAdd() string {
	return `
  # add new proxy server profile
  prw add my-proxy -p http://192.168.1.10:3128 -d "my proxy server" -n "127.0.0.1,localhost,something.com"
  
  # proxy server of my-proxy profile will be set for executing command
  prw -- kubectl get pods
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
				Desc:    flags.desc,
				Host:    flags.host,    // TODO: deprecate
				NoProxy: flags.noproxy, // TODO: deprecate
				Env:     []config.Env{},
			}
			profile.Env = config.ParseEnvFlagToEnv(flags.env)

			viper.Set(fmt.Sprintf("%v.%v", ConfigKeyProfile, profileName), profile)

			// TODO: study viper more. watch may not needed if viper.WriteConfig() reloads config after writing file.
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

	// set optional flag "profile". so that user can select proxy server without swithing proxy profile
	// selected profile by `use` command should be the profile if it is omitted
	cmd.Flags().StringVarP(&flags.host, "proxy", "p", "", "proxy host information with port number. e.g. http://<ip or domain>:<port number>")                    // TODO: deprecate
	cmd.MarkFlagRequired("proxy")                                                                                                                                 // TODO: deprecate
	cmd.Flags().StringVarP(&flags.noproxy, "noproxy", "n", "127.0.0.1,localhost", "comma seperated string of domains and ip addresses to be applied to no_proxy") // TODO: deprecate
	cmd.Flags().StringVarP(&flags.desc, "desc", "d", "", "description of proxy host profile")
	cmd.Flags().StringSliceVarP(&flags.env, "env", "e", []string{}, "usage string")

	return cmd
}
