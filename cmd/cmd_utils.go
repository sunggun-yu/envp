package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/sunggun-yu/envp/internal/config"
)

// CurrentProfile is function that returns config.Profile
// it checks args and return default profile if args has no profile name
// otherwise speicified profile will be return
func currentProfile(cfg *config.Config, args []string) (profile *config.NamedProfile, err error) {
	switch {
	case len(args) > 0:
		profile, err = cfg.Profile(args[0])
	default:
		profile, err = cfg.DefaultProfile()
	}
	return profile, err
}

// print command example
func printExample(cmd *cobra.Command) {
	fmt.Println("Example:")
	fmt.Println(cmd.Example)
}

// check the error type and print out command help for specific error types
func checkErrorAndPrintCommandExample(cmd *cobra.Command, err error) {
	switch err.(type) {
	case *config.DefaultProfileNotSetError:
		printExample(cmd)
	case *config.ProfileNameInputEmptyError:
		printExample(cmd)
	}
}
