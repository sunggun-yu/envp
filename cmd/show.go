package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
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

	cmd := &cobra.Command{
		Use:               "show profile-name [flags]",
		Short:             "Print the environment variables of profile",
		SilenceUsage:      true,
		ValidArgsFunction: ValidArgsProfileList,
		Example:           cmdExampleShow(),
		RunE: func(cmd *cobra.Command, args []string) error {

			profile, err := currentProfile(args)
			if err != nil {
				checkErrorAndPrintCommandExample(cmd, err)
				return err
			}

			fmt.Println("# profile:", profile.Name)
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
