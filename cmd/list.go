package cmd

import (
	"fmt"
	"sort"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/sunggun-yu/proxy-wrapper/internal/config"
)

func init() {
	rootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:     "list",
	Short:   "List all profiles",
	Aliases: []string{"ls"},
	RunE: func(cmd *cobra.Command, args []string) error {
		// current default profile name to compare
		defaultProfile := viper.GetString(ConfigKeyDefaultProfile)

		// build string array for profiles names to be sorted
		profiles, err := listProfiles()
		if err != nil {
			return err
		}

		if len(profiles) < 1 {
			fmt.Println("no profile is existing")
			return nil
		}

		// print profiles. mark default profile with *
		for _, p := range profiles {
			if p == defaultProfile {
				fmt.Println("*", p)
			} else {
				fmt.Println(" ", p)
			}
		}
		return nil
	},
}

// get profile names list
func listProfiles() ([]string, error) {
	// unmarshall profile sub section to get keys
	var c map[string]config.ProxyProfile

	// TODO: set the global var for profile sub section in the init of root command.
	ps := viper.Sub(ConfigKeyProfile)
	if ps == nil {
		return nil, nil
	}

	if err := ps.Unmarshal(&c); err != nil {
		return nil, err
	}

	// build string array for profiles names to be sorted
	profiles := []string{}
	for p := range c {
		profiles = append(profiles, p)
	}
	// sort the profiles
	sort.Strings(profiles)
	return profiles, nil
}
