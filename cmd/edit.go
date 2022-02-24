package cmd

import (
	"github.com/spf13/cobra"
	"github.com/sunggun-yu/envp/internal/config"
)

// flags struct for edit command
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
    -d 'updated profile desc' \
    -e 'NO_PROXY=127.0.0.1,localhost'
  `
}

// editCommand edit/update environment variable profile and it's envionment variables in the config file
func editCommand() *cobra.Command {
	var flags editFlags

	cmd := &cobra.Command{
		Use:          "edit profile-name [flags]",
		Short:        "Edit environment variable profile",
		SilenceUsage: true,
		Example:      cmdExampleEdit(),
		Args: cobra.MatchAll(
			arg0AsProfileName(),
			arg0NotExistingProfile(),
		),
		ValidArgsFunction: validArgsProfileList,
		RunE: func(cmd *cobra.Command, args []string) error {

			cfg, err := configFile.Read()
			if err != nil {
				return err
			}
			profile, err := currentProfile(cfg, args)
			if err != nil {
				checkErrorAndPrintCommandExample(cmd, err)
				return err
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

			if err := configFile.Save(); err != nil {
				return err
			}

			cmd.Println("Profile", profile.Name, "updated successfully")

			return nil
		},
	}

	cmd.Flags().StringVarP(&flags.desc, "desc", "d", "", "description of profile")
	cmd.Flags().StringArrayVarP(&flags.env, "env", "e", []string{}, "'VAR=VAL' format of string")

	return cmd
}
