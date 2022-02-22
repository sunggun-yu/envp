package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/sunggun-yu/envp/internal/config"
)

// flags struct for add command
type addFlags struct {
	desc string
	env  []string
}

// init
func init() {
	rootCmd.AddCommand(addCommand())
}

// example of add command
func cmdExampleAdd() string {
	return `
  envp add my-proxy \
    -d 'profile desc' \
    -e 'VAR=VAL' \
    -e HTTPS_PROXY=http://some-proxy:3128 \
    -e 'NO_PROXY=127.0.0.1,localhost' \
    -e 'DOCKER_HOST=ssh://myuser@some-server'
  `
}

// addCommand add/create environment variable profile and it's envionment variables in the config file
func addCommand() *cobra.Command {
	// add flags
	var flags addFlags

	// add command
	cmd := &cobra.Command{
		Use:          "add profile-name",
		Short:        "Add environment variable profile",
		SilenceUsage: true,
		Example:      cmdExampleAdd(),
		Args: cobra.MatchAll(
			arg0AsProfileName(),
			arg0ExistingProfile(),
		),
		RunE: func(cmd *cobra.Command, args []string) error {

			cfg, err := configFile.Read()
			if err != nil {
				return err
			}

			name := args[0]
			profile := config.Profile{
				Desc: flags.desc,
				Env:  []config.Env{},
			}
			profile.Env = config.ParseEnvFlagToEnv(flags.env)

			// set profile
			cfg.SetProfile(name, profile)

			// set profile as default profile if default is empty and no profile is existing
			if cfg.Default == "" {
				cfg.SetDefault(name)
			}

			if err := configFile.Save(); err != nil {
				return err
			}

			fmt.Println("Profile", name, "added successfully")

			return nil
		},
	}

	cmd.Flags().StringVarP(&flags.desc, "desc", "d", "", "description of profile")
	cmd.Flags().StringArrayVarP(&flags.env, "env", "e", []string{}, "'VAR=VAL' format of string")
	cmd.MarkFlagRequired("env")
	return cmd
}
