package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(useCommand())
}

// example of use command
func cmdExampleUse() string {
	return `
  # set default profile to profile-name
  envp use profile-name
  
  # env vars in the default profile will be set during command execution
  envp -- kubectl get pods
  `
}

// add command
func useCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "use profile-name",
		Short:        "Set default environment variable profile",
		SilenceUsage: true,
		Example:      cmdExampleUse(),
		Args: cobra.MatchAll(
			Arg0AsProfileName(),
			Arg0NotExistingProfile(),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			p := args[0]
			viper.Set(ConfigKeyDefaultProfile, p)
			if err := viper.WriteConfig(); err != nil {
				return fmt.Errorf("failed to updating the config file: %v", err.Error())
			}
			fmt.Println("Default profile is set to", viper.GetString(ConfigKeyDefaultProfile))
			return nil
		},
	}
	return cmd
}
