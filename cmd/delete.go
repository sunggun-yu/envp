package cmd

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(deleteCmd)
}

var deleteCmd = &cobra.Command{
	Use:          "delete profile-name",
	Long:         "Delete profile from config. please set default profile in case you delete default profile",
	Aliases:      []string{"del"},
	SilenceUsage: true,
	Example: `
  # delete profile
  prw delete my-proxy
  prw del another-profile
	`,
	Args: cobra.MatchAll(
		Arg0AsProfileName(),
		Arg0NotExistingProfile(),
	),
	RunE: func(cmd *cobra.Command, args []string) error {
		p := args[0]

		// use built-in function to delete key(profile) from map (profiles)
		delete(viper.Get(ConfigKeyProfile).(map[string]interface{}), p)

		// set default="" if default profile is being deleted
		if p == viper.GetString(ConfigKeyDefaultProfile) {
			viper.Set(ConfigKeyDefaultProfile, "")
			fmt.Println("WARN: Deleting default profile. please set default profile once it is deleted")
		}

		// TODO: study viper more. watch may not needed if viper.WriteConfig() reloads config after writing file.
		// watch config changes
		viper.WatchConfig()
		// wait for the config file update and verify profile is added or not
		rc := make(chan error, 1)
		// I think underlying of viper.OnConfiChange is goroutine. but just run it as goroutine just in case
		go viper.OnConfigChange(func(e fsnotify.Event) {
			if viper.Sub(ConfigKeyProfile).Get(p) != nil {
				rc <- fmt.Errorf("profile %v not deleted", p)
				return
			}
			fmt.Println("Profile", p, "deleted successfully:", e.Name)
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
