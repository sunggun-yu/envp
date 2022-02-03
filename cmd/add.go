package cmd

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/sunggun-yu/proxy-wrapper/internal/config"
)

var profileToAdd = config.ProxyProfile{
	// TODO: remove comment
	// prepare to accpet and export arbitrary environment variables. so that this tool doesn't need to be limited to proxy ;)
	Envs: map[string]string{},
}

func init() {
	// set optional flag "profile". so that user can select proxy server without swithing proxy profile
	// selected profile by `use` command should be the profile if it is omitted
	addCmd.Flags().StringVarP(&profileToAdd.Host, "proxy", "p", "", "proxy host information with port number. e.g. http://<ip or domain>:<port number>")
	addCmd.MarkFlagRequired("proxy")
	addCmd.Flags().StringVarP(&profileToAdd.Description, "desc", "d", "", "description of proxy host profile")
	addCmd.Flags().StringVarP(&profileToAdd.NoProxy, "noproxy", "n", "127.0.0.1,localhost", "comma seperated string of domains and ip addresses to be applied to no_proxy")
	rootCmd.AddCommand(addCmd)
}

var addCmd = &cobra.Command{
	Use:          "add [profile-name-with-no-space] [flags]",
	Short:        "add profile",
	SilenceUsage: true,
	Example: `
  # add new proxy server profile
  prw add my-proxy -p http://192.168.1.10:3128 -d "my proxy server" -n "127.0.0.1,localhost,google.com"
  
  # proxy server of my-proxy profile will be set for executing command
  prw -- kubectl get pods
	`,
	Args: cobra.MatchAll(
		Arg0AsProfileName(),
		Arg0ExistingProfile(),
	),
	RunE: func(cmd *cobra.Command, args []string) error {
		p := args[0]
		viper.Set(fmt.Sprintf("%v.%v", ConfigKeyProfile, p), profileToAdd)

		// TODO: study viper more. watch may not needed if viper.WriteConfig() reloads config after writing file.
		// watch config changes
		viper.WatchConfig()
		// wait for the config file update and verify profile is added or not
		rc := make(chan error, 1)
		// I think underlying of viper.OnConfiChange is goroutine. but just run it as goroutine just in case
		go viper.OnConfigChange(func(e fsnotify.Event) {
			// assuming
			if viper.Sub(ConfigKeyProfile).Get(p) == nil {
				rc <- fmt.Errorf("profile %v not added", p)
			}
			fmt.Println("profile", p, "added successfully:", e.Name)
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
