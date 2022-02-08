package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/sunggun-yu/envp/internal/config"
)

// flags struct for show command
type showFlags struct {
	export bool
}

func init() {
	rootCmd.AddCommand(showCommand())
}

// example of show command
func cmdExampleShow() string {
	return `
  envp show
  envp show profile-name
  envp show --export
  envp show -e
  `
}

// showCommand prints out all the environment variables of profile
func showCommand() *cobra.Command {
	var flags showFlags
	var profileName string

	cmd := &cobra.Command{
		Use:               "show profile-name [flags]",
		Short:             "Print the environment variables of profile",
		SilenceUsage:      true,
		ValidArgsFunction: ValidArgsProfileList,
		Example:           cmdExampleShow(),
		RunE: func(cmd *cobra.Command, args []string) error {

			if len(args) == 0 {
				profileName = viper.GetString(ConfigKeyDefaultProfile)
			} else {
				profileName = args[0]
			}

			var profile config.Profile
			// validate selected profile
			selected := configProfiles.Sub(profileName)
			if selected == nil {
				return fmt.Errorf("profile %v is not existing", profileName)
			}

			// unmarshal into Profile
			err := selected.Unmarshal(&profile)
			if err != nil {
				return fmt.Errorf("profile '%v' malformed configuration %e", profile, err)
			}

			if flags.export {
				fmt.Println("# you can export env vars of profile with following command")
				fmt.Println("# eval $(envp show --export)")
				fmt.Println("# eval $(envp show profile-name --export)")
				fmt.Println("")
			}
			for _, e := range profile.Env {
				if flags.export {
					fmt.Print("export ")
				}
				fmt.Println(fmt.Sprint(e.Name, "=", e.Value))
			}
			return nil
		},
	}

	cmd.Flags().BoolVarP(&flags.export, "export", "e", false, "show env vars with export command")
	return cmd
}
