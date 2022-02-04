package cmd

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/sunggun-yu/proxy-wrapper/internal/config"
)

var profileToEdit config.ProxyProfile

func init() {
	// set optional flag "profile". so that user can select proxy server without swithing proxy profile
	// selected profile by `use` command should be the profile if it is omitted
	editCmd.Flags().StringVarP(&profileToEdit.Host, "proxy", "p", "", "proxy host information with port number. e.g. http://<ip or domain>:<port number>")
	editCmd.Flags().StringVarP(&profileToEdit.Desc, "desc", "d", "", "description of proxy host profile")
	editCmd.Flags().StringVarP(&profileToEdit.NoProxy, "noproxy", "n", "127.0.0.1,localhost", "comma seperated string of domains and ip addresses to be applied to no_proxy")
	rootCmd.AddCommand(editCmd)
}

var editCmd = &cobra.Command{
	Use:          "edit [profile-name-with-no-space] [flags]",
	Short:        "edit profile",
	SilenceUsage: true,
	Example: `
  # edit exiting proxy server profile
  prw edit my-proxy -p http://192.168.1.10:3128 -d "my proxy server" -n "127.0.0.1,localhost,something-1.com"
  prw edit my-proxy -n "127.0.0.1,localhost,something-1.com,something-2.com"
	`,
	Args: cobra.MatchAll(
		Arg0AsProfileName(),
		Arg0NotExistingProfile(),
	),
	RunE: func(cmd *cobra.Command, args []string) error {
		p := args[0]
		var c map[string]config.ProxyProfile
		ps := viper.Sub(ConfigKeyProfile)
		if err := ps.Unmarshal(&c); err != nil {
			return err
		}

		// copy from existing desc if input is empty
		if profileToEdit.Desc == "" {
			profileToEdit.Desc = c[p].Desc
		}

		// copy from existing host if input is empty
		if profileToEdit.Host == "" {
			profileToEdit.Host = c[p].Host
		}

		// copy from existing noproxy if input is empty
		if profileToEdit.NoProxy == "" {
			profileToEdit.NoProxy = c[p].NoProxy
		}

		// overwrite the profile
		viper.Set(fmt.Sprintf("%v.%v", ConfigKeyProfile, p), profileToEdit)
		// watch config changes
		viper.WatchConfig()
		// wait for the config file update and verify profile is added or not
		rc := make(chan error, 1)
		// I think underlying of viper.OnConfiChange is goroutine. but just run it as goroutine just in case
		go viper.OnConfigChange(func(e fsnotify.Event) {
			// assuming
			if viper.Sub(ConfigKeyProfile).Get(p) == nil {
				rc <- fmt.Errorf("profile %v not added", p)
				return
			}
			fmt.Println("profile", p, "updated successfully:", e.Name)
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
